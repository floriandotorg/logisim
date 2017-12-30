package logisim

type ReadOnlyBus interface {
  OnChange(EventFunc)
  Read() uint64
  Width() uint8
  Branch(...uint8) ReadOnlyBus
}

type Bus interface {
  ReadOnlyBus
  Write(uint64)
  WriteableBranch(...uint8) Bus
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

func NewBusLiteral(width uint8, val uint64) ReadOnlyBus {
  return &bus{
    val: val,
    width: width,
  }
}

func (b *bus) Branch(pinMap ...uint8) ReadOnlyBus {
  return NewBranch(b, pinMap...)
}

func (b *bus) WriteableBranch(pinMap ...uint8) Bus {
  return NewBranch(b, pinMap...)
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
