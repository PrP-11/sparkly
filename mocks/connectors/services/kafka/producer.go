// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// ProducerService is an autogenerated mock type for the ProducerService type
type ProducerService struct {
	mock.Mock
}

type ProducerService_Expecter struct {
	mock *mock.Mock
}

func (_m *ProducerService) EXPECT() *ProducerService_Expecter {
	return &ProducerService_Expecter{mock: &_m.Mock}
}

// PushMessage provides a mock function with given fields: ctx, topic, body
func (_m *ProducerService) PushMessage(ctx context.Context, topic string, body interface{}) {
	_m.Called(ctx, topic, body)
}

// ProducerService_PushMessage_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'PushMessage'
type ProducerService_PushMessage_Call struct {
	*mock.Call
}

// PushMessage is a helper method to define mock.On call
//   - ctx context.Context
//   - topic string
//   - body interface{}
func (_e *ProducerService_Expecter) PushMessage(ctx interface{}, topic interface{}, body interface{}) *ProducerService_PushMessage_Call {
	return &ProducerService_PushMessage_Call{Call: _e.mock.On("PushMessage", ctx, topic, body)}
}

func (_c *ProducerService_PushMessage_Call) Run(run func(ctx context.Context, topic string, body interface{})) *ProducerService_PushMessage_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(interface{}))
	})
	return _c
}

func (_c *ProducerService_PushMessage_Call) Return() *ProducerService_PushMessage_Call {
	_c.Call.Return()
	return _c
}

func (_c *ProducerService_PushMessage_Call) RunAndReturn(run func(context.Context, string, interface{})) *ProducerService_PushMessage_Call {
	_c.Call.Return(run)
	return _c
}

// NewProducerService creates a new instance of ProducerService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewProducerService(t interface {
	mock.TestingT
	Cleanup(func())
}) *ProducerService {
	mock := &ProducerService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}