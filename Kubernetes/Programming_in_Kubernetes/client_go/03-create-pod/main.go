package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	//"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// 1. Connect to the cluster
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err.Error())
	}
	kubeconfig := filepath.Join(home, ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := "k8s-learning"
	podName := "demo-pod"

	// 2. Delete existing Pod if it exists
	_, err = clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err == nil {
		fmt.Printf("Pod %s exists. Deleting first...\n", podName)
		deletePolicy := metav1.DeletePropagationForeground
		err = clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{
			PropagationPolicy: &deletePolicy,
		})
		if err != nil {
			panic(err.Error())
		}
		// Brief wait for deletion
		time.Sleep(3 * time.Second)
	}

	// 3. Define the Pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podName,
			Namespace: namespace,
			Labels: map[string]string{
				"app":     "demo",
				"version": "v1",
				"owner":   "mick",
				"created": "by-client-go",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:alpine",
					Ports: []corev1.ContainerPort{
						{Name: "http", ContainerPort: 80, Protocol: corev1.ProtocolTCP},
					},
				},
			},
		},
	}

	// 4. Create the Pod
	fmt.Printf("Creating Pod %s...\n", podName)
	_, err = clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Pod created successfully.")

	// 5. Watch for the Pod to reach Running status
	fmt.Println("Watching for Pod to become Running...")

	watcher, err := clientset.CoreV1().Pods(namespace).Watch(context.TODO(), metav1.ListOptions{
		FieldSelector: fmt.Sprintf("metadata.name=%s", podName),
	})
	if err != nil {
		panic(err.Error())
	}
	defer watcher.Stop()

	for event := range watcher.ResultChan() {
		p, ok := event.Object.(*corev1.Pod)
		if !ok {
			continue
		}

		fmt.Printf("Pod phase: %s\n", p.Status.Phase)

		switch p.Status.Phase {
		case corev1.PodRunning:
			fmt.Println("🎉 Pod is now Running!")
			return
		case corev1.PodFailed, corev1.PodSucceeded:
			fmt.Printf("❌ Pod ended with phase: %s\n", p.Status.Phase)
			return
		}
	}
}
