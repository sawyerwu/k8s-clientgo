package types

import (
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"time"
)

type SwuPod struct {
	Namespace string
	Name      string
	Status    corev1.PodStatus
}

func (swuPod *SwuPod) CreatePod(clientset *kubernetes.Clientset) *corev1.Pod {
	var w watch.Interface
	pod := getNginxPod()
	result, err := clientset.CoreV1().Pods(swuPod.Namespace).Create(pod)
	if err != nil {
		panic(err.Error())
	}

	status := result.Status
	fmt.Println("pod current status", status.Phase)

	if w, err = clientset.CoreV1().Pods(swuPod.Namespace).Watch(metav1.ListOptions{
		Watch:           true,
		ResourceVersion: result.ResourceVersion,
	}); err != nil {
		panic(err.Error())
	}

	func() {
		for {
			select {
			case events, ok := <-w.ResultChan():
				if !ok {
					return
				}
				fmt.Println("event: ", events)
				result = events.Object.(*corev1.Pod)
				fmt.Println("Pod status:", result.Status.Phase)
				status = result.Status
				if result.Status.Phase != corev1.PodPending {
					w.Stop()
				}
			case <-time.After(10 * time.Second):
				fmt.Println("timeout to wait for pod active")
				w.Stop()
			}
		}
	}()

	return result
}

func getNginxPod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "nginx-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "nginx-c",
					Image: "nginx:alpine",
				},
			},
		},
	}
}
