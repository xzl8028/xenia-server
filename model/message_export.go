// Copyright (c) 2017-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package model

type MessageExport struct {
	TeamId          *string
	TeamName        *string
	TeamDisplayName *string

	ChannelId          *string
	ChannelName        *string
	ChannelDisplayName *string
	ChannelType        *string

	UserId    *string
	UserEmail *string
	Username  *string
	IsBot     bool

	PostId         *string
	PostCreateAt   *int64
	PostMessage    *string
	PostType       *string
	PostRootId     *string
	PostProps      *string
	PostOriginalId *string
	PostFileIds    StringArray
}
