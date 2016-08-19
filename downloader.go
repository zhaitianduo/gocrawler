package crawler

import (
	"errors"
	"gocrawler/base"
	"gocrawler/middleware"
	"net/http"
	"reflect"
)

var idGenerator middleware.IdGenerator = middleware.NewIdGenerator()

type PageDownloader interface {
	Id() uint32
	Download(req base.Request) (*base.Response, error)
}

type PageDownloaderPool interface {
	//get a downloader from pool
	Take() (PageDownloader, error)
	//return a downloader to pool
	Return(pd PageDownloader) error
	//get the total size of the pool
	Total() uint32
	//get the number of already used downloader
	Used() uint32
}

type myPageDownloader struct {
	id         uint32
	httpClient http.Client
}

func genDownloaderId() uint32 {
	return idGenerator.GetUint32()
}

func NewPageDownloader(client *http.Client) PageDownloader {
	if client == nil {
		client = new(http.Client)
	}
	return &myPageDownloader{
		id:         genDownloaderId(),
		httpClient: *client,
	}
}

func (m *myPageDownloader) Id() uint32 {
	return m.id
}

func (m *myPageDownloader) Download(req base.Request) (*base.Response, error) {
	// if m.httpClient == nil {
	// 	errMsg := "http client is not initialized!"
	// 	return nil, errors.New(errMsg)
	// }

	res, err := m.httpClient.Do(req.Get())
	if err != nil {
		return nil, err
	}
	httpRes := base.NewResponse(res, req.Depth())
	return httpRes, nil
}

type myPageDownloaderPool struct {
	pool  middleware.Pool
	etype reflect.Type
}

type GenPageDownloader func() PageDownloader

func NewPageDownloaderPool(
	total uint32,
	gen GenPageDownloader,
) (PageDownloaderPool, error) {
	genEntity := func() middleware.Entity { return gen() }
	etype := reflect.TypeOf(gen())
	pool, err := middleware.NewPool(total, etype, genEntity)
	if err != nil {
		return nil, err
	}
	myPool := &myPageDownloaderPool{
		pool:  pool,
		etype: etype,
	}
	return myPool, nil
}

func (m *myPageDownloaderPool) Take() (PageDownloader, error) {
	pdl, err := m.pool.Take()
	if err != nil {
		return nil, err
	}
	dl, ok := pdl.(PageDownloader)
	if !ok {
		errMsg := "Entity type doesn't match!"
		panic(errors.New(errMsg))
	}
	return dl, nil
}

func (m *myPageDownloaderPool) Return(p PageDownloader) error {
	return m.pool.Return(p)
}

func (m *myPageDownloaderPool) Total() uint32 {
	return m.pool.Total()
}

func (m *myPageDownloaderPool) Used() uint32 {
	return m.pool.Used()
}
