package zmqasyncserver

type ZMQOutMessage struct{
	Id []byte
	Message []byte
}

func (this *ZMQOutMessage) Init(id []byte, message []byte){
	this.Id = id
	this.Message = message
}