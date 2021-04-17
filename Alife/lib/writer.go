package lib

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Writer struct {
	finished 	chan bool
	ch       	chan []float64
	f           *os.File
	numOfAgents	int
}

func NewWriter(finished chan bool, ch chan []float64, name string, n int) *Writer {
	f, err := os.Create(name)
	check(err)
	return &Writer{
		finished 	: finished,
		ch		 	: ch,
		f 		 	: f,
		numOfAgents : n,
	}
}

func (w *Writer) Loop(){
	header := ""
	for i:=0;i<w.numOfAgents; i++ {
		header = header+"Agent_"+strconv.Itoa(i+1)
		if i != w.numOfAgents-1 {
			header = header + ","
		}
	}
	header = header + "\n"
	_, err1 := w.f.WriteString(header)
	check(err1)
	loop:
	for {
		data := make([][]float64, w.numOfAgents)
		for i:=0;i<w.numOfAgents; i++ {
			s := <-w.ch
			switch s {
			case nil:
				err := w.f.Close()
				check(err)
				break loop
			default:
				data[i] = s
			}
		}
		sort.Slice(data, func(i, j int) bool {
			// edge cases
			if len(data[i]) == 0 && len(data[j]) == 0 {
				return false // two empty slices - so one is not less than other i.e. false
			}
			if len(data[i]) == 0 || len(data[j]) == 0 {
				return len(data[i]) == 0 // empty slice listed "first" (change to != 0 to put them last)
			}

			// both slices len() > 0, so can test this now:
			return data[i][0] < data[j][0]
		})
		res := ""
		for i:=0;i<w.numOfAgents; i++ {
			res = res + "'[" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(data[i])), ","), "[]") + "]'"
			if i != w.numOfAgents-1 {
				res = res + ","
			}
		}
		res = res + "\n"
		_, err2 := w.f.WriteString(res)
		check(err2)
	}
	w.finished <- true
}


func check(e error) {
	if e != nil {
		panic(e)
	}
}