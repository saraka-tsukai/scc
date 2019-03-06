package scc

import (
	"errors"

	"github.com/docker/libkv/store"
	"github.com/smallnest/rpcx/codec"
)

var (
	ErrUnspecifiedCodec   = errors.New("unspecified codec")
	ErrUnspecifiedAppName = errors.New("unspecified app name")
	ErrStopWatchFailed    = errors.New("stop watch failed")
)

type Store interface {
	Get(key string, value interface{}) error
	Set(key string, value interface{}) error
	Watch(key string, cb func(data []byte)) error
	StopWatch(key string) error
	SetOption(options *store.Config) error
	SetCodec(codec Codec)
	GetCodec() Codec
}

type Codec codec.Codec
