package lru

import (
	"container/list"
	"os"
	"time"

	"github.com/esaiy/golang-lirs/simulator"
	"github.com/petar/GoLLRB/llrb"
)

const MAXFREQ = 1000

type LFU struct {
	CacheTuple
}

type CacheTuple struct {
	maxlen      int
	available   int
	totalaccess int
	hit         int
	miss        int
	pagefault   int
	write       int

	tlba    *llrb.LLRB
	freqArr [MAXFREQ]*list.List
}

func NewLRU(cacheSize int) *LFU {
	lfu := &LFU{
		CacheTuple: CacheTuple{
			maxlen:      cacheSize,
			available:   cacheSize,
			totalaccess: 0,
			hit:         0,
			miss:        0,
			pagefault:   0,
			write:       0,
			tlba:        llrb.New(),
			freqArr:     [MAXFREQ]*list.List{},
		},
	}
	for i := 0; i < MAXFREQ; i++ {
		lfu.freqArr[i] = list.New()
	}
	return lfu
}

func (lfu LFU) Get(trace simulator.Trace) (err error) {
	return nil
}
func (lfu LFU) PrintToFile(file *os.File, timeStart time.Time) (err error) {
	return nil
}
