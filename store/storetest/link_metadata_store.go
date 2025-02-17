// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package storetest

import (
	"net/http"
	"testing"
	"time"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/xzl8028/xenia-server/model"
	"github.com/xzl8028/xenia-server/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests are ran on the same store instance, so this provides easier unique, valid timestamps
var linkMetadataTimestamp int64 = 1546300800000

func getNextLinkMetadataTimestamp() int64 {
	linkMetadataTimestamp += int64(time.Hour) / 1000
	return linkMetadataTimestamp
}

func TestLinkMetadataStore(t *testing.T, ss store.Store) {
	t.Run("Save", func(t *testing.T) { testLinkMetadataStoreSave(t, ss) })
	t.Run("Get", func(t *testing.T) { testLinkMetadataStoreGet(t, ss) })
	t.Run("Types", func(t *testing.T) { testLinkMetadataStoreTypes(t, ss) })
}

func testLinkMetadataStoreSave(t *testing.T, ss store.Store) {
	t.Run("should save item", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_IMAGE,
			Data:      &model.PostImage{},
		}

		linkMetadata, err := ss.LinkMetadata().Save(metadata)

		require.Nil(t, err)
		assert.Equal(t, *metadata, *linkMetadata)
	})

	t.Run("should fail to save invalid item", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "",
			Timestamp: 0,
			Type:      "garbage",
			Data:      nil,
		}

		_, err := ss.LinkMetadata().Save(metadata)

		assert.NotNil(t, err)
	})

	t.Run("should save with duplicate URL and different timestamp", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_IMAGE,
			Data:      &model.PostImage{},
		}

		_, err := ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)

		metadata.Timestamp = getNextLinkMetadataTimestamp()

		linkMetadata, err := ss.LinkMetadata().Save(metadata)

		require.Nil(t, err)
		assert.Equal(t, *metadata, *linkMetadata)
	})

	t.Run("should save with duplicate timestamp and different URL", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_IMAGE,
			Data:      &model.PostImage{},
		}

		_, err := ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)

		metadata.URL = "http://example.com/another/page"

		linkMetadata, err := ss.LinkMetadata().Save(metadata)

		require.Nil(t, err)
		assert.Equal(t, *metadata, *linkMetadata)
	})

	t.Run("should not save with duplicate URL and timestamp, but should not return an error", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_IMAGE,
			Data:      &model.PostImage{},
		}

		linkMetadata, err := ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)
		assert.Equal(t, &model.PostImage{}, linkMetadata.Data)

		metadata.Data = &model.PostImage{Height: 10, Width: 20}

		linkMetadata, err = ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)
		assert.Equal(t, linkMetadata.Data, &model.PostImage{Height: 10, Width: 20})

		// Should return the original result, not the duplicate one
		linkMetadata, err = ss.LinkMetadata().Get(metadata.URL, metadata.Timestamp)
		require.Nil(t, err)
		assert.Equal(t, &model.PostImage{}, linkMetadata.Data)
	})
}

func testLinkMetadataStoreGet(t *testing.T, ss store.Store) {
	t.Run("should get value", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_IMAGE,
			Data:      &model.PostImage{},
		}

		_, err := ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)

		linkMetadata, err := ss.LinkMetadata().Get(metadata.URL, metadata.Timestamp)

		require.Nil(t, err)
		require.IsType(t, metadata, linkMetadata)
		assert.Equal(t, *metadata, *linkMetadata)
	})

	t.Run("should return not found with incorrect URL", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_IMAGE,
			Data:      &model.PostImage{},
		}

		_, err := ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)

		_, err = ss.LinkMetadata().Get("http://example.com/another_page", metadata.Timestamp)

		require.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
	})

	t.Run("should return not found with incorrect timestamp", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_IMAGE,
			Data:      &model.PostImage{},
		}

		_, err := ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)

		_, err = ss.LinkMetadata().Get(metadata.URL, getNextLinkMetadataTimestamp())

		require.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.StatusCode)
	})
}

func testLinkMetadataStoreTypes(t *testing.T, ss store.Store) {
	t.Run("should save and get image metadata", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_IMAGE,
			Data: &model.PostImage{
				Width:  123,
				Height: 456,
			},
		}

		received, err := ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)

		require.IsType(t, &model.PostImage{}, received.Data)
		assert.Equal(t, *(metadata.Data.(*model.PostImage)), *(received.Data.(*model.PostImage)))

		received, err = ss.LinkMetadata().Get(metadata.URL, metadata.Timestamp)
		require.Nil(t, err)

		require.IsType(t, &model.PostImage{}, received.Data)
		assert.Equal(t, *(metadata.Data.(*model.PostImage)), *(received.Data.(*model.PostImage)))
	})

	t.Run("should save and get opengraph data", func(t *testing.T) {
		og := &opengraph.OpenGraph{
			URL: "http://example.com",
			Images: []*opengraph.Image{
				{
					URL: "http://example.com/image.png",
				},
			},
		}

		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_OPENGRAPH,
			Data:      og,
		}

		received, err := ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)

		require.IsType(t, &opengraph.OpenGraph{}, received.Data)
		assert.Equal(t, *(metadata.Data.(*opengraph.OpenGraph)), *(received.Data.(*opengraph.OpenGraph)))

		received, err = ss.LinkMetadata().Get(metadata.URL, metadata.Timestamp)
		require.Nil(t, err)

		require.IsType(t, &opengraph.OpenGraph{}, received.Data)
		assert.Equal(t, *(metadata.Data.(*opengraph.OpenGraph)), *(received.Data.(*opengraph.OpenGraph)))
	})

	t.Run("should save and get nil", func(t *testing.T) {
		metadata := &model.LinkMetadata{
			URL:       "http://example.com",
			Timestamp: getNextLinkMetadataTimestamp(),
			Type:      model.LINK_METADATA_TYPE_NONE,
			Data:      nil,
		}

		received, err := ss.LinkMetadata().Save(metadata)
		require.Nil(t, err)
		assert.Nil(t, received.Data)

		received, err = ss.LinkMetadata().Get(metadata.URL, metadata.Timestamp)
		require.Nil(t, err)

		require.Nil(t, received.Data)
	})
}
