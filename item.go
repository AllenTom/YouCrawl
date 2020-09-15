package youcrawl

import "errors"

var (
	KeyNotContainError = errors.New("key not in item")
	TypeError          = errors.New("error type of item value")
)

type DefaultItem struct {
	Store map[string]interface{}
}

func (i *DefaultItem) GetValue(key string) (interface{}, error) {
	value, isExist := i.Store[key]
	if !isExist {
		return "", KeyNotContainError
	}
	return value, nil
}

func (i *DefaultItem) SetValue(key string, value interface{}) {
	i.Store[key] = value
}

func (i *DefaultItem) GetString(key string) (string, error) {
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

func (i *DefaultItem) GetInt(key string) (int, error) {
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

func (i *DefaultItem) GetFloat64(key string) (float64, error) {
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
