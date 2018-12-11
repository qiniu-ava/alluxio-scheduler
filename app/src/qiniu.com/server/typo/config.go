package typo

type Config struct {
	Kube struct {
		ConfigPath string `json:"config_path" config:"config_path"`
	} `json:"kube" config:"kube"`
	Namespace     string `json:"namespace" config:"namespace"`
	NodeListPath  string `json:"nodelist" config:"nodelist"`
	GroupListPath string `json:"grouplist" config:"grouplist"`
}

type GroupConfigItem struct {
	Name     string `json:"name" config:"name"`
	ImageTag string `json:"imagetag" config:"imagetag"`
}

type NodeConfigItem struct {
	Name string `json:"name" config:"name"`
	Ip   string `json:"ip" config:"ip"`
}

type LocalWorker struct {
	Group string `json:"group"`
	Node  string `json:"node"`
}

func (w *LocalWorker) Equal(v LocalWorker) bool {
	return w.Group == v.Group && w.Node == v.Node
}
