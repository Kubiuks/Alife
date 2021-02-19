package lib

import (
	"os"
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
	loop:
	for {
		for i:=0;i<w.numOfAgents; i++ {
			s := <-w.ch
			switch s {
			case "end":
				err := w.f.Close()
				check(err)
				break loop
			default:
				if i != w.numOfAgents-1 {
					s = s + ","
				}
				_, err := w.f.WriteString(s)
				check(err)
			}
		}
		_, err := w.f.WriteString("\n")
		check(err)
	}
	w.finished <- true
}