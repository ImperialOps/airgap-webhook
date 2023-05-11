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

	admissionReview.Response = admissionResponse
	admissionReview.SetGroupVersionKind(admissionReview.GroupVersionKind())
	admissionReview.Response.UID = admissionReview.Request.UID

	return admissionReview.AdmissionReview, nil
}

func (r *AdmissionReview) handleResource() error {
	s := (r.Request.Kind.Kind)
	switch s {
	case "Pod":
		return r.handlePodResource()
	case "Job":
		return r.handleJobResource()
	case "CronJob":
		return r.handleCronjobResource()
	case "Deployment":
		return r.handleDeploymentResource()
	case "DaemonSet":
		return r.handleDaemonsetResource()
	case "StatefulSet":
		return r.handleStatefulsetResource()
	case "ReplicaSet":
		return r.handleReplicasetResource()
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s.%s, not implemented", r.Request.Kind.Version, s))
	}
}

func (r *AdmissionReview) handlePodResource() error {
	s := (r.Request.Kind.Version)
	switch s {
	case "v1":
		return r.handlePodV1Resource()
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s.%s, not implemented", s, r.Request.Kind.Kind))
	}
}

func (r *AdmissionReview) handlePodV1Resource() error {
	rawRequest := r.Request.Object.Raw
	pod := corev1.Pod{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &pod); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&pod.Spec)
}

func (r *AdmissionReview) handleJobResource() error {
	s := (r.Request.Kind.Version)
	switch s {
	case "v1":
		return r.handleJobV1Resource()
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s.%s, not implemented", s, r.Request.Kind.Kind))
	}
}

func (r *AdmissionReview) handleJobV1Resource() error {
	rawRequest := r.Request.Object.Raw
	resource := batchv1.Job{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.Template.Spec)
}

func (r *AdmissionReview) handleCronjobResource() error {
	s := (r.Request.Kind.Version)
	switch s {
	case "v1":
		return r.handleCronjobV1Resource()
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s.%s, not implemented", s, r.Request.Kind.Kind))
	}
}

func (r *AdmissionReview) handleCronjobV1Resource() error {
	rawRequest := r.Request.Object.Raw
	resource := batchv1.CronJob{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.JobTemplate.Spec.Template.Spec)
}

func (r *AdmissionReview) handleDeploymentResource() error {
	s := (r.Request.Kind.Version)
	switch s {
	case "v1":
		return r.handleDeploymentV1Resource()
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s.%s, not implemented", s, r.Request.Kind.Kind))
	}
}

func (r *AdmissionReview) handleDeploymentV1Resource() error {
	rawRequest := r.Request.Object.Raw
	resource := appsv1.Deployment{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.Template.Spec)
}

func (r *AdmissionReview) handleDaemonsetResource() error {
	s := (r.Request.Kind.Version)
	switch s {
	case "v1":
		return r.handleDaemonsetV1Resource()
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s.%s, not implemented", s, r.Request.Kind.Kind))
	}
}

func (r *AdmissionReview) handleDaemonsetV1Resource() error {
	rawRequest := r.Request.Object.Raw
	resource := appsv1.DaemonSet{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.Template.Spec)
}

func (r *AdmissionReview) handleStatefulsetResource() error {
	s := (r.Request.Kind.Version)
	switch s {
	case "v1":
		return r.handleStatefulsetV1Resource()
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s.%s, not implemented", s, r.Request.Kind.Kind))
	}
}

func (r *AdmissionReview) handleStatefulsetV1Resource() error {
	rawRequest := r.Request.Object.Raw
	resource := appsv1.StatefulSet{}
	if _, _, err := deserializer.Decode(rawRequest, nil, &resource); err != nil {
		return NewApiError(http.StatusBadRequest, err.Error())
	}
	return r.handlePodSpec(&resource.Spec.Template.Spec)
}

func (r *AdmissionReview) handleReplicasetResource() error {
	s := (r.Request.Kind.Version)
	switch s {
	case "v1":
		return r.handleReplicasetV1Resource()
	default:
		return NewApiError(http.StatusNotImplemented, fmt.Sprintf("resource kind %s.%s, not implemented", s, r.Request.Kind.Kind))
	}
}

func (r *AdmissionReview) handleReplicasetV1Resource() error {
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
