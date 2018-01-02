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

	contents []uint64
}

type RamParams struct {
	addr         logisim.ReadOnlyBus
	data         logisim.Bus
	writeEnable  logisim.TriggerLine
	outputEnable logisim.TriggerLine
}

func NewRam(p RamParams) *Ram {
	ram := &Ram{
		addr:         p.addr,
		data:         p.data,
		writeEnable:  p.writeEnable,
		outputEnable: p.outputEnable,
		contents:     make([]uint64, 1<<uint64(p.addr.Width())),
	}
	p.addr.OnChange(ram.onAddrChange)
	p.outputEnable.OnRisingEdge(ram.onOutput)
	p.writeEnable.OnFallingEdge(ram.onInput)
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
	if r.outputEnable.IsHigh() {
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

type RegisterParams struct {
	in              logisim.ReadOnlyBus
	out             logisim.Bus
	writeEnable     logisim.TriggerLine
	outputEnable    logisim.TriggerLine
	incrementEnable logisim.TriggerLine
	reset           logisim.TriggerLine
	clockLine       logisim.TriggerLine
}

func NewRegister(p RegisterParams) *Register {
	register := &Register{
		in:              p.in,
		out:             p.out,
		writeEnable:     p.writeEnable,
		outputEnable:    p.outputEnable,
		incrementEnable: p.incrementEnable,
		reset:           p.reset,
		clockLine:       p.clockLine,
		max:             1 << p.out.Width(),
	}
	p.clockLine.OnRisingEdge(register.onInput)
	p.reset.OnRisingEdge(register.onReset)
	p.outputEnable.OnRisingEdge(register.onOutput)
	return register
}

func (r *Register) onReset() {
	r.val = 0

	if r.outputEnable.IsHigh() {
		r.out.Write(r.val)
	}
}

func (r *Register) onInput() {
	if !r.reset.IsHigh() {
		if r.writeEnable.IsHigh() {
			r.val = r.in.Read()
		}

		if r.outputEnable.IsHigh() {
			r.out.Write(r.val)
		}

		if r.incrementEnable.IsHigh() {
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
	clock := logisim.NewClock()
	clockLine := clock.GetClockLine()

	microInstNum := logisim.NewBus(7)
	microInstContents := make([]uint64, 1<<uint64(microInstNum.Width()))
	microInstContents[0] = 1 << MAR_WE
	microInstContents[1] = 1<<MEM_OE | 1<<IRE_WE
	microInstContents[4] = 1 << MIC_RE
	controlWord := logisim.NewBus(20)
	controlRom := NewRam(RamParams{
		addr:         microInstNum,
		data:         controlWord,
		writeEnable:  logisim.LOW(),
		outputEnable: logisim.HIGH(),
	})
	controlRom.write(0, microInstContents)

	data := logisim.NewBus(8)

	fmt.Print(controlRom)

	microInstNumSub := microInstNum.WriteableBranch(0, 1, 2)

	// Micro Ins
	NewRegister(RegisterParams{
		in:              nil,
		out:             microInstNumSub,
		writeEnable:     logisim.LOW(),
		outputEnable:    logisim.HIGH(),
		incrementEnable: logisim.HIGH(),
		reset:           controlWord.TriggerBranch(MIC_RE),
		clockLine:       clockLine,
	})

	memAddrBus := logisim.NewBus(8)

	// Mem Addr
	NewRegister(RegisterParams{
		in:              data,
		out:             memAddrBus,
		writeEnable:     controlWord.TriggerBranch(MAR_WE),
		outputEnable:    logisim.HIGH(),
		incrementEnable: logisim.LOW(),
		reset:           logisim.LOW(),
		clockLine:       clockLine,
	})

	// Ins
	NewRegister(RegisterParams{
		in:              data,
		out:             data,
		writeEnable:     controlWord.TriggerBranch(IRE_WE),
		outputEnable:    logisim.LOW(),
		incrementEnable: logisim.LOW(),
		reset:           logisim.LOW(),
		clockLine:       clockLine,
	})

	NewRam(RamParams{
		addr:         memAddrBus,
		data:         data,
		writeEnable:  controlWord.TriggerBranch(MEM_WE),
		outputEnable: controlWord.TriggerBranch(MEM_OE),
	})

	// dirty hack, todo: remove
	controlRom.onOutput()

	for n := 0; n < 10; n++ {
		clock.Tick()
		fmt.Printf("mI# %04b[%03b] CW %020b\n", microInstNum.Read()>>3, microInstNumSub.Read(), controlWord.Read())
	}
}
