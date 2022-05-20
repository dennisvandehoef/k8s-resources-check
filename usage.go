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

func getUsage(config *rest.Config, ns string) (map[string]Resource, error) {
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

	result := make(map[string]Resource)
	for _, pod := range metrics.Items {
		for _, container := range pod.Containers {
			if !strings.HasPrefix(pod.Name, container.Name) {
				continue
			}

			r := Resource{
				Cpu:    container.Usage.Cpu().MilliValue(),
				Memory: container.Usage.Memory().MilliValue(),
			}

			result[pod.Name] = r
		}
	}

	return result, nil
}
