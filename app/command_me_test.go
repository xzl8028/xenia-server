package app

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/xzl8028/xenia-server/model"
)

func TestMeProviderDoCommand(t *testing.T) {
	th := Setup(t).InitBasic()
	defer th.TearDown()

	mp := MeProvider{}

	msg := "hello"

	resp := mp.DoCommand(th.App, &model.CommandArgs{}, msg)

	assert.Equal(t, model.COMMAND_RESPONSE_TYPE_IN_CHANNEL, resp.ResponseType)
	assert.Equal(t, model.POST_ME, resp.Type)
	assert.Equal(t, "*"+msg+"*", resp.Text)
	assert.Equal(t, model.StringInterface{
		"message": msg,
	}, resp.Props)
}
