package main

import (
	"./logisim"
	"bytes"
	"fmt"
)

type Ram struct {
	addr         logisim.ReadOnlyBus
	data         logisim.Bus
	writeEnable  logisim.TriggerLine
	outputEnable logisim.TriggerLine
	clockLine    logisim.TriggerLine

	contents []uint64
}

func NewRam(addr logisim.ReadOnlyBus, data logisim.Bus, writeEnable, outputEnable, clockLine logisim.TriggerLine) *Ram {
	ram := &Ram{
		addr:         addr,
		data:         data,
		writeEnable:  writeEnable,
		outputEnable: outputEnable,
		clockLine:    clockLine,
		contents:     make([]uint64, 1<<uint64(addr.Width())),
	}
	addr.OnChange(ram.onAddrChange)
	outputEnable.OnRisingEdge(ram.onOutput)
	writeEnable.OnFallingEdge(ram.onInput)
	return ram
}

func (r *Ram) write(offset int, data []uint64) {
	if offset+len(data) > len(r.contents) {
		panic("This has gone too far!")
	}

	for n := 0; n < len(data); n++ {
		r.contents[n+offset] = data[n]
	}
}

func (r *Ram) onAddrChange(old uint64) {
	if r.outputEnable.Read() {
		r.onOutput()
	}
}

func (r *Ram) onInput() {
	r.contents[r.addr.Read()] = r.data.Read()
}

func (r *Ram) onOutput() {
	r.data.Write(r.contents[r.addr.Read()])
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
	in              logisim.ReadOnlyBus
	out             logisim.Bus
	writeEnable     logisim.TriggerLine
	outputEnable    logisim.TriggerLine
	incrementEnable logisim.TriggerLine
	reset           logisim.TriggerLine
	clockLine       logisim.TriggerLine

	max uint64
	val uint64
}

func NewRegister(in logisim.ReadOnlyBus, out logisim.Bus, writeEnable, outputEnable, incrementEnable, reset, clockLine logisim.TriggerLine) *Register {
	register := &Register{
		in:              in,
		out:             out,
		writeEnable:     writeEnable,
		outputEnable:    outputEnable,
		incrementEnable: incrementEnable,
		reset:           reset,
		clockLine:       clockLine,

		max: 1 << out.Width(),
	}
	clockLine.OnRisingEdge(register.onInput)
	reset.OnRisingEdge(register.onReset)
	outputEnable.OnRisingEdge(register.onOutput)
	return register
}

func (r *Register) onReset() {
	r.val = 0

	if r.outputEnable.Read() {
		r.out.Write(r.val)
	}
}

func (r *Register) onInput() {
	if !r.reset.Read() {
		if r.writeEnable.Read() {
			r.val = r.in.Read()
		}

		if r.outputEnable.Read() {
			r.out.Write(r.val)
		}

		if r.incrementEnable.Read() {
			r.val = (r.val + 1) % r.max
		}
	}
}

func (r *Register) onOutput() {
	r.out.Write(r.val)
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

	falseTrigger := falseBusBit.TriggerBranch(0)
	trueTrigger := trueBusBit.TriggerBranch(0)

	clock := logisim.NewClock()
	clockLine := clock.GetClockLine()

	microInstNum := logisim.NewBus(7)
	microInstContents := make([]uint64, 1<<uint64(microInstNum.Width()))
	microInstContents[0] = 1 << MAR_WE
	microInstContents[1] = 1<<MEM_OE | 1<<IRE_WE
	microInstContents[4] = 1 << MIC_RE
	controlWord := logisim.NewBus(20)
	controlRom := NewRam(microInstNum, controlWord, falseTrigger, trueTrigger, clockLine)
	controlRom.write(0, microInstContents)

	data := logisim.NewBus(8)

	fmt.Print(controlRom)

	microInstNumSub := microInstNum.WriteableBranch(0, 1, 2)

	// Micro Ins
	NewRegister(nil, microInstNumSub, falseTrigger, trueTrigger, trueTrigger, controlWord.TriggerBranch(MIC_RE), clockLine)

	memAddrBus := logisim.NewBus(8)
	// Mem Addr
	NewRegister(data, memAddrBus, controlWord.TriggerBranch(MAR_WE), trueTrigger, falseTrigger, falseTrigger, clockLine)

	// Ins
	NewRegister(data, data, controlWord.TriggerBranch(IRE_WE), falseTrigger, falseTrigger, falseTrigger, clockLine)

	NewRam(memAddrBus, data, controlWord.TriggerBranch(MEM_WE), controlWord.TriggerBranch(MEM_OE), clockLine)

	// TODO
	// - const VCC/GND
	// - parameter structs
	// - REMOVE NEXT LINE
	controlRom.onOutput()

	for n := 0; n < 10; n++ {
		clock.Tick()
		fmt.Printf("mI# %04b[%03b] CW %020b\n", microInstNum.Read()>>3, microInstNumSub.Read(), controlWord.Read())
	}
}
