// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/xzl8028/xenia-server/model"

// OAuthStore is an autogenerated mock type for the OAuthStore type
type OAuthStore struct {
	mock.Mock
}

// DeleteApp provides a mock function with given fields: id
func (_m *OAuthStore) DeleteApp(id string) *model.AppError {
	ret := _m.Called(id)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// GetAccessData provides a mock function with given fields: token
func (_m *OAuthStore) GetAccessData(token string) (*model.AccessData, *model.AppError) {
	ret := _m.Called(token)

	var r0 *model.AccessData
	if rf, ok := ret.Get(0).(func(string) *model.AccessData); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessData)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAccessDataByRefreshToken provides a mock function with given fields: token
func (_m *OAuthStore) GetAccessDataByRefreshToken(token string) (*model.AccessData, *model.AppError) {
	ret := _m.Called(token)

	var r0 *model.AccessData
	if rf, ok := ret.Get(0).(func(string) *model.AccessData); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessData)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(token)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAccessDataByUserForApp provides a mock function with given fields: userId, clientId
func (_m *OAuthStore) GetAccessDataByUserForApp(userId string, clientId string) ([]*model.AccessData, *model.AppError) {
	ret := _m.Called(userId, clientId)

	var r0 []*model.AccessData
	if rf, ok := ret.Get(0).(func(string, string) []*model.AccessData); ok {
		r0 = rf(userId, clientId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.AccessData)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(userId, clientId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetApp provides a mock function with given fields: id
func (_m *OAuthStore) GetApp(id string) (*model.OAuthApp, *model.AppError) {
	ret := _m.Called(id)

	var r0 *model.OAuthApp
	if rf, ok := ret.Get(0).(func(string) *model.OAuthApp); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OAuthApp)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(id)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAppByUser provides a mock function with given fields: userId, offset, limit
func (_m *OAuthStore) GetAppByUser(userId string, offset int, limit int) ([]*model.OAuthApp, *model.AppError) {
	ret := _m.Called(userId, offset, limit)

	var r0 []*model.OAuthApp
	if rf, ok := ret.Get(0).(func(string, int, int) []*model.OAuthApp); ok {
		r0 = rf(userId, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.OAuthApp)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, int, int) *model.AppError); ok {
		r1 = rf(userId, offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetApps provides a mock function with given fields: offset, limit
func (_m *OAuthStore) GetApps(offset int, limit int) ([]*model.OAuthApp, *model.AppError) {
	ret := _m.Called(offset, limit)

	var r0 []*model.OAuthApp
	if rf, ok := ret.Get(0).(func(int, int) []*model.OAuthApp); ok {
		r0 = rf(offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.OAuthApp)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(int, int) *model.AppError); ok {
		r1 = rf(offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAuthData provides a mock function with given fields: code
func (_m *OAuthStore) GetAuthData(code string) (*model.AuthData, *model.AppError) {
	ret := _m.Called(code)

	var r0 *model.AuthData
	if rf, ok := ret.Get(0).(func(string) *model.AuthData); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AuthData)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string) *model.AppError); ok {
		r1 = rf(code)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAuthorizedApps provides a mock function with given fields: userId, offset, limit
func (_m *OAuthStore) GetAuthorizedApps(userId string, offset int, limit int) ([]*model.OAuthApp, *model.AppError) {
	ret := _m.Called(userId, offset, limit)

	var r0 []*model.OAuthApp
	if rf, ok := ret.Get(0).(func(string, int, int) []*model.OAuthApp); ok {
		r0 = rf(userId, offset, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.OAuthApp)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, int, int) *model.AppError); ok {
		r1 = rf(userId, offset, limit)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetPreviousAccessData provides a mock function with given fields: userId, clientId
func (_m *OAuthStore) GetPreviousAccessData(userId string, clientId string) (*model.AccessData, *model.AppError) {
	ret := _m.Called(userId, clientId)

	var r0 *model.AccessData
	if rf, ok := ret.Get(0).(func(string, string) *model.AccessData); ok {
		r0 = rf(userId, clientId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessData)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(userId, clientId)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// PermanentDeleteAuthDataByUser provides a mock function with given fields: userId
func (_m *OAuthStore) PermanentDeleteAuthDataByUser(userId string) *model.AppError {
	ret := _m.Called(userId)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// RemoveAccessData provides a mock function with given fields: token
func (_m *OAuthStore) RemoveAccessData(token string) *model.AppError {
	ret := _m.Called(token)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(token)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// RemoveAllAccessData provides a mock function with given fields:
func (_m *OAuthStore) RemoveAllAccessData() *model.AppError {
	ret := _m.Called()

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func() *model.AppError); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// RemoveAuthData provides a mock function with given fields: code
func (_m *OAuthStore) RemoveAuthData(code string) *model.AppError {
	ret := _m.Called(code)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(string) *model.AppError); ok {
		r0 = rf(code)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// SaveAccessData provides a mock function with given fields: accessData
func (_m *OAuthStore) SaveAccessData(accessData *model.AccessData) (*model.AccessData, *model.AppError) {
	ret := _m.Called(accessData)

	var r0 *model.AccessData
	if rf, ok := ret.Get(0).(func(*model.AccessData) *model.AccessData); ok {
		r0 = rf(accessData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessData)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.AccessData) *model.AppError); ok {
		r1 = rf(accessData)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// SaveApp provides a mock function with given fields: app
func (_m *OAuthStore) SaveApp(app *model.OAuthApp) (*model.OAuthApp, *model.AppError) {
	ret := _m.Called(app)

	var r0 *model.OAuthApp
	if rf, ok := ret.Get(0).(func(*model.OAuthApp) *model.OAuthApp); ok {
		r0 = rf(app)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OAuthApp)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.OAuthApp) *model.AppError); ok {
		r1 = rf(app)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// SaveAuthData provides a mock function with given fields: authData
func (_m *OAuthStore) SaveAuthData(authData *model.AuthData) (*model.AuthData, *model.AppError) {
	ret := _m.Called(authData)

	var r0 *model.AuthData
	if rf, ok := ret.Get(0).(func(*model.AuthData) *model.AuthData); ok {
		r0 = rf(authData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AuthData)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.AuthData) *model.AppError); ok {
		r1 = rf(authData)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdateAccessData provides a mock function with given fields: accessData
func (_m *OAuthStore) UpdateAccessData(accessData *model.AccessData) (*model.AccessData, *model.AppError) {
	ret := _m.Called(accessData)

	var r0 *model.AccessData
	if rf, ok := ret.Get(0).(func(*model.AccessData) *model.AccessData); ok {
		r0 = rf(accessData)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AccessData)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.AccessData) *model.AppError); ok {
		r1 = rf(accessData)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// UpdateApp provides a mock function with given fields: app
func (_m *OAuthStore) UpdateApp(app *model.OAuthApp) (*model.OAuthApp, *model.AppError) {
	ret := _m.Called(app)

	var r0 *model.OAuthApp
	if rf, ok := ret.Get(0).(func(*model.OAuthApp) *model.OAuthApp); ok {
		r0 = rf(app)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OAuthApp)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.OAuthApp) *model.AppError); ok {
		r1 = rf(app)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}
