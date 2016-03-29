package main

import (
	fmt "fmt"
	_rand "math/rand"
	_atomic "sync/atomic"
	time "time"
)

type _readOp struct {
	_key  int
	_resp chan int
}
type _writeOp struct {
	_key  int
	_val  int
	_resp chan bool
}

func main() {
	var _ops int64 = 0
	_reads := make(chan *_readOp)
	_writes := make(chan *_writeOp)
	go func() {
		var _state = make(map[int]int)
		for {
			select {
			case _read := <-_reads:
				_read._resp <- _state[_read._key]
			case _write := <-_writes:
				_state[_write._key] = _write._val
				_write._resp <- true
			}
		}
	}()
	for _r := 0; _r < 100; _r++ {
		go func() {
			for {
				_read := &_readOp{_key: _rand.Intn(5), _resp: make(chan int)}
				_reads <- _read
				<-_read._resp
				_atomic.AddInt64(&_ops, 1)
			}
		}()
	}
	for _w := 0; _w < 10; _w++ {
		go func() {
			for {
				_write := &_writeOp{_key: _rand.Intn(5), _val: _rand.Intn(100), _resp: make(chan bool)}
				_writes <- _write
				<-_write._resp
				_atomic.AddInt64(&_ops, 1)
			}
		}()
	}
	time.Sleep(time.Second)
	_opsFinal := _atomic.LoadInt64(&_ops)
	fmt.Println("ops:", _opsFinal)

	_range := "abc" //when used with 让, Go keywords like "range" can be used as identifiers
	_range = "abcdefg"
	fmt.Printf("range: %v\n", _range)
}
