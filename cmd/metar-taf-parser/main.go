package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	metartafparser "github.com/ryansavara/metar-taf-parser"
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

func parseAuto(input string) (any, error) {
	trimmed := strings.TrimSpace(input)

	tryTAF := strings.HasPrefix(trimmed, "TAF") || strings.HasPrefix(trimmed, "taf")

	if tryTAF {
		taf, err := metartafparser.ParseTAF(input, nil)
		if err == nil {
			return taf, nil
		}
		metar, err := metartafparser.ParseMetar(input, nil)
		if err == nil && isMetarValid(metar) {
			return metar, nil
		}
		return nil, errInvalidInput
	}

	metar, err := metartafparser.ParseMetar(input, nil)
	if err == nil && isMetarValid(metar) {
		return metar, nil
	}

	taf, err := metartafparser.ParseTAF(input, nil)
	if err == nil {
		return taf, nil
	}

	if metar != nil {
		return nil, errInvalidInput
	}

	return nil, errInvalidInput
}

func isMetarValid(m *metartafparser.Metar) bool {
	if m == nil {
		return false
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
