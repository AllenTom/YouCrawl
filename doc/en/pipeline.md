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

## ItemChannelPipeline
This Pipeline supports sending the current Item to the specified Channel. 

In some scenarios, it is necessary to input tasks to the executing Engine, such as an Engine in daemon mode, and obtain results.

Example code:
```go
fun test () {
    e.Pool.AddTasks(&Task{
                Url:       "http://www.example.com",
                Context:   Context{
                    Item: DefaultItem{Store: map[string]interface{}{
                        ItemKeyChannelToken: "thisistoken",
                    }},
                },
            })

    resultChannel := make(chan interface{})
    channelPipeline.ChannelMapping.Store("thisistoken",resultChannel)
    result := <- resultChannel
    item := result.(DefaultItem)
    fmt.Println(item.GetString("title"))
}
```
First, generate a `ChannelToken` (ID can be reused), add the corresponding `Channel` to the pipeline Map; then provide the corresponding ʻID` when adding tasks. We can retrieve the result through the Channel.

**Note**: DefaultItem has implemented `ChannelPipelineToken` by default. For custom Item, you need to implement this interface to work properly.