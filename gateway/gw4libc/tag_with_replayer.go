// +build koala_replayer

package gw4libc

// #cgo CFLAGS: -DKOALA_REPLAYER -DPTHREAD -DPTHREAD_SINGLETHREADED_TIME
// #cgo CXXFLAGS: -DKOALA_REPLAYER -DPTHREAD -DPTHREAD_SINGLETHREADED_TIME --std=c++11 -Wno-ignored-attributes
import "C"