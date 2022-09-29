package again_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

type inputCall struct {
	ctx           context.Context
	fakeOperation *FakeOperation
}

func (currentCall inputCall) Returns(value int, err error) *FakeOperation {
	expectedCalls := currentCall.fakeOperation.expectedCalls[currentCall.ctx]
	expectedCalls = append(expectedCalls, call{
		input: currentCall.ctx,
		value: value,
		err:   err,
	})
	currentCall.fakeOperation.expectedCalls[currentCall.ctx] = expectedCalls
	return currentCall.fakeOperation
}

type call struct {
	input context.Context
	value int
	err   error
}
type FakeOperation struct {
	t             *testing.T
	times         int
	called        []call
	expectedCalls map[context.Context][]call
	allowAnyCall  bool
}

func NewFakeOperation(t *testing.T) *FakeOperation {
	t.Helper()
	return &FakeOperation{
		t:             t,
		expectedCalls: make(map[context.Context][]call),
	}
}

func (currentFakeOperator *FakeOperation) Run(context context.Context) (int, error) {
	expectedCalls, ok := currentFakeOperator.expectedCalls[context]
	require.True(
		currentFakeOperator.t,
		ok || currentFakeOperator.allowAnyCall,
		"Unexpected call for FakeOperation",
	)
	expectedCall := call{
		input: context,
		value: 23,
		err:   errors.New("default err"),
	}
	if !currentFakeOperator.allowAnyCall {
		require.NotZero(currentFakeOperator.t, expectedCalls)
		call := expectedCalls[0]
		expectedCall = call
		currentFakeOperator.expectedCalls[context] = expectedCalls[1:]
	}
	currentFakeOperator.called = append(currentFakeOperator.called, expectedCall)
	currentFakeOperator.times += 1
	return expectedCall.value, expectedCall.err
}

func (currentFakeOperator *FakeOperation) givenContext(ctx context.Context) inputCall {
	require.NotNil(currentFakeOperator.t, ctx)
	return inputCall{
		ctx:           ctx,
		fakeOperation: currentFakeOperator,
	}
}

func (currentFakeOperator FakeOperation) haveBeenCalled(times int) {
	require.Equal(currentFakeOperator.t, times, currentFakeOperator.times)
}
