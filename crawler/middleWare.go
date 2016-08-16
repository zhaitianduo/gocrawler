package crawler

import (
	"errors"
	"fmt"
	"sync"
)

type IdGenerator interface {
	GetUint32() uint32
}

type ChannelManagerStatus uint8

var defaultChannleLen = uint32(100)

const (
	CHANNEL_MANAGER_STATUS_UNINITIALIZED ChannelManagerStatus = 0
	CHANNEL_MANAGER_STATUS_INITIALIZED   ChannelManagerStatus = 1
	CHANNEL_MANAGER_STATUS_CLOSED        ChannelManagerStatus = 2
)

//use a map to return status info or error
var statusNameMap = map[ChannelManagerStatus]string{
	CHANNEL_MANAGER_STATUS_UNINITIALIZED: "uninitialized",
	CHANNEL_MANAGER_STATUS_INITIALIZED:   "initialized",
	CHANNEL_MANAGER_STATUS_CLOSED:        "closed",
}

type ChannelManager interface {
	//channelLen represents the initial length of channels
	//reset means whether or not to reset the channel manager
	Init(channelLen uint32, reset bool) bool
	//close the channel manager
	Close() bool
	//get the request channel
	ReqChan() (chan *Request, error)
	//get the response channel
	ResChan() (chan *Response, error)
	//get the item channel
	ItemChan() (chan *Item, error)
	//get the error channel
	ErrorChan() (chan error, error)
	//get the length of channel
	ChannelLen() uint32
	//get the status of channel manager
	Status() ChannelManagerStatus
	//get the summary info
	Summary() string
}

type myChannelManager struct {
	channelLen uint32
	reqChan    chan *Request
	resChan    chan *Response
	itemChan   chan *Item
	errorChan  chan error
	status     ChannelManagerStatus
	rwmutex    sync.RWMutex
}

func (m *myChannelManager) Init(channelLen uint32, reset bool) bool {
	if channelLen == 0 {
		panic(errors.New("The channel length is invalid!"))
	}
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()
	if m.status == CHANNEL_MANAGER_STATUS_INITIALIZED && !reset {
		return false
	}
	m.channelLen = channelLen
	m.reqChan = make(chan *Request, channelLen)
	m.resChan = make(chan *Response, channelLen)
	m.itemChan = make(chan *Item, channelLen)
	m.errorChan = make(chan error, channelLen)
	m.status = CHANNEL_MANAGER_STATUS_INITIALIZED
	return true
}

func (m *myChannelManager) Close() bool {
	m.rwmutex.Lock()
	defer m.rwmutex.Unlock()
	if m.status != CHANNEL_MANAGER_STATUS_INITIALIZED {
		return false
	}
	close(m.reqChan)
	close(m.resChan)
	close(m.itemChan)
	close(m.errorChan)
	m.status = CHANNEL_MANAGER_STATUS_CLOSED
	return true
}

//check status is used in get channel function(ReqChan, ResChan...), mutex will be used in get channel function to make sure it's thread safe
func (m *myChannelManager) checkInitializedStatus() error {
	if m.status == CHANNEL_MANAGER_STATUS_INITIALIZED {
		return nil
	} else {
		return errors.New(statusNameMap[m.status])
	}
}

func (m *myChannelManager) ReqChan() (chan *Request, error) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if err := m.checkInitializedStatus(); err != nil {
		return nil, err
	} else {
		return m.reqChan, nil
	}
}

func (m *myChannelManager) ResChan() (chan *Response, error) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if err := m.checkInitializedStatus(); err != nil {
		return nil, err
	} else {
		return m.resChan, nil
	}
}

func (m *myChannelManager) ItemChan() (chan *Item, error) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if err := m.checkInitializedStatus(); err != nil {
		return nil, err
	} else {
		return m.itemChan, nil
	}
}

func (m *myChannelManager) ErrorChan() (chan error, error) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if err := m.checkInitializedStatus(); err != nil {
		return nil, err
	} else {
		return m.errorChan, nil
	}
}

var chanSummaryTemplate = "status: %s, " +
	"request channel: %d/%d, " +
	"response channel: %d/%d, " +
	"item channel: %d/%d, " +
	"error channel: %d/%d, "

func (m *myChannelManager) Summary() string {
	summary := fmt.Sprintf(chanSummaryTemplate, statusNameMap[m.status], m.reqChan,
		len(m.reqChan), cap(m.reqChan),
		len(m.resChan), cap(m.resChan),
		len(m.itemChan), cap(m.itemChan),
		len(m.errorChan), cap(m.errorChan),
	)
	return summary
}

func (m *myChannelManager) ChannelLen() uint32 {
	return m.channelLen
}

func (m *myChannelManager) Status() ChannelManagerStatus {
	return m.status
}

func NewChannelManager(channelLen uint32, reset bool) ChannelManager {
	if channelLen == 0 {
		channelLen = defaultChannleLen
	}
	cm := new(myChannelManager)
	if !cm.Init(channelLen, reset) {
		fmt.Println("Failed to initialize channel manager!")
	}
	return cm
}

type Entity interface {
	Id() uint32
}

type Pool interface {
	Take() (*Entity, error)
	Return(entity *Entity) error
	Total() uint32
	Used() uint32
}

type StopSign interface {
	// set the stop sign, if the stop sign has been signed, return false
	Sign() bool
	//get the info whether the stop sign has been sent
	Signed() bool
	//reset the stop sign
	Reset()
	//code represents the code of stop sign processor
	Deal(code string) uint32
	//get the stop sign processed count from one stop sign processor
	DealCount(code string) uint32
	//get the total count of stop sign processed by all stop sign processors
	DealTotal() uint32
	//get the summary info
	Summary() string
}
