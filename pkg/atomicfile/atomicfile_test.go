package atomicfile

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/kr/pretty"
)

func TestWrite(t *testing.T) {
	t.Parallel()

	const (
		existingContent = "bananas"
		newContent      = "more bananas"
	)

	fileExists := func(path string) bool {
		_, err := os.Stat(path)
		return !os.IsNotExist(err)
	}

	mustExist := func(t *testing.T, path string) {
		t.Helper()
		if !fileExists(path) {
			t.Fatalf("file %v not found, should exist", path)
		}
	}
	mustNotExist := func(t *testing.T, path string) {
		t.Helper()
		if fileExists(path) {
			t.Fatalf("file %v exists, should not exist", path)
		}
	}
	mustBeContent := func(t *testing.T, path, want string) {
		t.Helper()
		got, err := ioutil.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Fatalf("got content %q, want %q", got, want)
		}
	}

	tests := []struct {
		name string

		// fn is called allowing a test to perform actions on the txn (such as
		// commit, or abort) and test the result.
		fn func(t *testing.T, txn *Write)
	}{
		{
			name: "commit",
			fn: func(t *testing.T, txn *Write) {
				err := txn.Commit()
				if err != nil {
					t.Errorf("got %v, want %v", pretty.Sprint(err), pretty.Sprint(nil))
				}

				mustExist(t, txn.targetPath)
				mustNotExist(t, txn.File.Name())
				mustBeContent(t, txn.targetPath, newContent)
			},
		},
		{
			name: "abort",
			fn: func(t *testing.T, txn *Write) {
				err := txn.Abort()
				if err != nil {
					t.Errorf("got %v, want %v", pretty.Sprint(err), pretty.Sprint(nil))
				}

				mustExist(t, txn.targetPath)
				mustNotExist(t, txn.File.Name())
				mustBeContent(t, txn.targetPath, existingContent)
			},
		},
		{
			name: "commit, commit errors",
			fn: func(t *testing.T, txn *Write) {
				if err := txn.Commit(); err != nil {
					t.Errorf("got %v, want %v", pretty.Sprint(err), pretty.Sprint(nil))
				}
				if err := txn.Commit(); err == nil {
					t.Error("expected error, got nil")
				}

				mustExist(t, txn.targetPath)
				mustNotExist(t, txn.File.Name())
				mustBeContent(t, txn.targetPath, newContent)
			},
		},
		{
			name: "close errors",
			fn: func(t *testing.T, txn *Write) {
				if err := txn.Close(); err != errCloseDisallowed {
					t.Errorf("got %v, want %v", pretty.Sprint(err), pretty.Sprint(nil))
				}

				mustExist(t, txn.targetPath)
				mustExist(t, txn.File.Name())
				mustBeContent(t, txn.targetPath, existingContent)

				_ = txn.Abort()
			},
		},
		{
			name: "commit, write errors",
			fn: func(t *testing.T, txn *Write) {
				if err := txn.Commit(); err != nil {
					t.Errorf("got %v, want %v", pretty.Sprint(err), pretty.Sprint(nil))
				}
				if _, err := txn.Write([]byte{42}); err == nil {
					t.Error("expected error, got nil")
				}

				mustExist(t, txn.targetPath)
				mustNotExist(t, txn.File.Name())
				mustBeContent(t, txn.targetPath, newContent)
			},
		},
		{
			name: "abort, abort errors",
			fn: func(t *testing.T, txn *Write) {
				if err := txn.Abort(); err != nil {
					t.Errorf("got %v, want %v", pretty.Sprint(err), pretty.Sprint(nil))
				}
				if err := txn.Abort(); err == nil {
					t.Error("expected error, got nil")
				}

				mustExist(t, txn.targetPath)
				mustNotExist(t, txn.File.Name())
				mustBeContent(t, txn.targetPath, existingContent)
			},
		},
		{
			name: "abort, write errors",
			fn: func(t *testing.T, txn *Write) {
				if err := txn.Abort(); err != nil {
					t.Errorf("got %v, want %v", pretty.Sprint(err), pretty.Sprint(nil))
				}
				if _, err := txn.Write([]byte{42}); err == nil {
					t.Error("expected error, got nil")
				}

				mustExist(t, txn.targetPath)
				mustNotExist(t, txn.File.Name())
				mustBeContent(t, txn.targetPath, existingContent)
			},
		},
	}

	for _, tt := range tests {
		var tt = tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// Create a temporary file as the txn destination file
			targetFile, err := ioutil.TempFile("", "testfile-*.json")
			if err != nil {
				t.Fatal(err)
			}

			// Write some context and close the file
			if _, err := io.WriteString(targetFile, existingContent); err != nil {
				t.Fatal(err)
			}
			_ = targetFile.Close()

			// Initialise a new write txn for targetFile
			txn, err := NewWrite(targetFile.Name())
			if err != nil {
				t.Fatal(err)
			}

			// Write some content to the txn file
			if _, err := io.WriteString(txn, newContent); err != nil {
				t.Fatal(err)
			}

			mustExist(t, txn.targetPath)
			mustExist(t, txn.File.Name())
			mustBeContent(t, txn.targetPath, existingContent)

			tt.fn(t, txn)
		})
	}
}
