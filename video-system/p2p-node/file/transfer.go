package file

import (
	"io"
	"os"
)

func SendFile(w io.Writer, path string) error {

	f, err := os.Open(path)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(w, f)

	return err
}
