// Code generated by counterfeiter. DO NOT EDIT.
package commandsfakes

import (
	"sync"

	"github.com/pivotal/hammer/commands"
)

type FakeUI struct {
	DisplayErrorStub        func(error)
	displayErrorMutex       sync.RWMutex
	displayErrorArgsForCall []struct {
		arg1 error
	}
	DisplayTextStub        func(string)
	displayTextMutex       sync.RWMutex
	displayTextArgsForCall []struct {
		arg1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeUI) DisplayError(arg1 error) {
	fake.displayErrorMutex.Lock()
	fake.displayErrorArgsForCall = append(fake.displayErrorArgsForCall, struct {
		arg1 error
	}{arg1})
	stub := fake.DisplayErrorStub
	fake.recordInvocation("DisplayError", []interface{}{arg1})
	fake.displayErrorMutex.Unlock()
	if stub != nil {
		fake.DisplayErrorStub(arg1)
	}
}

func (fake *FakeUI) DisplayErrorCallCount() int {
	fake.displayErrorMutex.RLock()
	defer fake.displayErrorMutex.RUnlock()
	return len(fake.displayErrorArgsForCall)
}

func (fake *FakeUI) DisplayErrorCalls(stub func(error)) {
	fake.displayErrorMutex.Lock()
	defer fake.displayErrorMutex.Unlock()
	fake.DisplayErrorStub = stub
}

func (fake *FakeUI) DisplayErrorArgsForCall(i int) error {
	fake.displayErrorMutex.RLock()
	defer fake.displayErrorMutex.RUnlock()
	argsForCall := fake.displayErrorArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeUI) DisplayText(arg1 string) {
	fake.displayTextMutex.Lock()
	fake.displayTextArgsForCall = append(fake.displayTextArgsForCall, struct {
		arg1 string
	}{arg1})
	stub := fake.DisplayTextStub
	fake.recordInvocation("DisplayText", []interface{}{arg1})
	fake.displayTextMutex.Unlock()
	if stub != nil {
		fake.DisplayTextStub(arg1)
	}
}

func (fake *FakeUI) DisplayTextCallCount() int {
	fake.displayTextMutex.RLock()
	defer fake.displayTextMutex.RUnlock()
	return len(fake.displayTextArgsForCall)
}

func (fake *FakeUI) DisplayTextCalls(stub func(string)) {
	fake.displayTextMutex.Lock()
	defer fake.displayTextMutex.Unlock()
	fake.DisplayTextStub = stub
}

func (fake *FakeUI) DisplayTextArgsForCall(i int) string {
	fake.displayTextMutex.RLock()
	defer fake.displayTextMutex.RUnlock()
	argsForCall := fake.displayTextArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeUI) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.displayErrorMutex.RLock()
	defer fake.displayErrorMutex.RUnlock()
	fake.displayTextMutex.RLock()
	defer fake.displayTextMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeUI) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ commands.UI = new(FakeUI)
