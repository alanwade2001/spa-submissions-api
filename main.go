package main

import "k8s.io/klog/v2"

func main() {
	klog.Infoln("starting spa-customer-api")

	serverAPI := InitialiseServerAPI()

	if err := serverAPI.Run(); err != nil {
		panic(err)
	}

}
