package model

import (
	"Alife/lib"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/mobile/event/size"
	"image"
	"log"
	"time"
)

type UI struct {
	width, height int
	ch            <-chan [][]interface{}

	grid *lib.GridWidget

	window screen.Window
	buffer screen.Buffer
}

var _ lib.UI = &UI{}
var _ lib.Grid = &UI{}

const (
	winWidth  = 800
	winHeight = 800
)

func NewUI(w, h int) *UI {
	return &UI{
		width:  w,
		height: h,
	}
}

func (ui *UI) Stop() {
	ui.window.Release()
}

func (ui *UI) AddGrid(ch <-chan [][]interface{}) {
	ui.grid = lib.NewGridWidget(ui.width+1, ui.height+1, winWidth, winHeight)
	ui.ch = ch
}

func (ui *UI) Loop() {
	driver.Main(func(s screen.Screen) {
		w, err := s.NewWindow(&screen.NewWindowOptions{
			Title:  "Alife Social Stress",
			Width:  winWidth,
			Height: winHeight,
		})
		if err != nil {
			log.Fatal(err)
		}
		ui.window = w
		defer w.Release()

		defer func() {
			if ui.buffer != nil {
				ui.buffer.Release()
			}
		}()

		if ui.ch != nil {
			go ui.readGridData()
		}
		for {
			switch e := w.NextEvent().(type) {
			case key.Event:
				if e.Code == key.CodeEscape || e.Code == key.CodeQ {
					return
				}
			case paint.Event:
				if ui.grid != nil {
					ui.grid.Draw(ui.buffer.RGBA())
				}
				ui.window.Upload(image.Point{}, ui.buffer, ui.buffer.Bounds())
				ui.window.Publish()
			case size.Event:
				if ui.buffer != nil {
					ui.buffer.Release()
				}
				b, err := s.NewBuffer(e.Size())
				if err != nil {
					log.Fatal(err)
				}
				ui.buffer = b
				if ui.grid != nil {
					ui.grid.Draw(ui.buffer.RGBA())
				}
			case error:
				log.Print(e)
			}
		}

	})
}

func (ui *UI) readGridData() {
	time.Sleep(1 * time.Second)
	for dump := range ui.ch {
		ui.grid.SetGrid(dump)
		if ui.window != nil {
			ui.window.Send(paint.Event{})
		}
	}
}
