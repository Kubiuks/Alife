package lib

// UI defines the minimal user interface type.
type UI interface {
	Stop()
	Loop()
}

type Grid interface {
	AddGrid(<-chan [][]interface{})
}