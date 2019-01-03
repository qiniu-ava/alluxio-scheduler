package executor

import (
	"fmt"
	"time"

	"qiniu.com/server/util/constants"

	"qiniu.com/server/util"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"qiniu.com/server/typo"
)

const DefaultListLimit = 10000

type PodHandler struct {
	kubeClientCache *kubernetes.Clientset
	conf            *typo.Config
	nodes           []typo.NodeConfigItem
	groups          []typo.GroupConfigItem
}

func (h *PodHandler) handle(t time.Time) {
	if err := h.beforeHandle(t); err != nil {
		util.Logger.Errorf("error when prepare for handling: %v", err)
		return
	}

	if err := h.updateWorkers(); err != nil {
		util.Logger.Errorf("failed to update workers, error: %v", err)
	}

	if err := h.afterHandle(t); err != nil {
		util.Logger.Errorf("%v", err)
	}
}

func (h *PodHandler) beforeHandle(t time.Time) error {
	// do setup staff
	if err := h.updateClient(); err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := h.updateNodes(); err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := h.updateGroups(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (h *PodHandler) updateNodes() error {
	nodes, err := util.GetNodes(h.conf.NodeListPath)
	if err != nil {
		util.Logger.Errorf("failed to get nodes: %v", err)
		return err
	}
	h.nodes = nodes
	return nil
}

func (h *PodHandler) updateGroups() error {
	groups, err := util.GetGroups(h.conf.GroupListPath)
	if err != nil {
		util.Logger.Errorf("failed to get groups: %v", err)
		return err
	}
	h.groups = groups
	return nil
}

func (h *PodHandler) afterHandle(t time.Time) error {
	// do clean staff

	return nil
}

func (h *PodHandler) updateClient() error {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", h.conf.Kube.ConfigPath)
	if err != nil {
		return fmt.Errorf("failed to get kube apiserver config, %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("failed to create kube apiserver client, %v", err)
	}
	h.kubeClientCache = clientset

	return nil
}

func (h *PodHandler) getNecessaryWorkers() (workers []typo.LocalWorker, err error) {
	pods, err := h.kubeClientCache.CoreV1().Pods(h.conf.Namespace).List(metav1.ListOptions{Limit: DefaultListLimit})
	if err != nil {
		util.Logger.Errorf("failed to get all pods in %s namespaces, error: %v", h.conf.Namespace, err)
		return
	}

	nodeWorkerMap := make(map[string][]string)

	for index := range pods.Items {
		pod := pods.Items[index]

		allowedPhase := false
		for _, phs := range constants.WORKER_ALLOW_PHASE {
			if pod.Status.Phase == phs {
				allowedPhase = true
				break
			}
		}

		// only pods in allowed phase need to be handle
		if !allowedPhase {
			continue
		}

		for _, vol := range pod.Spec.Volumes {
			if vol.FlexVolume != nil {
				group := vol.FlexVolume.Options["group"]

				// if group is not in group list, it should be trim to default
				gExsist := false
				for _, g := range h.groups {
					if g.Name == group {
						gExsist = true
					}
				}

				if !gExsist {
					group = "default"
				}

				util.Logger.Debugf("pod %s with flex volume in group %s, on node %s", pod.Name, group, pod.Spec.NodeName)
				if nodeWorkerMap[pod.Spec.NodeName] != nil {
					marked := false
					for _, gName := range nodeWorkerMap[pod.Spec.NodeName] {
						if gName == group {
							marked = true
						}
					}
					if !marked {
						nodeWorkerMap[pod.Spec.NodeName] = append(nodeWorkerMap[pod.Spec.NodeName], group)
					}
				} else {
					nodeWorkerMap[pod.Spec.NodeName] = []string{group}
				}
				break
			}
		}
	}

	workers = []typo.LocalWorker{}
	for node, groups := range nodeWorkerMap {
		for _, group := range groups {
			workers = append(workers, typo.LocalWorker{
				Node:  node,
				Group: group,
			})
		}
	}

	return
}

func (h *PodHandler) getLocalWorkers() (workers []typo.LocalWorker, err error) {
	pods, err := h.kubeClientCache.CoreV1().Pods(constants.WORKER_NAMESPACE).List(metav1.ListOptions{
		Limit:         DefaultListLimit,
		LabelSelector: fmt.Sprintf("%s=%s", constants.LOCAL_WORKER_LABEL_NAME, constants.LOCAL_WORKER_LABEL_VALUE),
	})

	if err != nil {
		return
	}

	workers = []typo.LocalWorker{}
	for _, pod := range pods.Items {
		group, node := util.ParseLocalWorkerName(pod.Name)
		workers = append(workers, typo.LocalWorker{Group: group, Node: node})
	}

	return
}

func (h *PodHandler) updateWorkers() error {
	necessaryWorkers, err := h.getNecessaryWorkers()
	if err != nil {
		return err
	}

	currentWorkers, err := h.getLocalWorkers()
	if err != nil {
		return err
	}

	podResolver := NewPodResolver(h.kubeClientCache)

	for _, nWorker := range necessaryWorkers {
		created := false
		for _, cWorker := range currentWorkers {
			if nWorker.Equal(cWorker) {
				created = true
			}
		}

		if !created {
			podResolver.create(PodCreateOptions{
				Node:   nWorker.Node,
				Group:  nWorker.Group,
				Nodes:  h.nodes,
				Groups: h.groups,
			})
		}
	}

	for _, cWorker := range currentWorkers {
		need := false
		for _, nWorker := range necessaryWorkers {
			if nWorker.Equal(cWorker) {
				need = true
			}
		}

		// TODO: we should try later, but not immediately
		if !need {
			podResolver.delete(PodDeleteOptions{
				Node:   cWorker.Node,
				Group:  cWorker.Group,
				Nodes:  h.nodes,
				Groups: h.groups,
			})
		}
	}

	return nil
}

func GetPodHandler(conf *typo.Config) TickerHandlerFunc {
	h := PodHandler{}
	h.conf = conf
	return h.handle
}
