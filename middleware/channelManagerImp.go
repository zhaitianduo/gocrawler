package middleware

import (
	"errors"
	"fmt"
	"gocrawler/base"
	"sync"
)

type myChannelManager struct {
	channelLen uint32
	reqChan    chan *base.Request
	resChan    chan *base.Response
	itemChan   chan *base.Item
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
	m.reqChan = make(chan *base.Request, channelLen)
	m.resChan = make(chan *base.Response, channelLen)
	m.itemChan = make(chan *base.Item, channelLen)
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

func (m *myChannelManager) ReqChan() (chan *base.Request, error) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if err := m.checkInitializedStatus(); err != nil {
		return nil, err
	} else {
		return m.reqChan, nil
	}
}

func (m *myChannelManager) ResChan() (chan *base.Response, error) {
	m.rwmutex.RLock()
	defer m.rwmutex.RUnlock()
	if err := m.checkInitializedStatus(); err != nil {
		return nil, err
	} else {
		return m.resChan, nil
	}
}

func (m *myChannelManager) ItemChan() (chan *base.Item, error) {
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
