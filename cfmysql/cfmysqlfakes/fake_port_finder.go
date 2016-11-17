// This file was generated by counterfeiter
package cfmysqlfakes

import (
	"sync"

	"github.com/andreasf/cf-mysql-plugin/cfmysql"
)

type FakePortFinder struct {
	GetPortStub        func() int
	getPortMutex       sync.RWMutex
	getPortArgsForCall []struct{}
	getPortReturns     struct {
		result1 int
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePortFinder) GetPort() int {
	fake.getPortMutex.Lock()
	fake.getPortArgsForCall = append(fake.getPortArgsForCall, struct{}{})
	fake.recordInvocation("GetPort", []interface{}{})
	fake.getPortMutex.Unlock()
	if fake.GetPortStub != nil {
		return fake.GetPortStub()
	} else {
		return fake.getPortReturns.result1
	}
}

func (fake *FakePortFinder) GetPortCallCount() int {
	fake.getPortMutex.RLock()
	defer fake.getPortMutex.RUnlock()
	return len(fake.getPortArgsForCall)
}

func (fake *FakePortFinder) GetPortReturns(result1 int) {
	fake.GetPortStub = nil
	fake.getPortReturns = struct {
		result1 int
	}{result1}
}

func (fake *FakePortFinder) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getPortMutex.RLock()
	defer fake.getPortMutex.RUnlock()
	return fake.invocations
}

func (fake *FakePortFinder) recordInvocation(key string, args []interface{}) {
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

var _ cfmysql.PortFinder = new(FakePortFinder)