package logisim

type ReadOnlyBus interface {
	OnChange(ChangeFunc)
	Read() uint64
	Width() uint8
	Branch(...uint8) ReadOnlyBus
	TriggerBranch(uint8) TriggerLine
}

type Bus interface {
	ReadOnlyBus
	Write(uint64)
	WriteableBranch(...uint8) Bus
}

type bus struct {
	val      uint64
	width    uint8
	onChange []ChangeFunc
}

func NewBus(width uint8) Bus {
	return &bus{
		val:   0,
		width: width,
	}
}

func NewBusLiteral(width uint8, val uint64) ReadOnlyBus {
	return &bus{
		val:   val,
		width: width,
	}
}

func (b *bus) Branch(pinMap ...uint8) ReadOnlyBus {
	return NewBranch(b, pinMap...)
}

func (b *bus) WriteableBranch(pinMap ...uint8) Bus {
	return NewBranch(b, pinMap...)
}

func (b *bus) TriggerBranch(pin uint8) TriggerLine {
	return NewBusTriggerLine(b, pin)
}

func (b *bus) OnChange(f ChangeFunc) {
	b.onChange = append(b.onChange, f)
}

func (b *bus) Read() uint64 {
	return b.val
}

func (b *bus) Write(val uint64) {
	if b.val != val {
		old := b.val
		b.val = val

		for _, f := range b.onChange {
			f(old)
		}
	}
}

func (b *bus) Width() uint8 {
	return b.width
}
