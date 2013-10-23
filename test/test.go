package test

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func loadFixtures(name string) []byte {
	var filename string
	pwd, _ := os.Getwd()
PathIterator:
	for {
		filename = pwd + "/test/fixtures/" + name + ".json"
		log.Println(filename)
		if _, err := os.Open(filename); err == nil {
			break PathIterator
		}
		if pwd == "/" {
			panic("Fixtures not found")
		}
		pwd = filepath.Dir(pwd)
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return b
}

func Load(name string, docs interface{}) {
	data := loadFixtures(name)
	err := json.Unmarshal(data, docs)
	if err != nil {
		panic(err)
	}
}
