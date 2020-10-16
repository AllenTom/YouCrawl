# Pipeline
pipeline 是处理Item的主要组件，对获取到的数据进行进一步加工处理。例如保存获取到的信息等。

## 自定义Pipeline

实现`youcrawl.Pipeline`接口即可。使用`e.AddPipelines(myPipeline)`添加
```go
type Pipeline interface {
	Process(item *Item, store GlobalStore) error
}
```

## 图片下载Pipeline
图片下载Pipeline会将Item中`downloadImgURLs`字段下所有的url进行下载。
```go
type ImageDownloadPipeline struct {
    // 获取保存文件夹路径
    GetStoreFileFolder func(item *Item, store GlobalStore) string
    // 获取保存文件名
    GetSaveFileName    func(item *Item, store GlobalStore, rawURL string) string
    // 最大同时下载
    MaxDownload        int
    // 下载请求时的Middleware
	Middlewares        []Middleware
}
```

## ItemChannelPipeline
这个Pipeline支持向指定的Channel发送当前Item  
在一些场景中，需要向正在执行的Engine，例如开启Daemon的Engine，输入任务，并将获取到结果返回。
如下：
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
首先，生成一个`ChannelToken`（这里可以复用ID），在pipeline的Map中添加相应的`Channel`；然后在添加任务时提供相应的`ID`。我们就可以通过Channel中取出结果。

**注意**：DefaultItem已经默认实现了`ChannelPipelineToken`的，对于自定义的Item，需要实现该接口才可以正常工作。
