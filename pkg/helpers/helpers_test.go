package helpers

import (
	"errors"
	"fmt"
	"testing"
)

func TestIsExistingDirOrFile(t *testing.T) {
	tests := []struct {
		name    string
		path    string
		want    bool
		wantErr bool
		err     error
	}{
		{
			name:    "HappyPathDir",
			path:    "happy-path-dir",
			want:    true,
			wantErr: false,
		},
		{
			name:    "HappyPathDir",
			path:    "happy-path-file",
			want:    false,
			wantErr: false,
		},
		{
			name:    "UnexistingDirOrFile",
			path:    "unexisting",
			want:    false,
			wantErr: true,
			err:     errors.New("stat testdata/unexisting: no such file or directory"),
		},
	}
	for _, tt := range tests {
		path := fmt.Sprintf("testdata/%s", tt.path)
		got, err := IsExistingDirOrFile(path)
		if err != nil && !tt.wantErr {
			t.Errorf("unexpected error %v", err)
		}
		if err != nil && tt.wantErr {
			if err.Error() != tt.err.Error() {
				t.Errorf("unexpected error, got: %v, want: %v", err.Error(), tt.err.Error())
			}
		}
		if got != tt.want {
			t.Errorf("got %v, want %v", got, tt.want)
		}
	}
}
