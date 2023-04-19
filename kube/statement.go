package kube

import (
	"reflect"

	"k8s.io/client-go/kubernetes"
)

type Statement struct {
	// A reference to Kubernetes API
	Client *kubernetes.Clientset

	// An object provided during updates
	Origin reflect.Value

	// An object (or list) provided during
	// the creation or while fetching objets
	Target reflect.Value

	// A reference to the API of a given
	// (target) object through reflection
	API KubeAPI
}

func (s *Statement) GetOrigin() interface{} {
	return s.Origin.Elem().Interface()
}

func (s *Statement) GetTarget() interface{} {
	return s.Target.Elem().Interface()
}

func (s *Statement) SetTarget(value interface{}) {
	s.Target.Elem().Set(reflect.ValueOf(value))
}
