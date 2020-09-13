# Post Process
It will be executed when all the requests are completed, for example: generating the log report of the current crawler, etc.

## Custom PostProcess
Implement the `PostProcess` interface.
```go
type PostProcess interface {
    Process(store GlobalStore) error
}
```

## OutPutJsonPostProcess

Output the `items` field in GlobalStore to the json file
```go
type OutputJsonPostProcess struct {
    // output path
	StorePath string
}
```