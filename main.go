package main

import (
	"os"

	"k8s.io/klog/v2"
)

func main() {
	klog.Infof("starting %s", os.Args[0])

	var serverAPI ServerAPI

	if os.Getenv("MOCK") != "" {
		klog.Infoln("Mocking server")
		serverAPI = InitialiseMockedServerAPI()
	} else {
		serverAPI = InitialiseServerAPI()
	}

	if err := serverAPI.Run(); err != nil {
		panic(err)
	}

}
