package main

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getResources(config *rest.Config, ns string) ([]Owner, error) {
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

	results := make(map[string]*Owner)

	for _, p := range podList.Items {
		for _, c := range p.Spec.Containers {
			if !strings.HasPrefix(p.Name, c.Name) {
				continue
			}
			pod := Pod{
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

			ownerName := "None"
			ownerKind := "Unknown"

			if len(p.OwnerReferences) > 0 {
				or := p.OwnerReferences[0]
				ownerName = or.Name
				ownerKind = or.Kind
			}

			owner := results[ownerName]
			if owner == nil {
				owner = &Owner{
					Name: ownerName,
					Kind: ownerKind,
				}
				results[ownerName] = owner
			}
			owner.Pods = append(owner.Pods, pod)
		}
	}

	// Drop the map before returning
	slResults := make([]Owner, 0, len(results))

	for _, value := range results {
		slResults = append(slResults, *value)
	}

	return slResults, nil
}
