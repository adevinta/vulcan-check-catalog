package manifest

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    Data
		wantErr bool
		err     error
	}{
		{
			name: "HappyPath",
			path: "happy-path.toml",
			want: Data{
				Description: "Check1 description",
				Timeout:     3600,
				AssetTypes:  []string{"WebAddress"},
				Options:     `{"depth": 2, "active": true, "url": "http://example.com"}`,
			},
		},
		{
			name:    "Malformed",
			path:    "malformed.toml",
			want:    Data{},
			wantErr: true,
			err:     errors.New("toml: line 2: expected '.' or '=', but got '\\n' instead"),
		},
		{
			name:    "Unexisting",
			path:    "unexisting.toml",
			want:    Data{},
			wantErr: true,
			err:     errors.New("open testdata/unexisting.toml: no such file or directory"),
		},
		{
			name:    "MalformedOptions",
			path:    "malformed-options.toml",
			want:    Data{},
			wantErr: true,
			err:     errors.New("options field is not a valid json string: invalid character 'd' looking for beginning of value"),
		},
	}
	for _, tt := range tests {
		path := fmt.Sprintf("testdata/%s", tt.path)
		got, err := Read(path)
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
