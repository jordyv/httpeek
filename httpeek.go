package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/antchfx/htmlquery"
	flag "github.com/spf13/pflag"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	filePath   string
	silent     bool
	xpathQuery string
	timeout    time.Duration

	client *http.Client

	output = os.Stdout
	errOut = os.Stderr
)

type outputLine struct {
	URL        string `json:"url"`
	StatusCode int    `json:"status_code"`
	Result     string `json:"result"`
}

func initHTTPClient() {
	client = http.DefaultClient
	client.Timeout = timeout
}

func main() {
	flag.StringVarP(&filePath, "file", "f", "", "Input file to use instead of stdin")
	flag.StringVarP(&xpathQuery, "query", "q", "//title", "XPath query to lookup in HTML output")
	flag.BoolVarP(&silent, "silent", "s", false, "Only output actual results")
	flag.DurationVarP(&timeout, "timeout", "t", time.Second*3, "Timeout for HTTP requests")
	flag.ErrHelp = errors.New("")
	flag.Parse()

	if silent {
		errOut, _ = os.Open(os.DevNull)
	}

	initHTTPClient()

	var inputFile *os.File
	var err error
	if filePath != "" {
		inputFile, err = os.Open(filePath)
		if err != nil {
			logError("Could not open file - %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		inputFile = os.Stdin
	}

	defer inputFile.Close()

	buf := bufio.NewReader(inputFile)
	for {
		line, _, err := buf.ReadLine()
		if err != nil {
			if err != io.EOF {
				logError("Error while reading line: %s\n", err.Error())
				continue
			}
			break
		}

		lineStr := strings.TrimSpace(string(line))
		_, err = url.Parse(lineStr)
		if err != nil {
			logError("Invalid URL '%s'\n", lineStr)
		}

		result, err := parse(lineStr)
		if err == nil {
			s, _ := json.Marshal(result)
			logOutput("%s\n", s)
		}
	}
}

func parse(line string) (*outputLine, error) {
	resp, err := client.Get(line)
	if err == nil {
		doc, err := htmlquery.Parse(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("could not parse HTML for '%s' - %s\n", line, err.Error())
		}
		resultString := ""
		resultNode := htmlquery.FindOne(doc, xpathQuery)
		if resultNode != nil {
			resultString = fmt.Sprintf("%+v", htmlquery.InnerText(resultNode))
		}
		return &outputLine{
			URL:        line,
			StatusCode: resp.StatusCode,
			Result:     resultString,
		}, nil
	}
	return nil, fmt.Errorf("could not access '%s' - %s\n", line, err.Error())
}

func logOutput(format string, arguments ...interface{}) {
	fmt.Fprintf(output, format, arguments...)
}

func logError(format string, arguments ...interface{}) {
	fmt.Fprintf(errOut, format, arguments...)
}
