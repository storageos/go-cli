// Package atomicfile exports a mechanism for transactional writes to a file.
package atomicfile

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

var errCloseDisallowed = errors.New("cannot call close on transactional write file, must commit or abort")

// Write adds transactional write behaviour to an *os.File.
//
// It extends the method set to include Commit() and Abort(), but blocks direct
// calls to Close().
type Write struct {
	*os.File

	targetPath string
}

// NewWrite creates a temporary file in the same directory as the destination
// file.
//
// The same directory is used for the temp file because cross-filesystem
// renames are not atomic and it is possible that the OS temp directory is
// mounted on a different filesystem to targetPath.
//
// The returned Write may be used as a normal file, except it must be committed
// to save to targetPath once any writing is finished.
func NewWrite(targetPath string) (*Write, error) {
	targetDir, fileName := filepath.Split(targetPath)

	tempFile, err := ioutil.TempFile(targetDir, fileName+"-txn-*.json")
	if err != nil {
		return nil, err
	}

	return &Write{
		File: tempFile,

		targetPath: targetPath,
	}, nil
}

// Abort closes the staged transaction file, removing it from disk. This
// discards any data written to the transaction file.
func (txn *Write) Abort() error {
	err := txn.File.Close()
	if err != nil {
		return err
	}

	err = os.Remove(txn.File.Name())
	if err != nil {
		return err
	}

	return nil
}

// Commit flushes any pending writes to the transaction file, atomically
// writes it to the target file path and then closes it.
func (txn *Write) Commit() error {
	err := txn.File.Sync()
	if err != nil {
		return err
	}

	err = os.Rename(txn.File.Name(), txn.targetPath)
	if err != nil {
		return err
	}

	err = txn.File.Close()
	if err != nil {
		return err
	}

	return nil
}

// Close will always error. Allowing consumers of a transaction to
// independently close the file breaks the transactional semantics.
func (txn *Write) Close() error {
	return errCloseDisallowed
}
