package main

import (
  "fmt"
  "bytes"
  "./logisim"
)

type Rom struct {
  addr logisim.ReadOnlyBus
  addrLatch uint64
  data logisim.Bus
  outputEnable logisim.ReadOnlyBus
  clk logisim.Clock

  contents []uint64
}

func NewRom(addr logisim.ReadOnlyBus, data logisim.Bus, outputEnable logisim.ReadOnlyBus, clk logisim.Clock, contents []uint64) *Rom {
  if outputEnable.Width() != 1 {
    panic("FU")
  }
  if len(contents) != 1 << uint64(addr.Width()) {
    panic("ROM and contents size differ")
  }
  rom := &Rom{
    addr: addr,
    data: data,
    outputEnable: outputEnable,
    clk: clk,
    contents: contents,
  }
  clk.OnWrite(rom.onWrite)
  return rom
}

func (r *Rom) onWrite() {
  status := r.outputEnable.Read()
  if status == 0x01 {
    r.data.Write(r.contents[r.addrLatch])
  }
}

func (r *Rom) onRead() {
  r.addrLatch = r.addr.Read()
}

type Ram struct {
  addr logisim.ReadOnlyBus
  addrLatch uint64
  data logisim.Bus
  writeEnable logisim.ReadOnlyBus
  outputEnable logisim.ReadOnlyBus
  clk logisim.Clock

  contents []uint64
}

func NewRam(addr logisim.ReadOnlyBus, data logisim.Bus, writeEnable, outputEnable logisim.ReadOnlyBus, clk logisim.Clock) *Ram {
  ram := &Ram{
    addr: addr,
    data: data,
    writeEnable: writeEnable,
    outputEnable: outputEnable,
    clk: clk,
    contents: make([]uint64, 1 << uint64(addr.Width())),
  }
  clk.OnWrite(ram.onWrite)
  clk.OnRead(ram.onRead)
  return ram
}

func (r *Ram) onWrite() {
  if r.outputEnable.Read() == 1 {
    r.data.Write(r.contents[r.addrLatch])
  }
}

func (r *Ram) onRead() {
  r.addrLatch = r.addr.Read()
  if r.writeEnable.Read() == 1 {
    r.contents[r.addrLatch] = r.data.Read()
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

type Register struct {
  in logisim.ReadOnlyBus
  out logisim.Bus
  writeEnable logisim.ReadOnlyBus
  outputEnable logisim.ReadOnlyBus
  incrementEnable logisim.ReadOnlyBus
  reset logisim.ReadOnlyBus
  clk logisim.Clock

  max uint64
  val uint64
}

func NewRegister(in logisim.ReadOnlyBus, out logisim.Bus, writeEnable logisim.ReadOnlyBus, outputEnable, incrementEnable logisim.ReadOnlyBus, reset logisim.ReadOnlyBus, clk logisim.Clock) *Register {
  register := &Register{
    in: in,
    out: out,
    writeEnable: writeEnable,
    outputEnable: outputEnable,
    incrementEnable: incrementEnable,
    reset: reset,
    clk: clk,

    max: 1 << out.Width(),
  }
  clk.OnWrite(register.onWrite)
  clk.OnRead(register.onRead)
  return register
}

func (r *Register) onWrite() {
  if (r.max == 8) {
    fmt.Printf("%v / %v\n", r.val, r.max)
  }
  if r.reset.Read() == 1 {
    r.val = 0
  }
  if r.outputEnable.Read() == 1 {
    r.out.Write(r.val)
  }
  if r.incrementEnable.Read() == 1 {
    if r.val >= r.max {
      r.val = 0
    } else {
      r.val++
    }
  }
}

func (r *Register) onRead() {
  if r.writeEnable.Read() == 1 {
    r.val = r.in.Read()
  }
}

// MEM = Memory (RAM)
// MAR = Memory Address Register
// MIC = Micro Instruction Counter
// IRE = Instruction Register

// OE = Output Enable
// WE = Write Enable
// RE = Reset

const MEM_OE = 0
const MEM_WE = 1

const MAR_WE = 2
const MIC_RE = 3
const IRE_WE = 4

func main() {
  falseBusBit := logisim.NewBusLiteral(1, 0)
  trueBusBit := logisim.NewBusLiteral(1, 1)

  clkLine := logisim.NewClock()
  clk := clkLine // Workaround

  microInstNum := logisim.NewBus(7)
  microInstContents := make([]uint64, 1 << uint64(microInstNum.Width()))
  microInstContents[0] = 1 << MAR_WE
  microInstContents[1] = 1 << MEM_OE | 1 << IRE_WE
  microInstContents[4] = 1 << MIC_RE
  //microInstContents := []uint64{1 << MAR_WE, 1 << MIC_RE}
  controlWord := logisim.NewBus(20)
  NewRom(microInstNum, controlWord, trueBusBit, clkLine, microInstContents)

  data := logisim.NewBus(8)

  /*  microInstNum.Write(1) // Set to dummy value to print whenit changes to 0
  microInstNum.OnChange(func () {
    fmt.Println(microInstNum.Read())
  })*/

  microInstNumSub := microInstNum.WriteableBranch(2, 1, 0)

  NewRegister(nil, microInstNumSub, falseBusBit, trueBusBit, trueBusBit, controlWord.Branch(MIC_RE), clkLine)

  memAddrBus := logisim.NewBus(8)
  NewRegister(data, memAddrBus, controlWord.Branch(MAR_WE), trueBusBit, falseBusBit, falseBusBit, clkLine)

  NewRegister(data, data, controlWord.Branch(IRE_WE), falseBusBit, falseBusBit, falseBusBit, clk)

  NewRam(memAddrBus, data, controlWord.Branch(MEM_WE), controlWord.Branch(MEM_OE), clkLine)
  fmt.Print(ram)

  data.Write(42)
  controlWord.Write(1 << MAR_WE)
  clk.Tick()
  controlWord.Write(0)

  data.Write(74)
  clk.Tick()

  fmt.Print(ram)

  controlWord.Write(1 << MEM_WE)
  fmt.Printf("CW %020b MEM_WE %v ram's %v\n", controlWord.Read(), we.Read(), ram.addrLatch)
  clk.Tick()

  fmt.Print(ram)

  clk.Ticks(10)
}
