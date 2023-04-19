package kube

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

type DeploymentAPI struct {
	Kubernetes *Kubernetes
}

func (d *DeploymentAPI) GetClient() v1.DeploymentInterface {
	namespace := apiv1.NamespaceDefault
	if ns := d.Kubernetes.Statement.GetTarget().(Deployment).Namespace; ns != "" {
		namespace = ns
	}
	return d.Kubernetes.Client.AppsV1().Deployments(namespace)
}

func NewDeploymentAPI(kubernetes *Kubernetes) *DeploymentAPI {
	return &DeploymentAPI{Kubernetes: kubernetes}
}

func (d *DeploymentAPI) Find() {
	stmt := d.Kubernetes.Statement

	switch stmt.GetTarget().(type) {
	case Deployment:
		deployment := stmt.GetTarget().(Deployment)
		deploymentSpec, _ := d.GetClient().Get(
			context.TODO(), deployment.Name, metav1.GetOptions{})

		deployment.Unmarshal(deploymentSpec)
		stmt.SetTarget(deployment)

	case []Deployment:
		stmt.SetTarget(
			[]Deployment{{Name: "matrix"}},
		)
	}
}

func (d *DeploymentAPI) Create() {
	stmt := d.Kubernetes.Statement

	deployment := stmt.GetTarget().(Deployment)
	deploymentSpec, _ := deployment.Marshal()

	result, _ := d.GetClient().Create(
		context.TODO(), &deploymentSpec, metav1.CreateOptions{})

	fmt.Println(result.GetObjectMeta().GetName())

	// return err
	stmt.SetTarget(deployment)
}

func (d *DeploymentAPI) Update() {}
func (d *DeploymentAPI) Upsert() {}
func (d *DeploymentAPI) Delete() {}
