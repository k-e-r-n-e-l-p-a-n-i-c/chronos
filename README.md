# Chronos

A custom Kubernetes Controller to watch Pods in all namespaces

### Introduction

Chronos is a kubernetes controller that I wrote to understand the internal workings of controllers. I followed a pattern that was widely used by many production ready custom controllers like Kubewatch and Nginx ingress. Important topics that you would understand are:

 1) Informers
 2) Work Queues

The command line is written using Cobra CLI Framework.

### Installation

After cloning the repository,you can run this program by executing:

```
#/>go run main.go -k {path to the kubeconfig}
                    or  
You can also choose to build the program and run it.  
                    or  
Build a docker image and deploy it as a pod in a Kubernetes cluster and it will use the Incluster Config file to communicate with API server.Create a role and rolebinding to provide access to Kubernetes Object store.
```

Check the help menu by passing -h or --help:
```
#/>go run main.go -h  
Kubernetes event collector and notifier  

Usage:  
  chronos [flags]  
  chronos [command]  

Available Commands:  
  help        Help about any command  
  version     chronos version  

Flags:  
  -h, --help                help for chronos  
  -k, --kubeconfig string   path to kubeconfig file  

Use "chronos [command] --help" for more information about a command.  
```
### What does it do

The current version of Chronos tracks the status of Pods in all namespaces and logs on the command line
whenever it detects a state change for a Pod.

### References

Find the below links that I refered to write the controller:
1) https://engineering.bitnami.com/articles/kubewatch-an-example-of-kubernetes-custom-controller.html
2) https://github.com/bitnami-labs/kubewatch/blob/master/pkg/controller/controller.go
