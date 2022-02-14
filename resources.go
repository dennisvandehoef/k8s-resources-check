package main

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type PodResources struct {
	Name      string
	Requested Resource
	Limit     Resource
}

type Resource struct {
	Cpu    int64
	Memory int64
}

func getResources(config *rest.Config, ns string) (map[string]PodResources, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	api := clientset.CoreV1()
	listOptions := metav1.ListOptions{
		LabelSelector: "",
		FieldSelector: "",
	}
	podList, err := api.Pods(ns).List(context.TODO(), listOptions)
	if err != nil {
		return nil, err
	}

	results := make(map[string]PodResources)

	for _, p := range podList.Items {
		for _, c := range p.Spec.Containers {
			if !strings.HasPrefix(p.Name, c.Name) {
				continue
			}
			pod := PodResources{
				Name: p.Name,
				Requested: Resource{
					Cpu:    c.Resources.Requests.Cpu().MilliValue(),
					Memory: c.Resources.Requests.Memory().MilliValue(),
				},
				Limit: Resource{
					Cpu:    c.Resources.Limits.Cpu().MilliValue(),
					Memory: c.Resources.Limits.Memory().MilliValue(),
				},
			}

			results[p.Name] = pod
		}
	}

	return results, nil
}
