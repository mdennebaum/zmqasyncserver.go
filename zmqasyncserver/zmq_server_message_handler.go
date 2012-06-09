package zmqasyncserver

type ZMQServerMessageHandler interface{
	Incoming(channel *ZMQChannel,message []byte)
	Error(err error)
}