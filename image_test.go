package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewImage(t *testing.T) {
	tests := []struct {
		image    string
		expected Image
	}{
		{"nginx",
			Image{
				registry:   "docker.io",
				repository: "nginx",
				tag:        "latest",
				digestHash: "",
				digest:     "",
			},
		},
		{"nginx@sha256:f2fee5c7194cbbfb9d2711fa5de094c797a42a51aa42b0c8ee8ca31547c872b1",
			Image{
				registry:   "docker.io",
				repository: "nginx",
				tag:        "latest",
				digestHash: "sha256",
				digest:     "f2fee5c7194cbbfb9d2711fa5de094c797a42a51aa42b0c8ee8ca31547c872b1",
			},
		},
		{"public.ecr.aws/lts/ubuntu:edge",
			Image{
				registry:   "public.ecr.aws",
				repository: "lts/ubuntu",
				tag:        "edge",
				digestHash: "",
				digest:     "",
			},
		},
		{"public.ecr.aws/nginx/nginx:stable-perl@sha256:1b624e3e6af841b907b1f5747b6f29ccb5ccb422f9e881eae82bd4b8b72cb7a1",
			Image{
				registry:   "public.ecr.aws",
				repository: "nginx/nginx",
				tag:        "stable-perl",
				digestHash: "sha256",
				digest:     "1b624e3e6af841b907b1f5747b6f29ccb5ccb422f9e881eae82bd4b8b72cb7a1",
			},
		},
	}

    for _, test := range tests {
        image := NewImage(test.image)
        assert.Equal(t, image, test.expected, "got %v, expected %v", image, test.expected)
    }
}
