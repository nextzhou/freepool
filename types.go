package freepool

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
	"unsafe"
)

type String reflect.StringHeader

func (s String) String() string {
	return *(*string)(unsafe.Pointer(&s))
}

func (s String) RefString() string {
	return *(*string)(unsafe.Pointer(&s))
}

func (s String) CloneString() string {
	data := make([]byte, s.Len)
	copy(data, []byte(s.RefString()))
	return *(*string)(unsafe.Pointer(&data))
}

func (s String) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.RefString()) // 字符串的 json 转义比较麻烦，为了保证兼容性还是直接 marshal 比较好
}

type Strings reflect.SliceHeader

func (ss Strings) RefStrings() []string {
	return *(*[]string)(unsafe.Pointer(&ss))
}

func (ss Strings) CloneStrings() []string {
	s := *(*[]String)(unsafe.Pointer(&ss))
	ret := make([]string, len(s))
	for i := range s {
		ret[i] = s[i].CloneString()
	}
	return ret
}

func (ss Strings) MarshalJSON() ([]byte, error) {
	return json.Marshal(ss.RefStrings())
}

type Bytes reflect.SliceHeader

func (b Bytes) Bytes() []byte {
	return *(*[]byte)(unsafe.Pointer(&b))
}

func (b Bytes) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Bytes())
}

type Ints reflect.SliceHeader

func (s Ints) Ints() []int {
	return *(*[]int)(unsafe.Pointer(&s))
}

func (s Ints) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Ints() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.Itoa(n))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Int8s reflect.SliceHeader

func (s Int8s) Int8s() []int8 {
	return *(*[]int8)(unsafe.Pointer(&s))
}

func (s Int8s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Int8s() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.Itoa(int(n)))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Int16s reflect.SliceHeader

func (s Int16s) Int16s() []int16 {
	return *(*[]int16)(unsafe.Pointer(&s))
}

func (s Int16s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Int16s() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.Itoa(int(n)))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Int32s reflect.SliceHeader

func (s Int32s) Int32s() []int32 {
	return *(*[]int32)(unsafe.Pointer(&s))
}

func (s Int32s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Int32s() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.Itoa(int(n)))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Int64s reflect.SliceHeader

func (s Int64s) Int64s() []int64 {
	return *(*[]int64)(unsafe.Pointer(&s))
}

func (s Int64s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Int64s() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.FormatInt(n, 10))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Uints reflect.SliceHeader

func (s Uints) Uints() []uint {
	return *(*[]uint)(unsafe.Pointer(&s))
}

func (s Uints) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Uints() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.FormatUint(uint64(n), 10))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Uint8s reflect.SliceHeader

func (s Uint8s) Uint8s() []uint8 {
	return *(*[]uint8)(unsafe.Pointer(&s))
}

func (s Uint8s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Uint8s() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.FormatUint(uint64(n), 10))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Uint16s reflect.SliceHeader

func (s Uint16s) Uint16s() []uint16 {
	return *(*[]uint16)(unsafe.Pointer(&s))
}

func (s Uint16s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Uint16s() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.FormatUint(uint64(n), 10))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Uint32s reflect.SliceHeader

func (s Uint32s) Uint32s() []uint32 {
	return *(*[]uint32)(unsafe.Pointer(&s))
}

func (s Uint32s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Uint32s() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.FormatUint(uint64(n), 10))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Uint64s reflect.SliceHeader

func (s Uint64s) Uint64s() []uint64 {
	return *(*[]uint64)(unsafe.Pointer(&s))
}

func (s Uint64s) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBufferString("[")
	for i, n := range s.Uint64s() {
		if i > 0 {
			buf.WriteByte(',')
		}
		buf.WriteString(strconv.FormatUint(uint64(n), 10))
	}
	buf.WriteByte(']')
	return buf.Bytes(), nil
}

type Float32s reflect.SliceHeader

func (s Float32s) Float32s() []float32 {
	return *(*[]float32)(unsafe.Pointer(&s))
}

func (s Float32s) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Float32s())
}

type Float64s reflect.SliceHeader

func (s Float64s) Float64s() []float64 {
	return *(*[]float64)(unsafe.Pointer(&s))
}

func (s Float64s) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.Float64s())
}
