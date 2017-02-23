package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

func dirMustExist(dest string) {
	err := os.MkdirAll(dest, 0644)
	if err != nil {
		log.Fatalf("Unable to create %v: %v", dest, err.Error())
	}
}

func mustGetPwd() string {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("could not get working directory: %v", err)
	}

	return dir
}

func mustReadInto(path string, dest interface{}) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("unable to read file %v: %v", path, err.Error())
	}

	err = json.Unmarshal(b, dest)
	if err != nil {
		log.Fatalf("failure processing file %v: %v", path, err.Error())
	}
}
