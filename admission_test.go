package main

import (
	"testing"

	"github.com/imperialops/airgap-webhook/admission"
	"github.com/stretchr/testify/assert"
)

var (
	v1Pod = []byte(`apiVersion: v1
kind: Pod
metadata:
  labels:
    name: webserver
  name: nginx-webserver
spec:
  initContainers:
  - name: init
    image: busybox:1.28
    command: ['sh', '-c', "echo test me"]
  containers:
    - image: nginx@sha256:f2fee5c7194cbbfb9d2711fa5de094c797a42a51aa42b0c8ee8ca31547c872b1
      name: weblatest
      ports:
        - containerPort: 80
          name: http
    - image: nginx:stable@sha256:f3c37d8a26f7a7d8a547470c58733f270bcccb7e785da17af81ec41576170da8
      name: webstable
      ports:
        - containerPort: 8080
          name: httpother
    - image: ghcr.io/stefanprodan/podinfo:6.3.6
      name: podinfo
      ports:
        - containerPort: 9898
          name: http`)
)

func TestHandleResource(t *testing.T) {
	tests := []struct {
		resource []byte
		expected []Image
	}{
        {v1Pod, []Image{
            {
                registry: "docker.io",
                repository: "busybox",
                tag: "1.28",
                digestHash: "",
                digest: "",
            },
            {
                registry: "docker.io",
                repository: "nginx",
                tag: "latest",
                digestHash: "sha256",
                digest: "f2fee5c7194cbbfb9d2711fa5de094c797a42a51aa42b0c8ee8ca31547c872b1",
            },
            {
                registry: "docker.io",
                repository: "nginx",
                tag: "stable",
                digestHash: "sha256",
                digest: "f3c37d8a26f7a7d8a547470c58733f270bcccb7e785da17af81ec41576170da8",
            },
            {
                registry: "ghcr.io",
                repository: "stefanprodan/podinfo",
                tag: "6.3.6",
                digestHash: "",
                digest: "",
            },
        }},
	}

	for _, test := range tests {
        bytes, err := admission.CreateAdmissionReviewRequest(test.resource, "create", "imperialops", []string{})
        admissionReview := MustAdmissionReview(bytes)
        assert.NoError(t, err)
        admissionReview.handleResource()
        assert.Equal(t, admissionReview.images, test.expected, "got %v, expected %v", admissionReview.images, test.expected)
	}
}
