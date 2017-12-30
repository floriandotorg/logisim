package logisim

import (
  "../bitarray"
)

type Bus interface {
  OnChange(EventFunc)
  Read() bitarray.BitArray
  Write(bitarray.BitArray)
  Size() uint8
}

type bus struct {
  val bitarray.BitArray
  onChange []EventFunc
}

func NewBus(width uint8) Bus {
  return &bus{
    val: bitarray.NewBitArray(width),
  }
}

func (p *bus) OnChange(f EventFunc) {
  p.onChange = append(p.onChange, f)
}

func (p *bus) Read() bitarray.BitArray {
  return p.val
}

func (p *bus) Write(val bitarray.BitArray) {
  if !p.val.Equals(val) {
    p.val = val

    for _, f := range p.onChange {
      f()
    }
  }
}

func (b *bus) Size() uint8 {
  return b.val.Size()
}
