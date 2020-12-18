package utils

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

//GetClient returns a kubernetes client
func GetClient(configpath string) (*kubernetes.Clientset, error) {

	if configpath == "" {
		logrus.Info("Using Incluster configuration")
	}
	config, err := clientcmd.BuildConfigFromFlags("", configpath)
	if err != nil {
		logrus.Fatalf("Error occured while reading kubeconfig:%v", err)
	}
	return kubernetes.NewForConfig(config)
}
