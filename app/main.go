package app

import (
	"github.com/falun/bundle/manifest"
)

func Main(pwd, dest, pkg string, mf manifest.Manifest) {
	WriteBundle(dest, pkg, mf)
	WriteLib(dest, pkg)
}
