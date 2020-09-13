# Pipeline
The pipeline, as a component for processing items, further processes the collected data. For example, save the collected information.
## Custom Pipeline

Just implement the `youcrawl.Pipeline` interface. Use `e.AddPipelines(myPipeline)` to add
```go
type Pipeline interface {
	Process(item *Item, store GlobalStore) error
}
```

## Image download Pipeline
The image download pipeline will download all the urls under the `downloadImgURLs` field in the Item.
```go
type ImageDownloadPipeline struct {
    // save folder
    GetStoreFileFolder func(item *Item, store GlobalStore) string
    // save filename
    GetSaveFileName    func(item *Item, store GlobalStore, rawURL string) string
    // maximum parallel download
    MaxDownload        int
    // middlewares for download image request
	Middlewares        []Middleware
}
```