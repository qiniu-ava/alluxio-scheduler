package constants

import (
	corev1 "k8s.io/api/core/v1"
)

const (
	CREATE_METRICS_CLIENT_FAILED int = 101
	CREATE_CONTROLLER_FAILED     int = 102
	WATCH_POD_FAILED             int = 103
)

var WORKER_ALLOW_PHASE = []corev1.PodPhase{corev1.PodPending, corev1.PodRunning}

const WORKER_IMAGE_PULL_SECRET = "atlab-images"

const WORKER_NAMESPACE = "ava"

const LOCAL_WORKER_LABEL_NAME = "ava.qiniu.com/application"
const LOCAL_WORKER_LABEL_VALUE = "alluxio-local-worker"
const LOCAL_RBD_NAME = "rbd-cachedisk"
