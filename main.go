package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
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

	// Sort all pods within an owner
	for _, o := range owners {
		sort.Slice(o.Pods, func(i, j int) bool {
			return o.Pods[i].maxRequestedUsage() > o.Pods[j].maxRequestedUsage()
		})
	}

	// sort owners according there higest pod
	sort.Slice(owners, func(i, j int) bool {
		return owners[i].Pods[0].maxRequestedUsage() > owners[j].Pods[0].maxRequestedUsage()
	})

	rowConfigAutoMerge := table.RowConfig{AutoMerge: true}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Pod name", "Parent", "Namespace", "CPU", "CPU", "CPU", "Memory", "Memory", "Memory"}, rowConfigAutoMerge)
	t.AppendHeader(table.Row{"", "", "", "used", "request", "limit", "used", "request", "limit"})

	for _, o := range owners {
		t.AppendSeparator()
		for _, p := range o.Pods {
			t.AppendRows([]table.Row{
				{
					p.Name,
					o.Kind,
					o.Namespace,
					p.Usage.Cpu,
					fmt.Sprintf("%d (%.2f%%)", p.Requested.Cpu, p.RequestedCpuUsage()),
					fmt.Sprintf("%d (%.2f%%)", p.Limit.Cpu, p.LimitCpuUsage()),
					fmt.Sprintf("%dMi", p.Usage.MemoryAsMebibyte()),
					fmt.Sprintf("%dMi (%.2f%%)", p.Requested.MemoryAsMebibyte(), p.RequestedMemUsage()),
					fmt.Sprintf("%dMi (%.2f%%)", p.Limit.MemoryAsMebibyte(), p.LimitMemUsage()),
				},
			})
		}
		t.AppendSeparator()
	}

	t.Render()
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
