# Plugin
Plugin 在Engine运行的时goruntime，可执行日志记录，WebAPI等。
Plugin 是可选的组件。
## 使用方法
```go
e.AddPlugins(yourPlugin)
```
## StatusOutputPlugin
组件提供一些基础的统计信息，包括未完成的计数、完成的计数、速度等。并将数据记录至GlobalStore


