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

## OutputCSVPostProcess
Output the `items` field in GlobalStore to the csv file
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