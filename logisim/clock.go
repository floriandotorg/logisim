package logisim

type Clock interface {
  Tick()
  Ticks(number uint64)
}

type clock struct {
  line TriggerLine
}

func NewClock(line TriggerLine) Clock {
  return &clock{
    line: line,
  }
}

func (c *clock) Tick() {
  c.Ticks(1)
}

func (c *clock) Ticks(number uint64) {
  for ; number > 0; number-- {
    c.line.Write(true)
    c.line.Write(false)
  }
}
