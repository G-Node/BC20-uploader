package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"testing"
)

func writeTmpFile(newfile string) error {
	file, err := os.Create(newfile)
	if err != nil {
		return err
	}
	return file.Close()
}

func checkDirFiles(tmpDir string, numfiles int) error {
	files, err := ioutil.ReadDir(tmpDir)
	if err != nil {
		return fmt.Errorf("Error reading directory %v", err)
	}
	if len(files) != numfiles {
		return fmt.Errorf("Found invalid number of files: %d", len(files))
	}
	return err
}

func TestRenameExistingFiles(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test_bc_rename")
	if err != nil {
		t.Fatalf("Error creating tmp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// test no issue on no file
	newfile := filepath.Join(tmpDir, "newfile.md")
	renameExistingFiles(newfile, 2)
	err = checkDirFiles(tmpDir, 0)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = writeTmpFile(newfile)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}

	// test valid file rename
	renameExistingFiles(newfile, 2)
	err = checkDirFiles(tmpDir, 1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = writeTmpFile(newfile)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}

	// test valid number of files kept after multiple renames
	renameExistingFiles(newfile, 2)
	err = writeTmpFile(newfile)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	renameExistingFiles(newfile, 2)
	err = writeTmpFile(newfile)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	renameExistingFiles(newfile, 2)
	err = checkDirFiles(tmpDir, 1)
	if err != nil {
		t.Fatalf(err.Error())
	}

	// test no rename on different file
	newfile = filepath.Join(tmpDir, "newfile2.md")
	err = writeTmpFile(newfile)
	if err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	renameExistingFiles(newfile, 2)
	err = checkDirFiles(tmpDir, 2)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestSaveFile(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "test_bc_save")
	if err != nil {
		t.Fatalf("Error creating tmp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// prepare mock request file content
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	mockwriter, err := writer.CreateFormFile("test-file", "testfile.md")
	if err != nil {
		t.Fatalf("Error creating form file: %v", err)
	}
	filecontent := "[{'content':'test'}]"
	_, err = mockwriter.Write([]byte(filecontent))
	if err != nil {
		t.Fatalf("Error writing file content: %v", err)
	}
	// yes, I am ignoring this error
	_ = writer.Close()
	req, _ := http.NewRequest("POST", "https://nowhere.com/", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// test saveFile with form file data
	content, _, err := req.FormFile("test-file")
	if err != nil {
		t.Fatalf("Error fetching request content: %v", err)
	}
	topath := filepath.Join(tmpDir, "upload.md")
	err = saveFile(content, topath)
	if err != nil {
		t.Fatalf("Error saving request content: %v", err)
	}

	// Check file content
	cont, err := ioutil.ReadFile(topath)
	if err != nil {
		t.Fatalf("Error reading output file: %v", err)
	}
	if string(cont) != filecontent {
		t.Fatalf("Unexpected file content: %s", cont)
	}
}
