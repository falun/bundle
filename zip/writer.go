package zip

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"log"

	"github.com/falun/gobundle/manifest"
)

func Save(files []manifest.File) []byte {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	for _, file := range files {
		if !file.Compress {
			continue
		}

		f, err := w.Create(file.Path)
		if err != nil {
			log.Fatal(err)
		}

		content, err := ioutil.ReadFile(file.Path)
		if err != nil {
			log.Fatal(err)
		}

		_, err = f.Write(content)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

	return buf.Bytes()
}
