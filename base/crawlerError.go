package base

import (
	"bytes"
)

type ErrorType string

const (
	DOWNLOADER_ERROR     ErrorType = "Downloader Error"
	ANALYZER_ERROR       ErrorType = "Analyzer Error"
	ITEM_PROCESSOR_ERROR ErrorType = "Item Processor Error"
)

type CrawlerError interface {
	Type() ErrorType
	Error() string
}

type myCrawlerError struct {
	errType    ErrorType
	errMsg     string
	fullErrMsg string
}

func (c *myCrawlerError) Type() ErrorType {
	return c.errType
}

func (c *myCrawlerError) Error() string {
	if c.fullErrMsg == "" {
		c.genFullMsg()
	}
	return c.fullErrMsg
}

func (c *myCrawlerError) genFullMsg() {
	var buf bytes.Buffer
	buf.WriteString("Crawler Error: ")
	if c.errType != "" {
		buf.WriteString(string(c.errType))
		buf.WriteString(":")
	}
	buf.WriteString(c.errMsg)
	c.fullErrMsg = buf.String()
}

func NewCrawlerError(errType ErrorType, errMsg string) CrawlerError {
	return &myCrawlerError{
		errType: errType,
		errMsg:  errMsg,
	}
}
