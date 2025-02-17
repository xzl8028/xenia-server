// Code generated by mockery v1.0.0. DO NOT EDIT.

// Regenerate this file using `make store-mocks`.

package mocks

import mock "github.com/stretchr/testify/mock"
import model "github.com/xzl8028/xenia-server/model"

// ClusterDiscoveryStore is an autogenerated mock type for the ClusterDiscoveryStore type
type ClusterDiscoveryStore struct {
	mock.Mock
}

// Cleanup provides a mock function with given fields:
func (_m *ClusterDiscoveryStore) Cleanup() *model.AppError {
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

// Delete provides a mock function with given fields: discovery
func (_m *ClusterDiscoveryStore) Delete(discovery *model.ClusterDiscovery) (bool, *model.AppError) {
	ret := _m.Called(discovery)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*model.ClusterDiscovery) bool); ok {
		r0 = rf(discovery)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.ClusterDiscovery) *model.AppError); ok {
		r1 = rf(discovery)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// Exists provides a mock function with given fields: discovery
func (_m *ClusterDiscoveryStore) Exists(discovery *model.ClusterDiscovery) (bool, *model.AppError) {
	ret := _m.Called(discovery)

	var r0 bool
	if rf, ok := ret.Get(0).(func(*model.ClusterDiscovery) bool); ok {
		r0 = rf(discovery)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(*model.ClusterDiscovery) *model.AppError); ok {
		r1 = rf(discovery)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// GetAll provides a mock function with given fields: discoveryType, clusterName
func (_m *ClusterDiscoveryStore) GetAll(discoveryType string, clusterName string) ([]*model.ClusterDiscovery, *model.AppError) {
	ret := _m.Called(discoveryType, clusterName)

	var r0 []*model.ClusterDiscovery
	if rf, ok := ret.Get(0).(func(string, string) []*model.ClusterDiscovery); ok {
		r0 = rf(discoveryType, clusterName)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.ClusterDiscovery)
		}
	}

	var r1 *model.AppError
	if rf, ok := ret.Get(1).(func(string, string) *model.AppError); ok {
		r1 = rf(discoveryType, clusterName)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*model.AppError)
		}
	}

	return r0, r1
}

// Save provides a mock function with given fields: discovery
func (_m *ClusterDiscoveryStore) Save(discovery *model.ClusterDiscovery) *model.AppError {
	ret := _m.Called(discovery)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.ClusterDiscovery) *model.AppError); ok {
		r0 = rf(discovery)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}

// SetLastPingAt provides a mock function with given fields: discovery
func (_m *ClusterDiscoveryStore) SetLastPingAt(discovery *model.ClusterDiscovery) *model.AppError {
	ret := _m.Called(discovery)

	var r0 *model.AppError
	if rf, ok := ret.Get(0).(func(*model.ClusterDiscovery) *model.AppError); ok {
		r0 = rf(discovery)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.AppError)
		}
	}

	return r0
}
