// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package model

const (
	POST_EMBED_IMAGE              PostEmbedType = "image"
	POST_EMBED_MESSAGE_ATTACHMENT PostEmbedType = "message_attachment"
	POST_EMBED_OPENGRAPH          PostEmbedType = "opengraph"
)

type PostEmbedType string

type PostEmbed struct {
	Type PostEmbedType `json:"type"`

	// The URL of the embedded content. Used for image and OpenGraph embeds.
	URL string `json:"url,omitempty"`

	// Any additional data for the embedded content. Only used for OpenGraph embeds.
	Data interface{} `json:"data,omitempty"`
}
