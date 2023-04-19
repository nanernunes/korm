package korm

import (
	"github.com/nanernunes/korm/kube"
)

type Config struct {
	Logger bool
}

func Open(address string, credentials kube.Credential, config *Config) (*kube.Kubernetes, error) {

	clientset, err := credentials.Login(address)
	if err != nil {
		return nil, err
	}

	kubernetes := &kube.Kubernetes{
		Client: clientset,
		Statement: &kube.Statement{
			Client: clientset,
		},
	}

	return kubernetes, nil

}
