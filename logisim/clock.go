package logisim

type Clock interface {
	Tick()
	TickNTimes(uint64)
	GetClockLine() TriggerLine
}

type clock struct {
	clockBus  Bus
	clockLine TriggerLine
}

func NewClock() Clock {
	clockBus := NewBus(1)

	return &clock{
		clockBus:  clockBus,
		clockLine: clockBus.TriggerBranch(0),
	}
}

func (c *clock) Tick() {
	c.clockBus.Write(0x01)
	c.clockBus.Write(0x00)
}

func (c *clock) TickNTimes(n uint64) {
	for ; n > 0; n-- {
		c.Tick()
	}
}

func (c *clock) GetClockLine() TriggerLine {
	return c.clockLine
}
