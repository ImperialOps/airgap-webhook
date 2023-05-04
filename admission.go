package main

import (
	"fmt"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchv1 "k8s.io/api/batch/v1"
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
		return r.handleJobResource()
	case "v1.CronJob":
		return r.handleCronjobResource()
	case "v1.Deployment":
		return r.handleDeploymentResource()
	case "v1.DaemonSet":
		return r.handleDaemonsetResource()
	case "v1.StatefulSet":
		return r.handleStatefulsetResource()
	case "v1.ReplicaSet":
		return r.handleReplicasetResource()
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s, not implemented", s))
	}
}

func (r *AdmissionReview) handlePodResource() error {
	rawRequest := r.Request.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &pod); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&pod.Spec)
}

func (r *AdmissionReview) handleJobResource() error {
	rawRequest := r.Request.Object.Raw
	resource := batchv1.Job{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.Template.Spec)
}

func (r *AdmissionReview) handleCronjobResource() error {
	rawRequest := r.Request.Object.Raw
	resource := batchv1.CronJob{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.JobTemplate.Spec.Template.Spec)
}

func (r *AdmissionReview) handleDeploymentResource() error {
	rawRequest := r.Request.Object.Raw
	resource := appsv1.Deployment{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.Template.Spec)
}

func (r *AdmissionReview) handleDaemonsetResource() error {
	rawRequest := r.Request.Object.Raw
	resource := appsv1.DaemonSet{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.Template.Spec)
}

func (r *AdmissionReview) handleStatefulsetResource() error {
	rawRequest := r.Request.Object.Raw
	resource := appsv1.StatefulSet{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.Template.Spec)
}

func (r *AdmissionReview) handleReplicasetResource() error {
	rawRequest := r.Request.Object.Raw
	resource := appsv1.ReplicaSet{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.Template.Spec)
}

func (r *AdmissionReview) handlePodSpec(spec *corev1.PodSpec) error {
	for _, container := range spec.InitContainers {
		r.images = append(r.images, NewImage(container.Image))
	}
	for _, container := range spec.Containers {
		r.images = append(r.images, NewImage(container.Image))
	}
	return nil
}
