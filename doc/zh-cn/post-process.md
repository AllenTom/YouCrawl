# Post Process
当所有的请求完成后会按照添加顺序完成一系列操作，例如：生成当前爬虫的日志报告等。

## 自定义PostProcess
实现 `PostProcess` 接口。
```go
type PostProcess interface {
    Process(store GlobalStore) error
}
```

## OutPutJsonPostProcess

将GlobalStore中字段为`items`输出至json文件
```go
type OutputJsonPostProcess struct {
    // 输出路径
	StorePath string
}
```