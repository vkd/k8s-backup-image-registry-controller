package registry

import (
	"context"
	"testing"
)

func TestBackup_CopyImage(t *testing.T) {
	tests := []struct {
		image    string
		prefix   string
		newImage string
	}{
		{"gcr.io/gc/pause", "my.reg/username", "my.reg/username/gcr-io-gc-pause:latest"},
		{"gcr.io/gc/pause", "my.reg", "my.reg/library/gcr-io-gc-pause:latest"},
		{"gcr.io/gc/pause", "", "index.docker.io/library/gcr-io-gc-pause:latest"},

		{"golang:alpine", "my.reg/username", "my.reg/username/golang:alpine"},
		{"my.reg/username/golang:alpine", "my.reg/username", ""},
	}
	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			ctx := context.Background()

			b := Backup{
				RegistryDomain: tt.prefix,
				Copier:         testOkCopier{},
			}

			newImage, ok, err := b.CopyImage(ctx, tt.image)
			if err != nil {
				t.Errorf("Backup.CopyImage() error = %v", err)
				return
			}
			if ok != (tt.newImage != "") {
				t.Errorf("Backup.CopyImage() gotOk = %v, want %q != \"\"", ok, tt.newImage)
			}

			if newImage != tt.newImage {
				t.Errorf("Backup.CopyImage() got = %q, want %q", newImage, tt.newImage)
			}
		})
	}
}

type testOkCopier struct{}

func (t testOkCopier) Copy(_, _ string) error { return nil }
