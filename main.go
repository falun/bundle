package main

import (
	"flag"
	"log"
	"path/filepath"

	"github.com/falun/gobundle/app"
	"github.com/falun/gobundle/manifest"
)

func main() {
	var (
		pkg  = "bundle"
		dest = "./bundle"
	)

	flag.StringVar(&pkg, "package", pkg, "the generated code will be part of this package")
	flag.StringVar(&dest, "dest", dest, "produce generated code into this directory")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		log.Fatalf("manifest file expected")
	}
	mfPath := args[0]
	mf := manifest.Manifest{}

	mustReadInto(mfPath, &mf)

	pwd := mustGetPwd()

	dest = filepath.Join(pwd, dest)
	dirMustExist(dest)

	app.Main(pwd, dest, pkg, mf)
}
