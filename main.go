package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dlnilsson/2go/togo"
	"gopkg.in/yaml.v3"
)

func main() {
	var (
		nested   = flag.Bool("nested", false, "nested structs")
		scanner  = bufio.NewScanner(os.Stdin)
		inputStr strings.Builder
	)
	flag.Parse()

	for scanner.Scan() {
		inputStr.WriteString(scanner.Text() + "\n")
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading standard input:", err)
		os.Exit(1)
	}

	var (
		input      = inputStr.String()
		data       any
		formatType = "json"
	)

	if json.Valid([]byte(input)) {
		if err := json.Unmarshal([]byte(input), &data); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to parse JSON input:", err)
			os.Exit(1)
		}
	} else {
		decoder := yaml.NewDecoder(bytes.NewReader([]byte(input)))
		if err := decoder.Decode(&data); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to parse YAML input:", err)
			os.Exit(1)
		}
		formatType = "yaml"
	}

	goCode, err := togo.ConvertToGoStructs(data, *nested, formatType)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to generate Go code:", err)
		os.Exit(1)
	}

	fmt.Fprintln(os.Stdout, goCode)
}
