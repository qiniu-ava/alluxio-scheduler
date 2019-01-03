package util

import (
	"fmt"
	"strings"

	"qiniu.com/server/typo"

	"qiniu.com/server/util/constants"

	corev1 "k8s.io/api/core/v1"
)

const volumePath = ""

func GetLocalWorkerName(group, node string) string {
	return fmt.Sprintf("alluxio-worker-%s-%s", group, node)
}

func ParseLocalWorkerName(name string) (group, node string) {
	slices := strings.Split(name, "-")
	group = slices[2]
	node = slices[3]
	return
}

func GetNodeSelector(node string) map[string]string {
	ns := make(map[string]string)
	ns["kubernetes.io/hostname"] = node
	return ns
}

func GetImage(groups []typo.GroupConfigItem, group string) string {
	for _, groupItem := range groups {
		if groupItem.Name == group {
			return fmt.Sprintf("reg-xs.qiniu.io/atlab/alluxio:%s", groupItem.ImageTag)
		}
	}

	defaultTag := ""
	for _, groupItem := range groups {
		if groupItem.Name == "default" {
			defaultTag = groupItem.ImageTag
		}
	}

	Logger.Warnf("group %s not found in group list, use default group")
	return fmt.Sprintf("reg-xs.qiniu.io/atlab/alluxio:%s", defaultTag)
}

func GetNodeIP(nodes []typo.NodeConfigItem, node string) string {
	for _, item := range nodes {
		if item.Name == node {
			return item.Ip
		}
	}
	Logger.Errorf("node %s not found in node list")
	return ""
}

func GetCfgName(group string) string {
	return fmt.Sprintf("alluxio-worker-%s-configmap", group)
}

func GetImagePullSecret() []corev1.LocalObjectReference {
	return []corev1.LocalObjectReference{corev1.LocalObjectReference{
		Name: constants.WORKER_IMAGE_PULL_SECRET,
	}}
}

func GetVolumePath(group string) string {
	return fmt.Sprintf("/disk-rbd/alluxio/workers/%s/cachedisk", group)
}

func GetVolumes(group string) []corev1.Volume {
	return []corev1.Volume{corev1.Volume{
		Name: constants.LOCAL_RBD_NAME,
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: GetVolumePath(group),
			},
		},
	}}
}
