package main

import (
	"flag"
	"fmt"
	"github.com/sawyerwu/k8s-clientgo/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

func main() {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	fmt.Println(*kubeconfig)

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// waiting for test
	types.WatchPod(clientset)

	nsObj := &types.SwuNamespace{
		Name: "demo",
	}
	ns := nsObj.CreateNamespace(clientset)
	fmt.Println(ns.Name)

	podObj := &types.SwuPod{
		Name:      "nginx-pod",
		Namespace: nsObj.Name,
	}
	pod := podObj.CreatePod(clientset)
	fmt.Println(pod.Name)

	/*pods, err := clientset.CoreV1().Pods("application").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}*/
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
