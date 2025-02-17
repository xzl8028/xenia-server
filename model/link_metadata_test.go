// Copyright (c) 2015-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package model

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/dyatlov/go-opengraph/opengraph"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const BigText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Vivamus maximus faucibus ex, vitae placerat neque feugiat ac. Nam tempus libero quis pellentesque feugiat. Cras tristique diam vel condimentum viverra. Proin molestie posuere leo. Nam pulvinar, ex quis tristique cursus, turpis ante commodo elit, a dapibus est ipsum id eros. Mauris tortor dolor, posuere ac velit vitae, faucibus viverra fusce."

func sampleImage(imageName string) *opengraph.Image {
	return &opengraph.Image{
		URL:       fmt.Sprintf("http://example.com/%s", imageName),
		SecureURL: fmt.Sprintf("https://example.com/%s", imageName),
		Type:      "png",
		Width:     32,
		Height:    32,
	}
}

func TestLinkMetadataIsValid(t *testing.T) {
	for _, test := range []struct {
		Name     string
		Metadata *LinkMetadata
		Expected bool
	}{
		{
			Name: "should be valid image metadata",
			Metadata: &LinkMetadata{
				URL:       "http://example.com",
				Timestamp: 1546300800000,
				Type:      LINK_METADATA_TYPE_IMAGE,
				Data:      &PostImage{},
			},
			Expected: true,
		},
		{
			Name: "should be valid opengraph metadata",
			Metadata: &LinkMetadata{
				URL:       "http://example.com",
				Timestamp: 1546300800000,
				Type:      LINK_METADATA_TYPE_OPENGRAPH,
				Data:      &opengraph.OpenGraph{},
			},
			Expected: true,
		},
		{
			Name: "should be valid with no metadata",
			Metadata: &LinkMetadata{
				URL:       "http://example.com",
				Timestamp: 1546300800000,
				Type:      LINK_METADATA_TYPE_NONE,
				Data:      nil,
			},
			Expected: true,
		},
		{
			Name: "should be invalid because of empty URL",
			Metadata: &LinkMetadata{
				Timestamp: 1546300800000,
				Type:      LINK_METADATA_TYPE_IMAGE,
				Data:      &PostImage{},
			},
			Expected: false,
		},
		{
			Name: "should be invalid because of empty timestamp",
			Metadata: &LinkMetadata{
				URL:  "http://example.com",
				Type: LINK_METADATA_TYPE_IMAGE,
				Data: &PostImage{},
			},
			Expected: false,
		},
		{
			Name: "should be invalid because of unrounded timestamp",
			Metadata: &LinkMetadata{
				URL:       "http://example.com",
				Timestamp: 1546300800001,
				Type:      LINK_METADATA_TYPE_IMAGE,
				Data:      &PostImage{},
			},
			Expected: false,
		},
		{
			Name: "should be invalid because of invalid type",
			Metadata: &LinkMetadata{
				URL:       "http://example.com",
				Timestamp: 1546300800000,
				Type:      "garbage",
				Data:      &PostImage{},
			},
			Expected: false,
		},
		{
			Name: "should be invalid because of empty data",
			Metadata: &LinkMetadata{
				URL:       "http://example.com",
				Timestamp: 1546300800000,
				Type:      LINK_METADATA_TYPE_IMAGE,
			},
			Expected: false,
		},
		{
			Name: "should be invalid because of mismatched data and type, image type and opengraph data",
			Metadata: &LinkMetadata{
				URL:       "http://example.com",
				Timestamp: 1546300800000,
				Type:      LINK_METADATA_TYPE_IMAGE,
				Data:      &opengraph.OpenGraph{},
			},
			Expected: false,
		},
		{
			Name: "should be invalid because of mismatched data and type, opengraph type and image data",
			Metadata: &LinkMetadata{
				URL:       "http://example.com",
				Timestamp: 1546300800000,
				Type:      LINK_METADATA_TYPE_OPENGRAPH,
				Data:      &PostImage{},
			},
			Expected: false,
		},
		{
			Name: "should be invalid because of mismatched data and type, image type and random data",
			Metadata: &LinkMetadata{
				URL:       "http://example.com",
				Timestamp: 1546300800000,
				Type:      LINK_METADATA_TYPE_OPENGRAPH,
				Data:      &Channel{},
			},
			Expected: false,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			err := test.Metadata.IsValid()

			if test.Expected {
				assert.Nil(t, err)
			} else {
				assert.NotNil(t, err)
			}
		})
	}
}

func TestLinkMetadataDeserializeDataToConcreteType(t *testing.T) {
	t.Run("should convert []byte to PostImage", func(t *testing.T) {
		image := &PostImage{
			Height: 400,
			Width:  500,
		}

		metadata := &LinkMetadata{
			Type: LINK_METADATA_TYPE_IMAGE,
			Data: []byte(image.ToJson()),
		}

		require.IsType(t, []byte{}, metadata.Data)

		err := metadata.DeserializeDataToConcreteType()

		assert.Nil(t, err)
		assert.IsType(t, &PostImage{}, metadata.Data)
		assert.Equal(t, *image, *metadata.Data.(*PostImage))
	})

	t.Run("should convert string to OpenGraph", func(t *testing.T) {
		og := &opengraph.OpenGraph{
			URL:         "http://example.com",
			Description: "Hello, world!",
			Images: []*opengraph.Image{
				{
					URL: "http://example.com/image.png",
				},
			},
		}

		b, err := json.Marshal(og)

		require.Nil(t, err)

		metadata := &LinkMetadata{
			Type: LINK_METADATA_TYPE_OPENGRAPH,
			Data: b,
		}

		require.IsType(t, []byte{}, metadata.Data)

		err = metadata.DeserializeDataToConcreteType()

		assert.Nil(t, err)
		assert.IsType(t, &opengraph.OpenGraph{}, metadata.Data)
		assert.Equal(t, *og, *metadata.Data.(*opengraph.OpenGraph))
	})

	t.Run("should ignore data of the correct type", func(t *testing.T) {
		metadata := &LinkMetadata{
			Type: LINK_METADATA_TYPE_OPENGRAPH,
			Data: 1234,
		}

		err := metadata.DeserializeDataToConcreteType()

		assert.Nil(t, err)
	})

	t.Run("should ignore an invalid type", func(t *testing.T) {
		metadata := &LinkMetadata{
			Type: "garbage",
			Data: "garbage",
		}

		err := metadata.DeserializeDataToConcreteType()

		assert.Nil(t, err)
	})

	t.Run("should return error for invalid data", func(t *testing.T) {
		metadata := &LinkMetadata{
			Type: LINK_METADATA_TYPE_IMAGE,
			Data: "garbage",
		}

		err := metadata.DeserializeDataToConcreteType()

		assert.NotNil(t, err)
	})
}

func TestFloorToNearestHour(t *testing.T) {
	assert.True(t, isRoundedToNearestHour(FloorToNearestHour(1546346096000)))
}

func TestTruncateText(t *testing.T) {
	t.Run("Shouldn't affect strings smaller than 300 characters", func(t *testing.T) {
		assert.Equal(t, utf8.RuneCountInString(truncateText("abc")), 3, "should be 3")
	})
	t.Run("Shouldn't affect empty strings", func(t *testing.T) {
		assert.Equal(t, utf8.RuneCountInString(truncateText("")), 0, "should be empty")
	})
	t.Run("Truncates string to 300 + 5", func(t *testing.T) {
		assert.Equal(t, utf8.RuneCountInString(truncateText(BigText)), 305, "should be 300 chars + 5")
	})
	t.Run("Truncated text ends in elipsis", func(t *testing.T) {
		assert.True(t, strings.HasSuffix(truncateText(BigText), "[...]"))
	})
}

func TestFirstNImages(t *testing.T) {
	t.Run("when empty, return an empty one", func(t *testing.T) {
		empty := make([]*opengraph.Image, 0)
		assert.Exactly(t, firstNImages(empty, 1), empty, "Should be the same element")
	})
	t.Run("when it contains one element, return the same array", func(t *testing.T) {
		one := []*opengraph.Image{sampleImage("image.png")}
		assert.Exactly(t, firstNImages(one, 1), one, "Should be the same element")
	})
	t.Run("when it contains more than one element and asking for only one, return the first one", func(t *testing.T) {
		two := []*opengraph.Image{sampleImage("image.png"), sampleImage("notme.png")}
		assert.True(t, strings.HasSuffix(firstNImages(two, 1)[0].URL, "image.png"), "Should be the image element")
	})
	t.Run("when it contains less than asked, return the original", func(t *testing.T) {
		two := []*opengraph.Image{sampleImage("image.png"), sampleImage("notme.png")}
		assert.Equal(t, two, firstNImages(two, 10), "should be the same pointer")
	})

	t.Run("asking for negative images", func(t *testing.T) {
		six := []*opengraph.Image{
			sampleImage("image.png"),
			sampleImage("another.png"),
			sampleImage("yetanother.jpg"),
			sampleImage("metoo.gif"),
			sampleImage("fifth.ico"),
			sampleImage("notme.tiff"),
		}
		assert.Len(t, firstNImages(six, -10), MAX_IMAGES, "On negative, go for defaults")
	})

}

func TestTruncateOpenGraph(t *testing.T) {
	og := opengraph.OpenGraph{
		Type:             "something",
		URL:              "http://myawesomesite.com",
		Title:            BigText,
		Description:      BigText,
		Determiner:       BigText,
		SiteName:         BigText,
		Locale:           "[EN-en]",
		LocalesAlternate: []string{"[EN-ca]", "[ES-es]"},
		Images: []*opengraph.Image{
			sampleImage("image.png"),
			sampleImage("another.png"),
			sampleImage("yetanother.jpg"),
			sampleImage("metoo.gif"),
			sampleImage("fifth.ico"),
			sampleImage("notme.tiff")},
		Audios:  []*opengraph.Audio{&opengraph.Audio{}},
		Videos:  []*opengraph.Video{&opengraph.Video{}},
		Article: &opengraph.Article{},
		Book:    &opengraph.Book{},
		Profile: &opengraph.Profile{},
	}
	result := TruncateOpenGraph(&og)
	assert.Nil(t, result.Article, "No article stored")
	assert.Nil(t, result.Book, "No book stored")
	assert.Nil(t, result.Profile, "No profile stored")
	assert.Len(t, result.Images, 5, "Only the first 5 images")
	assert.Len(t, result.Audios, 0, "No audios stored")
	assert.Len(t, result.Videos, 0, "No videos stored")
	assert.Len(t, result.LocalesAlternate, 0, "No alternate locales stored")
	assert.Equal(t, result.Determiner, "", "No determiner stored")
	assert.Equal(t, utf8.RuneCountInString(result.Title), 305, "Title text is truncated")
	assert.Equal(t, utf8.RuneCountInString(result.Description), 305, "Description text is truncated")
	assert.Equal(t, utf8.RuneCountInString(result.SiteName), 305, "SiteName text is truncated")
}
