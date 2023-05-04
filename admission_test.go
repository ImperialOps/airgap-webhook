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
	v1Job = []byte(`apiVersion: batch/v1
kind: Job
metadata:
  name: pi
spec:
  template:
    spec:
      containers:
      - name: pi
        image: perl:5.34.0
        command: ["perl",  "-Mbignum=bpi", "-wle", "print bpi(2000)"]
      restartPolicy: Never
  backoffLimit: 4`)
	v1CronJob = []byte(`apiVersion: batch/v1
kind: CronJob
metadata:
  name: hello
spec:
  schedule: "* * * * *"
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: hello
            image: busybox:1.28
            imagePullPolicy: IfNotPresent
            command:
            - /bin/sh
            - -c
            - date; echo Hello from the Kubernetes cluster
          restartPolicy: OnFailure`)
	v1Deployment = []byte(`apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: public.ecr.aws/nginx/nginx:stable-perl@sha256:1b624e3e6af841b907b1f5747b6f29ccb5ccb422f9e881eae82bd4b8b72cb7a1
        ports:
        - containerPort: 80`)
	v1Daemonset = []byte(`apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: fluentd-elasticsearch
  namespace: kube-system
  labels:
    k8s-app: fluentd-logging
spec:
  selector:
    matchLabels:
      name: fluentd-elasticsearch
  template:
    metadata:
      labels:
        name: fluentd-elasticsearch
    spec:
      tolerations:
      - key: node-role.kubernetes.io/control-plane
        operator: Exists
        effect: NoSchedule
      - key: node-role.kubernetes.io/master
        operator: Exists
        effect: NoSchedule
      containers:
      - name: fluentd-elasticsearch
        image: quay.io/fluentd_elasticsearch/fluentd:v2.5.2
        resources:
          limits:
            memory: 200Mi
          requests:
            cpu: 100m
            memory: 200Mi
        volumeMounts:
        - name: varlog
          mountPath: /var/log
      terminationGracePeriodSeconds: 30
      volumes:
      - name: varlog
        hostPath:
          path: /var/log`)
	v1StatefulSet = []byte(`apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: web
spec:
  selector:
    matchLabels:
      app: nginx # has to match .spec.template.metadata.labels
  serviceName: "nginx"
  replicas: 3 # by default is 1
  minReadySeconds: 10 # by default is 0
  template:
    metadata:
      labels:
        app: nginx # has to match .spec.selector.matchLabels
    spec:
      terminationGracePeriodSeconds: 10
      containers:
      - name: nginx
        image: registry.k8s.io/nginx-slim:0.8
        ports:
        - containerPort: 80
          name: web
        volumeMounts:
        - name: www
          mountPath: /usr/share/nginx/html
  volumeClaimTemplates:
  - metadata:
      name: www
    spec:
      accessModes: [ "ReadWriteOnce" ]
      storageClassName: "my-storage-class"
      resources:
        requests:
          storage: 1Gi`)
	v1ReplicaSet = []byte(`apiVersion: apps/v1
kind: ReplicaSet
metadata:
  name: frontend
  labels:
    app: guestbook
    tier: frontend
spec:
  # modify replicas according to your case
  replicas: 3
  selector:
    matchLabels:
      tier: frontend
  template:
    metadata:
      labels:
        tier: frontend
    spec:
      containers:
      - name: php-redis
        image: gcr.io/google_samples/gb-frontend:v3`)
)

func TestHandleResource(t *testing.T) {
    tests := []struct {
        resource []byte
        expected []Image
    }{
        {v1Pod, []Image{
            {
                registry:   "docker.io",
                repository: "busybox",
                tag:        "1.28",
                digestHash: "",
                digest:     "",
            },
            {
                registry:   "docker.io",
                repository: "nginx",
                tag:        "latest",
                digestHash: "sha256",
                digest:     "f2fee5c7194cbbfb9d2711fa5de094c797a42a51aa42b0c8ee8ca31547c872b1",
            },
            {
                registry:   "docker.io",
                repository: "nginx",
                tag:        "stable",
                digestHash: "sha256",
                digest:     "f3c37d8a26f7a7d8a547470c58733f270bcccb7e785da17af81ec41576170da8",
            },
            {
                registry:   "ghcr.io",
                repository: "stefanprodan/podinfo",
                tag:        "6.3.6",
                digestHash: "",
                digest:     "",
            },
        }},
        {v1Job, []Image{
            {
                registry:   "docker.io",
                repository: "perl",
                tag:        "5.34.0",
                digestHash: "",
                digest:     "",
            },
        }},
        {v1CronJob, []Image{
            {
                registry:   "docker.io",
                repository: "busybox",
                tag:        "1.28",
                digestHash: "",
                digest:     "",
            },
        }},
        {v1Deployment, []Image{
            {
                registry:   "public.ecr.aws",
                repository: "nginx/nginx",
                tag:        "stable-perl",
                digestHash: "sha256",
                digest:     "1b624e3e6af841b907b1f5747b6f29ccb5ccb422f9e881eae82bd4b8b72cb7a1",
            },
        }},
        {v1Daemonset, []Image{
            {
                registry:   "quay.io",
                repository: "fluentd_elasticsearch/fluentd",
                tag:        "v2.5.2",
                digestHash: "",
                digest:     "",
            },
        }},
        {v1StatefulSet, []Image{
            {
                registry:   "registry.k8s.io",
                repository: "nginx-slim",
                tag:        "0.8",
                digestHash: "",
                digest:     "",
            },
        }},
        {v1ReplicaSet, []Image{
            {
                registry:   "gcr.io",
                repository: "google_samples/gb-frontend",
                tag:        "v3",
                digestHash: "",
                digest:     "",
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
