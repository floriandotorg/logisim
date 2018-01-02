package logisim

type Const interface {
	ReadOnlyBus
	TriggerLine
}

type constStruct struct {
	val uint64
}

func (*constStruct) OnChange(ChangeFunc)             {}
func (c *constStruct) Read() uint64                  { return c.val }
func (*constStruct) Width() uint8                    { return 0 }
func (*constStruct) Branch(...uint8) ReadOnlyBus     { return nil }
func (*constStruct) TriggerBranch(uint8) TriggerLine { return nil }
func (*constStruct) OnRisingEdge(EventFunc)          {}
func (*constStruct) OnFallingEdge(EventFunc)         {}
func (c *constStruct) IsHigh() bool                  { return c.val == 1 }

func HIGH() Const {
	return &constStruct{
		val: 1,
	}
}

func LOW() Const {
	return &constStruct{
		val: 0,
	}
}
