package registry

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/name"
)

const RegistryRepoDelimeter = "/"

type Backup struct {
	RegistryDomain string
	Copier         interface {
		Copy(src, dst string) error
	}
}

func NewBackupRegistry(domain string) Backup {
	return Backup{
		RegistryDomain: strings.TrimSuffix(domain, "/"),
		Copier:         CraneCopier{},
	}
}

func (b Backup) CopyImage(ctx context.Context, image string) (newImage string, ok bool, _ error) {
	var err error

	registry, namespace, _ := strings.Cut(b.RegistryDomain, RegistryRepoDelimeter)
	if registry == "" {
		registry = name.DefaultRegistry
	}
	if namespace == "" {
		namespace = "library" // default namespace
	}

	newImagePrefix := registry + RegistryRepoDelimeter + namespace + RegistryRepoDelimeter

	if strings.HasPrefix(image, newImagePrefix) {
		return "", false, nil
	}

	newImage = image
	newImage = strings.ReplaceAll(newImage, ".", "-")
	newImage = strings.ReplaceAll(newImage, RegistryRepoDelimeter, "-")
	newImage = newImagePrefix + newImage

	newImg, err := name.ParseReference(newImage)
	if err != nil {
		return "", false, fmt.Errorf("parse new image %q: %w", newImage, err)
	}

	out := newImg.Name()

	err = b.Copier.Copy(image, out)
	if err != nil {
		return "", false, fmt.Errorf("copy new image %q: %w", out, err)
	}

	return out, true, nil
}

type CraneCopier struct{}

func (c CraneCopier) Copy(src, dst string) error {
	return crane.Copy(src, dst, crane.WithAuthFromKeychain(authn.DefaultKeychain))
}
