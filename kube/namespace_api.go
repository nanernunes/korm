package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type NamespaceAPI struct {
	Kubernetes *Kubernetes
	Client     v1.NamespaceInterface
}

func NewNamespaceAPI(kubernetes *Kubernetes) *NamespaceAPI {
	client := kubernetes.Client.CoreV1().Namespaces()
	return &NamespaceAPI{Kubernetes: kubernetes, Client: client}
}

func (n *NamespaceAPI) Find() {
	stmt := n.Kubernetes.Statement

	switch stmt.GetTarget().(type) {
	case Namespace:
		namespace := stmt.GetTarget().(Namespace)
		namespaceSpec, _ := n.Client.Get(context.TODO(), namespace.Name, metav1.GetOptions{})

		namespace.Unmarshal(namespaceSpec)
		stmt.SetTarget(namespace)

	case []Namespace:
		namespaces := stmt.GetTarget().([]Namespace)
		namespacesSpec, _ := n.Client.List(context.TODO(), metav1.ListOptions{})

		for _, n := range namespacesSpec.Items {
			var namespace Namespace
			namespace.Unmarshal(&n)
			namespaces = append(namespaces, namespace)
		}
		stmt.SetTarget(namespaces)
	}

}

func (n *NamespaceAPI) Create() {
	stmt := n.Kubernetes.Statement

	namespace := stmt.GetTarget().(Namespace)
	namespaceSpec, _ := namespace.Marshal()

	n.Client.Create(context.TODO(), &namespaceSpec, metav1.CreateOptions{})
	// need to return the created data

}

func (n *NamespaceAPI) Update() {
	stmt := n.Kubernetes.Statement

	old := stmt.GetTarget().(Namespace)
	new := stmt.GetOrigin().(Namespace)

	namespaceSpec, _ := n.Client.Get(context.TODO(), old.Name, metav1.GetOptions{})

	if new.Labels != nil {
		namespaceSpec.SetLabels(new.Labels)
	}

	if updated, err := n.Client.Update(context.TODO(), namespaceSpec, metav1.UpdateOptions{}); err == nil {
		new.Unmarshal(updated)
		stmt.SetTarget(new)
	}

}

func (n *NamespaceAPI) Upsert() {

}

func (n *NamespaceAPI) Delete() {
	stmt := n.Kubernetes.Statement

	namespace := stmt.GetTarget().(Namespace)
	n.Client.Delete(context.TODO(), namespace.Name, metav1.DeleteOptions{})
}
