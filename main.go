package main

import (
	"context"
	"log"
	"os"
	"path"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

func must1[T any](x T, err error) T {
	must(err)
	return x
}

func createRestConfig() (*rest.Config, error) {
	if os.Getenv("KUBERNETES_SERVICE_HOST") != "" {

		return rest.InClusterConfig()
	} else if kubeconfig, ok := os.LookupEnv("KUBECONFIG"); ok {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		home := homedir.HomeDir()

		return clientcmd.BuildConfigFromFlags("", path.Join(home, ".kube", "config"))
	}
}

func createClient() (*kubernetes.Clientset, error) {
	config, err := createRestConfig()
	if err != nil {
		return nil, err
	}

	return kubernetes.NewForConfig(config)
}

func main() {
	client := must1(createClient())

	pods := must1(client.CoreV1().Pods("default").List(context.Background(), v1.ListOptions{}))

	for _, pod := range pods.Items {
		log.Printf("Pod: %s\n", pod.Name)
	}
}
