package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

func main() {
	ns := "sales"
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	resources := []PodResources{}

	err = processNamespace(config, ns, &resources)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("pod (ns)| CPU use/request/limit | MEM use/request/limit")
	for _, r := range resources {
		fmt.Printf("%s (%s) | %dm/%dm/%dm | %dM/%dM/%dM\n",
			r.Name,
			r.Namespace,
			r.Usage.Cpu,
			r.Requested.Cpu,
			r.Limit.Cpu,
			r.Usage.Memory/(1000*1000*1000),
			r.Requested.Memory/(1000*1000*1000),
			r.Limit.Memory/(1000*1000*1000),
		)
	}
}

func processNamespace(config *rest.Config, ns string, resources *[]PodResources) error {
	reservedResources, err := getResources(config, ns)
	if err != nil {
		return err
	}
	usage, err := getUsage(config, ns)
	if err != nil {
		return err
	}

	for k, pu := range usage {
		rr := reservedResources[k]

		r := PodResources{
			Name:      k,
			Namespace: ns,
			Usage:     pu.Usage,
			Requested: rr.Requested,
			Limit:     rr.Limit,
		}

		*resources = append(*resources, r)
	}

	return nil
}
