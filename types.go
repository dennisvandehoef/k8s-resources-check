package main

import "math"

type Owner struct {
	Name      string
	Namespace string
	Kind      string
	Pods      []Pod
}

type Pod struct {
	Name      string
	Requested Resource
	Limit     Resource
	Usage     Resource
}

type Resource struct {
	Cpu    int64
	Memory int64
}

const (
	BibyteFactor   = float64(1.048576)
	MegabyteFactor = 1000000000 // 1000 * 1000 * 1000
)

func (r *Resource) MemoryAsMegabyte() int64 {
	return r.Memory / MegabyteFactor
}

func (r *Resource) MemoryAsMebibyte() int64 {
	return int64(float64(r.MemoryAsMegabyte()) / BibyteFactor)
}

func (pr *Pod) RequestedMemUsage() float64 {
	if pr.Requested.Memory == 0.0 {
		return 0.0
	}
	return (float64(pr.Usage.Memory) / float64(pr.Requested.Memory)) * 100.0
}

func (pr *Pod) LimitMemUsage() float64 {
	if pr.Limit.Memory == 0.0 {
		return 0.0
	}
	return (float64(pr.Usage.Memory) / float64(pr.Limit.Memory)) * 100.0
}

func (pr *Pod) RequestedCpuUsage() float64 {
	if pr.Requested.Cpu == 0.0 {
		return 0.0
	}
	return (float64(pr.Usage.Cpu) / float64(pr.Requested.Cpu)) * 100.0
}

func (pr *Pod) LimitCpuUsage() float64 {
	if pr.Limit.Cpu == 0.0 {
		return 0.0
	}
	return (float64(pr.Usage.Cpu) / float64(pr.Limit.Cpu)) * 100.0
}

func (pr *Pod) maxRequestedUsage() float64 {
	return math.Max(pr.RequestedCpuUsage(), pr.RequestedMemUsage())
}
