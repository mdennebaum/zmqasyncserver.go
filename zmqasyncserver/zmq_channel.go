package zmqasyncserver

type ZMQChannel struct{
	Id []byte
	outgoing *ZMQServerOutgoing
}

func (this *ZMQChannel) Init(id []byte, out *ZMQServerOutgoing) {
	this.Id = id
	this.outgoing = out
}

func (this *ZMQChannel) Send(message []byte) {
	this.outgoing.Send(this.Id, message)
}
