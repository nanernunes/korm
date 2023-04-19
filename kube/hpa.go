package kube

type HPA struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type HPAAPI struct {
	Kubernetes *Kubernetes
	Client     interface{}
}

func NewHPAAPI(kubernetes *Kubernetes) *HPAAPI {
	return &HPAAPI{Kubernetes: kubernetes}
}

func (h *HPAAPI) Find()   {}
func (h *HPAAPI) Create() {}
func (h *HPAAPI) Update() {}
func (h *HPAAPI) Upsert() {}
func (h *HPAAPI) Delete() {}
