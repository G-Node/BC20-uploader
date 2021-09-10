package main

import (
	"log"
	"net/http"
	"text/template"
)

func success(w http.ResponseWriter, data map[string]interface{}) {
	tmpl, err := PrepareTemplate(SuccessTmpl)
	if err != nil {
		failure(w, http.StatusInternalServerError, data, "Submission success but error occured. Please contact...")
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := tmpl.Execute(w, &data); err != nil {
		failure(w, http.StatusInternalServerError, data, "Submission success but error occured. Please contact...")
		return
	}
}

func failure(w http.ResponseWriter, status int, data map[string]interface{}, message string) {
	tmpl, err := PrepareTemplate(FailureTmpl)
	if err != nil {
		w.Write([]byte(message))
		return
	}

	errData := map[string]interface{}{
		"Message": message,
	}
	// Handle conference page link and support email in the page header
	if data != nil {
		errData = data
		errData["Message"] = message
	}

	w.WriteHeader(status)
	if err := tmpl.Execute(w, &errData); err != nil {
		log.Printf("Error rendering fail page: %v", err)
		return
	}

}

func emailfailure(w http.ResponseWriter, status int, data map[string]interface{}, message string) {
	tmpl, err := PrepareTemplate(EmailFailTmpl)
	if err != nil {
		w.Write([]byte(message))
		return
	}

	errData := map[string]interface{}{
		"Message": message,
	}
	// Handle conference page link and support email in the page header
	if data != nil {
		errData = data
		errData["Message"] = message
	}

	w.WriteHeader(status)
	if err := tmpl.Execute(w, &errData); err != nil {
		log.Printf("Error rendering email fail page: %v", err)
		return
	}
}

func PrepareTemplate(contentTemplate string) (*template.Template, error) {
	tmpl := template.New("layout")
	tmpl, err := tmpl.Parse(Layout)
	if err != nil {
		return nil, err
	}
	return tmpl.Parse(contentTemplate)
}
