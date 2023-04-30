package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

type AdmissionReview struct {
	admissionv1.AdmissionReview
	images []Image
}

func NewAdmissionReview(b []byte) (*AdmissionReview, error) {
	codecs := serializer.NewCodecFactory(runtime.NewScheme())
	deserializer := codecs.UniversalDeserializer()

	// Decode the bytes
	admissionReview := &AdmissionReview{}
	if _, _, err := deserializer.Decode(b, nil, admissionReview); err != nil {
		return &AdmissionReview{}, err
	}

	admissionReview.images = []Image{}
	return admissionReview, nil
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

func handleAdmissionReview(i []byte) (admissionv1.AdmissionReview, error) {
	// Decode the request body
	admissionReview, err := NewAdmissionReview(i)
	if err != nil {
		return admissionv1.AdmissionReview{}, NewApiError(http.StatusBadRequest, err.Error())
	}

	admissionReview.images, err = handleResource(admissionReview)
	if err != nil {
		return admissionv1.AdmissionReview{}, err
	}

	// Construct the response, which is just an AdmissionReview.
	admissionResponse := &admissionv1.AdmissionResponse{}
	admissionResponse.Allowed = true

	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admissionResponse
	admissionReviewResponse.SetGroupVersionKind(admissionReview.GroupVersionKind())
	admissionReviewResponse.Response.UID = admissionReview.Request.UID

	return admissionReviewResponse, nil
}

func handleResource(review *AdmissionReview) ([]Image, error) {
	codecs := serializer.NewCodecFactory(runtime.NewScheme())
	deserializer := codecs.UniversalDeserializer()
	rawRequest := review.Request.Object.Raw

	s := (review.Request.Kind.Version + "." + review.Request.Kind.Kind)
	switch s {
	case "v1.Pod":
		pod := corev1.Pod{}
		if _, _, err := deserializer.Decode(rawRequest, nil, &pod); err != nil {
			return []Image{}, NewApiError(http.StatusBadRequest, err.Error())
		}
		return handlePodResource(&pod)
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
		return []Image{}, NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s, not implemented", s))
	}
	return []Image{}, nil
}

func handlePodResource(pod *corev1.Pod) ([]Image, error) {
	log.Printf("got pod %s", pod.Name)
	return []Image{}, nil
}

func handleDeploymentResource(deployment *appsv1.Deployment) ([]Image, error) {
	return []Image{}, nil
}
