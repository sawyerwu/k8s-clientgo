package types

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type SwuNamespace struct {
	Name string
}

func (swuNs *SwuNamespace) CreateNamespace(clientset *kubernetes.Clientset) *corev1.Namespace {
	result, err := clientset.CoreV1().Namespaces().Get(swuNs.Name, metav1.GetOptions{})
	if err != nil {
		if errors.IsNotFound(err) {
			_, err := clientset.CoreV1().Namespaces().Create(&corev1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: swuNs.Name}})
			if err != nil {
				panic(err.Error())
			}
		}
	}

	return result
}
