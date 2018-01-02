package logisim

type TriggerLine interface {
	OnChange(EventFunc)
	OnRisingEdge(EventFunc)
	Read() bool
	Write(bool)
}

type triggerLine struct {
	val          bool
	onChange     []EventFunc
	onRisingEdge []EventFunc
}

func NewTriggerLine() TriggerLine {
	return &triggerLine{}
}

func (t *triggerLine) OnRisingEdge(f EventFunc) {
	t.onRisingEdge = append(t.onRisingEdge, f)
}

func (t *triggerLine) OnChange(f EventFunc) {
	t.onChange = append(t.onChange, f)
}

func (t *triggerLine) Read() bool {
	return t.val
}

func (t *triggerLine) Write(val bool) {
	if t.val != val {
		t.val = val

		for _, f := range t.onChange {
			f()
		}

		if t.val {
			for _, f := range t.onRisingEdge {
				f()
			}
		}
	}
}
