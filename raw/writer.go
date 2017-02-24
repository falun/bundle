package raw

import (
	"io/ioutil"
	"log"

	"github.com/falun/gobundle/manifest"
)

func Save(files []manifest.File) map[string][]byte {
	out := map[string][]byte{}

	for _, file := range files {
		if file.Compress {
			continue
		}

		content, err := ioutil.ReadFile(file.Path)
		if err != nil {
			log.Fatal(err)
		}

		out[file.Path] = content
	}

	return out
}
