package fsm

import "testing"

import "github.com/sdegutis/assert"

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

  rules := []Rule{
    {From: "locked", Event: "coin", To: "unlocked", Action: "token_inc"},
    {From: "locked", Event: OnEntry, Action: "enter"},
    {From: "locked", Event: Default, To: "locked", Action: "default"},
    {From: "unlocked", Event: "turn", To: "locked"},
    {From: "unlocked", Event: OnExit, Action: "exit"},
  }

  tm := NewStateMachine(rules, &delegate)

  var e *Error

  assert.Equals(t, tm.CurrentState, "locked")
  assert.Equals(t, delegate.count, 0)
  assert.Equals(t, int(delegate.char), 0)

  e = tm.Process("coin", 'i')
  assert.True(t, e == nil)
  assert.Equals(t, tm.CurrentState, "unlocked")
  assert.Equals(t, delegate.count, 1)
  assert.Equals(t, delegate.char, 'i')

  e = tm.Process("foobar", 'i')
  assert.True(t, e != nil)
  assert.Equals(t, e.BadEvent, "foobar")
  assert.Equals(t, e.InState, "unlocked")
  assert.Equals(t, e.Error(), "state machine error: cannot find rule for event [foobar] when in state [unlocked]\n")
  assert.Equals(t, tm.CurrentState, "unlocked")
  assert.Equals(t, delegate.count, 1)

  e = tm.Process("turn", 'q')
  assert.True(t, e == nil)
  assert.Equals(t, tm.CurrentState, "locked")
  assert.Equals(t, delegate.count, 1)
  assert.Equals(t, delegate.entered, 8)

  e = tm.Process("random", 'p')
  assert.True(t, e == nil)
  assert.Equals(t, tm.CurrentState, "locked")
  assert.Equals(t, delegate.entered, 88)
}
