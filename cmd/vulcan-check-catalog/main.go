package main

import (
	"encoding/json"
	"flag"
	"io"
	"os"
	"strings"

	"github.com/adevinta/vulcan-check-catalog/pkg/helpers"
	"github.com/adevinta/vulcan-check-catalog/pkg/model"
	"github.com/adevinta/vulcan-check-catalog/pkg/processor"
	log "github.com/sirupsen/logrus"
)

const (
	OKExitCode = 0
	KOExitCode = 1
)

func main() {
	os.Exit(realMain(os.Stdout))
}

func realMain(out io.Writer) int {
	// Setup the logger.
	var logger = log.New()
	logger.SetOutput(out)
	logger.SetLevel(log.DebugLevel)
	logger.SetFormatter(&log.TextFormatter{})

	// Read and validate flags and arguments.
	var registryURL, tag, checktypesURLList, output string
	flag.StringVar(&registryURL, "registry-url", "", "Docker image registry base URL")
	flag.StringVar(&tag, "tag", "", "Docker image tag")
	flag.StringVar(&checktypesURLList, "checktypes-url-list", "", "Checktypes URL list")
	flag.StringVar(&output, "output", "", "Output file path")

	flag.Parse()

	if registryURL == "" {
		flag.PrintDefaults()
		logger.Errorf("-registry-url flag is mandatory")
		return KOExitCode
	}

	path := flag.Arg(0)
	if flag.Arg(0) == "" {
		logger.Errorf("path argument not provided")
		return KOExitCode
	}
	isDir, err := helpers.IsExistingDirOrFile(path)
	if err != err {
		logger.Errorf("path [%s] does not exist", path)
		return KOExitCode
	}
	if !isDir {
		logger.Errorf("path [%s] is not a directory", path)
		return KOExitCode
	}

	// Process.
	p := processor.New(
		path,
		registryURL,
		tag,
		logger,
	)
	ct, err := p.Run()
	if err != nil {
		logger.Errorf("check processor failed: %s", err)
		return KOExitCode
	}

	// Fetch and merge checktypes list.
	for _, ctURL := range strings.Split(checktypesURLList, ",") {
		if ctURL == "" {
			continue
		}
		ctr, err := model.FetchChecktypesFromURL(ctURL)
		if err != nil {
			logger.Errorf("can't fetch checks from %s URL: %s", ctURL, err)
			return KOExitCode
		}
		ct = model.MergeChecktypes(ct, ctr)
	}

	// Generate output.
	b, err := json.MarshalIndent(ct, "", "\t")
	if err != nil {
		logger.Errorf("can't marshal check catalog json: %s", err)
		return KOExitCode
	}
	if output == "" {
		out.Write(b)
		return OKExitCode
	}

	f, err := os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logger.Errorf("can't open file to write output: %s", err)
		return KOExitCode
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		logger.Errorf("can't write output to file: %s", err)
		return KOExitCode
	}

	return OKExitCode
}
