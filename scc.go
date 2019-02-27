package scc

import (
	"errors"
)

var (
	ErrUnspecifiedCodec   = errors.New("unspecified codec")
	ErrUnspecifiedAppName = errors.New("unspecified app name")
)

type Store interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}) error
	Watch(key string, cb func(data []byte, decoder Decoder) (stop bool)) error
}

type Decoder interface {
	Decode(data []byte, i interface{}) error
}
