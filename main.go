package main

import (
	"./logisim"
	"bytes"
	"fmt"
)

type Rom struct {
	addr         logisim.ReadOnlyBus
	addrLatch    uint64
	data         logisim.Bus
	outputEnable logisim.ReadOnlyBus
	clk          logisim.Clock

	contents []uint64
}

func NewRom(addr logisim.ReadOnlyBus, data logisim.Bus, outputEnable logisim.ReadOnlyBus, clk logisim.Clock, contents []uint64) *Rom {
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
	addr         logisim.ReadOnlyBus
	addrLatch    uint64
	data         logisim.Bus
	writeEnable  logisim.ReadOnlyBus
	outputEnable logisim.ReadOnlyBus
	clk          logisim.Clock

	contents []uint64
}

func NewRam(addr logisim.ReadOnlyBus, data logisim.Bus, writeEnable, outputEnable logisim.ReadOnlyBus, clk logisim.Clock) *Ram {
	ram := &Ram{
		addr:         addr,
		data:         data,
		writeEnable:  writeEnable,
		outputEnable: outputEnable,
		clk:          clk,
		contents:     make([]uint64, 1<<uint64(addr.Width())),
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

// OE = Output Enable
// WE = Write Enable

const MEM_OE = 0
const MEM_WE = 1

func main() {
	falseBusBit := logisim.NewBusLiteral(1, 0)
	trueBusBit := logisim.NewBusLiteral(1, 1)

	clkLine := logisim.NewClock()
	clk := clkLine // Workaround

	microInstNum := logisim.NewBus(7)
	microInstContents := make([]uint64, 1<<uint64(microInstNum.Width()))
	microInstContents[0] = 1 << MAR_WE
	microInstContents[1] = 1<<MEM_OE | 1<<IRE_WE
	microInstContents[4] = 1 << MIC_RE
	//microInstContents := []uint64{1 << MAR_WE, 1 << MIC_RE}
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
