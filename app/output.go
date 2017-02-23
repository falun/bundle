package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"log"
	"path/filepath"

	"github.com/falun/bundle/manifest"
	"github.com/falun/bundle/raw"
	"github.com/falun/bundle/zip"
)

func getBundle(pkg string, files []manifest.File) string {
	compressedBytes := zip.Save(files)
	rawBytesMap := raw.Save(files)

	mkByteStr := func(bs []byte) string {
		byteStr := ""
		for i, b := range bs {
			if i != 0 {
				byteStr += ", "
			}
			byteStr += fmt.Sprintf("%d", b)
		}
		return byteStr
	}

	compressedStr := mkByteStr(compressedBytes)
	compressedComment := "// compressed is a compressed version of the following files:\n"
	for _, f := range files {
		if !f.Compress { continue }
		compressedComment += fmt.Sprintf("\t//  - %v\n", f.Path)
	}

	rawStr := ""
	if 0 != len(rawBytesMap) {
		rawStr = "\n"
	}
	for k, bs := range rawBytesMap {
		rawStr += fmt.Sprintf("\t\t%q: []byte{%v},", k, mkByteStr(bs))
	}

	return fmt.Sprintf(`package %s

const (
	%v	compressed = []byte{%v}

	// raw is a map from file path to contents
	raw = map[string][]byte{%v}
)
`, pkg, compressedComment, compressedStr, rawStr)
}

func WriteBundle(dest, pkg string, mf manifest.Manifest) {
	contents := getBundle(pkg, mf.Files)

	err := ioutil.WriteFile(
		filepath.Join(dest, "contents.go"),
		[]byte(contents),
		os.FileMode(644))

	if err != nil {
		log.Fatalf("unable to write bundle contents: %v", err.Error())
	}
}
