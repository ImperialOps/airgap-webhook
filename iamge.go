package main

import "strings"

type Image struct {
	registry   string `json:"registry"`
	repository string `json:"repository"`
	tag        string `json:"tag"`
	digestHash string `json:"hash"`
	digest     string `json:"digest"`
}

func NewImage(i string) Image {
	image := Image{}
	if i == "" {
		return image
	}

	// hangle digest
	dirtyTag, dirtyDigest, foundDigest := strings.Cut(i, "@")
	if foundDigest {
		digestParts := strings.Split(dirtyDigest, ":")
		image.digestHash = digestParts[0]
		image.digest = digestParts[1]
	} else {
		image.digestHash = ""
		image.digest = ""
	}

	// handle tag
	name, tag, foundTag := strings.Cut(dirtyTag, ":")
	if foundTag {
		image.tag = tag
	} else {
		image.tag = "latest"
	}

	// handle registry and repository
	nameParts := strings.Split(name, "/")
	if strings.Contains(nameParts[0], ".") {
		image.registry = nameParts[0]
		image.repository = strings.Join(nameParts[1:], "/")
	} else {
		image.registry = "docker.io"
		image.repository = strings.Join(nameParts, "/")
	}
	return image
}
