package main

import (
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
