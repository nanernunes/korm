package kube

import (
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Namespace struct {
	Name   string            `json:"name"`
	Labels map[string]string `json:"labels"`
}

func (n *Namespace) Marshal() (apiv1.Namespace, error) {
	namespace := apiv1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:   n.Name,
			Labels: n.Labels,
		},
	}

	return namespace, nil
}

func (n *Namespace) Unmarshal(namespace *apiv1.Namespace) error {
	n.Name = namespace.Name
	n.Labels = namespace.Labels
	return nil
}
