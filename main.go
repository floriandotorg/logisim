package main

import (
  "fmt"
  "./logisim"
)

type Ram struct {
  addr logisim.Bus
  data logisim.Bus
  ctrl logisim.Bus
  clk logisim.TriggerLine

  contents []uint64
}

func NewRam(addr logisim.Bus, data logisim.Bus, ctrl logisim.Bus, clk logisim.TriggerLine) *Ram {
  if ctrl.Width() != 2 {
    panic("FU")
  }
  ram := &Ram{
    addr: addr,
    data: data,
    ctrl: ctrl,
    clk: clk,
    contents: make([]uint64, 1 << uint64(addr.Width())),
  }
  clk.OnRisingEdge(ram.onTick)
  return ram
}

func (r *Ram) onTick() {
  status := r.ctrl.Read()
  addr := r.addr.Read()
  if status == 0x01 {
    r.data.Write(r.contents[addr])
  } else if status == 0x02 {
    r.contents[addr] = r.data.Read()
  fmt.Println(r.contents)
  }
}

func tick(clk logisim.TriggerLine) {
  clk.Write(true)
  clk.Write(false)
}

func (r *Ram) Print() {
  for i := 0; i < len(r.contents); i += 0 {
    fmt.Printf("%08x:", i)
    for j := 0; j < 16; j++ {
      fmt.Printf(" %02x", r.contents[i])
      i++
    }
    fmt.Println()
  }
  fmt.Println()
}
}

func main() {
  clk := logisim.NewTriggerLine()
  addr := logisim.NewBus(4)
  data := logisim.NewBus(8)
  ctrl := logisim.NewBus(2)
  ram := NewRam(addr, data, ctrl, clk)

  ram.Print()

  addr.Write(0)
  data.Write(74)

  tick(clk)

  // todo
  // clock object
  // Bus Print as Stringer
  // Ram Print as Stringer

  ram.Print()

  ctrl.Write(0x02)
  tick(clk)

  ram.Print()
}
