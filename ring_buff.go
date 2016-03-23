package main

import (
	"fmt"
	//"io/ioutil"
	//"container/list"
	//"runtime"
	//"time"
)

type Object []interface{}
type Ringbuffer struct {
	write uint32
	read  uint32
	array Object
	mask  uint32
	size  uint32
	flag  uint8
}

func Ring_buffer_create(n uint32) *Ringbuffer {

	var sz uint32
	var rb *Ringbuffer

	if n < 1 || n > 30 {
		return nil
	}

	rb = new(Ringbuffer)
	if rb == nil {
		return nil
	}

	sz = (1 << n) + 1
	rb.mask = (1 << n) - 1
	rb.size = sz
	rb.flag = 1
	rb.array = make(Object, sz)
	if rb.array == nil {
		return nil
	}
	return rb
}

func Ring_buffer_size(rb *Ringbuffer) uint32 {

	size := (rb.write - rb.read) & rb.mask
	return (size)
}

func Ring_buffer_is_empty(rb *Ringbuffer) int {

	if rb.write == rb.read {
		return 1
	}
	return 0
}

func Ring_buffer_is_full(rb *Ringbuffer) int {

	if ((rb.write + 1) & rb.mask) == rb.read {
		return 1
	}
	return 0
}

func Ring_buffer_get(rb *Ringbuffer) (ptr interface{}) {

	if rb.write == rb.read {
		return nil
	}
	ptr = rb.array[rb.read]
	rb.read = (rb.read + 1) & rb.mask
	return ptr
}

func Ring_buffer_put(rb *Ringbuffer, ptr interface{}) int {

	if ((rb.write + 1) & rb.mask) == rb.read {
		return -1
	}

	rb.array[rb.write] = ptr
	rb.write = (rb.write + 1) & rb.mask
	return 0
}

func Ring_buffer_puts(start int, pos int, end int, ptr interface{}, slave_work_que Slave_work_que) int {
	var i int
	for i = start; i < end; i++ {
		if Ring_buffer_put(slave_work_que[pos].slave_ring, ptr) == 0 {
			break
		}
		if pos == end-1 {
			pos = start
		} else {
			pos++
		}
	}
	if i != end {
		return 0
	} else {
		return -1
	}
}
