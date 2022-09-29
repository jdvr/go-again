package internal_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

type inputCall struct {
	ctx           context.Context
	fakeOperation *FakeOperation
}

func (currentCall inputCall) Returns(err error) *FakeOperation {
	currentCall.fakeOperation.expectedCalls[currentCall.ctx] = call{
		input:  currentCall.ctx,
		result: err,
	}
	return currentCall.fakeOperation
}

type call struct {
	input  context.Context
	result error
}
type FakeOperation struct {
	t             *testing.T
	times         int
	called        []call
	expectedCalls map[context.Context]call
}

func NewFakeOperation(t *testing.T) *FakeOperation {
	t.Helper()
	return &FakeOperation{
		t:             t,
		expectedCalls: make(map[context.Context]call),
	}
}

func (currentFakeOperator *FakeOperation) Run(context context.Context) error {
	expectedCall, ok := currentFakeOperator.expectedCalls[context]
	require.True(currentFakeOperator.t, ok, "Unexpected call for FakeOperation")
	currentFakeOperator.called = append(currentFakeOperator.called, expectedCall)
	currentFakeOperator.times += 1
	return expectedCall.result
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
