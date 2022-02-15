package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	resources := []PodResources{}
	namespaces, err := getNamespaces(config)
	if err != nil {
		log.Fatal(err)
	}

	for _, ns := range namespaces {
		err = processNamespace(config, ns, &resources)
		if err != nil {
			log.Fatal(err)
		}
	}

	sort.Slice(resources, func(i, j int) bool {
		return resources[i].maxRequestedUsage() > resources[j].maxRequestedUsage()
	})

	fmt.Println("pod (ns)| CPU use/request/limit = request%/limi% | MEM use/request/limit = request%/limi%")
	for _, r := range resources {
		fmt.Printf("%s (%s) | %dm/%dm/%dm = %.2f%%/%.2f%% | %dM/%dM/%dM = %.2f%%/%.2f%%\n",
			r.Name,
			r.Namespace,
			r.Usage.Cpu,
			r.Requested.Cpu,
			r.Limit.Cpu,
			r.RequestedCpuUsage(),
			r.LimitCpuUsage(),
			r.Usage.Memory/(1000*1000*1000),
			r.Requested.Memory/(1000*1000*1000),
			r.Limit.Memory/(1000*1000*1000),
			r.RequestedMemUsage(),
			r.LimitMemUsage(),
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
