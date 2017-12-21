package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type handler struct {
	Unmarshal func([]byte, interface{}) error
	Marshal   func(interface{}) ([]byte, error)
}

func usage() {
	fmt.Println("gota <input> <output>")
}

func ext(file string) string {
	return strings.TrimPrefix(filepath.Ext(file), ".")
}

func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(2)
	}

	handlers := map[string]handler{
		"json": handler{json.Unmarshal, json.Marshal},
	}

	var storage map[string]interface{}

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

	// TODO read from the real file
	err := unmarshaler.Unmarshal([]byte(`{"hi": true}`), &storage)

	if err != nil {
		fmt.Errorf("Error: unable to decode %s file %s: %v", extIn, fileIn, err)
		os.Exit(2)
	}

	contents, err := marshaler.Marshal(storage)

	if err != nil {
		fmt.Errorf("Error: unable to encode %s to %s: %v", extIn, extOut, err)
		os.Exit(2)
	}

	// TODO write to real file
	fmt.Println(string(contents))
}
