package logisim

type branch struct {
	parent Bus
	pinMap []uint8
	mask   uint64
}

func NewBranch(parent Bus, pinMap ...uint8) Bus {
	var mask uint64

	for _, p := range pinMap {
		if p >= parent.Width() {
			panic("you idiot")
		}

		mask |= 1 << p
	}

	return &branch{
		parent: parent,
		pinMap: pinMap,
		mask:   mask,
	}
}

func (b *branch) Branch(pinMap ...uint8) ReadOnlyBus {
	return NewBranch(b, pinMap...)
}

func (b *branch) TriggerBranch(pin uint8) TriggerLine {
	return NewBusTriggerLine(b, pin)
}

func (b *branch) WriteableBranch(pinMap ...uint8) Bus {
	return NewBranch(b, pinMap...)
}

func (b *branch) OnChange(f ChangeFunc) {
	b.parent.OnChange(f)
}

func (b *branch) Read() uint64 {
	var result uint64
	val := b.parent.Read()

	for n, p := range b.pinMap {
		result |= ((val >> p) & 0x01) << uint8(n)
	}

	return result
}

func (b *branch) Write(newVal uint64) {
	var val uint64

	for n, p := range b.pinMap {
		val |= ((newVal >> uint8(n)) & 0x01) << p
	}

	b.parent.Write((b.parent.Read() & ^b.mask) | val)
}

func (b *branch) Width() uint8 {
	return uint8(len(b.pinMap))
}
