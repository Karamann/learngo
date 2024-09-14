// Copyright © 2018 Inanc Gumus
// Learn Go Programming Course
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/
//
// For more tutorials  : https://learngoprogramming.com
// In-person training  : https://www.linkedin.com/in/inancgumus/
// Follow me on twitter: https://twitter.com/inancgumus

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

// Detect returns the files that have a valid header (file signature).
// A valid header is determined by the format.
func Detect(format string, filenames []string) (valids []string, err error) {
	defer func() {
		if rerr := recover(); rerr != nil {
			err = fmt.Errorf("cannot detect: %v", rerr)
			fmt.Println("Error:", err)
		}
	}()

	valids = detect(format, filenames)
	if len(valids) > 0 {
		fmt.Println("Valid files detected:", valids)
	} else {
		fmt.Println("No valid files found.")
	}
	return
}

func detect(format string, filenames []string) (valids []string) {
	header, err := headerOf(format)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	buf := make([]byte, len(header))

	for _, filename := range filenames {
		if read(filename, buf) != nil {
			continue
		}

		if bytes.Equal([]byte(header), buf) {
			valids = append(valids, filename)
		}
	}
	return
}

// headerOf returns the file signature (magic number) for a given format.
// If the format is unsupported, it returns an error.
func headerOf(format string) (string, error) {
	switch format {
	case "png":
		return "\x89PNG\r\n\x1a\n", nil
	case "jpg":
		return "\xff\xd8\xff", nil
	case "gif":
		return "GIF89a", nil
	case "pdf":
		return "%PDF-", nil
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// read reads len(buf) bytes to buf from a file.
func read(filename string, buf []byte) error {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", filename, "-", err)
		return err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		fmt.Println("Error getting file info:", filename, "-", err)
		return err
	}

	if fi.Size() < int64(len(buf)) {
		fmt.Printf("File %s is smaller than expected header size.\n", filename)
		return fmt.Errorf("file size < len(buf)")
	}

	_, err = io.ReadFull(file, buf)
	if err != nil {
		fmt.Println("Error reading file:", filename, "-", err)
	}
	return err
}

func main() {
	// Set up command-line flags.
	format := flag.String("format", "png", "File format to detect (png, jpg, gif, pdf)")
	flag.Parse()
	filenames := flag.Args()

	// Detect valid files.
	if len(filenames) == 0 {
		fmt.Println("Please provide at least one filename to check.")
		os.Exit(1)
	}

	_, err := Detect(*format, filenames)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
