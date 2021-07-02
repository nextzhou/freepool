package freepool

import (
	"fmt"
	"math/bits"
	"reflect"
	"unsafe"
)

func New(pageSize int) FreePool {
	pool := new(freePool) // 一定要确保分配在堆上，因为其他地方会持有它的 uintptr 指针
	pool.pageSize = pageSize
	pool.oversizedPages = make(map[int]*page)
	pool.wrappedPageSize = pageSize + pageExtraOffset
	return pool
}

type Conf struct {
	PageSize  int
	Callbacks *struct {
		OnOversize      func(size int)
		OnNewChunk      func()
		OnPageLink      func(head int, length int)
		OnCallbackPanic func(recover interface{}, stack []byte)
	}
}

func (fp *freePool) Page() Page {
	return Page{
		pool: uintptr(unsafe.Pointer(fp)),
	}
}

type FreePool interface {
	//Allocatable
	Page() Page
	Stats() StatsInfo
}

func (fp *freePool) Stats() StatsInfo {
	return StatsInfo{fp: fp}
}

type Allocatable interface {
	String(string) String
	Strings([]string) Strings
	Bytes([]byte) Bytes
	Ints([]int) Ints
	Int8s([]int8) Int8s
	Int16s([]int16) Int16s
	Int32s([]int32) Int32s
	Int64s([]int64) Int64s
	Uints([]uint) Uints
	Uint8s([]uint8) Uint8s
	Uint16s([]uint16) Uint16s
	Uint32s([]uint32) Uint32s
	Uint64s([]uint64) Uint64s
	Float32s([]float32) Float32s
	Float64s([]float64) Float64s
}

type Page struct {
	pool     uintptr
	headAddr uintptr
}

var _ Allocatable = &Page{}

func (p *Page) String(s string) String {
	b := p.Bytes([]byte(s))
	return String{Data: b.Data, Len: b.Len}
}

func (p *Page) Strings(ss []string) Strings {
	slice := make([]String, len(ss))
	for i, s := range ss {
		slice[i] = p.String(s)
	}
	// 把 slice 也得放到 chunk 里
	ret := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&slice)), int(unsafe.Sizeof(String{})))
	return Strings(ret)
}

//go:nocheckptr
func (p *Page) Bytes(bytes []byte) Bytes {
	size := len(bytes)
	if size == 0 {
		return Bytes{Data: 0, Len: 0, Cap: 0}
	}
	pool := (*freePool)(unsafe.Pointer(p.pool))
	var selectedPage *page
	ptrToTail := &p.headAddr
	var dataAddr uintptr
	if size > pool.pageSize { // 超尺寸
		selectedPage, dataAddr = newOversizedPage(bytes)
		// onPageSize
		pool.overflowIdx--
		selectedPage.idx = pool.overflowIdx
		pool.oversizedPages[selectedPage.idx] = selectedPage
		pool.roughlyAllocBytes += int(selectedPage.pageOffset)
		for addr := p.headAddr; addr != 0; {
			page := (*page)(unsafe.Pointer(addr))
			ptrToTail = &page.nextPageAddr
			addr = page.nextPageAddr
		}
	} else { // 寻找一下放得下的 page
		maxOffset := pool.wrappedPageSize - alignedSize(size)
		for addr := p.headAddr; addr != 0; {
			page := (*page)(unsafe.Pointer(addr))
			if int(page.pageOffset) <= maxOffset {
				selectedPage = page
				break
			}
			addr = page.nextPageAddr
			ptrToTail = &page.nextPageAddr
		}
		if selectedPage == nil { // 没有放得下的 page, 新申请一个
			//TODO onPageLink
			selectedPage = pool.initNextIdlePage()
		}
		dataAddr = selectedPage.write(bytes)
	}
	*ptrToTail = uintptr(unsafe.Pointer(selectedPage))
	return Bytes{Data: dataAddr, Len: size, Cap: size}
}

func (p *Page) unsafeSlice(s *reflect.SliceHeader, elemSize int) reflect.SliceHeader {
	bytes := reflect.SliceHeader{
		Data: s.Data,
		Len:  s.Len * elemSize,
		Cap:  s.Len * elemSize,
	}
	ret := reflect.SliceHeader(p.Bytes(*(*[]byte)(unsafe.Pointer(&bytes))))
	ret.Len, ret.Cap = s.Len, s.Len
	return ret
}

func (p *Page) Ints(s []int) Ints {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), int(unsafe.Sizeof(int(0))))
	return Ints(h)
}

func (p *Page) Int8s(s []int8) Int8s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 1)
	return Int8s(h)
}

func (p *Page) Int16s(s []int16) Int16s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 2)
	return Int16s(h)
}

func (p *Page) Int32s(s []int32) Int32s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 4)
	return Int32s(h)
}

func (p *Page) Int64s(s []int64) Int64s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 8)
	return Int64s(h)
}

func (p *Page) Uints(s []uint) Uints {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), int(unsafe.Sizeof(uint(0))))
	return Uints(h)
}

func (p *Page) Uint8s(s []uint8) Uint8s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 1)
	return Uint8s(h)
}

func (p *Page) Uint16s(s []uint16) Uint16s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 2)
	return Uint16s(h)
}

func (p *Page) Uint32s(s []uint32) Uint32s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 4)
	return Uint32s(h)
}

func (p *Page) Uint64s(s []uint64) Uint64s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 8)
	return Uint64s(h)
}

func (p *Page) Float32s(s []float32) Float32s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 4)
	return Float32s(h)
}

func (p *Page) Float64s(s []float64) Float64s {
	h := p.unsafeSlice((*reflect.SliceHeader)(unsafe.Pointer(&s)), 8)
	return Float64s(h)
}

//go:nocheckptr
func (p *Page) Release() {
	pool := (*freePool)(unsafe.Pointer(p.pool))
	pool.releasePage(p.headAddr)
	p.headAddr = 0
}

type freePool struct {
	pageSize          int
	wrappedPageSize   int // page & extra info
	overflowIdx       int
	idlePagesPerChunk []uint16 // 记录每个chunk中可用page数量
	chunks            []*chunk
	oversizedPages    map[int]*page // 记录大于 page size 的超尺寸 page
	roughlyAllocBytes int
}

const pagesPerChunk = 64 * 64

//go:nocheckptr
func (fp *freePool) releasePage(pageAddr uintptr) {
	for pageAddr != 0 {
		page := (*page)(unsafe.Pointer(pageAddr))
		if page.idx < 0 {
			delete(fp.oversizedPages, page.idx)
		} else {
			chunkIdx, pageIdxInChunk := page.idx/pagesPerChunk, page.idx%pagesPerChunk
			chunk := fp.chunks[chunkIdx]
			fp.idlePagesPerChunk[chunkIdx] -= chunk.releasePage(pageIdxInChunk)
		}
		pageAddr = page.nextPageAddr
		page.nextPageAddr = 0
	}
}

func (fp *freePool) initNextIdlePage() *page {
	for i, chunk := range fp.chunks {
		if fp.idlePagesPerChunk[i] == 0 {
			continue
		}
		page := chunk.initNextIdlePage(fp.wrappedPageSize)
		if page != nil {
			fp.idlePagesPerChunk[page.idx/pagesPerChunk]--
			return page
		}
		panic(fmt.Errorf("BUG: expected %d usable pages, but not found", fp.idlePagesPerChunk[i]))
	}
	chunk := newChunk(len(fp.chunks), fp.wrappedPageSize)
	fp.roughlyAllocBytes += fp.wrappedPageSize * pagesPerChunk
	fp.chunks = append(fp.chunks, chunk)
	fp.idlePagesPerChunk = append(fp.idlePagesPerChunk, pagesPerChunk-1)
	return chunk.initNextIdlePage(fp.wrappedPageSize)
}

// 每一个 chunk 内包含 64 * 64 = 4096 个 page
type chunk struct {
	pageFlags [64]uint64 // 记录pages的使用情况，0 为已占用，1为未占用。每个page占用1bit
	// 后续是 pages 的实际内容。虽然没写定义但实际上是连续内存分配
}

func newChunk(chunkIdx int, wrappedPageSize int) *chunk {
	size := 64*int(unsafe.Sizeof(uint64(0))) + wrappedPageSize*pagesPerChunk
	bytes := make([]byte, size)
	chunk := (*chunk)(unsafe.Pointer(&bytes[0]))
	for i := range chunk.pageFlags {
		chunk.pageFlags[i] = ^uint64(0) // 全部 bit 都置为 1
	}
	for i := 0; i < pagesPerChunk; i++ {
		chunk.page(i, wrappedPageSize).idx = chunkIdx*pagesPerChunk + i
	}
	return chunk
}

func (c *chunk) releasePage(pageIdxInChunk int) uint16 {
	flagIdx := pageIdxInChunk / 64
	flagBitIdx := pageIdxInChunk % 64
	mask := pageFlagMasks[flagBitIdx]
	flags := c.pageFlags[flagIdx]
	preOnes := bits.OnesCount64(flags)
	flags |= mask // 标记位置1
	c.pageFlags[flagIdx] = flags
	return uint16(preOnes - bits.OnesCount64(flags))
}

func (c *chunk) initNextIdlePage(wrappedPageSize int) *page {
	//TODO optimize: remove for loop
	for i, flags := range c.pageFlags {
		idleIdx := bits.TrailingZeros64(flags)
		if idleIdx < 64 {
			c.pageFlags[i] &= ^pageFlagMasks[idleIdx] // 标记位置0
			page := c.page(i*64+idleIdx, wrappedPageSize)
			page.pageOffset = uintptr(pageExtraOffset)
			return page
		}
	}
	return nil
}

//go:nocheckptr
func (c *chunk) page(idxInChunk int, wrappedPageSize int) *page {
	addr := uintptr(unsafe.Pointer(c)) + 64*8 + uintptr(wrappedPageSize*idxInChunk)
	return (*page)(unsafe.Pointer(addr))
}

type page struct {
	idx          int // freepool 中第几个 page, 方便 release 时快速定位
	nextPageAddr uintptr
	pageOffset   uintptr // 已使用的位置
}

func (p *page) id() int {
	if p == nil {
		return 0
	}
	if p.idx < 0 {
		return p.idx
	}
	return p.idx + 1 // 固定尺寸 page 的 id 为 idx +1, 把 0 留给未分配空间的 page
}

//go:nocheckptr
func (p *page) nextPage() *page {
	return (*page)(unsafe.Pointer(p.nextPageAddr))
}

func (p *page) write(bytes []byte) uintptr {
	size := len(bytes)
	dstH := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(p)) + p.pageOffset,
		Len:  size,
		Cap:  size,
	}
	copy(*(*[]byte)(unsafe.Pointer(&dstH)), bytes)
	p.pageOffset += uintptr(alignedSize(size))
	return dstH.Data
}

func newOversizedPage(bytes []byte) (*page, uintptr) {
	pageMem := make([]byte, len(bytes)+pageExtraOffset)
	copy(pageMem[pageExtraOffset:], bytes)
	page := (*page)(unsafe.Pointer(&pageMem[0]))
	page.pageOffset = uintptr(alignedSize(len(pageMem)))
	return page, uintptr(unsafe.Pointer(&pageMem[pageExtraOffset]))
}

const wordSize = int(unsafe.Sizeof(int(0)))
const pageExtraOffset = wordSize * 3

var pageFlagMasks = func() [64]uint64 {
	var masks [64]uint64
	for i := 0; i < 64; i++ {
		masks[i] = 1 << i
	}
	return masks
}()

func alignedSize(size int) int {
	return ((size + wordSize - 1) / wordSize) * wordSize
}
