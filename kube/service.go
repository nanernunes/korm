package kube

import (
	"encoding/json"
	"fmt"
	"time"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type Service struct {
	Name        string            `json:"name"`
	Namespace   string            `json:"namespace"`
	AppSelector string            `json:"-"`
	Labels      map[string]string `json:"labels,omitempty"`
	Annotations map[string]string `json:"annotations,omitempty"`
	Ports       []int32

	CreatedAt time.Time `json:"created_at"`
}

func (s *Service) ToJSON() string {
	response, _ := json.Marshal(s)
	return string(response)
}

func (s *Service) Marshal() (apiv1.Service, error) {
	var ports []apiv1.ServicePort
	for _, p := range s.Ports {
		ports = append(ports, apiv1.ServicePort{
			Name:       fmt.Sprintf("port-%d", p),
			Protocol:   apiv1.ProtocolTCP,
			Port:       p,
			TargetPort: intstr.FromInt(int(p)),
		})
	}
	service := apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: apiv1.ServiceSpec{
			Selector: map[string]string{
				"app": s.AppSelector,
			},
			Ports: ports,
		},
	}
	// fmt.Println(service)

	return service, nil
}

func (s *Service) Unmarshal(service *apiv1.Service) error {
	s.Name = service.GetName()
	s.Namespace = service.GetNamespace()
	s.CreatedAt = service.GetCreationTimestamp().Time

	s.Labels = make(map[string]string)
	if labels := service.GetLabels(); labels != nil {
		s.Labels = labels
	}

	s.Annotations = make(map[string]string)
	if annotations := service.GetAnnotations(); annotations != nil {
		s.Annotations = annotations
	}

	return nil
}
