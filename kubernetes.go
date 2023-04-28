package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
)

type Image struct {
	registry   string `json:"registry"`
	repository string `json:"repository"`
	tag        string `json:"tag"`
	digest     string `json:"digest"`
}

func handleAdmissionReview(r *http.Request) (admissionv1.AdmissionReview, error) {
	// Get the body data, which will be the AdmissionReview
	// content for the request.
	var body []byte
	if r.Body != nil {
		requestData, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return admissionv1.AdmissionReview{}, NewApiError(http.StatusBadRequest, err.Error())
		}
		body = requestData
	}

	codecs := serializer.NewCodecFactory(runtime.NewScheme())
	deserializer := codecs.UniversalDeserializer()

	// Decode the request body
	admissionReview := &admissionv1.AdmissionReview{}
	if _, _, err := deserializer.Decode(body, nil, admissionReview); err != nil {
		return admissionv1.AdmissionReview{}, NewApiError(http.StatusBadRequest, err.Error())
	}

    _, err := handleResource(admissionReview)
    if err != nil {
        return admissionv1.AdmissionReview{}, err
    }

	admissionResponse := &admissionv1.AdmissionResponse{}
	admissionResponse.Allowed = true

	// Construct the response, which is just an AdmissionReview.
	var admissionReviewResponse admissionv1.AdmissionReview
	admissionReviewResponse.Response = admissionResponse
	admissionReviewResponse.SetGroupVersionKind(admissionReview.GroupVersionKind())
	admissionReviewResponse.Response.UID = admissionReview.Request.UID

	return admissionReviewResponse, nil
}

func handleResource(review *admissionv1.AdmissionReview) ([]Image, error) {
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

