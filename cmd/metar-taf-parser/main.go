package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	metartafparser "github.com/ryansavara/go-metar-taf-parser"
)

var errInvalidInput = errors.New("invalid input")

func main() {
	input := strings.Join(os.Args[1:], " ")
	if input == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading stdin:", err)
			os.Exit(1)
		}
		input = strings.TrimSpace(string(data))
	}

	if input == "" {
		fmt.Fprintln(os.Stderr, "No input provided")
		os.Exit(1)
	}

	result, err := parseAuto(input)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Invalid input: neither METAR nor TAF could be parsed")
		os.Exit(1)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	err = enc.Encode(result)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error encoding JSON:", err)
		os.Exit(1)
	}
}

var tafValidityRe = regexp.MustCompile(`\b\d{4}/\d{4}\b`)

func hasTAFValidityPeriod(input string) bool {
	return tafValidityRe.MatchString(input)
}

func parseAuto(input string) (any, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" || trimmed != strings.ToUpper(trimmed) {
		return nil, errInvalidInput
	}

	tokens := strings.Fields(trimmed)
	if len(tokens) == 0 {
		return nil, errInvalidInput
	}

	firstToken := tokens[0]

	if firstToken == "TAF" {
		taf, err := metartafparser.ParseTAF(input, nil)
		if err == nil {
			return taf, nil
		}
		return nil, errInvalidInput
	}

	if firstToken == "METAR" || firstToken == "SPECI" {
		metar, err := metartafparser.ParseMetar(input, nil)
		if err == nil && isMetarValid(metar) {
			return metar, nil
		}
		return nil, errInvalidInput
	}

	metar, err := metartafparser.ParseMetar(input, nil)
	if err == nil && isMetarValid(metar) {
		if hasTAFValidityPeriod(trimmed) {
			if taf, tafErr := metartafparser.ParseTAF(input, nil); tafErr == nil {
				return taf, nil
			}
		}
		return metar, nil
	}

	taf, err := metartafparser.ParseTAF(input, nil)
	if err == nil {
		return taf, nil
	}

	return nil, errInvalidInput
}

func isMetarValid(m *metartafparser.Metar) bool {
	if m == nil {
		return false
	}
	if m.Nil != nil && *m.Nil {
		return true
	}
	if m.Canceled != nil && *m.Canceled {
		return true
	}
	if m.Wind != nil {
		return true
	}
	if m.Visibility != nil {
		return true
	}
	if len(m.Clouds) > 0 {
		return true
	}
	if len(m.WeatherConditions) > 0 {
		return true
	}
	return false
}
