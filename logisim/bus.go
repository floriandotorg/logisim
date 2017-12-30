package logisim

import (
)

type Bus interface {
  OnChange(EventFunc)
  Read() uint64
  Write(uint64)
  Width() uint8
}

type bus struct {
  val uint64
  width uint8
  onChange []EventFunc
}

func NewBus(width uint8) Bus {
  return &bus{
    val: 0,
    width: width,
  }
}

func (p *bus) OnChange(f EventFunc) {
  p.onChange = append(p.onChange, f)
}

func (p *bus) Read() uint64 {
  return p.val
}

func (p *bus) Write(val uint64) {
  if p.val != val {
    p.val = val

    for _, f := range p.onChange {
      f()
    }
  }
}

func (b *bus) Width() uint8 {
  return b.width
}
