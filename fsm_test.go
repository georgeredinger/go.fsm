package fsm

import "testing"

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

	tm := NewStateMachine(&delegate,
		Transition{From: "locked", Event: "coin", To: "unlocked", Action: "token_inc"},
		Transition{From: "locked", Event: OnEntry, Action: "enter"},
		Transition{From: "locked", Event: Default, To: "locked", Action: "default"},
		Transition{From: "unlocked", Event: "turn", To: "locked"},
		Transition{From: "unlocked", Event: OnExit, Action: "exit"},
	)

	var e Error

	if !(tm.currentState.From == "locked") {
		t.Errorf("state machine failure")
	}
	if !(delegate.count == 0) {
		t.Errorf("state machine failure")
	}
	if !(delegate.char == 0) {
		t.Errorf("state machine failure")
	}

	e = tm.Process("coin", 'i')
	if !(e == nil) {
		t.Errorf("state machine failure")
	}
	if !(tm.currentState.From == "unlocked") {
		t.Errorf("state machine failure")
	}
	if !(delegate.count == 1) {
		t.Errorf("state machine failure")
	}
	if !(delegate.char == 'i') {
		t.Errorf("state machine failure")
	}

	e = tm.Process("foobar", 'i')
	if !(e != nil) {
		t.Errorf("state machine failure")
	}
	if !(e.BadEvent() == "foobar") {
		t.Errorf("state machine failure")
	}
	if !(e.InState() == "unlocked") {
		t.Errorf("state machine failure")
	}
	if !(e.Error() == "state machine error: cannot find transition for event [foobar] when in state [unlocked]\n") {
		t.Errorf("state machine failure")
	}
	if !(tm.currentState.From == "unlocked") {
		t.Errorf("state machine failure")
	}
	if !(delegate.count == 1) {
		t.Errorf("state machine failure")
	}
	if !(delegate.char == 'i') {
		t.Errorf("state machine failure")
	}

	e = tm.Process("turn", 'q')
	if !(e == nil) {
		t.Errorf("state machine failure")
	}
	if !(tm.currentState.From == "locked") {
		t.Errorf("state machine failure")
	}
	if !(delegate.count == 1) {
		t.Errorf("state machine failure")
	}
	if !(delegate.entered == 8) {
		t.Errorf("state machine failure, %d", delegate.entered)
	}

	e = tm.Process("random", 'p')
	if !(e == nil) {
		t.Errorf("state machine failure")
	}
	if !(tm.currentState.From == "locked") {
		t.Errorf("state machine failure")
	}
	if !(delegate.entered == 88) {
		t.Errorf("state machine failure, %d", delegate.entered)
	}
}
