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

## OutputCSVPostProcess
将GlobalStore中字段为`items`输出至csv文件

```go
type OutputCSVPostProcessOption struct {
    // output path.
    // if not provided,use `./output.csv` as default value
    OutputPath string
    // with header.
    // default : false
    WithHeader bool
    // key to write
    // if not provided,will write all key
    Keys []string
    // key to csv column name.
    // if not provide,use key name as csv column name
    KeysMapping map[string]string
    // if value not exist in item.
    // by default,use empty string
    NotExistValue string
}
```