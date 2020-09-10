package youcrawl

import "errors"

var (
	KeyNotContainError = errors.New("key not in item")
	TypeError          = errors.New("error type of item value")
)

type Item struct {
	Store map[string]interface{}
}

func (i *Item) GetValue(key string) (interface{}, error) {
	value, isExist := i.Store[key]
	if !isExist {
		return "", KeyNotContainError
	}
	return value, nil
}

func (i *Item) SetValue(key string, value interface{}) {
	i.Store[key] = value
}

func (i *Item) GetString(key string) (string, error) {
	rawValue, err := i.GetValue(key)
	if err != nil {
		return "", err
	}
	value, ok := rawValue.(string)
	if !ok {
		return "", TypeError
	}
	return value, nil
}

func (i *Item) GetInt(key string) (int, error) {
	rawValue, err := i.GetValue(key)
	if err != nil {
		return 0, err
	}
	value, ok := rawValue.(int)
	if !ok {
		return 0, TypeError
	}
	return value, nil
}

func (i *Item) GetFloat64(key string) (float64, error) {
	rawValue, err := i.GetValue(key)
	if err != nil {
		return 0, err
	}
	value, ok := rawValue.(float64)
	if !ok {
		return 0, TypeError
	}
	return value, nil
}
