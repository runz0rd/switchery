package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/nats-io/nats.go"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	versionBlue  = "blue"
	versionGreen = "green"
	version      = "version"
)

func main() {
	ctx := context.Background()
	// create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	// create the clientset
	cs, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	// nats
	nc, err := nats.Connect("nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	// Subscribe
	if _, err := nc.Subscribe("foo", func(m *nats.Msg) {
		// make sure get, list, update permissions for service are avialable on k8s
		log.Printf("got message %q with data %q", m.Subject, string(m.Data))
		data := strings.Split(string(m.Data), "/")
		if err := updateService(ctx, cs, data[0], data[1]); err != nil {
			log.Println(err)
		}
	}); err != nil {
		log.Fatal(err)
	}
	log.Println("waiting for nats messages")
	select {}
}

// update service resource label beetween blue and green
func updateService(ctx context.Context, cs *kubernetes.Clientset, namespace, service string) error {
	// namespace must be provided if going after the service by name
	svc, err := cs.CoreV1().Services(namespace).Get(ctx, service, v1.GetOptions{})
	if err != nil {
		return err
	}
	// switch between green and blue
	selectors := svc.Spec.Selector
	currentVersion, ok := selectors[version]
	if !ok {
		return fmt.Errorf("no version selector on service %v/%v", svc.Namespace, svc.Name)
	}
	// switch
	selectors[version] = versionBlue
	if currentVersion == versionBlue {
		selectors[version] = versionGreen
	}
	svc.Spec.Selector = selectors
	if _, err = cs.CoreV1().Services(namespace).Update(ctx, svc, v1.UpdateOptions{}); err != nil {
		return err
	}
	log.Printf("switched service %v/%v version from %v to %v", namespace, service, currentVersion, selectors[version])
	return nil
}
