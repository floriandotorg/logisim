package main

import (
	"./logisim"
	"bytes"
	"fmt"
)

type Rom struct {
	addr         logisim.ReadOnlyBus
	data         logisim.Bus
	outputEnable logisim.ReadOnlyBus
	clk          logisim.TriggerLine

	contents []uint64
}

func NewRom(addr logisim.ReadOnlyBus, data logisim.Bus, outputEnable logisim.ReadOnlyBus, clk logisim.TriggerLine, contents []uint64) *Rom {
	if outputEnable.Width() != 1 {
		panic("FU")
	}
	if len(contents) != 1<<uint64(addr.Width()) {
		panic("ROM and contents size differ")
	}
	rom := &Rom{
		addr:         addr,
		data:         data,
		outputEnable: outputEnable,
		clk:          clk,
		contents:     contents,
	}
	clk.OnRisingEdge(rom.onTick)
	return rom
}

func (r *Rom) onTick() {
	status := r.outputEnable.Read()
	addr := r.addr.Read()
	if status == 0x01 {
		r.data.Write(r.contents[addr])
	}
}

type Ram struct {
	addr         logisim.ReadOnlyBus
	data         logisim.Bus
	writeEnable  logisim.ReadOnlyBus
	outputEnable logisim.ReadOnlyBus
	clk          logisim.TriggerLine

	contents []uint64
}

func NewRam(addr logisim.ReadOnlyBus, data logisim.Bus, writeEnable, outputEnable logisim.ReadOnlyBus, clk logisim.TriggerLine) *Ram {
	ram := &Ram{
		addr:         addr,
		data:         data,
		writeEnable:  writeEnable,
		outputEnable: outputEnable,
		clk:          clk,
		contents:     make([]uint64, 1<<uint64(addr.Width())),
	}
	clk.OnRisingEdge(ram.onTick)
	return ram
}

func (r *Ram) onTick() {
	addr := r.addr.Read()
	if r.outputEnable.Read() == 1 {
		r.data.Write(r.contents[addr])
	} else if r.writeEnable.Read() == 1 {
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

// OE = Output Enable
// WE = Write Enable

const MEM_OE = 0
const MEM_WE = 1

func main() {
	clkLine := logisim.NewTriggerLine()
	clk := logisim.NewClock(clkLine)

	controlWord := logisim.NewBus(20)
	addr := logisim.NewBus(7)
	data := logisim.NewBus(8)

	ram := NewRam(addr, data, controlWord.Branch(MEM_WE), controlWord.Branch(MEM_OE), clkLine)

	fmt.Print(ram)

	addr.Write(42)
	data.Write(74)
	clk.Tick()

	fmt.Print(ram)

	controlWord.Write(1 << MEM_WE)
	clk.Tick()

	fmt.Print(ram)
}
