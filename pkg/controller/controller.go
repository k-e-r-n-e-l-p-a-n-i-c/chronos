package controller

import (
	"fmt"
	"time"

	"github.com/arunprasadmudaliar/chronos/pkg/utils"
	"github.com/sirupsen/logrus"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

type event struct {
	key          string
	eventType    string
	resourceType string
}

type controller struct {
	client   kubernetes.Interface
	informer cache.SharedIndexInformer
	queue    workqueue.RateLimitingInterface
}

func Start(config string) {
	kc, err := utils.GetClient(config)
	if err != nil {
		logrus.Fatal(err)
	}

	factory := informers.NewSharedInformerFactory(kc, 0)
	informer := factory.Core().V1().Pods().Informer()

	/* 	var ctx context.Context
	   	podInformer := cache.NewSharedIndexInformer(
	   		&cache.ListWatch{
	   			ListFunc: func(options meta.ListOptions) (runtime.Object, error) {
	   				return kc.CoreV1().Pods(meta.NamespaceAll).List(ctx, options)
	   			},
	   			WatchFunc: func(options meta.ListOptions) (watch.Interface, error) {
	   				return kc.CoreV1().Pods(meta.NamespaceAll).Watch(ctx, options)
	   			},
	   		},
	   		&v1.Pod{},
	   		0,
	   		cache.Indexers{},
	   	) */
	c := newController(kc, informer)
	stopCh := make(chan struct{})
	defer close(stopCh)

	c.Run(stopCh)

}

func newController(kc kubernetes.Interface, informer cache.SharedIndexInformer) *controller {
	q := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	var event event
	var err error
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			event.key, err = cache.MetaNamespaceKeyFunc(obj)
			event.eventType = "create"
			if err == nil {
				q.Add(event)
			}
			logrus.Infof("Event received of type [%s] for [%s]", event.eventType, event.key)
		},
		UpdateFunc: func(old, new interface{}) {
			event.key, err = cache.MetaNamespaceKeyFunc(old)
			event.eventType = "update"
			if err == nil {
				q.Add(event)
			}
			logrus.Infof("Event received of type [%s] for [%s]", event.eventType, event.key)
		},
		DeleteFunc: func(obj interface{}) {
			event.key, err = cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			event.eventType = "delete"
			if err == nil {
				q.Add(event)
			}
			logrus.Infof("Event received of type [%s] for [%s]", event.eventType, event.key)
		},
	})

	return &controller{
		client:   kc,
		informer: informer,
		queue:    q,
	}
}

func (c *controller) Run(stopper <-chan struct{}) {
	// don't let panics crash the process
	defer utilruntime.HandleCrash()
	// make sure the work queue is shutdown which will trigger workers to end
	defer c.queue.ShutDown()

	logrus.Info("Starting Chronos...")

	go c.informer.Run(stopper)

	logrus.Info("Synchronizing events...")
	// wait for the caches to synchronize before starting the worker

	if !cache.WaitForCacheSync(stopper, c.informer.HasSynced) {
		utilruntime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		logrus.Info("synchronization failed...")
		return
	}

	logrus.Info("synchronization complete!")

	// runWorker will loop until "something bad" happens.  The .Until will
	// then rekick the worker after one second
	wait.Until(c.runWorker, time.Second, stopper)
}

func (c *controller) HasSynced() bool {
	return c.informer.HasSynced()
}

func (c *controller) runWorker() {
	for c.processNextItem() {
		// continue looping
	}
}

func (c *controller) processNextItem() bool {
	e, term := c.queue.Get()

	if term {
		return false
	}

	err := c.processItem(e.(event))
	if err == nil {
		// No error, reset the ratelimit counters
		c.queue.Forget(e)
		return true
	}
	return true
}

func (c *controller) processItem(e event) error {
	obj, _, err := c.informer.GetIndexer().GetByKey(e.key)
	if err != nil {
		return fmt.Errorf("Error fetching object with key %s from store: %v", e.key, err)
	}

	//Use a switch clause instead
	logrus.Infof("Chronos has processed 1 event of type [%s] for object [%s]", e.eventType, obj)

	return nil
}

//Alternate method
/* func Start(config string) {
	kc, err := utils.GetClient(config)
	if err != nil {
		logrus.Fatal(err)
	}

	factory := informers.NewSharedInformerFactory(kc, 0)
	informer := factory.Core().V1().Events().Informer()
	stopper := make(chan struct{})
	defer close(stopper)
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*v1.Event)
			logrus.Infof("AddEvent : %s %s %s %s %s %s", mObj.Type, mObj.InvolvedObject.Kind, mObj.Namespace, mObj.Reason, mObj.Message, mObj.Series)
		},
		DeleteFunc: func(obj interface{}) {
			// "k8s.io/apimachinery/pkg/apis/meta/v1" provides an Object
			// interface that allows us to get metadata easily
			mObj := obj.(*v1.Event)
			logrus.Infof("DelEvent : %s %s %s %s %s %s", mObj.Type, mObj.InvolvedObject.Kind, mObj.Namespace, mObj.Reason, mObj.Message, mObj.Series)
		},
	})

	informer.Run(stopper)
} */
