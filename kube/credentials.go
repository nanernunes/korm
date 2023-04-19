package kube

import (
	"os"
	"path/filepath"
	"strings"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Credential interface {
	Login(string) (*kubernetes.Clientset, error)
}

type KubeConfigCredential struct {
	Path string
}

func (k KubeConfigCredential) GetPath() string {
	if strings.HasPrefix(k.Path, "~/") {
		dirname, _ := os.UserHomeDir()
		return filepath.Join(dirname, k.Path[2:])
	}

	return k.Path
}

func (k KubeConfigCredential) Login(address string) (*kubernetes.Clientset, error) {

	config, err := clientcmd.BuildConfigFromFlags(address, k.GetPath())
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil

}

type ServiceAccountCredential struct {
	Name  string
	Token string
}

func (s ServiceAccountCredential) Login(address string) (*kubernetes.Clientset, error) {
	return nil, nil
}

type UserPasswordCredential struct {
	Username string
	Password string
}

func (u UserPasswordCredential) Login(address string) (*kubernetes.Clientset, error) {
	return nil, nil
}
