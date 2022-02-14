package main

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	metricsapi "k8s.io/metrics/pkg/apis/metrics"
	metricsv1beta1api "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	metricsclientset "k8s.io/metrics/pkg/client/clientset/versioned"
)

type PodUsage struct {
	Name   string
	Cpu    int64
	Memory int64
}

func getUsage(config *rest.Config, ns string) (map[string]PodUsage, error) {
	metricsClient, err := metricsclientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	versionedMetrics, err := metricsClient.MetricsV1beta1().PodMetricses(ns).List(context.TODO(), metav1.ListOptions{LabelSelector: "", FieldSelector: ""})
	if err != nil {
		return nil, err
	}
	metrics := &metricsapi.PodMetricsList{}
	err = metricsv1beta1api.Convert_v1beta1_PodMetricsList_To_metrics_PodMetricsList(versionedMetrics, metrics, nil)
	if err != nil {
		return nil, err
	}

	result := make(map[string]PodUsage)
	for _, pod := range metrics.Items {
		for _, container := range pod.Containers {
			if !strings.HasPrefix(pod.Name, container.Name) {
				continue
			}

			pod := PodUsage{
				Name:   pod.Name,
				Cpu:    container.Usage.Cpu().MilliValue(),
				Memory: container.Usage.Memory().MilliValue(),
			}

			result[pod.Name] = pod
		}
	}

	return result, nil
}
