package main

import (
  "fmt"
  "bytes"
  "./logisim"
)

type Rom struct {
  addr logisim.ReadOnlyBus
  data logisim.Bus
  re logisim.ReadOnlyBus
  clk logisim.TriggerLine

  contents []uint64
}

func NewRom(addr logisim.ReadOnlyBus, data logisim.Bus, re logisim.ReadOnlyBus, clk logisim.TriggerLine, contents []uint64) *Rom {
  if re.Width() != 1 {
    panic("FU")
  }
  if len(contents) != 1 << uint64(addr.Width()) {
    panic("ROM and contents size differ")
  }
  rom := &Rom{
    addr: addr,
    data: data,
    re: re,
    clk: clk,
    contents: contents,
  }
  clk.OnRisingEdge(rom.onTick)
  return rom
}

func (r *Rom) onTick() {
  status := r.re.Read()
  addr := r.addr.Read()
  if status == 0x01 {
    r.data.Write(r.contents[addr])
  }
}

type Ram struct {
  addr logisim.ReadOnlyBus
  data logisim.Bus
  ctrl logisim.ReadOnlyBus
  clk logisim.TriggerLine

  contents []uint64
}

func NewRam(addr logisim.ReadOnlyBus, data logisim.Bus, ctrl logisim.ReadOnlyBus, clk logisim.TriggerLine) *Ram {
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
  }
}

func (r *Ram) String() string {
  output := bytes.Buffer{}
  for i := 0; i < len(r.contents); i += 0 {
    output.WriteString(fmt.Sprintf("%08x:", i))
    for j := 0; j < 16; j++ {
      output.WriteString(fmt.Sprintf(" %02x", r.contents[i]))
      i++
    }
    output.WriteString("\n")
  }
  output.WriteString("\n")
  return output.String()
}

func main() {
  clkLine := logisim.NewTriggerLine()
  clk := logisim.NewClock(clkLine)
  addr := logisim.NewBus(7)
  data := logisim.NewBus(8)
  ctrl := logisim.NewBus(2)
  ram := NewRam(addr, data, ctrl, clkLine)

  fmt.Print(ram)

  addr.Write(42)
  data.Write(74)
  clk.Tick()

  fmt.Print(ram)

  ctrl.Write(0x02)
  clk.Tick()

  fmt.Print(ram)
}
