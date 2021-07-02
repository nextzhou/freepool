// +build freepool_cgo_alloc

package freepool

// #include <stdlib.h>
import "C"

func alloc() []byte {
	//TODO
	return nil
}