package kube

import (
	"encoding/json"
	"time"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Deployment struct {
	Name          string            `json:"name"`
	Namespace     string            `json:"namespace"`
	Labels        map[string]string `json:"labels,omitempty"`
	Annotations   map[string]string `json:"annotations,omitempty"`
	Replicas      int               `json:"replicas"`
	Containers    []Container       `json:"containers"`
	NodeSelectors map[string]string `json:"node_selectors,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

func (d *Deployment) ToJSON() string {
	response, _ := json.Marshal(d)
	return string(response)
}

func (d *Deployment) Marshal() (appsv1.Deployment, error) {

	var kubeContainers []apiv1.Container
	for _, container := range d.Containers {
		kubeContainers = append(kubeContainers, container.Marshal())
	}

	var replicas int32 = int32(d.Replicas)
	if replicas == 0 {
		replicas = 1
	}

	deployment := appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      d.Name,
			Namespace: d.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": d.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": d.Name,
					},
					Annotations: d.Annotations,
				},
				Spec: apiv1.PodSpec{
					Containers: kubeContainers,
				},
			},
		},
	}

	if d.NodeSelectors != nil {
		deployment.Spec.Template.Spec.NodeSelector = d.NodeSelectors
	}

	return deployment, nil
}

func (d *Deployment) Unmarshal(deployment *appsv1.Deployment) error {
	d.Name = deployment.GetName()
	d.Namespace = deployment.GetNamespace()
	d.CreatedAt = deployment.GetCreationTimestamp().Time
	d.Replicas = int(*deployment.Spec.Replicas)

	d.Labels = make(map[string]string)
	if labels := deployment.GetLabels(); labels != nil {
		d.Labels = labels
	}

	d.Annotations = make(map[string]string)
	if annotations := deployment.GetAnnotations(); annotations != nil {
		d.Annotations = annotations
	}

	for _, c := range deployment.Spec.Template.Spec.Containers {
		container := Container{}
		container.Unmarshal(c)
		d.Containers = append(d.Containers, container)
	}

	d.NodeSelectors = make(map[string]string)
	if nodeSelectors := deployment.Spec.Template.Spec.NodeSelector; nodeSelectors != nil {
		d.NodeSelectors = nodeSelectors
	}

	return nil
}
