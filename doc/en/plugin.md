# Plugin
Plugin is goruntime, use for logging, WebAPI, etc. when the Engine is running.
Plugin is optional component.
## Useage
```go
e.AddPlugins(yourPlugin)
```
## StatusOutputPlugin
The component provides some basic statistical information, including unrequest count, completed count, speed, etc. And save to GlobalStore
