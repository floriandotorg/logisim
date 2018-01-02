package logisim

type EventFunc func()
type ChangeFunc func(uint64)
type TickFunc func() ChangeFunc
