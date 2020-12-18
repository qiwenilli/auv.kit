package config

import (
	"fmt"
	"sync"
	"testing"
)

func TestA(t *testing.T) {

	s := &sync.WaitGroup{}
	for i := 0; i < 200; i++ {
		s.Add(1)
		go func(i int) {
			defer s.Done()

			v := fmt.Sprintf("%v", i)
			err := Put(fmt.Sprintf("key%d", i), v)
			t.Log("----", i, err)
		}(i)
	}
	s.Wait()

	for i := 0; i < 20; i++ {
		s.Add(1)
		go func(i int) {
			defer s.Done()

			k := fmt.Sprintf("key%v", i)
			v := fmt.Sprintf("%v", i)
			Set(fmt.Sprintf(k, i), v)

			t.Log("---->", Int64(k)+100)
		}(i)
	}
	s.Wait()
}
