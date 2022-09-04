package search

type DataSearchServer interface {
    Config(args ...interface{}) error
    Start(args ...interface{}) error
    AddHandle(handler UrlHandler)
}

type UrlHandler struct {
    Url string
    Method string
    ReqBody interface{}
    ReqParam []string
    GetHandleFunc func(parameter ...interface{}) ([]interface{}, error)
    PostHandleFunc func(body []byte) ([]interface{}, error)
}