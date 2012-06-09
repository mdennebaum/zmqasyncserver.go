package zmqasyncserver

import(
	zmq "github.com/alecthomas/gozmq"
)

type ZMQServerOutgoing struct{
	Messages chan *ZMQOutMessage
	Context zmq.Context
	Connection string
}

func (this *ZMQServerOutgoing) Init(context zmq.Context, connection string){
	this.Messages = make(chan *ZMQOutMessage)
	this.Context = context
	this.Connection = connection
}

func (this *ZMQServerOutgoing) Start(){
	go this.Run()
}

func (this *ZMQServerOutgoing) Run(){
	socket,_ := this.Context.NewSocket(zmq.DEALER)
	defer socket.Close()
	
	socket.Connect(this.Connection)
	
	for {
		message := <- this.Messages
		socket.Send(message.Id, zmq.SNDMORE)
		socket.Send(message.Message, 0)
	}
}

func (this *ZMQServerOutgoing) Send(id []byte,message []byte){
	outMessage := new(ZMQOutMessage)
	outMessage.Init(id, message)
	this.Messages <- outMessage
}