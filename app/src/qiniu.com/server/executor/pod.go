package executor

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"qiniu.com/server/typo"
	"qiniu.com/server/util"
	"qiniu.com/server/util/constants"
)

const WorkerNamespace = "ava"
const maxRetryTimes = 3

type PodResolver struct {
	client *kubernetes.Clientset
}

func NewPodResolver(client *kubernetes.Clientset) PodResolver {
	// r := rest.Config{}
	// util.Logger.Errorf("%v", r)
	return PodResolver{
		client: client,
	}
}

type PodCreateOptions struct {
	Group  string
	Node   string
	Nodes  []typo.NodeConfigItem
	Groups []typo.GroupConfigItem
}

type PodDeleteOptions struct {
	Group  string
	Node   string
	Nodes  []typo.NodeConfigItem
	Groups []typo.GroupConfigItem
}

func (p *PodResolver) create(options PodCreateOptions) {
	pod := corev1.Pod{}
	pod.Name = util.GetLocalWorkerName(options.Group, options.Node)
	pod.Namespace = WorkerNamespace
	pod.Labels = map[string]string{
		constants.LOCAL_WORKER_LABEL_NAME: constants.LOCAL_WORKER_LABEL_VALUE,
	}
	podSpec := corev1.PodSpec{
		NodeSelector:     util.GetNodeSelector(options.Node),
		HostNetwork:      true,
		DNSPolicy:        corev1.DNSClusterFirstWithHostNet,
		Containers:       p.getContainers(options),
		ImagePullSecrets: util.GetImagePullSecret(),
		Volumes:          util.GetVolumes(options.Group),
	}
	pod.Spec = podSpec

	retryTimes := maxRetryTimes
retry:
	_, err := p.client.Core().Pods(WorkerNamespace).Create(&pod)
	if err != nil {
		if retryTimes > 0 {
			util.Logger.Warnf("failed to create local worker for group %s in node %s, times to retry: %f", options.Group, options.Node, retryTimes)
			retryTimes--
			goto retry
		} else {
			util.Logger.Errorf("failed to create local worker for group %s in node %s after %d", options.Group, options.Node, maxRetryTimes)
		}
	} else {
		util.Logger.Infof("creating local worker for group %s in node %s", options.Group, options.Node)
	}
}

func (p *PodResolver) getContainers(options PodCreateOptions) []corev1.Container {
	envs := []corev1.EnvVar{
		corev1.EnvVar{Name: "ALLUXIO_LOCALITY_NODE", Value: util.GetNodeIP(options.Nodes, options.Node)},
		corev1.EnvVar{Name: "ALLUXIO_WORKER_HOSTNAME", Value: util.GetNodeIP(options.Nodes, options.Node)},
	}
	envFromCfgmap := corev1.EnvFromSource{
		ConfigMapRef: &corev1.ConfigMapEnvSource{
			LocalObjectReference: corev1.LocalObjectReference{
				Name: util.GetCfgName(options.Group),
			},
		},
	}

	resourceLimit := map[corev1.ResourceName]resource.Quantity{}
	resourceRequest := map[corev1.ResourceName]resource.Quantity{}

	cpuRequest, _ := resource.ParseQuantity("1")
	cpuLimit, _ := resource.ParseQuantity("5")
	memRequest, _ := resource.ParseQuantity("8Gi")
	memLimit, _ := resource.ParseQuantity("20Gi")

	resourceLimit[corev1.ResourceCPU] = cpuLimit
	resourceLimit[corev1.ResourceMemory] = memLimit
	resourceRequest[corev1.ResourceCPU] = cpuRequest
	resourceRequest[corev1.ResourceMemory] = memRequest

	container := corev1.Container{
		Name:    "worker",
		Command: []string{"/entrypoint.sh", "worker", "--no-format"},
		Image:   util.GetImage(options.Groups, options.Group),
		Env:     envs,
		EnvFrom: []corev1.EnvFromSource{envFromCfgmap},
		Resources: corev1.ResourceRequirements{
			Requests: resourceRequest,
			Limits:   resourceLimit,
		},
		VolumeMounts: []corev1.VolumeMount{
			corev1.VolumeMount{
				Name:      constants.LOCAL_RBD_NAME,
				MountPath: util.GetVolumePath(options.Group),
			},
		},
	}
	return []corev1.Container{container}
}

func (p *PodResolver) delete(options PodDeleteOptions) {
	retryTimes := maxRetryTimes
retry:
	err := p.client.Core().Pods(WorkerNamespace).Delete(util.GetLocalWorkerName(options.Group, options.Node), &metav1.DeleteOptions{})
	if err != nil {
		if retryTimes > 0 {
			util.Logger.Warnf("failed to delete local worker for group %s in node %s, times to retry: %d", options.Group, options.Node, retryTimes)
			retryTimes--
			goto retry
		} else {
			util.Logger.Errorf("failed to delete local worker for group %s in node %s after %d", options.Group, options.Node, maxRetryTimes)
		}
	} else {
		util.Logger.Infof("deleting local worker for group %s in node %s", options.Group, options.Node)
	}
}
