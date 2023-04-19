package kube

type KubeAPI interface {
	Find()
	Create()
	Update()
	Upsert()
	Delete()
}
