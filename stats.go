package freepool

import "unsafe"

type StatsInfo struct {
	fp *freePool
}

func (si StatsInfo) GetRoughlyAllocBytes() int {
	return si.fp.roughlyAllocBytes
}

func (si StatsInfo) Pages() int {
	return len(si.fp.oversizedPages) + len(si.fp.chunks)*pagesPerChunk
}

func (si StatsInfo) PagesInuse() int {
	return si.Pages() - si.PagesIdle()
}

func (si StatsInfo) PagesIdle() int {
	idle := 0
	for _, n := range si.fp.idlePagesPerChunk {
		idle += int(n)
	}
	return idle
}

func (si StatsInfo) PagesOversized() int {
	return len(si.fp.oversizedPages)
}

func (si StatsInfo) PagesFixedsized() int {
	return len(si.fp.chunks) * pagesPerChunk
}

//go:nocheckptr
func (si StatsInfo) PagesLink(p Page) (id []int, usedBytes []int) {
	for page := (*page)(unsafe.Pointer(p.headAddr)); page != nil; page = page.nextPage() {
		id = append(id, page.id())
		usedBytes = append(usedBytes, int(page.pageOffset)-pageExtraOffset)
	}
	return
}

func (si StatsInfo) WalkPages(f func(pageID int, nextPageID int, usedBytes int) bool) {
	fp := si.fp
	for _, page := range fp.oversizedPages {
		f(page.id(), page.nextPage().id(), int(page.pageOffset)-pageExtraOffset)
	}
	for i, chunk := range fp.chunks {
		pagesToWalk := pagesPerChunk - int(fp.idlePagesPerChunk[i])
		for j := 0; j < pagesPerChunk && pagesToWalk > 0; j++ {
			flagIdx := j / 64
			flagBitIdx := j % 64
			if chunk.pageFlags[flagIdx]&pageFlagMasks[flagBitIdx] == 0 {
				pagesToWalk--
				page := chunk.page(j, fp.wrappedPageSize)
				f(page.id(), page.nextPage().id(), int(page.pageOffset)-pageExtraOffset)
			}
		}
	}
}
