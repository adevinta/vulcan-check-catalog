package main

import (
	"bytes"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
)

var (
	timeRegexp = regexp.MustCompile(`time="(.*?)"`)
	portRegexp = regexp.MustCompile(`:\d+/`)
)

func startHTTPTestServer() (string, *httptest.Server) {
	fs := http.FileServer(http.Dir(fmt.Sprintf("testdata/published")))
	s := httptest.NewServer(fs)
	publishedCheckCatalogURL := fmt.Sprintf("%s/published/", s.URL)
	return publishedCheckCatalogURL, s
}

func TestMain(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		flags          []string
		mergeFlag      bool
		mergePath      string
		wantExitCode   int
		outputPath     string
		wantOutputPath string
	}{
		{
			name:           "NoArgumentsNorFlags",
			wantExitCode:   1,
			wantOutputPath: "testdata/no-arguments-nor-flags.out",
		},
		{
			name:           "NoArguments",
			wantExitCode:   1,
			flags:          []string{"-registry-url", "https://example.com"},
			wantOutputPath: "testdata/no-arguments.out",
		},
		{
			name:           "WrongCheckFolder",
			wantExitCode:   1,
			flags:          []string{"-registry-url", "https://example.com"},
			args:           []string{"testdata/dummy"},
			wantOutputPath: "testdata/wrong-check-folder.out",
		},
		{
			name:           "NotExistingCheckFolder",
			wantExitCode:   1,
			flags:          []string{"-registry-url", "https://example.com"},
			args:           []string{"testdata/not-existing"},
			wantOutputPath: "testdata/not-existing-check-folder.out",
		},
		{
			name:           "HappyPath",
			wantExitCode:   0,
			flags:          []string{"-registry-url", "https://example.com"},
			args:           []string{"testdata/happy-path"},
			wantOutputPath: "testdata/happy-path.out",
		},
		{
			name:           "HappyPathWithTag",
			wantExitCode:   0,
			flags:          []string{"-registry-url", "https://example.com", "-tag", "2.5.3"},
			args:           []string{"testdata/happy-path"},
			wantOutputPath: "testdata/happy-path-with-tag.out",
		},
		{
			name:           "HappyPathToFile",
			wantExitCode:   0,
			flags:          []string{"-registry-url", "https://example.com", "-tag", "master", "-output", "testdata/output/happy-path.json"},
			args:           []string{"testdata/happy-path"},
			outputPath:     "testdata/output/happy-path.json",
			wantOutputPath: "testdata/happy-path-to-file.out",
		},
		{
			name:           "HappyPathWithMerge",
			wantExitCode:   0,
			flags:          []string{"-registry-url", "https://example.com", "-tag", "experimental"},
			mergeFlag:      true,
			mergePath:      "catalog.json",
			args:           []string{"testdata/happy-path"},
			wantOutputPath: "testdata/happy-path-with-merge.out",
		},
		{
			name:           "WrongMergeURL",
			wantExitCode:   1,
			flags:          []string{"-registry-url", "https://example.com", "-tag", "experimental"},
			mergeFlag:      true,
			mergePath:      "non-existing.json",
			args:           []string{"testdata/happy-path"},
			wantOutputPath: "testdata/wrong-merge-url.out",
		},
	}
	for _, tt := range tests {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()

		t.Run(tt.name, func(t *testing.T) {
			fs := http.FileServer(http.Dir(fmt.Sprintf("testdata/published")))
			s := httptest.NewServer(fs)
			defer s.Close()
			flag.CommandLine = flag.NewFlagSet(tt.name, flag.ExitOnError)
			flags := append([]string{tt.name}, tt.flags...)
			if tt.mergeFlag {
				flags = append(flags, "-checktypes-url-list", fmt.Sprintf("%s/%s", s.URL, tt.mergePath))
			}
			os.Args = append(flags, tt.args...)
			var buf bytes.Buffer
			gotExitCode := realMain(&buf)
			if gotExitCode != tt.wantExitCode {
				t.Errorf("unexpected exit code: got %v, want %v", gotExitCode, tt.wantExitCode)
			}
			// Remove time strings from test output.
			gotOutput := strings.TrimSpace(timeRegexp.ReplaceAllString(buf.String(), ""))
			// Remove port strings from test output.
			gotOutput = portRegexp.ReplaceAllString(gotOutput, "/")
			if tt.outputPath != "" {
				outputBytes, err := os.ReadFile(tt.outputPath)
				if err != nil {
					t.Errorf("unexpected error reading output file %s: %s", tt.outputPath, err)
				}
				defer os.RemoveAll(tt.outputPath)
				gotOutput = string(outputBytes)
			}
			wantOutput, err := os.ReadFile(tt.wantOutputPath)
			if err != nil {
				t.Errorf("unexpected error reading testdata file %s: %s", tt.wantOutputPath, err)
			}
			wantOutputStr := string(wantOutput)
			if string(wantOutputStr) != gotOutput {
				t.Errorf("unexpected output: got \n%v\n\n, want \n%v\n", gotOutput, wantOutputStr)
			}
		})
	}
}
