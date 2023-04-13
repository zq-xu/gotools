package websocket

type Message struct {
	Id      string
	Group   string
	Message []byte
}

type GroupMessage struct {
	Group   string
	Message []byte
}

type BroadCastMessage struct {
	Message []byte
}
