package cnf

import (
	"flag"
	"fmt"
	"os"
)

var (
	Namespace     string
	AllNamespaces bool = false
)

const usage = `Example usage:
  k8s-resources-check -all-namespaces

Options:
  -n, --Namespace NAMESPACE     The namespace to report for
  --all-namespaces              Report for all namespaces`

func FromFlags() {
	flag.StringVar(&Namespace, "namespace", "", "")
	flag.StringVar(&Namespace, "n", "", "")
	flag.BoolVar(&AllNamespaces, "all-namespaces", false, "")

	flag.Usage = func() { fmt.Println(usage) }

	flag.Parse()

	validate()
}

func validate() {
	if !AllNamespaces {
		if len(Namespace) <= 0 {
			configError("No namespace is given!")
		}
	}
}

func configError(error string) {
	fmt.Println("ERROR: " + error)
	fmt.Println("")
	fmt.Println(usage)
	os.Exit(1)
}
