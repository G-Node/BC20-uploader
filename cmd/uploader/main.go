package main

import (
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
	"strings"
	"text/template"

	"github.com/G-Node/tonic/templates"
	"github.com/G-Node/tonic/tonic/web"
	"gopkg.in/yaml.v2"
)

var (
	appname string = "uploader"
	build   string
	commit  string
	verstr  string
)

type Config struct {
	// Port to listen on
	Port uint16
	// Directory for saving uploaded files
	UploadDirectory string
	// File containing user info with passwords
	UsersFile string
	// True if video upload is enabled
	Videos bool
}

func defaultConfig() *Config {
	return &Config{
		Port:            3000,
		UploadDirectory: "uploads",
		UsersFile:       "userlist.json",
		Videos:          false,
	}
}

func readConfig(configFileName string) *Config {
	configFile, err := os.Open(configFileName)
	if err != nil {
		log.Printf("[os.Open] Error reading config file %q: %s", configFileName, err.Error())
		os.Exit(1)
	}
	defer configFile.Close()

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Printf("[ioutil.ReadAll] Error reading config file %q: %s", configFileName, err.Error())
		os.Exit(1)
	}

	config := new(Config)
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Printf("[yaml.Unmarshall] Error reading config file (%q): %s", configFileName, err.Error())
		os.Exit(1)
	}
	// create upload directory (if it doesn't exist)
	os.MkdirAll(config.UploadDirectory, 0777)
	return config
}

// writeConfig writes the default configuration values to the specified file.
func writeConfig(cfgFileName string) {
	// using fmt.Print for error messages here since it's run interactively and
	// the log-style formatting with timestamps makes it noisy.
	cfgYml, err := yaml.Marshal(defaultConfig())
	if err != nil {
		fmt.Printf("Error marshalling default config: %s\n", err.Error())
		os.Exit(1)
	}

	cfgFile, err := os.Create(cfgFileName)
	if err != nil {
		fmt.Printf("Error creating config file: %s\n", err.Error())
		os.Exit(1)
	}
	defer cfgFile.Close()

	if _, err := cfgFile.Write(cfgYml); err != nil {
		fmt.Printf("Error writing default config: %s\n", err.Error())
		os.Exit(1)
	}
}

func prompt(msg string) string {
	var response string
	fmt.Printf("%s: ", msg)
	fmt.Scanln(&response)
	return response
}

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

func NewUploader(cfg *Config) *Uploader {
	uploader := new(Uploader)
	uploader.Config = cfg
	// TODO: Use configured port when tonic.Web supports it
	srv := web.New()
	srv.Router.HandleFunc("/", uploader.renderForm).Methods("GET")
	srv.Router.HandleFunc("/submit", uploader.submit).Methods("POST")
	srv.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	uploader.Web = srv
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
	tmpl := template.New("layout")
	tmpl, err := tmpl.Parse(Layout)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal error: Please contact an administrator")
		return
	}
	tmpl, err = tmpl.Parse(Form)

	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "Internal error: Please contact an administrator")
		return
	}

	formOpts := map[string]interface{}{
		"videos": uploader.Config.Videos,
	}

	if err := tmpl.Execute(w, formOpts); err != nil {
		log.Printf("Failed to render form: %v", err)
	}
}

func (uploader *Uploader) submit(w http.ResponseWriter, r *http.Request) {
	log.Print("Submission received")
	err := r.ParseMultipartForm(1048576) // 1 MiB max mem
	if err != nil {
		// 500
		log.Printf("Failed to parse form: %v", err.Error())
		return
	}
	postValues := r.PostForm

	passcode := postValues.Get("passcode")
	if passcode == "" {
		// 401
		log.Printf("ERROR: empty passcode")
		return
	}
	user, err := uploader.getUserInfo(passcode)
	if err != nil {
		// Check error message if unauthorised or server error and return appropriate response
		log.Printf("ERROR: %v", err.Error())
		return
	}

	log.Printf("User %q", user.Name)

	fileBasename := fmt.Sprintf("%s_%s", user.Session, strings.ReplaceAll(user.Name, " ", "_")) // TODO: Sanitize names
	os.MkdirAll(uploader.Config.UploadDirectory, 0777)
	saveUploadedFile := func(file multipart.File, header *multipart.FileHeader) {
		ext := filepath.Ext(header.Filename)
		fname := fmt.Sprintf("%s%s", fileBasename, ext)
		log.Printf("Writing file %q", fname)
		if err := saveFile(file, filepath.Join(uploader.Config.UploadDirectory, fname)); err != nil {
			log.Printf("ERROR: %v", err.Error())
			return
		}
	}

	// Save poster pdf
	posterFile, posterHeader, err := r.FormFile("poster")
	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		return
	}
	saveUploadedFile(posterFile, posterHeader)

	if uploader.Config.Videos {
		// Save video file
		videoFile, videoHeader, err := r.FormFile("video")
		if err != nil {
			log.Printf("ERROR: %v", err.Error())
			return
		}
		saveUploadedFile(videoFile, videoHeader)
	}

	videoURL := r.PostForm.Get("video_url")
	if videoURL != "" {
		fname := fmt.Sprintf("%s.url", fileBasename)
		urlfile, err := os.Create(filepath.Join(uploader.Config.UploadDirectory, fname))
		if err != nil {
			log.Printf("ERROR: %v", err.Error())
			return
		}
		defer urlfile.Close()
		log.Printf("Writing file %q", fname)
		if _, err := urlfile.WriteString(videoURL); err != nil {
			log.Printf("ERROR: %v", err.Error())
			return
		}
	}
}

func (uploader *Uploader) getUserInfo(passcode string) (*BCUser, error) {
	users, err := loadUserList(uploader.Config.UsersFile)
	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		return nil, err
	}

	for _, user := range users {
		if user.Passcode == passcode {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("Passcode did not match")

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

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)

	tmpl := template.New("layout")
	tmpl, err := tmpl.Parse(templates.Layout)
	if err != nil {
		tmpl = template.New("content")
	}
	tmpl, err = tmpl.Parse(templates.Fail)
	if err != nil {
		w.Write([]byte(message))
		return
	}
	errinfo := struct {
		StatusCode int
		StatusText string
		Message    string
	}{
		status,
		http.StatusText(status),
		message,
	}
	if err := tmpl.Execute(w, &errinfo); err != nil {
		log.Printf("Error rendering fail page: %v", err)
	}
}

type BCUser struct {
	Name     string
	Session  string
	Title    string
	Passcode string
}

func loadUserList(fname string) ([]BCUser, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	users := make([]BCUser, 0, 100)
	if err := json.Unmarshal(fileData, &users); err != nil {
		return nil, err
	}

	return users, nil
}
