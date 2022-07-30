package search

type DataSearchServer interface {
    Config(args ...interface{}) error
    Start(args ...interface{}) error
    AddHandle(handler UrlHandler)
}

type UrlHandler struct {
    Url string
    Parameter string
    SearchString func(string) ([]interface{}, error)
}