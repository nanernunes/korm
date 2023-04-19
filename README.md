# KORM

[![Go](https://github.com/nanernunes/korm/actions/workflows/go.yml/badge.svg)](https://github.com/nanernunes/korm/actions/workflows/go.yml)
[![version](https://img.shields.io/github/tag/nanernunes/korm.svg)](https://github.com/nanernunes/korm/releases/latest)
[![GoDoc](https://godoc.org/github.com/korm?status.png)](https://godoc.org/github.com/nanernunes/korm)
[![license](https://img.shields.io/github/license/nanernunes/korm.svg)](../LICENSE.md)
[![LoC](https://tokei.rs/b1/github/nanernunes/korm?category=lines)](https://github.com/nanernunes/korm)
[![codecov](https://codecov.io/gh/nanernunes/korm/branch/master/graph/badge.svg)](https://codecov.io/gh/nanernunes/korm)

The Lite ORM library for Kubernetes, aims to be developer friendly.

> **Completely inspired by GORM (Go ORM for Databases, thanks devs :heart:)**

### Installation

```bash
go get github.com/nanernunes/korm
```

### Usage

```go
package main

import (
    "github.com/nanernunes/korm"
    "github.com/nanernunes/korm/kube"
)

func main() {
    k8s, err := korm.Open(
        "https://kubernetes.docker.internal:6443",
        kube.KubeConfigCredential{Path: "~/.kube/config"},
        &korm.Config{}
    )

    if err != nil {
      panic("failed to connect kubernetes")
    }

    // Create
    k8s.Create(&kube.Namespace{Name: "korm"})

    // Find - all namespaces
    var namespaces []kube.Namespace
    k8s.Find(&namespaces)

    // Find - a namespace with name korm
    namespace := kube.Namespace{Name: "korm"}
    k8s.Find(&namespace)

    // Update - namespace's labels
    labels := map[string]string{"framework": "korm"}
    k8s.Find(&namespace).Update(&kube.Namespace{Labels: labels})

    // Delete - a namespace
    k8s.Delete(&namespace)
}
```

### Authentication Types

```go
    var credentials kube.Credential

    // Based on Kubernetes' config file
    credentials = kube.KubeConfigCredential{
        Path: "~/.kube/config",
    }

    // Based on Kubernetes' Service Account
    credentials = kube.ServiceAccountCredential{
        Account: "",
        Token:   "",
    }

    // Based on Kubernetes' RBAC User and Password
    credentials = kube.UserPasswordCredential{
        Username: "",
        Password: "",
    }

    korm.Open("kubernetes-address", credentials, &korm.Config{})
```

### Query

```go
    var deployments []kube.Deployment
    result := k8s.Find(&deployments)

    result.RowsAffected // returns found records count, equals `len(deployments)`
    result.Error        // returns error or nil

    // check error ErrRecordNotFound
    errors.Is(result.Error, korm.ErrRecordNotFound)

    // k8s methods fill the object or collection with its data
    for _, deploy := range deployments {
        fmt.Println(deploy.Name)
        fmt.Println(deploy.Labels)
    }
```

### Creating a Namespace

```go
    namespace := kube.Namespace{
        Name:   "korm",
        Labels: map[string]string{
            "framework": "korm",
            "language": "Go",
        }
    }

    k8s.Create(&namespace)
```

### Creating a Deployment

```go
    deployment := kube.Deployment{
        Name:      "helloworld",
        Namespace: "korm",
        Containers: []kube.Container{
            {
                Name:  "helloworld",
                Image: "nginx:latest",
                Ports: []kube.Port{
                    {
                        Port: 80,
                    },
                },
                Environments: []kube.Environment{
                    {
                        Name:  "HELLO",
                        Value: "WORLD",
                    },
                },
            },
        }
    }

    k8s.Create(&deployment)
```

### Creating a Service

```go
    service := kube.Service{
        Name:      "svc-helloworld",
        Namespace: "korm",
    }

    k8s.Create(&service)
```

### Creating a HPA

```go
    hpa := kube.HPA{
        Name:      "helloworld",
        Namespace: "korm",
    }

    k8s.Create(&hpa)
```

### Updating objects

```go
    // Updating fields in a fetched object
    deploy := kube.Deployment{Name: "helloworld"}

    if tx := k8s.Find(&deploy); tx.Error == nil {
        changed := deploy
        changed.Containers[0].Image = "nginx:1.17.1"

        // After fetching an object we can call Update on it
        // with a fully or partialy filled new object and both
        // data will be merged in a final object
        tx.Update(changed)
    }


    // Updating fields without cloning the object
    labels := map[string]string{"framework": "korm"}

    namespace := kube.Namespace{Name: "korm"}

    if tx := k8s.Find(&namespace); tx.Error == nil {
        tx.Update(&kube.Namespace{Labels: labels})
    }

    // Or with no Error handling
    k8s.Find(&namespace).Update(&kube.Namespace{Labels: labels})

```

## Contributing

Contributions to KORM are welcome and appreciated. If you would like to contribute, please open an issue or submit a pull request.

## License

KORM is licensed under the MIT License. See the [LICENSE](LICENSE) file for more information.
