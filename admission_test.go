package main

import (
	"testing"

	"github.com/anderseknert/kube-review/pkg/admission"
	"github.com/stretchr/testify/assert"
	admissionv1 "k8s.io/api/admission/v1"
)

var (
	v1Pod = `apiVersion: v1
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
          name: http`
)

func TestHandleResource(t *testing.T) {
	tests := []struct {
		admissionRequest *admissionv1.AdmissionReview
		expected         []Image
	}{
		{newAdmissionRequest(v1Pod), []Image{}},
	}

	for _, test := range tests {
		_ = test.admissionRequest
		_, err := NewClientset(test.path)
		if err != nil && !test.err {
			t.Errorf("got err: %s, on path: %s", err.Error(), test.path)
		}
	}
	assert.Equal(t, true, true, "ok")
}
