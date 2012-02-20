package fsm

import "testing"
import "runtime"
import "path"

type testTokenMachineDelegate struct {
  count   int
  char    rune
  entered int
}

func (delegate *testTokenMachineDelegate) StateMachineCallback(action string, args []interface{}) {
  switch action {
  case "token_inc":
    delegate.count++
    delegate.char = args[0].(rune)
  case "enter":
    delegate.entered++
  case "exit":
    delegate.entered = 7
  case "default":
    delegate.entered = 88
  }
}

func TestTokenMachine(t *testing.T) {
  var delegate testTokenMachineDelegate

  rules := []StateMachineRule{
    {From: "locked", Event: "coin", To: "unlocked", Action: "token_inc"},
    {From: "locked", Event: OnEntry, Action: "enter"},
    {From: "locked", Event: Default, To: "locked", Action: "default"},
    {From: "unlocked", Event: "turn", To: "locked"},
    {From: "unlocked", Event: OnExit, Action: "exit"},
  }

  tm := NewStateMachine(rules, &delegate)

  var e Error

  assertEquals(t, tm.CurrentState(), "locked")
  assertEquals(t, delegate.count, 0)
  assertEquals(t, int(delegate.char), 0)

  e = tm.Process("coin", 'i')
  assertEquals(t, e, nil)
  assertEquals(t, tm.CurrentState(), "unlocked")
  assertEquals(t, delegate.count, 1)
  assertEquals(t, delegate.char, 'i')

  e = tm.Process("foobar", 'i')
  assertEquals(t, e == nil, false)
  assertEquals(t, e.BadEvent(), "foobar")
  assertEquals(t, e.InState(), "unlocked")
  assertEquals(t, e.Error(), "state machine error: cannot find rule for event [foobar] when in state [unlocked]\n")
  assertEquals(t, tm.CurrentState(), "unlocked")
  assertEquals(t, delegate.count, 1)

  e = tm.Process("turn", 'q')
  assertEquals(t, e, nil)
  assertEquals(t, tm.CurrentState(), "locked")
  assertEquals(t, delegate.count, 1)
  assertEquals(t, delegate.entered, 8)

  e = tm.Process("random", 'p')
  assertEquals(t, e, nil)
  assertEquals(t, tm.CurrentState(), "locked")
  assertEquals(t, delegate.entered, 88)
}

func assertEquals(t *testing.T, got interface{}, expected interface{}) {
  if got != expected {
    _, file, line, _ := runtime.Caller(1)
    t.Errorf("___ [%s:%d] state machine failure; got %v but expected %v", path.Base(file), line, got, expected)
  }
}
