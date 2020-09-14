package youcrawl

import (
	"encoding/csv"
	"fmt"
	"os"
)

type PostProcess interface {
	Process(store GlobalStore) error
}

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
type OutputCSVPostProcess struct {
	option OutputCSVPostProcessOption
}

func NewOutputCSVPostProcess(option OutputCSVPostProcessOption) *OutputCSVPostProcess {
	return &OutputCSVPostProcess{
		option: option,
	}
}

func (o *OutputCSVPostProcess) Process(store GlobalStore) error {
	data := store.GetValue("items")
	if data == nil {
		return nil
	}
	output := data.([]map[string]interface{})
	file, err := os.Create(o.option.OutputPath)
	defer file.Close()
	if err != nil {
		return err
	}

	csvRows := make([][]string, 0)
	//scan keys
	keys := o.option.Keys
	if keys == nil {
		keys := make([]string, 0)
		seenKeys := make(map[string]bool)
		for _, item := range output {
			for key := range item {
				_, hasSeen := seenKeys[key]
				if !hasSeen {
					keys = append(keys, key)
					seenKeys[key] = true
				}
			}
		}
	}

	if o.option.WithHeader {
		if o.option.KeysMapping == nil {
			csvRows = append(csvRows, keys)
		} else {
			colHeader := make([]string, 0)
			for _, key := range keys {
				name, exist := o.option.KeysMapping[key]
				if !exist {
					colHeader = append(colHeader, key)
				} else {
					colHeader = append(colHeader, name)
				}
			}
			csvRows = append(csvRows, colHeader)
		}
	}

	for _, item := range output {
		row := make([]string, 0)
		for _, key := range keys {
			value, exist := item[key]
			if !exist {
				row = append(row, o.option.NotExistValue)
			} else {
				row = append(row, fmt.Sprintf("%v", value))
			}

		}
		csvRows = append(csvRows, row)
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(csvRows)
	if err != nil {
		return err
	}
	return nil
}
