package main

import (
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestPrepareTemplate checks all used templates for
// valid go syntax.
func TestPrepareTemplate(t *testing.T) {
	_, err := PrepareTemplate(Form)
	if err != nil {
		t.Fatalf("Failed to prepare 'Form': %v", err)
	}

	_, err = PrepareTemplate(SuccessTmpl)
	if err != nil {
		t.Fatalf("Failed to prepare 'SuccessTmpl': %v", err)
	}

	_, err = PrepareTemplate(FailureTmpl)
	if err != nil {
		t.Fatalf("Failed to prepare 'FailureTmpl': %v", err)
	}

	_, err = PrepareTemplate(EmailFormTmpl)
	if err != nil {
		t.Fatalf("Failed to prepare 'EmailFormTmpl': %v", err)
	}

	_, err = PrepareTemplate(EmailSubmitTmpl)
	if err != nil {
		t.Fatalf("Failed to prepare 'EmailSubmitTmpl': %v", err)
	}

	_, err = PrepareTemplate(EmailFailTmpl)
	if err != nil {
		t.Fatalf("Failed to prepare 'EmailFailTmpl': %v", err)
	}
}

func TestSuccess(t *testing.T) {
	// test parsing valid success page content
	pDat := &BCPoster{
		Session:        "session",
		AbstractNumber: "absnum",
		Authors:        "auth",
		Title:          "title",
		Topic:          "topic",
		ID:             "id",
		UploadKey:      "upkey",
		Abstract:       "text",
	}
	valDat := map[string]interface{}{
		"UserData":          pDat,
		"PDFPath":           "posterPath",
		"VideoURL":          "videoURL",
		"PosterHash":        "posterHash",
		"supportemail":      "supportEmail",
		"conferencepageurl": "conferencePageURL",
	}
	w := httptest.NewRecorder()
	success(w, valDat)
	if w.Result().StatusCode != 200 {
		t.Fatalf("Invalid header on success page: %v", w.Result().StatusCode)
	}
	// check that the content has been properly parsed into the template
	res := w.Body.String()
	contentCheck := [...]string{
		pDat.Session, pDat.AbstractNumber, pDat.Authors,
		pDat.Title, pDat.Topic, pDat.ID, pDat.Abstract, 
		"posterPath", "videoURL", "posterHash",
		"supportEmail", "conferencePageURL",
	}
	containsAll := true
	var missing string
	for _, item := range contentCheck {
		if !strings.Contains(res, item) {
			containsAll = false
			missing = fmt.Sprintf("%s %s", missing, item)
		}
	}
	if !containsAll {
		t.Fatalf("Success page content is missing: %s", missing)
	}

	// test parsing empty success page content
	w = httptest.NewRecorder()
	success(w, map[string]interface{}{})
	if w.Result().StatusCode != 200 {
		t.Fatalf("Invalid header on success page: %v", w.Result().StatusCode)
	}

	// test parsing invalid success page content
	invalDat := map[string]interface{}{
		"something":  "completely different",
		"PosterHash": nil,
	}
	w = httptest.NewRecorder()
	success(w, invalDat)
	if w.Result().StatusCode != 200 {
		t.Fatalf("Invalid header on success page: %v", w.Result().StatusCode)
	}
}
