package main

import (
	"bufio"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/G-Node/tonic/tonic/web"
)

var (
	appname = "uploader"
	build   string
	commit  string
	verstr  string
)

func init() {
	if build == "" {
		build = "(dev)"
		commit = "???"
	}

	verstr = fmt.Sprintf("%s build %s [%s]", appname, build, commit)
}

// Uploader is the main service struct.
type Uploader struct {
	Config *Config
	Web    *web.Server
}

// NewUploader creates and initializes the Uploader server.
func NewUploader(cfg *Config) *Uploader {
	uploader := new(Uploader)
	uploader.Config = cfg
	// TODO: Use configured port when tonic.Web supports it
	srv := web.New()
	srv.Router.HandleFunc("/", uploader.renderForm).Methods("GET")
	srv.Router.HandleFunc("/submit", uploader.submit).Methods("POST")
	srv.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	srv.Router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(cfg.UploadDirectory))))
	srv.Router.HandleFunc("/uploademail", uploader.uploademail).Methods("GET")
	srv.Router.HandleFunc("/submitemail", uploader.submitemail).Methods("POST")
	uploader.Web = srv

	// Increase timeouts
	srv.Server.WriteTimeout = time.Minute * 10
	srv.Server.ReadTimeout = time.Minute
	srv.Server.IdleTimeout = time.Minute * 2
	return uploader
}

func main() {
	log.Print(verstr)
	help := flag.Bool("help", false, "help")
	writeConfigFlag := flag.Bool("write-config", false, "write default configuration to file (use --config to specify file location)")
	configFile := flag.String("config", "config", "config file")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if *writeConfigFlag {
		if _, err := os.Stat(*configFile); err == nil {
		PROMPT:
			for {
				resp := prompt(fmt.Sprintf("File %q already exists. Overwite? [yN]", *configFile))
				switch resp {
				case "": // no response; treat as no
					fallthrough
				case "N": // case insensitive match
					fallthrough
				case "n":
					fmt.Println("Cancelled")
					os.Exit(0)
				case "Y": // case insensitive match
					fallthrough
				case "y":
					break PROMPT
				default:
					continue
				}
			}

		}
		fmt.Printf("Writing default configuration to %q\n", *configFile)
		writeConfig(*configFile)
		os.Exit(0)
	}

	log.Printf("Loading configuration from %q", *configFile)
	config := readConfig(*configFile)
	log.Printf("%+v", config)
	uploader := NewUploader(config)
	uploader.Web.Start()
	log.Printf("Listening on port %d", config.Port)
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	uploader.Web.Stop()
}

func (uploader *Uploader) renderForm(w http.ResponseWriter, r *http.Request) {
	baseTemplateData := map[string]interface{}{
		"supportemail":      uploader.Config.SupportEmail,
		"conferencepageurl": uploader.Config.ConferencePageURL,
	}

	tmpl, err := PrepareTemplate(Form)
	if err != nil {
		failure(w, http.StatusInternalServerError, baseTemplateData, "Internal error: Please contact an administrator")
		return
	}

	submission := true
	if closedate, err := time.Parse("2006-01-02", uploader.Config.SubmissionClosedDate); err != nil {
		log.Println("Could not parse submission closing date; submission is open")
	} else {
		submission = time.Now().Before(closedate)
	}

	formOpts := map[string]interface{}{
		"submission":        submission,
		"videos":            uploader.Config.Videos,
		"viduploadurl":      uploader.Config.VideoUploadURL,
		"conferencepageurl": uploader.Config.ConferencePageURL,
		"supportemail":      uploader.Config.SupportEmail,
		"closedtext":        uploader.Config.SubmissionClosedText,
		"closedtextvid":     uploader.Config.SubmissionClosedVideoText,
	}

	if err := tmpl.Execute(w, formOpts); err != nil {
		log.Printf("Failed to render form: %v", err)
	}
}

func renameExistingFiles(path string, nversions int) {
	log.Printf("Checking for older versions of %s", path)
	ext := filepath.Ext(path)
	basename := strings.TrimSuffix(path, ext)

	fileExists := func(fname string) bool {
		if _, err := os.Stat(fname); err == nil {
			return true
		}
		return false
	}

	nthFilename := func(n int) string {
		return fmt.Sprintf("%s-v%d%s", basename, n, ext)
	}

	// delete ver = nversions - 1 (oldest to keep) if it exists
	oldestVer := nthFilename(nversions - 1)
	if fileExists(oldestVer) {
		log.Printf("Deleting old file %s", oldestVer)
		err := os.Remove(oldestVer)
		if err != nil {
			log.Printf("Error removing file %s: %s", oldestVer, err.Error())
		}
	}

	for n := nversions - 2; n > 0; n-- {
		// for each one that exists, move it up one version
		nthVer := nthFilename(n)
		if fileExists(nthVer) {
			nthPlusOne := nthFilename(n + 1)
			log.Printf("Renaming old file %s -> %s", nthVer, nthPlusOne)
			err := os.Rename(nthVer, nthPlusOne)
			if err != nil {
				log.Printf("Error renaming file %s: %s", nthVer, err.Error())
			}
		}
	}

	// check for base file (no version suffix)
	if fileExists(path) {
		oneVer := nthFilename(1)
		log.Printf("Renaming old file %s -> %s", path, oneVer)
		err := os.Rename(path, oneVer)
		if err != nil {
			log.Printf("Error renaming file %s: %s", path, err.Error())
		}
	}
}

func (uploader *Uploader) submit(w http.ResponseWriter, r *http.Request) {
	baseTemplateData := map[string]interface{}{
		"supportemail":      uploader.Config.SupportEmail,
		"conferencepageurl": uploader.Config.ConferencePageURL,
	}

	log.Print("Submission received")
	err := r.ParseMultipartForm(1048576) // 1 MiB max mem
	if err != nil {
		// 500
		log.Printf("Failed to parse form: %v", err.Error())
		failure(w, http.StatusInternalServerError, baseTemplateData, "An internal error occurred.")
		return
	}
	postValues := r.PostForm

	passcode := postValues.Get("passcode")
	if passcode == "" {
		// 401
		log.Printf("ERROR: empty passcode")
		failure(w, http.StatusUnauthorized, baseTemplateData, "Empty passcode")
		return
	}
	user, err := uploader.getUserInfo(passcode)
	if err != nil {
		// Check error message if unauthorised or server error and return appropriate response
		log.Printf("ERROR: %v", err.Error())
		failure(w, http.StatusUnauthorized, baseTemplateData, "Unauthorised: Incorrect passcode")
		return
	}

	log.Printf("User %q", user.Authors)

	fileBasename := user.ID
	err = os.MkdirAll(uploader.Config.UploadDirectory, 0777)
	if err != nil {
		log.Printf("ERROR handling upload directory: %v", err.Error())
		failure(w, http.StatusInternalServerError, baseTemplateData, "Poster upload failed")
		return
	}
	saveUploadedFile := func(file multipart.File, header *multipart.FileHeader) string {
		ext := filepath.Ext(header.Filename)
		fname := fmt.Sprintf("%s%s", fileBasename, ext)
		log.Printf("Writing file %q", fname)
		targetPath := filepath.Join(uploader.Config.UploadDirectory, fname)
		renameExistingFiles(targetPath, uploader.Config.KeepVersions)
		if err := saveFile(file, targetPath); err != nil {
			log.Printf("ERROR: %v", err.Error())
			failure(w, http.StatusInternalServerError, baseTemplateData, fmt.Sprintf("File upload (%s) failed", ext))
			return ""
		}
		return targetPath
	}

	// Save poster pdf
	posterFile, posterHeader, err := r.FormFile("poster")
	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		failure(w, http.StatusInternalServerError, baseTemplateData, "Poster upload failed")
		return
	}
	posterPath := saveUploadedFile(posterFile, posterHeader)
	posterHash, err := sha1File(posterPath)
	if err != nil {
		log.Printf("Failed to hash file upload %q: %s", posterPath, err.Error())
	}
	log.Printf("PDF file saved: %s (%s)", posterPath, posterHash)

	if uploader.Config.Videos {
		// Save video file
		// account for the case that file upload is not required and the formfile can be empty.
		// simply continue in this case.
		videoFile, videoHeader, err := r.FormFile("video")
		if err != nil && err != http.ErrMissingFile {
			log.Printf("ERROR: %v", err.Error())
			failure(w, http.StatusInternalServerError, baseTemplateData, "Video upload failed")
			return
		} else if err == http.ErrMissingFile {
			log.Print("No video provided")
		} else {
			videoPath := saveUploadedFile(videoFile, videoHeader)
			log.Printf("Video file saved: %s", videoPath)
		}
	}

	videoURL := r.PostForm.Get("video_url")
	if videoURL != "" {
		fname := fmt.Sprintf("%s.url", fileBasename)
		urlTargetPath := filepath.Join(uploader.Config.UploadDirectory, fname)
		renameExistingFiles(urlTargetPath, uploader.Config.KeepVersions)
		urlfile, err := os.Create(urlTargetPath)
		if err != nil {
			log.Printf("ERROR: %v", err.Error())
			failure(w, http.StatusInternalServerError, baseTemplateData, "Form submission failed")
			return
		}
		defer urlfile.Close()
		if _, err := urlfile.WriteString(videoURL); err != nil {
			log.Printf("ERROR: %v", err.Error())
			failure(w, http.StatusInternalServerError, baseTemplateData, "Form submission failed")
			return
		}
		log.Printf("URL file saved: %s (%s)", fname, videoURL)
	}

	submittedData := map[string]interface{}{
		"UserData":          user,
		"PDFPath":           posterPath,
		"VideoURL":          videoURL,
		"PosterHash":        posterHash,
		"supportemail":      uploader.Config.SupportEmail,
		"conferencepageurl": uploader.Config.ConferencePageURL,
	}
	success(w, submittedData)
}

func (uploader *Uploader) getUserInfo(key string) (*BCPoster, error) {
	users, err := loadUserList(uploader.Config.PostersInfoFile)
	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		return nil, err
	}

	for _, user := range users {
		if user.UploadKey == key {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("passcode did not match")
}

func (uploader *Uploader) uploademail(w http.ResponseWriter, r *http.Request) {
	baseTemplateData := map[string]interface{}{
		"supportemail":      uploader.Config.SupportEmail,
		"conferencepageurl": uploader.Config.ConferencePageURL,
	}

	tmpl, err := PrepareTemplate(EmailFormTmpl)
	if err != nil {
		log.Printf("Error rendering email form page: %v", err)
		failure(w, http.StatusInternalServerError, baseTemplateData, "Form cannot be displayed")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, &baseTemplateData)
	if err != nil {
		log.Printf("Error rendering email form page: %v", err)
		failure(w, http.StatusInternalServerError, baseTemplateData, "Form cannot be displayed")
	}
}

func (uploader *Uploader) submitemail(w http.ResponseWriter, r *http.Request) {
	var filename = uploader.Config.WhitelistFile
	var password = uploader.Config.WhitelistPW

	content := r.FormValue("content")
	pwd := r.FormValue("password")

	// Prepare minimal template information
	baseTemplateData := map[string]interface{}{
		"supportemail":      uploader.Config.SupportEmail,
		"conferencepageurl": uploader.Config.ConferencePageURL,
	}

	// In case of an invalid password redirect back to the upload form
	if pwd != password {
		log.Print("ERROR Invalid password received")
		failure(w, http.StatusUnauthorized, baseTemplateData, "Unauthorised: Incorrect password")
		return
	}
	log.Print("INFO Received whitelist email form")

	// Sanitize input and split on whitespaces, comma and semicolon
	rstring := regexp.MustCompile(`[\s,;]+`)
	sanstring := rstring.ReplaceAllString(content, " ")
	contentslice := strings.Split(sanstring, " ")

	mailmap := make(map[string]interface{})
	// The file is created below if it does not exist yet
	if _, err := os.Stat(filename); err == nil {
		// Read file lines to map for duplicate entry exclusion
		datafile, err := os.Open(filename)
		if err != nil {
			log.Printf("ERROR Could not open whitelist email file: '%v'", err.Error())
			failure(w, http.StatusInternalServerError, baseTemplateData, "Form submission failed")
			return
		}
		fileScanner := bufio.NewScanner(datafile)

		// Populate data map
		for fileScanner.Scan() {
			mailmap[fileScanner.Text()] = nil
		}

		// No defer close since the same file is opened again and truncated below
		err = datafile.Close()
		if err != nil {
			log.Printf("Error closing whitelist email file: %v", err)
		}
	}

	log.Printf("Loaded %d entries from existing file", len(mailmap))

	// Reconcile stored and new data; make sure new content is lower case before it is hashed
	for _, v := range contentslice {
		mailmap[sha1String(strings.ToLower(v))] = nil
	}

	log.Printf("New entries added. Total entries: %d", len(mailmap))

	// Truncate output file and write all data to it
	outfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ERROR Could not open outfile for writing: '%v'", err)
		failure(w, http.StatusInternalServerError, baseTemplateData, "Form submission failed")
		return
	}
	defer outfile.Close()

	for k := range mailmap {
		if k == "" {
			continue
		}
		_, err = fmt.Fprintln(outfile, k)
		if err != nil {
			log.Printf("ERROR Could not write content '%s' to whitelist email file: '%v'", k, err)
		}
	}

	tmpl, err := PrepareTemplate(EmailSubmitTmpl)
	if err != nil {
		log.Printf("Error rendering email submission page: %v", err)
		failure(w, http.StatusInternalServerError, baseTemplateData, "Form submission failed")
		return
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, &baseTemplateData)
	if err != nil {
		log.Printf("Error rendering email submission page: %v", err)
		failure(w, http.StatusInternalServerError, baseTemplateData, "Form submission failed")
	}
	log.Printf("Saved email hashes to %q", filename)
}

func saveFile(file multipart.File, target string) error {
	buf := make([]byte, 1024)
	outfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer outfile.Close()
	for {
		rn, err := file.Read(buf)
		if rn > 0 {
			// write the new bytes to file before checking EOF or error
			wn, err := outfile.Write(buf[:rn])
			if err != nil {
				return err
			}
			if wn != rn {
				log.Printf("Error: read %d bytes but wrote %d", rn, wn)
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error while reading uploaded file for %q: %v", target, err.Error())
			return err
		}
	}
	return nil
}

func sha1String(content string) string {
	hasher := sha1.New()
	_, err := io.WriteString(hasher, content)
	if err != nil {
		log.Printf("Error writing sha1String: %v", err)
	}
	hash := hasher.Sum(nil)
	encoded := hex.EncodeToString(hash[:])

	return encoded
}

func sha1File(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha1.New()

	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	hash := hasher.Sum(nil)
	encoded := hex.EncodeToString(hash[:])
	return encoded, nil

}

// BCPoster represents a conference poster item
type BCPoster struct {
	Session        string
	AbstractNumber string `json:"abstract_number"`
	Authors        string
	Title          string
	Topic          string
	ID             string
	UploadKey      string `json:"upload_key"`
	Abstract       string
}

func loadUserList(fname string) ([]BCPoster, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	users := make([]BCPoster, 0, 100)
	if err := json.Unmarshal(fileData, &users); err != nil {
		return nil, err
	}

	return users, nil
}
