package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

// This is how we're able to read and write to different formats. The shared
// and standard format is Go's `interface{}` which is standard in Go's
// encoding packages.
type handler struct {
	Unmarshal func([]byte, interface{}) error
	Marshal   func(interface{}) ([]byte, error)
}

var (
	handlers map[string]handler
)

func usage() {
	fmt.Println("gota <input> <output>")
}

func ext(file string) string {
	return strings.TrimPrefix(filepath.Ext(file), ".")
}

func init() {
	handlers = map[string]handler{
		"json": handler{json.Unmarshal, json.Marshal},
		"yaml": handler{yaml.Unmarshal, yaml.Marshal},
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func isError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}

	return (err != nil)
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(2)
	}

	fileIn, fileOut := os.Args[1], os.Args[2]
	extIn, extOut := ext(fileIn), ext(fileOut)

	fmt.Printf("Converting %s file %s into %s and saving it to %s.\n",
		extIn, fileIn, extOut, fileOut)

	if extIn == "" {
		fmt.Println("Error: missing extension on input file")
	}

	if extOut == "" {
		fmt.Println("Error: missing extension on output file")
	}

	if extIn == "" || extOut == "" {
		os.Exit(2)
	}

	unmarshaler, validIn := handlers[extIn]
	marshaler, validOut := handlers[extOut]

	if !validIn {
		fmt.Printf("Error: invalid input format: %s\n", extIn)
	}

	if !validOut {
		fmt.Printf("Error: invalid output format: %s\n", extOut)
	}

	if !validIn || !validOut {
		os.Exit(2)
	}

	file, error := ioutil.ReadFile(fileIn)
	check(error)
	var storage map[string]interface{}
	err := unmarshaler.Unmarshal([]byte(string(file)), &storage)

	if err != nil {
		fmt.Println(fmt.Errorf("Error: unable to decode %s file %s: %v", extIn, fileIn, err))
		os.Exit(2)
	}
	contents, err := marshaler.Marshal(storage)

	if err != nil {
		fmt.Printf("Error: unable to encode %s to %s: %v", extIn, extOut, err)
		os.Exit(2)
	}
	
	
	writeErr := ioutil.WriteFile(fileOut, contents, 0644)
	check(writeErr)

	fmt.Println(string(contents))
}
