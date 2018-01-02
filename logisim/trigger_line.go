package logisim

type TriggerLine interface {
	OnRisingEdge(EventFunc)
	OnFallingEdge(EventFunc)
	Read() bool
}

type triggerLine struct {
	parent        ReadOnlyBus
	mask          uint64
	onRisingEdge  []EventFunc
	onFallingEdge []EventFunc
}

func NewTriggerLine() TriggerLine {
	return NewBusTriggerLine(NewBus(1), 0)
}

func NewBusTriggerLine(bus ReadOnlyBus, pin uint8) TriggerLine {
	triggerLine := &triggerLine{
		parent: bus,
		mask:   1 << pin,
	}
	triggerLine.parent.OnChange(triggerLine.onChange)
	return triggerLine
}

func (t *triggerLine) OnRisingEdge(f EventFunc) {
	t.onRisingEdge = append(t.onRisingEdge, f)
}

func (t *triggerLine) OnFallingEdge(f EventFunc) {
	t.onFallingEdge = append(t.onFallingEdge, f)
}

func (t *triggerLine) onChange(old uint64) {
	if t.parent.Read()&t.mask != old&t.mask {
		if t.Read() {
			for _, f := range t.onRisingEdge {
				f()
			}
		} else {
			for _, f := range t.onFallingEdge {
				f()
			}
		}
	}
}

func (t *triggerLine) Read() bool {
	return t.parent.Read()&t.mask != 0x00
}
