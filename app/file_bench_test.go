// Copyright (c) 2018-present Xenia, Inc. All Rights Reserved.
// See License.txt for license information.

package app

import (
	"bytes"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"io"
	"io/ioutil"
	"math/rand"
	"testing"
	"time"

	"github.com/xzl8028/xenia-server/mlog"
	"github.com/xzl8028/xenia-server/model"
)

var randomJPEG []byte
var randomGIF []byte
var zero10M = make([]byte, 10*1024*1024)
var rgba *image.RGBA

func prepareTestImages(tb testing.TB) {
	if rgba != nil {
		return
	}

	// Create a random image (pre-seeded for predictability)
	rgba = image.NewRGBA(image.Rectangle{
		image.Point{0, 0},
		image.Point{2048, 2048},
	})
	_, err := rand.New(rand.NewSource(1)).Read(rgba.Pix)
	if err != nil {
		tb.Fatal(err)
	}

	// Encode it as JPEG and GIF
	buf := &bytes.Buffer{}
	err = jpeg.Encode(buf, rgba, &jpeg.Options{Quality: 50})
	if err != nil {
		tb.Fatal(err)
	}
	randomJPEG = buf.Bytes()

	buf = &bytes.Buffer{}
	err = gif.Encode(buf, rgba, nil)
	if err != nil {
		tb.Fatal(err)
	}
	randomGIF = buf.Bytes()
}

func BenchmarkUploadFile(b *testing.B) {
	prepareTestImages(b)
	th := Setup(b).InitBasic()
	defer th.TearDown()
	// disable logging in the benchmark, as best we can
	th.App.Log.SetConsoleLevel(mlog.LevelError)
	teamId := model.NewId()
	channelId := model.NewId()
	userId := model.NewId()

	mb := func(i int) int {
		return (i + 512*1024) / (1024 * 1024)
	}

	files := []struct {
		title string
		ext   string
		data  []byte
	}{
		{fmt.Sprintf("random-%dMb-gif", mb(len(randomGIF))), ".gif", randomGIF},
		{fmt.Sprintf("random-%dMb-jpg", mb(len(randomJPEG))), ".jpg", randomJPEG},
		{fmt.Sprintf("zero-%dMb", mb(len(zero10M))), ".zero", zero10M},
	}

	file_benchmarks := []struct {
		title string
		f     func(b *testing.B, n int, data []byte, ext string)
	}{
		{
			title: "raw-ish DoUploadFile",
			f: func(b *testing.B, n int, data []byte, ext string) {
				info1, err := th.App.DoUploadFile(time.Now(), teamId, channelId,
					userId, fmt.Sprintf("BenchmarkDoUploadFile-%d%s", n, ext), data)
				if err != nil {
					b.Fatal(err)
				}
				th.App.Srv.Store.FileInfo().PermanentDelete(info1.Id)
				th.App.RemoveFile(info1.Path)

			},
		},
		{
			title: "raw UploadFileX Content-Length",
			f: func(b *testing.B, n int, data []byte, ext string) {
				info, aerr := th.App.UploadFileX(channelId,
					fmt.Sprintf("BenchmarkUploadFileTask-%d%s", n, ext),
					bytes.NewReader(data),
					UploadFileSetTeamId(teamId),
					UploadFileSetUserId(userId),
					UploadFileSetTimestamp(time.Now()),
					UploadFileSetContentLength(int64(len(data))),
					UploadFileSetRaw())
				if aerr != nil {
					b.Fatal(aerr)
				}
				th.App.Srv.Store.FileInfo().PermanentDelete(info.Id)
				th.App.RemoveFile(info.Path)
			},
		},
		{
			title: "raw UploadFileX chunked",
			f: func(b *testing.B, n int, data []byte, ext string) {
				info, aerr := th.App.UploadFileX(channelId,
					fmt.Sprintf("BenchmarkUploadFileTask-%d%s", n, ext),
					bytes.NewReader(data),
					UploadFileSetTeamId(teamId),
					UploadFileSetUserId(userId),
					UploadFileSetTimestamp(time.Now()),
					UploadFileSetContentLength(-1),
					UploadFileSetRaw())
				if aerr != nil {
					b.Fatal(aerr)
				}
				th.App.Srv.Store.FileInfo().PermanentDelete(info.Id)
				th.App.RemoveFile(info.Path)
			},
		},
		{
			title: "image UploadFiles",
			f: func(b *testing.B, n int, data []byte, ext string) {
				resp, err := th.App.UploadFiles(teamId, channelId, userId,
					[]io.ReadCloser{ioutil.NopCloser(bytes.NewReader(data))},
					[]string{fmt.Sprintf("BenchmarkDoUploadFiles-%d%s", n, ext)},
					[]string{},
					time.Now())
				if err != nil {
					b.Fatal(err)
				}
				th.App.Srv.Store.FileInfo().PermanentDelete(resp.FileInfos[0].Id)
				th.App.RemoveFile(resp.FileInfos[0].Path)
			},
		},
		{
			title: "image UploadFileX Content-Length",
			f: func(b *testing.B, n int, data []byte, ext string) {
				info, aerr := th.App.UploadFileX(channelId,
					fmt.Sprintf("BenchmarkUploadFileTask-%d%s", n, ext),
					bytes.NewReader(data),
					UploadFileSetTeamId(teamId),
					UploadFileSetUserId(userId),
					UploadFileSetTimestamp(time.Now()),
					UploadFileSetContentLength(int64(len(data))))
				if aerr != nil {
					b.Fatal(aerr)
				}
				th.App.Srv.Store.FileInfo().PermanentDelete(info.Id)
				th.App.RemoveFile(info.Path)
			},
		},
		{
			title: "image UploadFileX chunked",
			f: func(b *testing.B, n int, data []byte, ext string) {
				info, aerr := th.App.UploadFileX(channelId,
					fmt.Sprintf("BenchmarkUploadFileTask-%d%s", n, ext),
					bytes.NewReader(data),
					UploadFileSetTeamId(teamId),
					UploadFileSetUserId(userId),
					UploadFileSetTimestamp(time.Now()),
					UploadFileSetContentLength(int64(len(data))))
				if aerr != nil {
					b.Fatal(aerr)
				}
				th.App.Srv.Store.FileInfo().PermanentDelete(info.Id)
				th.App.RemoveFile(info.Path)
			},
		},
	}

	for _, file := range files {
		for _, fb := range file_benchmarks {
			b.Run(file.title+"-"+fb.title, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					fb.f(b, i, file.data, file.ext)
				}
			})
		}
	}
}
