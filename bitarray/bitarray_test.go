package bitarray

import (
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestBitArrayNew(t *testing.T) {
  assert.Equal(t, uint64(0), NewBitArray(1).Get(), "New")
  assert.Equal(t, uint64(1), NewBitArrayWithVal(1, 1).Get(), "NewWithVal")
  assert.Equal(t, uint64(1), NewBitArrayWithVal(1, 0xfffffff).Get(), "NewWithVal (overflow)")
}

func TestBitArrayStringer(t *testing.T) {
  assert.Equal(t, "[0]", fmt.Sprintf("%v", NewBitArray(1)), "Stringer (binary)")
  assert.Equal(t, "[0x00]", fmt.Sprintf("%v", NewBitArray(5)), "Stringer (hex)")
}

func TestBitArrayEquals(t *testing.T) {
  assert.Equal(t, true, NewBitArray(1).Equals(NewBitArray(1)), "Equals (0)")
  assert.Equal(t, true, NewBitArrayWithVal(1, 1).Equals(NewBitArrayWithVal(1, 1)), "Equals (1)")
  assert.Equal(t, true, NewBitArrayWithVal(4, 0xf).Equals(NewBitArrayWithVal(4, 0xf)), "Equals (0xf)")
  assert.NotEqual(t, true, NewBitArrayWithVal(4, 0xf).Equals(NewBitArray(1)), "Equals (different size)")
}

func TestBitArrayOr(t *testing.T) {
  bitarr := NewBitArrayWithVal(8, 0xf0).Or(NewBitArrayWithVal(8, 0x0f))
  assert.Equal(t, uint64(0xff), bitarr.Get(), "Or")
}

func TestBitArraySetBit(t *testing.T) {
  assert.Equal(t, uint64(0x1), NewBitArray(1).SetBit(0).Get(), "SetBit")
}

func TestBitArrayClearBit(t *testing.T) {
  assert.Equal(t, uint64(0), NewBitArrayWithVal(1, 1).ClearBit(0).Get(), "ClearBit")
}
