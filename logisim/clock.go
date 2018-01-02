package logisim

type Clock interface {
	OnWrite(EventFunc)
	OnRead(EventFunc)
	Tick()
	Ticks(uint64)
}

type clock struct {
	onWrite []EventFunc
	onRead  []EventFunc
}

func NewClock() Clock {
	return &clock{}
}

func (c *clock) OnWrite(f EventFunc) {
	c.onWrite = append(c.onWrite, f)
}

func (c *clock) OnRead(f EventFunc) {
	c.onRead = append(c.onRead, f)
}

func (c *clock) Tick() {
	c.Ticks(1)
}

func (c *clock) Ticks(number uint64) {
	for ; number > 0; number-- {
		for _, f := range c.onWrite {
			f()
		}

		for _, f := range c.onRead {
			f()
		}
	}
}
