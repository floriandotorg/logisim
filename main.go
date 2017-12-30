package main

import (
  "fmt"
  "strings"
  "strconv"
  "./bitarray"
)

type pinMap map[string]uint64
type stateType map[string]bitarray.BitArray
type changeSet stateType
type connectionArray map[string][]string
type changeFuncType func (changeSet)
type triggerType int
const (
  EVERY_CHANGE triggerType = 0
  RISING_EDGE  triggerType = 1
  FALLING_EDGE triggerType = 2
)
type triggerMap map[triggerType][]string

type comp interface {
  getName() string
  getPins() pinMap
}

type triggeredComp interface {
  comp
  getTrigger() triggerMap
  run(state stateType) changeSet
}

type triggeringComp interface {
  comp
  setChangeFunc(changeFunc changeFuncType)
}

// func (s* stateType) String() string {
//   buffer := bytes.Buffer{}
//   buffer.WriteString("{\n")
//   for name, val := range *s {
//     buffer.WriteString(fmt.Sprintf("  %s: %v,\n", name, val))
//   }
//   buffer.WriteString("}")
//   return buffer.String()
// }

// func (c* changeSet) String() string {
//   return (*stateType)(c).String()
// }

type engine struct {
  state stateType
  pinToState map[string]string
  triggeredComponents map[triggerType]map[string][]triggeredComp
  runQueue []changeSet
}

func newEngine(components []comp, connections connectionArray) *engine {
  engine := new(engine)
  engine.init(components, connections)
  return engine
}

func (e *engine) init(components []comp, connections connectionArray) {
  pins := make(pinMap)
  for _, comp := range components {
    for pinName, pinWidth := range comp.getPins() {
      pins[fmt.Sprintf("%s.%s", comp.getName(), pinName)] = pinWidth
    }
  }

  e.state = make(stateType)
  e.pinToState = make(map[string]string)
  for name, connection := range connections {
    bitarr := bitarray.NewBitArray(pins[connection[0]])
    e.state[name] = bitarr

    for _, pinName := range connection {
      if pins[pinName] != bitarr.Size() {
        panic(fmt.Sprintf("bus width not matching: %s (%d), pin: %s (%d)", name, bitarr.Size(), pinName, pins[pinName]))
      }
      e.pinToState[pinName] = name
    }
  }

  e.triggeredComponents = map[triggerType]map[string][]triggeredComp{
    EVERY_CHANGE: map[string][]triggeredComp{},
    RISING_EDGE: map[string][]triggeredComp{},
    FALLING_EDGE: map[string][]triggeredComp{},
  }

  for _, comp := range components {
    if tcomp, ok := comp.(triggeredComp); ok {
      for ttype, pins := range tcomp.getTrigger() {
        for _, pin := range pins {
          e.triggeredComponents[ttype][fmt.Sprintf("%s.%s", tcomp.getName(), pin)] = append(e.triggeredComponents[ttype][pin], tcomp)
          // if e.triggeredComponents[ttype][pin] == nil {
          //   e.triggeredComponents[ttype][pin] = []triggeredComp{tcomp}
          // }
        }
      }
    }
  }

  for _, comp := range components {
    if tcomp, ok := comp.(triggeringComp); ok {
      tcomp.setChangeFunc(func(pinChangeSet changeSet) {
        chgSet := changeSet{}
        for name, val := range pinChangeSet {
          chgSet[e.pinToState[fmt.Sprintf("%s.%s", tcomp.getName(), name)]] = val
        }
        e.runQueue = append(e.runQueue, chgSet)
      })
    }
  }

  e.runQueue = []changeSet{}

  fmt.Printf("Initial State\n%v\n", e.state)
}

func mergeStates(a, b stateType) {
  for name, val := range b {
    if _, ok := a[name]; ok {
      a[name] = a[name].Or(val)
    } else {
      a[name] = val
    }
  }
}

func (e *engine) runComponents(comps []triggeredComp) changeSet {
  changeSets := []changeSet{}
  for _, comp := range comps {
    state := make(stateType)
    for name, stateName := range e.pinToState {
      if (strings.Contains(name, comp.getName())) {
        state[strings.Replace(name, fmt.Sprintf("%s.", comp.getName()), "", -1)] = e.state[stateName]
      }
    }

    if chgSet := comp.run(state); chgSet != nil {
      changeSets = append(changeSets, chgSet)
    }
  }

  master := changeSet{}
  for _, chgSet := range changeSets {
    mergeStates(stateType(master), stateType(chgSet))
  }

  masterState := changeSet{}
  for name, val := range master {
    masterState[e.pinToState[name]] = val
  }

  for name, val := range masterState {
    if (e.state[name].Equals(val)) {
      delete(masterState, name)
    }
  }

  return masterState
}

func (e *engine) step(chgSet changeSet) {
  fmt.Printf("Step changeSet\n%v\n", chgSet)

  for name, val := range chgSet {
    if (e.state[name].Equals(val)) {
      delete(chgSet, name)
    }
  }

  if len(chgSet) > 0 {
    pins := []string{}
    for name := range chgSet {
      for pin, state := range e.pinToState {
        if state == name {
          pins = append(pins, pin)
        }
      }
    }

    triggeredComponents := []triggeredComp{}
    for _, pin := range pins {
      if trigger := e.triggeredComponents[EVERY_CHANGE][pin]; trigger != nil {
        triggeredComponents = append(triggeredComponents, trigger...)
      }
    }

    mergeStates(e.state, stateType(chgSet))

    master := e.runComponents(triggeredComponents)

    if len(master) > 0 {
      mergeStates(e.state, stateType(master))
      fmt.Printf("Master\n%v", master)
      fmt.Printf("State\n%v", e.state)
      e.step(master)
    }
  }
}

func (e *engine) run() {
  for len(e.runQueue) > 0 {
    e.step(e.runQueue[0])
    e.runQueue = e.runQueue[1:]
  }
}

type dipSwitch struct {
  name string
  pins pinMap
  changeFunc changeFuncType
}

func newDipSwitch(name string, n int) *dipSwitch {
  swt := new(dipSwitch)
  swt.name = name
  swt.pins = make(pinMap)
  for ;n > 0; n-- {
    swt.pins[strconv.Itoa(n)] = 1
  }
  return swt
}

func (d *dipSwitch) getName() string {
  return fmt.Sprintf("dip_switch.%s", d.name)
}

func (d *dipSwitch) getPins() pinMap {
  return d.pins
}

func (d *dipSwitch) setChangeFunc(changeFunc changeFuncType) {
  d.changeFunc = changeFunc
}

func (d *dipSwitch) setSwt(n int, to bool) {
  chgSet := changeSet{}
  if to {
    chgSet[strconv.Itoa(n)] = bitarray.NewBitArrayWithVal(1, 1)
  } else {
    chgSet[strconv.Itoa(n)] = bitarray.NewBitArray(1)
  }
  d.changeFunc(chgSet)
}

type ledArray struct {
  name string
  pins pinMap
}

func newLedArray(name string, n int) *ledArray {
  leds := new(ledArray)
  leds.name = name
  leds.pins = make(pinMap)
  for ;n > 0; n-- {
    leds.pins[strconv.Itoa(n)] = 1
  }
  return leds
}

func (l *ledArray) getName() string {
  return fmt.Sprintf("led_array.%s", l.name)
}

func (l *ledArray) getPins() pinMap {
  return l.pins
}

func (l *ledArray) getTrigger() triggerMap {
  triggers := make([]string, len(l.pins))
  n := 0
  for name := range l.pins {
    triggers[n] = name
    n++
  }
  return triggerMap{
    EVERY_CHANGE: triggers,
  }
}

func (l *ledArray) run(state stateType) changeSet {
  fmt.Printf("LED\n%v\n", state)
  return nil
}

func main() {
  a := newDipSwitch("out", 4)
  b := newLedArray("in", 4)

  conn := connectionArray{
    "line1": []string{"dip_switch.out.1", "led_array.in.1"},
    "line2": []string{"dip_switch.out.2", "led_array.in.2"},
    "line3": []string{"dip_switch.out.3", "led_array.in.3"},
    "line4": []string{"dip_switch.out.4", "led_array.in.4"},
  }
  engine := newEngine([]comp{a,b}, conn)

  a.setSwt(1, true)
  engine.run()
}
