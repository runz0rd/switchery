package main

import (
	"context"
	"fmt"
	"log"

	"github.com/foomo/keel/config"
	"github.com/nats-io/nats.go"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	versionBlue  = "blue"
	versionGreen = "green"
	versionKey   = "version"
)

func main() {
	config := config.Config()
	natsAddr := config.GetString("nats.addr")
	namespace := config.GetString("namespace")
	service := config.GetString("service")

	ctx := context.Background()
	// create the in-cluster config
	cc, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal(err)
	}
	// create the clientset
	cs, err := kubernetes.NewForConfig(cc)
	if err != nil {
		log.Fatal(err)
	}
	// nats
	nc, err := nats.Connect(natsAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	// Subscribe
	if _, err := nc.Subscribe(fmt.Sprintf("%v/%v", namespace, service), func(m *nats.Msg) {
		// make sure get, list, update permissions for service are avialable on k8s
		log.Printf("got message %q with data %q", m.Subject, string(m.Data))
		if err := setVersion(ctx, cs, namespace, service, string(m.Data)); err != nil {
			log.Println(err)
		}
	}); err != nil {
		log.Fatal(err)
	}
	log.Println("waiting for nats messages")
	select {}
}

// switch service resource label beetween blue and green
func setVersion(ctx context.Context, cs *kubernetes.Clientset, namespace, service, version string) error {
	// namespace must be provided if going after the service by name
	svc, err := cs.CoreV1().Services(namespace).Get(ctx, service, v1.GetOptions{})
	if err != nil {
		return err
	}
	// switch between green and blue
	selectors := svc.Spec.Selector
	currentVersion, ok := selectors[versionKey]
	if !ok {
		return fmt.Errorf("no %v selector on service %v/%v", versionKey, svc.Namespace, svc.Name)
	}
	switch version {
	case "switch":
		selectors[versionKey] = versionBlue
		if currentVersion == versionBlue {
			selectors[versionKey] = versionGreen
		}
	case "blue", "green":
		selectors[versionKey] = version
	default:
		return fmt.Errorf("invalid version %q provided", version)
	}
	svc.Spec.Selector = selectors
	if _, err = cs.CoreV1().Services(namespace).Update(ctx, svc, v1.UpdateOptions{}); err != nil {
		return err
	}
	log.Printf("switched service %v/%v version from %v to %v", namespace, service, currentVersion, selectors[versionKey])
	return nil
}
