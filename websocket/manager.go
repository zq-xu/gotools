package websocket

import (
	"sync"

	"github.com/zq-xu/gotools/logx"
)

var GlobalManager = manager{
	Group:            make(map[string]map[string]*client),
	Register:         make(chan *client, 128),
	Unregister:       make(chan *client, 128),
	GroupMessage:     make(chan *GroupMessage, 128),
	Message:          make(chan *Message, 128),
	BroadCastMessage: make(chan *BroadCastMessage, 128),
	groupCount:       0,
	clientCount:      0,
}

type Manager interface {
	RegisterClient(client *client)
	UnregisterClient(client *client)
}

type manager struct {
	Group            map[string]map[string]*client
	groupCount       uint
	clientCount      uint
	Lock             sync.Mutex
	Register         chan *client
	Unregister       chan *client
	Message          chan *Message
	GroupMessage     chan *GroupMessage
	BroadCastMessage chan *BroadCastMessage
}

// Send send data to the certain client
func (mgr *manager) Send(id string, group string, message []byte) {
	data := &Message{
		Id:      id,
		Group:   group,
		Message: message,
	}
	mgr.Message <- data
}

// SendGroup send data to all the clients in the group
func (mgr *manager) SendGroup(group string, message []byte) {
	data := &GroupMessage{
		Group:   group,
		Message: message,
	}
	mgr.GroupMessage <- data
}

func (mgr *manager) SendAll(message []byte) {
	data := &BroadCastMessage{
		Message: message,
	}
	mgr.BroadCastMessage <- data
}

func (mgr *manager) RegisterClient(client *client) {
	mgr.Register <- client
}

func (mgr *manager) UnregisterClient(client *client) {
	mgr.Unregister <- client
}

func (mgr *manager) LenGroup() uint {
	return mgr.groupCount
}

func (mgr *manager) LenClient() uint {
	return mgr.clientCount
}

func (mgr *manager) Info() map[string]interface{} {
	mgrInfo := make(map[string]interface{})
	mgrInfo["groupLen"] = mgr.LenGroup()
	mgrInfo["clientLen"] = mgr.LenClient()
	mgrInfo["chanRegisterLen"] = len(mgr.Register)
	mgrInfo["chanUnregisterLen"] = len(mgr.Unregister)
	mgrInfo["chanMessageLen"] = len(mgr.Message)
	mgrInfo["chanGroupMessageLen"] = len(mgr.GroupMessage)
	mgrInfo["chanBroadCastMessageLen"] = len(mgr.BroadCastMessage)
	return mgrInfo
}

// Start mgr
func (mgr *manager) Start() {
	logx.Logger.Errorf("websocket manage start")

	for {
		select {
		case cli := <-mgr.Register:
			mgr.consumeRegister(cli)
		case cli := <-mgr.Unregister:
			mgr.consumeUnregister(cli)
		case msg := <-mgr.Message:
			mgr.consumeMessage(msg)
		case msg := <-mgr.GroupMessage:
			mgr.consumeGroupMessage(msg)
		case msg := <-mgr.BroadCastMessage:
			mgr.consumeBroadCastMessage(msg)
		}
	}
}

func (mgr *manager) consumeRegister(cli *client) {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	logx.Logger.Infof("client [%s] connect", cli.Id)
	logx.Logger.Infof("register client [%s] to group [%s]", cli.Id, cli.Group)

	if mgr.Group[cli.Group] == nil {
		mgr.Group[cli.Group] = make(map[string]*client)
		mgr.groupCount += 1
	}

	mgr.Group[cli.Group][cli.Id] = cli
	mgr.clientCount += 1
}

func (mgr *manager) consumeUnregister(cli *client) {
	mgr.Lock.Lock()
	defer mgr.Lock.Unlock()

	logx.Logger.Infof("unregister client [%s] from group [%s]", cli.Id, cli.Group)

	groupMap, ok := mgr.Group[cli.Group]
	if !ok {
		return
	}

	_, ok = groupMap[cli.Id]
	if !ok {
		return
	}

	close(cli.Message)
	delete(mgr.Group[cli.Group], cli.Id)

	mgr.clientCount -= 1
	if len(mgr.Group[cli.Group]) == 0 {
		//log.Printf("delete empty group [%s]", client.Group)
		delete(mgr.Group, cli.Group)
		mgr.groupCount -= 1
	}
}

func (mgr *manager) consumeMessage(msg *Message) {
	groupMap, ok := mgr.Group[msg.Group]
	if !ok {
		return
	}

	conn, ok := groupMap[msg.Id]
	if !ok {
		return
	}

	go func() { conn.Message <- msg.Message }()
}

func (mgr *manager) consumeGroupMessage(msg *GroupMessage) {
	groupMap, ok := mgr.Group[msg.Group]
	if !ok {
		return
	}

	for _, conn := range groupMap {
		go func() { conn.Message <- msg.Message }()
	}
}

func (mgr *manager) consumeBroadCastMessage(msg *BroadCastMessage) {
	for _, v := range mgr.Group {
		for _, conn := range v {
			go func() { conn.Message <- msg.Message }()
		}
	}
}
