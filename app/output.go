package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/falun/gobundle/manifest"
	"github.com/falun/gobundle/raw"
	"github.com/falun/gobundle/zip"
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
		if !f.Compress {
			continue
		}
		compressedComment += fmt.Sprintf("\t//  - %v\n", f.Path)
	}

	rawStr := ""
	for k, bs := range rawBytesMap {
		rawStr += fmt.Sprintf("\n\t\t%q: []byte{%v},\n", k, mkByteStr(bs))
	}
	if rawStr != "" {
		rawStr += "\t"
	}

	return fmt.Sprintf(`package %s

var (
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

func WriteLib(dest, pkg string) {
	libWithPkg := []byte(fmt.Sprintf(bundleLib, pkg))
	err := ioutil.WriteFile(filepath.Join(dest, "bundle.go"), libWithPkg, 644)
	if err != nil {
		log.Fatalf("unable to write bundle library: %v", err.Error())
	}
}

var bundleLib = `package %s

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
)

func writeFile(path string, contents []byte) error {
	err := os.MkdirAll(filepath.Dir(path), 744)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path, contents, 644)
	if err != nil {
		return err
	}

	return nil
}

func installCompressed(dir string) error {
	compressedBytes := bytes.NewReader(compressed)
	reader, err := zip.NewReader(compressedBytes, int64(compressedBytes.Len()))
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		c, err := file.Open()
		if err != nil {
			return err
		}

		contents, err := ioutil.ReadAll(c)
		if err != nil {
			return err
		}

		err = writeFile(filepath.Join(dir, file.Name), contents)
		if err != nil {
			return err
		}

		c.Close()
	}

	return nil
}

func installRaw(dir string) error {
	for path, contents := range raw {
		err := writeFile(filepath.Join(dir, path), contents)
		if err != nil {
			return err
		}
	}

	return nil
}

func Install(dir string) error {
	if err := installRaw(dir); err != nil {
		return err
	}
	if err := installCompressed(dir); err != nil {
		return err
	}

	return nil
}`
