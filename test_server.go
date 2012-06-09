package main

import (
	"runtime"
	"zmqasyncserver"
)

type TestHandler struct{}

func (this *TestHandler) Incoming(channel *zmqasyncserver.ZMQChannel, message []byte) {
	channel.Send(message)
}

func (this *TestHandler) Error(err error) {
	
}

func main(){
	runtime.GOMAXPROCS(runtime.NumCPU())
	apiserver := new(zmqasyncserver.ZMQServer)
	handler := new(TestHandler)
	exit := make(chan bool)
	go func(){
		apiserver.Listen(8988, handler, true)		
	}()
	<-exit
}