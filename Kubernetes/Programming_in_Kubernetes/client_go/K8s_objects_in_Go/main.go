package main

import (
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func main() {
	// Construct a Pod object in memory
	pod := &corev1.Pod{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Pod",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo-pod",
			Namespace: "k8s-learning",
			Labels: map[string]string{
				"app":     "demo",
				"version": "v1",
				"owner":   "mick",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx",
					Image: "nginx:alpine",
					Ports: []corev1.ContainerPort{
						{
							Name:          "http",
							ContainerPort: 80,
							Protocol:      corev1.ProtocolTCP,
						},
					},
				},
			},
		},
	}

	// Convert the Go object into nice YAML
	yamlBytes, err := yaml.Marshal(pod)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling to YAML: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(yamlBytes))
}
