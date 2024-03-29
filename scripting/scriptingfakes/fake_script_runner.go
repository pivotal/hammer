// Code generated by counterfeiter. DO NOT EDIT.
package scriptingfakes

import (
	"sync"

	"github.com/pivotal/hammer/scripting"
)

type FakeScriptRunner struct {
	RunScriptStub        func([]string, []string, bool) error
	runScriptMutex       sync.RWMutex
	runScriptArgsForCall []struct {
		arg1 []string
		arg2 []string
		arg3 bool
	}
	runScriptReturns struct {
		result1 error
	}
	runScriptReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeScriptRunner) RunScript(arg1 []string, arg2 []string, arg3 bool) error {
	var arg1Copy []string
	if arg1 != nil {
		arg1Copy = make([]string, len(arg1))
		copy(arg1Copy, arg1)
	}
	var arg2Copy []string
	if arg2 != nil {
		arg2Copy = make([]string, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.runScriptMutex.Lock()
	ret, specificReturn := fake.runScriptReturnsOnCall[len(fake.runScriptArgsForCall)]
	fake.runScriptArgsForCall = append(fake.runScriptArgsForCall, struct {
		arg1 []string
		arg2 []string
		arg3 bool
	}{arg1Copy, arg2Copy, arg3})
	stub := fake.RunScriptStub
	fakeReturns := fake.runScriptReturns
	fake.recordInvocation("RunScript", []interface{}{arg1Copy, arg2Copy, arg3})
	fake.runScriptMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakeScriptRunner) RunScriptCallCount() int {
	fake.runScriptMutex.RLock()
	defer fake.runScriptMutex.RUnlock()
	return len(fake.runScriptArgsForCall)
}

func (fake *FakeScriptRunner) RunScriptCalls(stub func([]string, []string, bool) error) {
	fake.runScriptMutex.Lock()
	defer fake.runScriptMutex.Unlock()
	fake.RunScriptStub = stub
}

func (fake *FakeScriptRunner) RunScriptArgsForCall(i int) ([]string, []string, bool) {
	fake.runScriptMutex.RLock()
	defer fake.runScriptMutex.RUnlock()
	argsForCall := fake.runScriptArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeScriptRunner) RunScriptReturns(result1 error) {
	fake.runScriptMutex.Lock()
	defer fake.runScriptMutex.Unlock()
	fake.RunScriptStub = nil
	fake.runScriptReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeScriptRunner) RunScriptReturnsOnCall(i int, result1 error) {
	fake.runScriptMutex.Lock()
	defer fake.runScriptMutex.Unlock()
	fake.RunScriptStub = nil
	if fake.runScriptReturnsOnCall == nil {
		fake.runScriptReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.runScriptReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeScriptRunner) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.runScriptMutex.RLock()
	defer fake.runScriptMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeScriptRunner) recordInvocation(key string, args []interface{}) {
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

var _ scripting.ScriptRunner = new(FakeScriptRunner)
