package main

import (
	"fmt"

	"github.com/saraka-tsukai/scc"
)

func main() {
	var (
		s   scc.Store
		err error
	)
	if s, err = scc.NewConsulStore("example1", []string{"http://127.0.0.1:8500"}); err != nil {
		fmt.Println(err)
		return
	}

	s.Set("test.key", 1)
	s.Watch("test.key", func(data []byte, decode scc.Decoder) bool {
		var tmp int
		decode.Decode(data, &tmp)
		fmt.Println("data:", tmp)
		return false
	})
	t := 0
	s.Get("test.key", &t)
	fmt.Println("get:", t)
	fmt.Scanln()
}
