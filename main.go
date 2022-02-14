package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	ns := "sales"
	kubeconfig := filepath.Join(
		os.Getenv("HOME"), ".kube", "config",
	)
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	resources, err := getResources(config, ns)
	if err != nil {
		log.Fatal(err)
	}
	usage, err := getUsage(config, ns)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("pod | CPU use/request/limit | MEM use/request/limit")
	for k, pu := range usage {
		r := resources[k]
		fmt.Printf("%s | %dm/%dm/%dm | %dM/%dM/%dM\n",
			k,
			pu.Cpu,
			r.Requested.Cpu,
			r.Limit.Cpu,
			pu.Memory/(1000*1000*1000),
			r.Requested.Memory/(1000*1000*1000),
			r.Limit.Memory/(1000*1000*1000),
		)
	}
}
