package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
)

func main() {
	f, err := ioutil.TempFile("", "foo.txt")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("f.Name():", f.Name())
	fmt.Println("filepath.Dir(f.Name()):", filepath.Dir(f.Name()))
	a, err := filepath.Abs(filepath.Dir(f.Name()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("filepath.Abs(filepath.Dir(f.Name())):", a)
}
