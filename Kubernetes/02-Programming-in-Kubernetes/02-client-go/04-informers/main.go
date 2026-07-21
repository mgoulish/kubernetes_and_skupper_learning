package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// Connect to cluster
	home, _ := os.UserHomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	clientset, _ := kubernetes.NewForConfig(config)

	// Create a shared informer factory (for namespace k8s-learning)
	factory := informers.NewSharedInformerFactoryWithOptions(clientset, 30*time.Second, informers.WithNamespace("k8s-learning"))

	// Create a Pod informer
	podInformer := factory.Core().V1().Pods()

	// Add event handlers
	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("[ADD] Pod %s created\n", pod.Name)
		},
		UpdateFunc: func(old, new interface{}) {
			//oldPod := old.(*corev1.Pod)
			newPod := new.(*corev1.Pod)
			fmt.Printf("[UPDATE] Pod %s changed (phase: %s)\n", newPod.Name, newPod.Status.Phase)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			fmt.Printf("[DELETE] Pod %s deleted\n", pod.Name)
		},
	})

	// Start the informer (this starts the watch)
	stopCh := make(chan struct{})
	factory.Start(stopCh)
	fmt.Println("Pod informer started. Watching for changes...")

	// Wait forever (Ctrl+C to stop)
	<-stopCh
}
