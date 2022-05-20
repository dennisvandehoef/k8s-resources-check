package main

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func getNamespaces(config *rest.Config) ([]string, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	api := clientset.CoreV1()
	list, err := api.Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var result []string

	for _, n := range list.Items {
		result = append(result, n.Name)
	}

	return result, nil
}
