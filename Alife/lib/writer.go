package lib

import (
	"os"
	"sort"
	"strconv"
)

type Writer struct {
	finished 	chan bool
	ch       	chan string
	f           *os.File
	numOfAgents	int
}

func NewWriter(finished chan bool, ch chan string, name string, n int) *Writer {
	f, err := os.Create("data/"+name)
	check(err)
	return &Writer{
		finished 	: finished,
		ch		 	: ch,
		f 		 	: f,
		numOfAgents : n,
	}
}

func check(e error) {
	if e != nil {
		panic(e)
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
		data := make([]string, w.numOfAgents)
		for i:=0;i<w.numOfAgents; i++ {
			s := <-w.ch
			switch s {
			case "end":
				err := w.f.Close()
				check(err)
				break loop
			default:
				data[i] = s
			}
		}
		sort.Strings(data)
		res := ""
		for i:=0;i<w.numOfAgents; i++ {
			res = res + data[i]
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