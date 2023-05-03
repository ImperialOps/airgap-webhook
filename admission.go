package main

import (
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

var (
	codecs       = serializer.NewCodecFactory(runtime.NewScheme())
	deserializer = codecs.UniversalDeserializer()
)

type AdmissionReview struct {
	admissionv1.AdmissionReview
	images []Image
}

func NewAdmissionReview(b []byte) (*AdmissionReview, error) {
	// Decode the bytes
	admissionReview := &AdmissionReview{}
	if _, _, err := deserializer.Decode(b, nil, admissionReview); err != nil {
		return &AdmissionReview{}, NewApiError(http.StatusBadRequest, err.Error())
	}

	admissionReview.images = []Image{}
	return admissionReview, nil
}

func MustAdmissionReview(b []byte) *AdmissionReview {
    admissionReview, err := NewAdmissionReview(b)
    if err != nil {
        log.Panicf("could not create admission review: %s", err)
    }
	return admissionReview
}

type Image struct {
	registry   string `json:"registry"`
	repository string `json:"repository"`
	tag        string `json:"tag"`
	digest     string `json:"digest"`
}

func NewImage(i string) Image {
	image := Image{}
	if i == "" {
		return image
	}

	return image
}

func handleAdmissionReview(b []byte) (admissionv1.AdmissionReview, error) {
	// Decode the request body
	admissionReview, err := NewAdmissionReview(b)
	if err != nil {
		return admissionv1.AdmissionReview{}, err
	}

	err = admissionReview.handleResource()
	if err != nil {
		return admissionv1.AdmissionReview{}, err
	}

	// TODO test with our AdmissionReview
	// Construct the response, which is just an AdmissionReview.
	admissionResponse := &admissionv1.AdmissionResponse{}
	admissionResponse.Allowed = true

	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admissionResponse
	admissionReviewResponse.SetGroupVersionKind(admissionReview.GroupVersionKind())
	admissionReviewResponse.Response.UID = admissionReview.Request.UID

	return admissionReviewResponse, nil
}

func (r *AdmissionReview) handleResource() error {
	s := (r.Request.Kind.Version + "." + r.Request.Kind.Kind)
	switch s {
	case "v1.Pod":
		return r.handlePodResource()
	case "v1.Job":
		_ = ""
	case "v1.CronJob":
		_ = ""
	case "v1.Deployment":
		_ = ""
	case "v1.Daemonset":
		_ = ""
	case "v1.StatefulSet":
		_ = ""
	case "v1.ReplicaSet":
		_ = ""
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s, not implemented", s))
	}
	return nil
}

func (r *AdmissionReview) handlePodResource() error {
	rawRequest := r.Request.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &pod); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}

	log.Printf("got pod %s", pod.Name)
	return nil
}

func (r *AdmissionReview) handleDeploymentResource() error {
	rawRequest := r.Request.Object.Raw
	deployment := appsv1.Deployment{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &deployment); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}

	log.Printf("got deployment %s", deployment.Name)
	return nil
}