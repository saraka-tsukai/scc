package main

import (
	"encoding/json"
	"fmt"

	"github.com/saraka/scc"
)

func main() {
	var (
		s   scc.Store
		err error
		t   interface{}
	)
	// if s, err = scc.NewZookeeperStore("example1", []string{"127.0.0.1:2181"}); err != nil {
	if s, err = scc.NewConsulStore("example1", []string{"127.0.0.1:8500"}); err != nil {
		fmt.Println(err)
		return
	}

	json.Unmarshal([]byte(`{"a":1}`), &t)
	s.Set("test.key", t)
	s.Watch("test.key", func(data []byte) {
		var tmp interface{}
		json.Unmarshal(data, &tmp)
		fmt.Println("data:", tmp)
	})
	s.Get("test.key", &t)
	fmt.Println("test.key", t)
	var input string
	for {
		fmt.Scanln(&input)
		json.Unmarshal([]byte(input), &t)
		s.Set("test.key", t)
	}
}
