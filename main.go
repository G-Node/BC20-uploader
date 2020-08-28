package main

import (
	"encoding/json"
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
)

func main() {
	log.Println("Starting")
	srv := web.New()
	srv.Router.HandleFunc("/", renderForm).Methods("GET")
	srv.Router.HandleFunc("/submit", submit).Methods("POST")
	srv.Router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))
	srv.Start()
	log.Println("Ready to engage")
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt)
	<-sigchan
	srv.Stop()
}

func renderForm(w http.ResponseWriter, r *http.Request) {
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

	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("Failed to render form: %v", err)
	}
}

func submit(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Printf("Failed to parse form: %v", err)
	}
	postValues := r.PostForm

	for key, value := range postValues {
		log.Printf("%s: %s", key, value)
	}
	posterFile, posterHeader, err := r.FormFile("poster")
	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		return
	}
	log.Printf("Poster filename: %s", posterHeader.Filename)

	users, err := loadUserList("./userlist.json")
	if err != nil {
		log.Printf("ERROR: %v", err.Error())
		return
	}

	fname := ""
	for _, user := range users {
		if user.Passcode == postValues.Get("passcode") {
			fname = fmt.Sprintf("%s_%s.pdf", user.Session, strings.ReplaceAll(user.Name, " ", "_")) // TODO: Sanitize names
		}
	}
	if fname == "" {
		log.Print("Unauthorised!!!")
		return
	}
	os.MkdirAll("uploads", 0777)
	if err := saveFile(posterFile, filepath.Join("uploads", fname)); err != nil {
		log.Printf("ERROR: %v", err.Error())
		return
	}
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
