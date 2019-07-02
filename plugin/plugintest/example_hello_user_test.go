package plugintest_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/plugin"
	"github.com/xzl8028/xenia-server/plugin/plugintest"
)

type HelloUserPlugin struct {
	plugin.XeniaPlugin
}

func (p *HelloUserPlugin) ServeHTTP(context *plugin.Context, w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("Xenia-User-Id")
	user, err := p.API.GetUser(userId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.API.LogError(err.Error())
		return
	}

	fmt.Fprintf(w, "Welcome back, %s!", user.Username)
}

func Example() {
	t := &testing.T{}
	user := &model.User{
		Id:       model.NewId(),
		Username: "billybob",
	}

	api := &plugintest.API{}
	api.On("GetUser", user.Id).Return(user, nil)
	defer api.AssertExpectations(t)

	helpers := &plugintest.Helpers{}
	defer helpers.AssertExpectations(t)

	p := &HelloUserPlugin{}
	p.SetAPI(api)
	p.SetHelpers(helpers)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	r.Header.Add("Xenia-User-Id", user.Id)
	p.ServeHTTP(&plugin.Context{}, w, r)
	body, err := ioutil.ReadAll(w.Result().Body)
	require.NoError(t, err)
	assert.Equal(t, "Welcome back, billybob!", string(body))
}
