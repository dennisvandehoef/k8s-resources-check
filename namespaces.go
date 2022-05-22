package main

import (
	"context"
	"fmt"

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

func processNamespace(config *rest.Config, ns string, owners *[]Owner) error {
	fmt.Println("Processing namespace: " + ns)

	reservedResourceOwners, err := getResources(config, ns)
	if err != nil {
		return err
	}

	usage, err := getUsage(config, ns)
	if err != nil {
		return err
	}

	for _, o := range reservedResourceOwners {
		o.Namespace = ns
		for i, p := range o.Pods {
			o.Pods[i].Usage = usage[p.Name]
		}

		*owners = append(*owners, o)
	}
	return nil
}
