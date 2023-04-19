package kube

import (
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type PullPolicy string

const (
	PullAlways       PullPolicy = "Always"
	PullNever        PullPolicy = "Never"
	PullIfNotPresent PullPolicy = "IfNotPresent"
)

type Resource struct {
	Requests string `json:"requests"`
	Limits   string `json:"limits"`
}

type Port struct {
	Type string `json:"type"`
	From int32  `json:"from"`
	To   int32  `json:"to"`
}

func (p *Port) Marshal() *apiv1.ContainerPort {
	return &apiv1.ContainerPort{
		Name:          "http",
		Protocol:      apiv1.ProtocolTCP,
		HostPort:      p.From,
		ContainerPort: p.To,
	}
}

func (p *Port) Unmarshal(port *apiv1.ContainerPort) {
	p.Type = string(port.Protocol)
	p.From = port.HostPort
	p.To = port.ContainerPort
}

type Environment struct {
	Type  string `json:"type"`
	From  string `json:"from"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (e *Environment) Marshal() *apiv1.EnvVar {
	return &apiv1.EnvVar{
		Name:  e.Name,
		Value: e.Value,
	}
}

func (e *Environment) Unmarshal(env *apiv1.EnvVar) {
	e.Type = "value"
	e.From = "value"
	e.Name = env.Name
	e.Value = env.Value
}

type Probe struct {
	TimeoutSeconds      int32 `json:"timeout_seconds"`
	InitialDelaySeconds int32 `json:"initial_delay_seconds"`
	PeriodSeconds       int32 `json:"period_seconds"`
	SuccessThreshold    int32 `json:"success_threshold"`
	FailureThreshold    int32 `json:"failure_threshold"`
}

func (p *Probe) Marshal() *apiv1.Probe {
	return &apiv1.Probe{
		TimeoutSeconds:      p.TimeoutSeconds,
		InitialDelaySeconds: p.InitialDelaySeconds,
		PeriodSeconds:       p.PeriodSeconds,
		SuccessThreshold:    p.SuccessThreshold,
		FailureThreshold:    p.FailureThreshold,
	}

}

func (p *Probe) Unmarshal(probe apiv1.Probe) {
	p.TimeoutSeconds = probe.TimeoutSeconds
	p.InitialDelaySeconds = probe.InitialDelaySeconds
	p.PeriodSeconds = probe.PeriodSeconds
	p.SuccessThreshold = probe.SuccessThreshold
	p.FailureThreshold = probe.FailureThreshold
}

type Container struct {
	Name           string        `json:"name"`
	Image          string        `json:"image"`
	PullPolicy     PullPolicy    `json:"pull_policy"`
	Ports          []Port        `json:"ports,omitempty"`
	Environments   []Environment `json:"environments,omitempty"`
	LivenessProbe  *Probe        `json:"liveness,omitempty"`
	ReadinessProbe *Probe        `json:"readiness,omitempty"`
	CPU            *Resource     `json:"cpu"`
	Memory         *Resource     `json:"memory"`
}

func (c *Container) Marshal() (container apiv1.Container) {
	container = apiv1.Container{
		Name:  c.Name,
		Image: c.Image,
		Resources: apiv1.ResourceRequirements{
			Requests: apiv1.ResourceList{},
			Limits:   apiv1.ResourceList{},
		},
	}

	for _, port := range c.Ports {
		container.Ports = append(container.Ports, *port.Marshal())
	}

	for _, env := range c.Environments {
		container.Env = append(container.Env, *env.Marshal())
	}

	if c.CPU != nil {
		container.Resources.Requests["cpu"] = resource.MustParse(c.CPU.Requests)
		container.Resources.Limits["cpu"] = resource.MustParse(c.CPU.Limits)
	}

	if c.Memory != nil {
		container.Resources.Requests["memory"] = resource.MustParse(c.Memory.Requests)
		container.Resources.Limits["memory"] = resource.MustParse(c.Memory.Limits)
	}

	if c.LivenessProbe != nil {
		container.LivenessProbe = c.LivenessProbe.Marshal()
	}

	if c.ReadinessProbe != nil {
		container.ReadinessProbe = c.ReadinessProbe.Marshal()
	}

	return
}

func (c *Container) Unmarshal(container apiv1.Container) {
	c.Name = container.Name
	c.Image = container.Image

	switch container.ImagePullPolicy {
	case apiv1.PullAlways:
		c.PullPolicy = PullAlways
	case apiv1.PullNever:
		c.PullPolicy = PullNever
	case apiv1.PullIfNotPresent:
		c.PullPolicy = PullIfNotPresent
	}

	c.Ports = make([]Port, 0)
	c.Environments = make([]Environment, 0)

	for _, port := range container.Ports {
		c.Ports = append(c.Ports, Port{
			Type: string(port.Protocol),
			From: port.HostPort,
			To:   port.ContainerPort,
		})
	}

	for _, env := range container.Env {
		c.Environments = append(c.Environments, Environment{
			Type:  "value",
			From:  "",
			Name:  env.Name,
			Value: env.Value,
		})
	}

	if liveness := container.LivenessProbe; liveness != nil {
		c.LivenessProbe.Unmarshal(*liveness)
	}

	if readiness := container.ReadinessProbe; readiness != nil {
		c.ReadinessProbe.Unmarshal(*readiness)
	}

	c.CPU = &Resource{
		Requests: fmt.Sprint(container.Resources.Requests.Cpu().Value()),
		Limits:   fmt.Sprint(container.Resources.Limits.Cpu().Value()),
	}

	c.Memory = &Resource{
		Requests: fmt.Sprint(container.Resources.Requests.Memory().Value()),
		Limits:   fmt.Sprint(container.Resources.Limits.Memory().Value()),
	}

}
