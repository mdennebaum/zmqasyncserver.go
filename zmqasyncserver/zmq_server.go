package zmqasyncserver

import (
	"fmt"
	zmq "github.com/alecthomas/gozmq"
)

type ZMQServer struct{
	Port int
	Handler ZMQServerMessageHandler
	Stopped bool
	Context zmq.Context
	PollingTimeout int
}

func (this *ZMQServer) Init(){
	this.Port = 8653
	this.Stopped = true
	this.PollingTimeout = 1000*1000
}

/*
starts the listener.
@param threaded should this be started in a new thread? if true, will return immediately, if false will never return.
*/
func (this *ZMQServer) Listen(port int,handler ZMQServerMessageHandler,threaded bool){
	this.Port = port
	this.Handler = handler
	if threaded {
		go this.Run()
	}else{
		this.Run()
	}
}

func (this *ZMQServer) Run(){
	context,_ := zmq.NewContext()
	defer context.Close()
	
	frontend,_ := context.NewSocket(zmq.ROUTER)
	defer frontend.Close()
	

	frontend.Bind(fmt.Sprintf("tcp://*:%d",this.Port))

	backend,_ := context.NewSocket(zmq.DEALER)
	defer backend.Close()
	
	backend.Bind("inproc://serverbackend")

	//set up the outgoing
	outgoing := new(ZMQServerOutgoing)
	outgoing.Init(context, "inproc://serverbackend")
	outgoing.Start()

	this.Stopped = false

	items := zmq.PollItems{
        zmq.PollItem{Socket: frontend, zmq.Events: zmq.POLLIN},
        zmq.PollItem{Socket: backend, zmq.Events: zmq.POLLIN},
    }
    
	println(fmt.Sprintf("Server Listening on : %d",this.Port))

	for {
		if this.Stopped{
			println("stopping server...")
			
			//set the socket linger options
			frontend.SetSockOptInt(zmq.LINGER,0)
			backend.SetSockOptInt(zmq.LINGER,0)
			
			//close our sockets
			frontend.Close()
			backend.Close()
			
			//close our context
			context.Close()
			return
		}
		
		//poll the items
		_, _ = zmq.Poll(items, -1)

		switch {
	        case items[0].REvents&zmq.POLLIN != 0:
	            for {
		            //  Process task
		            id,_ := frontend.Recv(0)
		            message,_ := frontend.Recv(0)
		           	more, err := frontend.GetSockOptUInt64(zmq.RCVMORE)
					if err != nil {
						println(err)
						return
					}
					
					channel := new(ZMQChannel)
					channel.Init(id,outgoing)
					this.HandleIncoming(channel,message)
	            	
	            	if more == 0 {
						break
					}
	            }
	        case items[1].REvents&zmq.POLLIN != 0:
	            for {
		            //  Process task
		            id,_ := backend.Recv(0)
		            message,_ := backend.Recv(0)
		           	
		           	more, err := backend.GetSockOptUInt64(zmq.RCVMORE)
					if err != nil {
						return
					}
					
					frontend.Send(id, zmq.SNDMORE)
					frontend.Send(message,0)
	            	if more == 0 {
						break
					}
	            }
        }
	}
}


/*
closes and cleans up this server

Currently this is an asynch call, it returns immediately, but it may take a 
second or three to actually clean up.  

TODO: block until the operation completes.
*/
func (this *ZMQServer) Close(){
	this.Stopped = true
}

func (this *ZMQServer) HandleIncoming (channel *ZMQChannel,message []byte){
	this.Handler.Incoming(channel, message)
}