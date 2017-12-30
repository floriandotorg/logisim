package bitarray

import (
  "fmt"
  "strconv"
)

type BitArray interface {
  Set(uint64) BitArray
  SetString(string) BitArray
  SetBit(uint8) BitArray
  ClearBit(uint8) BitArray
  Get() uint64
  Size() uint8
  Or(BitArray) BitArray
  Equals(BitArray) bool
}

type bitArray struct {
  val uint64
  size uint8
}

func NewBitArray(size uint8) BitArray {
  return &bitArray{
    size: size,
  }
}

func NewBitArrayWithVal(size uint8, val uint64) BitArray {
  return NewBitArray(size).Set(val)
}

func (b* bitArray) bitmask() uint64 {
  return ^uint64(0) >> (64 - b.size)
}

func (b* bitArray) Set(val uint64) BitArray {
  return &bitArray{
    val: val & b.bitmask(),
    size: b.size,
  }
}

func (b* bitArray) SetString(val string) BitArray {
  if val, err := strconv.ParseUint(val, 16, 64); err == nil {
    return &bitArray{
      val: val & b.bitmask(),
      size: b.size,
    }
  }

  return nil
}

func (b* bitArray) SetBit(bit uint8) BitArray {
  if bit >= b.size {
    panic("out of range")
  }

  return &bitArray{
    val: b.val | (1 << bit),
    size: b.size,
  }
}

func (b* bitArray) ClearBit(bit uint8) BitArray {
  if bit >= b.size {
    panic("out of range")
  }

  return &bitArray{
    val: b.val & ^(1 << bit),
    size: b.size,
  }
}

func (b* bitArray) Get() uint64 {
  return b.val
}

func (b* bitArray) Size() uint8 {
  return b.size
}

func (b* bitArray) Or(other BitArray) BitArray {
  if ba, ok := other.(*bitArray); ok {
    if b.size != ba.size {
      panic("trying Or on BitArray of different size")
    }

    return &bitArray{
      size: b.size,
      val: (b.val | ba.val) & b.bitmask(),
    }
  }

  return nil
}

func (b* bitArray) Equals(other BitArray) bool {
  if ba, ok := other.(*bitArray); ok {
    return b.size == ba.size && b.val == ba.val
  }

  return false
}

func (b* bitArray) String() string {
  if b.size > 1 {
    return fmt.Sprintf("[0x%02X]", b.val)
  } else {
    return fmt.Sprintf("[%d]", b.val)
  }
}
