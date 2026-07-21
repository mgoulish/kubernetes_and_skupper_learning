package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	createFlag := flag.Bool("create", false, "Create the demo pod")
	deleteFlag := flag.Bool("delete", false, "Delete the demo pod")
	flag.Parse()

	if !*createFlag && !*deleteFlag {
		fmt.Println("Usage: go run . -create  OR  go run . -delete")
		return
	}

	// Connect to cluster
	home, _ := os.UserHomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, _ := clientcmd.BuildConfigFromFlags("", kubeconfig)
	clientset, _ := kubernetes.NewForConfig(config)

	namespace := "k8s-learning"
	podName := "demo-pod"

	if *deleteFlag {
		fmt.Printf("Deleting Pod %s...\n", podName)
		err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{
			PropagationPolicy: func() *metav1.DeletionPropagation { p := metav1.DeletePropagationForeground; return &p }(),
		})
		if err != nil {
			fmt.Printf("Error deleting: %v\n", err)
		} else {
			fmt.Println("✅ Pod deleted successfully.")
		}
		return
	}

	if *createFlag {
		pod := &corev1.Pod{
			ObjectMeta: metav1.ObjectMeta{
				Name:      podName,
				Namespace: namespace,
				Labels: map[string]string{
					"app":     "demo",
					"owner":   "mick",
					"created": "by-client-go",
				},
			},
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name:  "nginx",
						Image: "nginx:alpine",
						Ports: []corev1.ContainerPort{{ContainerPort: 80}},
					},
				},
			},
		}

		fmt.Printf("Creating Pod %s...\n", podName)
		_, err := clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
		if err != nil {
			fmt.Printf("Error creating: %v\n", err)
		} else {
			fmt.Println("✅ Pod created successfully.")
		}
	}
}
