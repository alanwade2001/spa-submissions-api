package main

import (
	"os"

	"github.com/alanwade2001/spa-submissions-api/types"
	"k8s.io/klog/v2"
)

func main() {
	klog.Infof("starting %s", os.Args[0])

	var serverAPI types.ServerAPI

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
