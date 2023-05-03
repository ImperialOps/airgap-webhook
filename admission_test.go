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
  containers:
    - image: nginx
      name: webserver
      ports:
        - containerPort: 80
          name: http
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
                repository: "nginx",
                tag: "latest",
                digest: "",
            },
            {
                registry: "ghcr.io",
                repository: "stefanprodan/podinfo",
                tag: "6.3.6",
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
