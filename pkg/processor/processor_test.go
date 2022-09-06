package processor

import (
	"errors"
	"reflect"
	"testing"

	"github.com/adevinta/vulcan-check-catalog/pkg/model"
	log "github.com/sirupsen/logrus"
)

var logger = log.New()
var nilSlice []string

func TestNew(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		registryBaseURL string
		imageTag        string
		logger          *log.Logger
		want            Runner
	}{
		{
			name:            "DefaultValues",
			path:            "/tmp",
			registryBaseURL: "",
			imageTag:        "",
			logger:          logger,
			want: Runner{
				logger:          logger,
				path:            "/tmp",
				registryBaseURL: "",
				imageTag:        DefaultImageTag,
			},
		},
		{
			name:            "CustomValues",
			path:            "/tmp",
			registryBaseURL: "example.com",
			imageTag:        "latest",
			logger:          logger,
			want: Runner{
				logger:          logger,
				path:            "/tmp",
				registryBaseURL: "example.com",
				imageTag:        "latest",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.path, tt.registryBaseURL, tt.imageTag, tt.logger)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		registryBaseURL string
		imageTag        string
		logger          *log.Logger
		want            model.Checktypes
		wantErr         bool
		err             error
	}{
		{
			name:            "HappyPath",
			path:            "testdata/happy-path/",
			registryBaseURL: "example.com",
			imageTag:        "stable",
			logger:          logger,
			want: model.Checktypes{
				Checktype: []model.Checktype{
					{
						Name:         "check1",
						Description:  "Check1 description",
						Image:        "example.com/check1:stable",
						Timeout:      3600,
						Options:      map[string]interface{}{"active": true, "depth": 2.0, "url": "http://example.com"},
						RequiredVars: nilSlice,
						Assets:       []string{"WebAddress"},
					},
					{
						Name:         "check2",
						Description:  "Check2 description",
						Image:        "example.com/check2:stable",
						Timeout:      60,
						RequiredVars: []string{"var1", "var2"},
						Assets:       []string{"Hostname"},
					},
				},
			},
		},
		{
			name:            "WithoutAssetTypes",
			path:            "testdata/without-asset-types/",
			registryBaseURL: "example.com",
			imageTag:        "stable",
			logger:          logger,
			want: model.Checktypes{
				Checktype: []model.Checktype{
					{
						Name:         "check3",
						Description:  "Check3 description",
						Image:        "example.com/check3:stable",
						Timeout:      120,
						RequiredVars: []string{"var1"},
						Assets:       []string{},
					},
				},
			},
		},
		{
			name:            "UnexistingPath",
			path:            "testdata/unexisting-path/",
			registryBaseURL: "example.com",
			imageTag:        "stable",
			logger:          logger,
			want:            model.Checktypes{},
			wantErr:         true,
			err:             errors.New("open testdata/unexisting-path/: no such file or directory"),
		},
		{
			name:            "EmptyFolder",
			path:            "testdata/empty-folder/",
			registryBaseURL: "example.com",
			imageTag:        "stable",
			logger:          logger,
			want:            model.Checktypes{},
		},
		{
			name:            "FolderWithoutCheckManifest",
			path:            "testdata/no-check-folder",
			registryBaseURL: "example.com",
			imageTag:        "stable",
			logger:          logger,
			want:            model.Checktypes{},
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := New(tt.path, tt.registryBaseURL, tt.imageTag, tt.logger)
			got, err := r.Run()
			if err != nil && !tt.wantErr {
				t.Errorf("unexpected error %v", err)
			}
			if err != nil && tt.wantErr {
				if err.Error() != tt.err.Error() {
					t.Errorf("unexpected error, got: %v, want: %v", err.Error(), tt.err.Error())
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("\ngot:\n%#v\nwant:\n%#v", got, tt.want)
			}
		})
	}
}
