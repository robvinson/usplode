package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"path"
	"strings"
)

func main() {
	// Command-line flags
	filePath := flag.String("file", "", "Path to the file containing URLs (optional)")
	depth := flag.Int("depth", 1, "Depth of paths to extract")
	flag.Parse()
	*depth += 1 // accomodate for host

	// Input source (file or stdin)
	var input io.Reader
	if *filePath != "" {
		file, err := os.Open(*filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
	} else {
		input = os.Stdin
	}

	scanner := bufio.NewScanner(input)
	seen := make(map[string]bool)

	for scanner.Scan() {
		u, err := url.Parse(scanner.Text())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Invalid URL: %v\n", err)
			continue
		}

		parts := strings.Split(u.Path, "/")
		if len(parts) > 0 && parts[len(parts)-1] == "" {
			parts = parts[:len(parts)-1] // Remove trailing slash part
		}

		for i := 1; i <= *depth && i <= len(parts); i++ {
			u.Path = path.Join(parts[:i]...) + "/"
			u.RawQuery = ""
			u.Fragment = ""
			dir := u.String()
			if !seen[dir] {
				fmt.Println(dir)
				seen[dir] = true
			}
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}
