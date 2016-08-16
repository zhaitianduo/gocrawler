package crawler

type PageDownloader interface {
	Id() uint32
	Download(req *Request) (*Response, error)
}

type PageDownloaderPool interface {
	//get a downloader from pool
	Take() (*PageDownloader, error)
	//return a downloader to pool
	Return(pd *PageDownloader) error
	//get the total size of the pool
	Total() uint32
	//get the number of already used downloader
	Used() uint32
}
