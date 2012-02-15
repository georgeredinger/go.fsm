package fsm

import "testing"

type testTokenMachineContext struct {
	count   int
	char    rune
	entered int
}

func (ctx *testTokenMachineContext) StateMachineCallback(action string, args []interface{}) {
	switch action {
	case "token_inc":
		ctx.count++
		ctx.char = args[0].(rune)
	case "enter":
		ctx.entered++
	case "exit":
		ctx.entered = 7
	case "default":
		ctx.entered = 88
	}
}

func TestTokenMachine(t *testing.T) {
	var ctx testTokenMachineContext

	tm := NewStateMachine(&ctx,
		Transition{From: "locked", Event: "coin", To: "unlocked", Action: "token_inc"},
		Transition{From: "locked", Event: OnEntry, Action: "enter"},
		Transition{From: "locked", Event: Default, To: "locked", Action: "default"},
		Transition{From: "unlocked", Event: "turn", To: "locked"},
		Transition{From: "unlocked", Event: OnExit, Action: "exit"},
	)

	if !(tm.currentState.From == "locked") {
		t.Errorf("state machine failure")
	}
	if !(ctx.count == 0) {
		t.Errorf("state machine failure")
	}
	if !(ctx.char == 0) {
		t.Errorf("state machine failure")
	}

	tm.Process("coin", 'i')
	if !(tm.currentState.From == "unlocked") {
		t.Errorf("state machine failure")
	}
	if !(ctx.count == 1) {
		t.Errorf("state machine failure")
	}
	if !(ctx.char == 'i') {
		t.Errorf("state machine failure")
	}

	tm.Process("turn", 'q')
	if !(tm.currentState.From == "locked") {
		t.Errorf("state machine failure")
	}
	if !(ctx.count == 1) {
		t.Errorf("state machine failure")
	}
	if !(ctx.entered == 8) {
		t.Errorf("state machine failure, %d", ctx.entered)
	}

	tm.Process("random", 'p')
	if !(tm.currentState.From == "locked") {
		t.Errorf("state machine failure")
	}
	if !(ctx.entered == 88) {
		t.Errorf("state machine failure, %d", ctx.entered)
	}
}
