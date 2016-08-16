package middleware

import (
	"gocrawler/base"
)

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

var chanSummaryTemplate = "status: %s, " +
	"request channel: %d/%d, " +
	"response channel: %d/%d, " +
	"item channel: %d/%d, " +
	"error channel: %d/%d, "

type ChannelManager interface {
	//channelLen represents the initial length of channels
	//reset means whether or not to reset the channel manager
	Init(channelLen uint32, reset bool) bool
	//close the channel manager
	Close() bool
	//get the request channel
	ReqChan() (chan *base.Request, error)
	//get the response channel
	ResChan() (chan *base.Response, error)
	//get the item channel
	ItemChan() (chan *base.Item, error)
	//get the error channel
	ErrorChan() (chan error, error)
	//get the length of channel
	ChannelLen() uint32
	//get the status of channel manager
	Status() ChannelManagerStatus
	//get the summary info
	Summary() string
}
