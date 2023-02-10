// Package util contains utility functions.
package util_test

import (
	"testing"

	"github.com/nextmv-io/sdk/util"
)

func TestIsPseudoVersion(t *testing.T) {
	tests := []struct {
		version string
		want    bool
	}{
		{
			version: "v1.0.1-0.20210101000000-000000000000",
			want:    true,
		},
		{
			version: "v0.21.1-0.20230209001616-14e922147acd",
			want:    true,
		},
		{
			version: "v0.21.4",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			if got := util.IsPseudoVersion(tt.version); got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBaseOfPseudoVersion(t *testing.T) {
	tests := []struct {
		version string
		want    string
		err     string
	}{
		{
			version: "v1.0.1-0.20210101000000-000000000000",
			want:    "v1.0.0",
		},
		{
			version: "v0.21.16-0.20230209001616-14e922147acd",
			want:    "v0.21.15",
		},
		{
			version: "v0.21.4",
			err:     "version v0.21.4 is not a pseudo version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.version, func(t *testing.T) {
			got, err := util.GetBaseOfPseudoVersion(tt.version)
			if err != nil {
				if err.Error() != tt.err {
					t.Errorf("got %v, want %v", err, tt.err)
				}
				return
			}
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
