package main

import (
  "fmt"
  "./bitarray"
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
  if ctrl.Size() != 2 {
    panic("FU")
  }
  ram := &Ram{
    addr: addr,
    data: data,
    ctrl: ctrl,
    clk: clk,
    contents: make([]uint64, 1 << uint64(addr.Size())),
  }
  clk.OnRisingEdge(ram.onTick)
  return ram
}

func (r *Ram) onTick() {
  status := r.ctrl.Read().Get()
  addr := r.addr.Read().Get()
  if status == 0x01 {
    // DAS muss schöner gehen!
    //r.data.Write(bitarray.NewBitArrayWithVal(r.data.Size(), r.contents[addr]))
    fick_dich_als_workaround_for_now(r.data, r.contents[addr])
  } else if status == 0x02 {
    r.contents[addr] = r.data.Read().Get()
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

func fick_dich_als_workaround_for_now(bus logisim.Bus, val uint64) {
  bus.Write(bitarray.NewBitArrayWithVal(bus.Size(), val))
}

func main() {
  clk := logisim.NewTriggerLine()
  addr := logisim.NewBus(4)
  data := logisim.NewBus(8)
  ctrl := logisim.NewBus(2)
  ram := NewRam(addr, data, ctrl, clk)

  ram.Print()

  //addr.Write(42)
  fick_dich_als_workaround_for_now(addr, 0)
  //data.Write(74)
  fick_dich_als_workaround_for_now(data, 74)

  tick(clk)

  // todo
  // git init && git add -A && git ci -m "inital"
  // bitarray weg
  // clock object
  // Bus Print as Stringer
  // Ram Print as Stringer

  ram.Print()

  //ctrl.Write(0x01)
  fick_dich_als_workaround_for_now(ctrl, 0x02)
  tick(clk)

  ram.Print()
}
