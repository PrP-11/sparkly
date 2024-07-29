// Code generated by mockery v2.43.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	pkg "prp.com/sparkly/internal/pkg"

	primitive "go.mongodb.org/mongo-driver/bson/primitive"

	time "time"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

type Repository_Expecter struct {
	mock *mock.Mock
}

func (_m *Repository) EXPECT() *Repository_Expecter {
	return &Repository_Expecter{mock: &_m.Mock}
}

// GetPopularPosts provides a mock function with given fields: ctx, duration, limit
func (_m *Repository) GetPopularPosts(ctx context.Context, duration time.Duration, limit int) ([]primitive.M, error) {
	ret := _m.Called(ctx, duration, limit)

	if len(ret) == 0 {
		panic("no return value specified for GetPopularPosts")
	}

	var r0 []primitive.M
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration, int) ([]primitive.M, error)); ok {
		return rf(ctx, duration, limit)
	}
	if rf, ok := ret.Get(0).(func(context.Context, time.Duration, int) []primitive.M); ok {
		r0 = rf(ctx, duration, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]primitive.M)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, time.Duration, int) error); ok {
		r1 = rf(ctx, duration, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Repository_GetPopularPosts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetPopularPosts'
type Repository_GetPopularPosts_Call struct {
	*mock.Call
}

// GetPopularPosts is a helper method to define mock.On call
//   - ctx context.Context
//   - duration time.Duration
//   - limit int
func (_e *Repository_Expecter) GetPopularPosts(ctx interface{}, duration interface{}, limit interface{}) *Repository_GetPopularPosts_Call {
	return &Repository_GetPopularPosts_Call{Call: _e.mock.On("GetPopularPosts", ctx, duration, limit)}
}

func (_c *Repository_GetPopularPosts_Call) Run(run func(ctx context.Context, duration time.Duration, limit int)) *Repository_GetPopularPosts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(time.Duration), args[2].(int))
	})
	return _c
}

func (_c *Repository_GetPopularPosts_Call) Return(_a0 []primitive.M, _a1 error) *Repository_GetPopularPosts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Repository_GetPopularPosts_Call) RunAndReturn(run func(context.Context, time.Duration, int) ([]primitive.M, error)) *Repository_GetPopularPosts_Call {
	_c.Call.Return(run)
	return _c
}

// Insert provides a mock function with given fields: ctx, activity
func (_m *Repository) Insert(ctx context.Context, activity pkg.PostInteraction) error {
	ret := _m.Called(ctx, activity)

	if len(ret) == 0 {
		panic("no return value specified for Insert")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, pkg.PostInteraction) error); ok {
		r0 = rf(ctx, activity)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Repository_Insert_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Insert'
type Repository_Insert_Call struct {
	*mock.Call
}

// Insert is a helper method to define mock.On call
//   - ctx context.Context
//   - activity pkg.PostInteraction
func (_e *Repository_Expecter) Insert(ctx interface{}, activity interface{}) *Repository_Insert_Call {
	return &Repository_Insert_Call{Call: _e.mock.On("Insert", ctx, activity)}
}

func (_c *Repository_Insert_Call) Run(run func(ctx context.Context, activity pkg.PostInteraction)) *Repository_Insert_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(pkg.PostInteraction))
	})
	return _c
}

func (_c *Repository_Insert_Call) Return(_a0 error) *Repository_Insert_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Repository_Insert_Call) RunAndReturn(run func(context.Context, pkg.PostInteraction) error) *Repository_Insert_Call {
	_c.Call.Return(run)
	return _c
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
