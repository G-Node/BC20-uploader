package main

import (
	"fmt"
	"io/ioutil"
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
