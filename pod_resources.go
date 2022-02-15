package main

import "math"

type PodResources struct {
	Name      string
	Namespace string
	Requested Resource
	Limit     Resource
	Usage     Resource
}

type Resource struct {
	Cpu    int64
	Memory int64
}

func (pr *PodResources) RequestedMemUsage() float64 {
	if pr.Requested.Memory == 0.0 {
		return 0.0
	}
	return (float64(pr.Usage.Memory) / float64(pr.Requested.Memory)) * 100.0
}

func (pr *PodResources) LimitMemUsage() float64 {
	if pr.Limit.Memory == 0.0 {
		return 0.0
	}
	return (float64(pr.Usage.Memory) / float64(pr.Limit.Memory)) * 100.0
}

func (pr *PodResources) RequestedCpuUsage() float64 {
	if pr.Requested.Cpu == 0.0 {
		return 0.0
	}
	return (float64(pr.Usage.Cpu) / float64(pr.Requested.Cpu)) * 100.0
}

func (pr *PodResources) LimitCpuUsage() float64 {
	if pr.Limit.Cpu == 0.0 {
		return 0.0
	}
	return (float64(pr.Usage.Cpu) / float64(pr.Limit.Cpu)) * 100.0
}

func (pr *PodResources) maxRequestedUsage() float64 {
	return math.Max(pr.RequestedCpuUsage(), pr.RequestedMemUsage())
}
