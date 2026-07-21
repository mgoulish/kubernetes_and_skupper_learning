package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	home, _ := os.UserHomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	dynClient, _ := dynamic.NewForConfig(config)

	// Define the GVR (Group, Version, Resource) for Pods
	podGVR := schema.GroupVersionResource{
		Group:    "",
		Version:  "v1",
		Resource: "pods",
	}

	// List Pods using dynamic client
	pods, err := dynClient.Resource(podGVR).Namespace("k8s-learning").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Found %d pods:\n", len(pods.Items))
	for _, p := range pods.Items {
		fmt.Printf("- %s\n", p.GetName())
	}
}
