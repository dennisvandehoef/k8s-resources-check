package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	owners := []Owner{}
	namespaces, err := getNamespaces(config)
	if err != nil {
		log.Fatal(err)
	}

	for _, ns := range namespaces {
		err = processNamespace(config, ns, &owners)
		if err != nil {
			log.Fatal(err)
		}
	}

	// sort.Slice(resources, func(i, j int) bool {
	// 	return resources[i].maxRequestedUsage() > resources[j].maxRequestedUsage()
	// })

	fmt.Println("controller name (type) namespace")
	fmt.Println("- pod | CPU use/request/limit = request%/limit% | MEM use/request/limit = request%/limit%")
	fmt.Println("--------------------------------------")
	for _, o := range owners {
		fmt.Printf("%s (%s) %s:\n", o.Name, o.Kind, o.Namespace)
		for _, p := range o.Pods {
			fmt.Printf("- %s | %dm/%dm/%dm = %.2f%%/%.2f%% | %dMi/%dMi/%dMi = %.2f%%/%.2f%%\n",
				p.Name,
				p.Usage.Cpu,
				p.Requested.Cpu,
				p.Limit.Cpu,
				p.RequestedCpuUsage(),
				p.LimitCpuUsage(),
				p.Usage.MemoryAsMebibyte(),
				p.Requested.MemoryAsMebibyte(),
				p.Limit.MemoryAsMebibyte(),
				p.RequestedMemUsage(),
				p.LimitMemUsage(),
			)
		}

	}
}

func processNamespace(config *rest.Config, ns string, owners *[]Owner) error {
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
		for _, p := range o.Pods {
			p.Usage = usage[p.Name].Usage
		}

		*owners = append(*owners, o)
	}
	return nil
}
