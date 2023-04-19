package kube

import (
	"fmt"
	"reflect"

	"k8s.io/client-go/kubernetes"
)

type Kubernetes struct {
	Client    *kubernetes.Clientset
	Statement *Statement

	Error        error
	RowsAffected int64
}

func (k *Kubernetes) newTransaction(target interface{}, origin interface{}) (tx *Kubernetes) {
	tx = &Kubernetes{
		Client: k.Client,
		Statement: &Statement{
			Client: k.Client,
		},
	}

	if origin == nil {
		tx.Statement.Target = reflect.ValueOf(target)
	} else {
		tx.Statement.Target = k.Statement.Target
	}

	tx.Statement.Origin = reflect.ValueOf(origin)

	entity := target
	if entity == nil {
		entity = origin
	}

	switch t := entity.(type) {
	case *Namespace, *[]Namespace:
		tx.Statement.API = NewNamespaceAPI(tx)

	case *Deployment, *[]Deployment:
		tx.Statement.API = NewDeploymentAPI(tx)

	case *Service, *[]Service:
		tx.Statement.API = NewServiceAPI(tx)

	case *HPA, *[]HPA:
		tx.Statement.API = NewHPAAPI(tx)

	default:
		panic(fmt.Sprintf("unsupported type: %s", t))
	}

	return
}

func (k *Kubernetes) Find(entity interface{}) (tx *Kubernetes) {
	tx = k.newTransaction(entity, nil)
	tx.Statement.API.Find()
	return
}

func (k *Kubernetes) Create(entity interface{}) (tx *Kubernetes) {
	tx = k.newTransaction(entity, nil)
	tx.Statement.API.Create()
	return
}

func (k *Kubernetes) Update(entity interface{}) (tx *Kubernetes) {
	// Only in this transaction, the entity is an object with new
	// values (named Origin) we'll use to update the previous value
	// inside the tx.Target defined by a Find, Create, Upsert, etc.
	tx = k.newTransaction(nil, entity)
	tx.Statement.API.Update()
	return
}

func (k *Kubernetes) Upsert(entity interface{}) (tx *Kubernetes) {
	tx = k.newTransaction(entity, nil)
	tx.Statement.API.Upsert()
	return
}

func (k *Kubernetes) Delete(entity interface{}) (tx *Kubernetes) {
	tx = k.newTransaction(entity, nil)
	tx.Statement.API.Delete()
	return
}
