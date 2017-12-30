package logisim

type Bus interface {
  OnChange(EventFunc)
  Read() uint64
  Write(uint64)
  Width() uint8
  Branch(PinMap) Bus
}

type bus struct {
  val uint64
  width uint8
  onChange []EventFunc
}

func NewBus(width uint8) Bus {
  return &bus{
    val: 0,
    width: width,
  }
}

func (b *bus) Branch(pinMap PinMap) Bus {
  return NewBranch(b, pinMap)
}

func (b *bus) OnChange(f EventFunc) {
  b.onChange = append(b.onChange, f)
}

func (b *bus) Read() uint64 {
  return b.val
}

func (b *bus) Write(val uint64) {
  if b.val != val {
    b.val = val

    for _, f := range b.onChange {
      f()
    }
  }
}

func (b *bus) Width() uint8 {
  return b.width
}
