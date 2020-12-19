package utils

import (
	"github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

//GetClient returns a kubernetes client
func GetClient(configpath string) (*kubernetes.Clientset, error) {

	if configpath == "" {
		logrus.Info("Using Incluster configuration")
		config, err := rest.InClusterConfig()
		if err != nil {
			logrus.Fatalf("Error occured while reading incluster kubeconfig:%v", err)
			return nil, err
		}
		return kubernetes.NewForConfig(config)
	}

	logrus.Infof("Using configuration file:%s", configpath)
	config, err := clientcmd.BuildConfigFromFlags("", configpath)
	if err != nil {
		logrus.Fatalf("Error occured while reading kubeconfig:%v", err)
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
