package model

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestMergeChectypes(t *testing.T) {
	tests := []struct {
		name string
		ctda Checktypes
		ctdb Checktypes
		// Be aware that reflect.DeepEqual order in slices matters.
		want Checktypes
	}{
		{
			name: "MergeBEmpty",
			ctda: Checktypes{[]Checktype{
				{
					Name: "check1",
				},
			}},
			ctdb: Checktypes{},
			want: Checktypes{[]Checktype{
				{
					Name: "check1",
				},
			}},
		},
		{
			name: "MergeAEmpty",
			ctda: Checktypes{},
			ctdb: Checktypes{[]Checktype{
				{
					Name: "check1",
				},
			}},
			want: Checktypes{[]Checktype{
				{
					Name: "check1",
				},
			}},
		},
		{
			name: "MergeAandBEmpty",
			ctda: Checktypes{},
			ctdb: Checktypes{},
			want: Checktypes{},
		},
		{
			name: "MergeAandBNoDuplicates",
			ctda: Checktypes{[]Checktype{
				{
					Name: "check1",
				},
			}},
			ctdb: Checktypes{[]Checktype{
				{
					Name: "check2",
				},
			}},
			want: Checktypes{[]Checktype{
				{
					Name: "check2",
				},
				{
					Name: "check1",
				},
			}},
		},
		{
			name: "MergeAandBDuplicates",
			ctda: Checktypes{[]Checktype{
				{
					Name:        "check1",
					Description: "A",
				},
			}},
			ctdb: Checktypes{[]Checktype{
				{
					Name:        "check1",
					Description: "B",
				},
			}},
			// Be aware that reflect.DeepEqual order in slices matters.
			want: Checktypes{[]Checktype{
				{
					Name:        "check1",
					Description: "A",
				},
			}},
		},
		{
			name: "MergeAandBDuplicatesWithMoreData",
			ctda: Checktypes{[]Checktype{
				{
					Name:        "check1",
					Description: "A",
				},
				{
					Name:        "check2",
					Description: "A",
				},
			}},
			ctdb: Checktypes{[]Checktype{
				{
					Name:        "check1",
					Description: "B",
				},
				{
					Name:        "check3",
					Description: "B",
				},
			}},
			// Be aware that reflect.DeepEqual order in slices matters.
			want: Checktypes{[]Checktype{
				{
					Name:        "check1",
					Description: "A",
				},
				{
					Name:        "check3",
					Description: "B",
				},
				{
					Name:        "check2",
					Description: "A",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeChecktypes(tt.ctda, tt.ctdb)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchChecktypesFromURL(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    Checktypes
		wantErr bool
		err     error
	}{
		{
			name: "HappyPath",
			path: "happy-path.json",
			want: Checktypes{
				[]Checktype{
					{
						Name:         "check1",
						Description:  "A",
						Image:        "example.com/check1:stable",
						Timeout:      3600,
						RequiredVars: nil,
						Assets:       []string{"Hostname"},
					},
					{
						Name:         "check2",
						Description:  "B",
						Image:        "example.com/check2:stable",
						Timeout:      0,
						RequiredVars: nil,
						Assets:       []string{"Hostname", "IP"},
					},
				},
			},
		},
		{
			name:    "Malformed",
			path:    "malformed.json",
			want:    Checktypes{},
			wantErr: true,
			err:     errors.New("invalid character 'm' looking for beginning of value"),
		},
		{
			name:    "WrongURL",
			path:    "wrong.json",
			want:    Checktypes{},
			wantErr: true,
			err:     errors.New("unexpected status code fetching checktypes from URL"),
		},
	}
	for _, tt := range tests {
		fs := http.FileServer(http.Dir(fmt.Sprintf("testdata/")))
		s := httptest.NewServer(fs)
		defer s.Close()
		url := fmt.Sprintf("%s/%s", s.URL, tt.path)
		got, err := FetchChecktypesFromURL(url)
		if err != nil && !tt.wantErr {
			t.Errorf("unexpected error %v", err)
		}
		if err != nil && tt.wantErr {
			if err.Error() != tt.err.Error() {
				t.Errorf("unexpected error, got: %v, want: %v", err.Error(), tt.err.Error())
			}
		}
		if !reflect.DeepEqual(got, tt.want) {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}
