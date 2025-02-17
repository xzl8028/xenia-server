// Copyright (c) 2016-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"strings"

	"github.com/xzl8028/xenia-server/utils/jsonutils"
)

const (
	COMMAND_RESPONSE_TYPE_IN_CHANNEL = "in_channel"
	COMMAND_RESPONSE_TYPE_EPHEMERAL  = "ephemeral"
)

type CommandResponse struct {
	ResponseType   string             `json:"response_type"`
	Text           string             `json:"text"`
	Username       string             `json:"username"`
	ChannelId      string             `json:"channel_id"`
	IconURL        string             `json:"icon_url"`
	Type           string             `json:"type"`
	Props          StringInterface    `json:"props"`
	GotoLocation   string             `json:"goto_location"`
	TriggerId      string             `json:"trigger_id"`
	Attachments    []*SlackAttachment `json:"attachments"`
	ExtraResponses []*CommandResponse `json:"extra_responses"`
}

func (o *CommandResponse) ToJson() string {
	b, _ := json.Marshal(o)
	return string(b)
}

func CommandResponseFromHTTPBody(contentType string, body io.Reader) (*CommandResponse, error) {
	if strings.TrimSpace(strings.Split(contentType, ";")[0]) == "application/json" {
		return CommandResponseFromJson(body)
	}
	if b, err := ioutil.ReadAll(body); err == nil {
		return CommandResponseFromPlainText(string(b)), nil
	}
	return nil, nil
}

func CommandResponseFromPlainText(text string) *CommandResponse {
	return &CommandResponse{
		Text: text,
	}
}

func CommandResponseFromJson(data io.Reader) (*CommandResponse, error) {
	b, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	var o CommandResponse
	err = json.Unmarshal(b, &o)
	if err != nil {
		return nil, jsonutils.HumanizeJsonError(err, b)
	}

	o.Attachments = StringifySlackFieldValue(o.Attachments)

	if o.ExtraResponses != nil {
		for _, resp := range o.ExtraResponses {
			resp.Attachments = StringifySlackFieldValue(resp.Attachments)
		}
	}

	return &o, nil
}
