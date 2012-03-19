// Finite State Machines, in idiomatic Go1.
//
// Here is the basic API:
//
//     sm := []fsm.Rule{
//
//       { From: "locked",    Event: "coin",     To: "unlocked",  Action: "token_inc" },
//       { From: "locked",    Event: OnEntry,                     Action: "enter" },
//       { From: "locked",    Event: Default,    To: "locked",    Action: "default" },
//
//       { From: "unlocked",  Event: "turn",     To: "locked",    },
//       { From: "unlocked",  Event: OnExit,                      Action: "exit" },
//
//     }
//
//     sm := fsm.NewStateMachine(rules, &delegate)
//
//     sm.Process("coin")
//     sm.Process("turn", optionalArg, ...)
//     sm.Process("break")
//
// For a more complete usage, see the test file.
package fsm

import "fmt"

const (
  OnEntry = "ON_ENTRY"
  OnExit  = "ON_EXIT"
  Default = "DEFAULT"
)

type Rule struct {
  From   string
  Event  string
  To     string
  Action string
}

// 'action' corresponds to what's in a Rule
type Delegate interface {
  StateMachineCallback(action string, args []interface{})
}

type StateMachine struct {
  CurrentState string // it seems silly to alter this directly, but not necessarily erroneous

  rules        []Rule
  delegate     Delegate
}

type Error struct {
  BadEvent string
  InState string
}

func (e Error) Error() string {
  return fmt.Sprintf("state machine error: cannot find rule for event [%s] when in state [%s]\n", e.BadEvent, e.InState)
}

// Use this in conjunction with Rule literals, keeping
// in mind that To may be omitted for actions, and Action may
// always be omitted. See the overview above for an example.
func NewStateMachine(rules []Rule, delegate Delegate) StateMachine {
  return StateMachine{delegate: delegate, rules: rules, CurrentState: rules[0].From}
}

func (m *StateMachine) Process(event string, args ...interface{}) *Error {
  trans := m.findTransMatching(m.CurrentState, event)
  if trans == nil {
    trans = m.findTransMatching(m.CurrentState, Default)
  }

  if trans == nil {
    return &Error{
      BadEvent: event,
      InState: m.CurrentState,
    }
  }

  changing_states := trans.From != trans.To

  if changing_states {
    m.runAction(m.CurrentState, OnExit, args)
  }

  if trans.Action != "" {
    m.delegate.StateMachineCallback(trans.Action, args)
  }

  if changing_states {
    m.runAction(trans.To, OnEntry, args)
  }

  m.CurrentState = trans.To

  return nil
}

func (m *StateMachine) findTransMatching(fromState string, event string) *Rule {
  for _, v := range m.rules {
    if v.From == fromState && v.Event == event {
      return &v
    }
  }
  return nil
}

func (m *StateMachine) runAction(state string, event string, args []interface{}) {
  if trans := m.findTransMatching(state, event); trans != nil && trans.Action != "" {
    m.delegate.StateMachineCallback(trans.Action, args)
  }
}
