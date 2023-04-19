package kube

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type ServiceAPI struct {
	Kubernetes *Kubernetes
}

func (s *ServiceAPI) GetClient() v1.ServiceInterface {
	namespace := apiv1.NamespaceDefault
	if ns := s.Kubernetes.Statement.GetTarget().(Service).Namespace; ns != "" {
		namespace = ns
	}
	return s.Kubernetes.Client.CoreV1().Services(namespace)
}

func NewServiceAPI(kubernetes *Kubernetes) *ServiceAPI {
	return &ServiceAPI{Kubernetes: kubernetes}
}

func (s *ServiceAPI) Find() {
	stmt := s.Kubernetes.Statement

	switch stmt.GetTarget().(type) {
	case Service:
		service := stmt.GetTarget().(Service)
		serviceSpec, _ := s.GetClient().Get(
			context.TODO(), service.Name, metav1.GetOptions{})

		service.Unmarshal(serviceSpec)
		stmt.SetTarget(service)

	case []Service:
		stmt.SetTarget(
			[]Service{{Name: "matrix"}},
		)
	}
}

func (s *ServiceAPI) Create() {
	stmt := s.Kubernetes.Statement

	service := stmt.GetTarget().(Service)
	serviceSpec, _ := service.Marshal()

	result, err := s.GetClient().Create(
		context.TODO(), &serviceSpec, metav1.CreateOptions{})

	if err != nil {
		fmt.Println(service.Namespace, service.Name, " -> ", err)
	}

	fmt.Println(result.GetObjectMeta().GetName())

	// return err
	stmt.SetTarget(service)
}

func (s *ServiceAPI) Update() {}
func (s *ServiceAPI) Upsert() {}
func (s *ServiceAPI) Delete() {}
