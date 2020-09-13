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