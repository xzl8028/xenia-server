package app

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCondenseSiteURL(t *testing.T) {
	require.Equal(t, "", condenseSiteURL(""))
	require.Equal(t, "xenia.com", condenseSiteURL("xenia.com"))
	require.Equal(t, "xenia.com", condenseSiteURL("xenia.com/"))
	require.Equal(t, "chat.xenia.com", condenseSiteURL("chat.xenia.com"))
	require.Equal(t, "chat.xenia.com", condenseSiteURL("chat.xenia.com/"))
	require.Equal(t, "xenia.com/subpath", condenseSiteURL("xenia.com/subpath"))
	require.Equal(t, "xenia.com/subpath", condenseSiteURL("xenia.com/subpath/"))
	require.Equal(t, "chat.xenia.com/subpath", condenseSiteURL("chat.xenia.com/subpath"))
	require.Equal(t, "chat.xenia.com/subpath", condenseSiteURL("chat.xenia.com/subpath/"))

	require.Equal(t, "xenia.com:8080", condenseSiteURL("http://xenia.com:8080"))
	require.Equal(t, "xenia.com:8080", condenseSiteURL("http://xenia.com:8080/"))
	require.Equal(t, "chat.xenia.com:8080", condenseSiteURL("http://chat.xenia.com:8080"))
	require.Equal(t, "chat.xenia.com:8080", condenseSiteURL("http://chat.xenia.com:8080/"))
	require.Equal(t, "xenia.com:8080/subpath", condenseSiteURL("http://xenia.com:8080/subpath"))
	require.Equal(t, "xenia.com:8080/subpath", condenseSiteURL("http://xenia.com:8080/subpath/"))
	require.Equal(t, "chat.xenia.com:8080/subpath", condenseSiteURL("http://chat.xenia.com:8080/subpath"))
	require.Equal(t, "chat.xenia.com:8080/subpath", condenseSiteURL("http://chat.xenia.com:8080/subpath/"))
}
