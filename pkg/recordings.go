package pkg

import (
	"sync"

	"github.com/montanaflynn/stats"
)

type Recordings struct {
	sync.Mutex
	values []float64
}

func NewRecordings() *Recordings {
	return &Recordings{}
}

func (r *Recordings) Add(value float64) {
	r.Lock()
	defer r.Unlock()
	r.values = append(r.values, value)
}

func (r *Recordings) GetMedian() float64 {
	r.Lock()
	defer r.Unlock()
	median, _ := stats.Median(r.values)
	return median
}

func (r *Recordings) Reset() {
	r.Lock()
	defer r.Unlock()
	r.values = []float64{}
}
