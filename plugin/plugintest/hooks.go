// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make plugin-mocks`.

package plugintest

import http "net/http"
import io "io"
import mock "github.com/stretchr/testify/mock"
import model "github.com/xzl8028/xenia-server/model"
import plugin "github.com/xzl8028/xenia-server/plugin"

// Hooks is an autogenerated mock type for the Hooks type
type Hooks struct {
	mock.Mock
}

// ChannelHasBeenCreated provides a mock function with given fields: c, channel
func (_m *Hooks) ChannelHasBeenCreated(c *plugin.Context, channel *model.Channel) {
	_m.Called(c, channel)
}

// ExecuteCommand provides a mock function with given fields: c, args
func (_m *Hooks) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	ret := _m.Called(c, args)

	var r0 *model.CommandResponse
	if rf, ok := ret.Get(0).(func(*plugin.Context, *model.CommandArgs) *model.CommandResponse); ok {
		r0 = rf(c, args)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.CommandResponse)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*plugin.Context, *model.CommandArgs) *model.AppError); ok {
		r1 = rf(c, args)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// FileWillBeUploaded provides a mock function with given fields: c, info, file, output
func (_m *Hooks) FileWillBeUploaded(c *plugin.Context, info *model.FileInfo, file io.Reader, output io.Writer) (*model.FileInfo, string) {
	ret := _m.Called(c, info, file, output)

	var r0 *model.FileInfo
	if rf, ok := ret.Get(0).(func(*plugin.Context, *model.FileInfo, io.Reader, io.Writer) *model.FileInfo); ok {
		r0 = rf(c, info, file, output)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.FileInfo)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(*plugin.Context, *model.FileInfo, io.Reader, io.Writer) string); ok {
		r1 = rf(c, info, file, output)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// Implemented provides a mock function with given fields:
func (_m *Hooks) Implemented() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MessageHasBeenPosted provides a mock function with given fields: c, post
func (_m *Hooks) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	_m.Called(c, post)
}

// MessageHasBeenUpdated provides a mock function with given fields: c, newPost, oldPost
func (_m *Hooks) MessageHasBeenUpdated(c *plugin.Context, newPost *model.Post, oldPost *model.Post) {
	_m.Called(c, newPost, oldPost)
}

// MessageWillBePosted provides a mock function with given fields: c, post
func (_m *Hooks) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	ret := _m.Called(c, post)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(*plugin.Context, *model.Post) *model.Post); ok {
		r0 = rf(c, post)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(*plugin.Context, *model.Post) string); ok {
		r1 = rf(c, post)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// MessageWillBeUpdated provides a mock function with given fields: c, newPost, oldPost
func (_m *Hooks) MessageWillBeUpdated(c *plugin.Context, newPost *model.Post, oldPost *model.Post) (*model.Post, string) {
	ret := _m.Called(c, newPost, oldPost)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(*plugin.Context, *model.Post, *model.Post) *model.Post); ok {
		r0 = rf(c, newPost, oldPost)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(*plugin.Context, *model.Post, *model.Post) string); ok {
		r1 = rf(c, newPost, oldPost)
	} else {
		r1 = ret.Get(1).(string)
	}

	return r0, r1
}

// OnActivate provides a mock function with given fields:
func (_m *Hooks) OnActivate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnConfigurationChange provides a mock function with given fields:
func (_m *Hooks) OnConfigurationChange() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// OnDeactivate provides a mock function with given fields:
func (_m *Hooks) OnDeactivate() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ServeHTTP provides a mock function with given fields: c, w, r
func (_m *Hooks) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	_m.Called(c, w, r)
}

// UserHasBeenCreated provides a mock function with given fields: c, user
func (_m *Hooks) UserHasBeenCreated(c *plugin.Context, user *model.User) {
	_m.Called(c, user)
}

// UserHasJoinedChannel provides a mock function with given fields: c, channelMember, actor
func (_m *Hooks) UserHasJoinedChannel(c *plugin.Context, channelMember *model.ChannelMember, actor *model.User) {
	_m.Called(c, channelMember, actor)
}

// UserHasJoinedTeam provides a mock function with given fields: c, teamMember, actor
func (_m *Hooks) UserHasJoinedTeam(c *plugin.Context, teamMember *model.TeamMember, actor *model.User) {
	_m.Called(c, teamMember, actor)
}

// UserHasLeftChannel provides a mock function with given fields: c, channelMember, actor
func (_m *Hooks) UserHasLeftChannel(c *plugin.Context, channelMember *model.ChannelMember, actor *model.User) {
	_m.Called(c, channelMember, actor)
}

// UserHasLeftTeam provides a mock function with given fields: c, teamMember, actor
func (_m *Hooks) UserHasLeftTeam(c *plugin.Context, teamMember *model.TeamMember, actor *model.User) {
	_m.Called(c, teamMember, actor)
}

// UserHasLoggedIn provides a mock function with given fields: c, user
func (_m *Hooks) UserHasLoggedIn(c *plugin.Context, user *model.User) {
	_m.Called(c, user)
}

// UserWillLogIn provides a mock function with given fields: c, user
func (_m *Hooks) UserWillLogIn(c *plugin.Context, user *model.User) string {
	ret := _m.Called(c, user)

	var r0 string
	if rf, ok := ret.Get(0).(func(*plugin.Context, *model.User) string); ok {
		r0 = rf(c, user)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
