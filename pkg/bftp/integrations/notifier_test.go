// Copyright 2020 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package integrations

import (
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/bin-work/go-example/pkg/bftp"
	"github.com/bin-work/go-example/pkg/bftp/driver/file"

	"github.com/jlaffaye/ftp"
	"github.com/stretchr/testify/assert"
)

var (
	_ bftp.Filter = &mockNotifier{}
)

type mockNotifier struct {
	actions []string
	lock    sync.Mutex
}

func (m *mockNotifier) BeforeLoginUser(ctx *bftp.Context, userName string) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "BeforeLoginUser")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) BeforePutFile(ctx *bftp.Context, dstPath string) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "BeforePutFile")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) BeforeDeleteFile(ctx *bftp.Context, dstPath string) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "BeforeDeleteFile")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) BeforeChangeCurDir(ctx *bftp.Context, oldCurDir, newCurDir string) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "BeforeChangeCurDir")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) BeforeCreateDir(ctx *bftp.Context, dstPath string) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "BeforeCreateDir")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) BeforeDeleteDir(ctx *bftp.Context, dstPath string) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "BeforeDeleteDir")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) BeforeDownloadFile(ctx *bftp.Context, dstPath string) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "BeforeDownloadFile")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) AfterUserLogin(ctx *bftp.Context, userName, password string, passMatched bool, err error) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "AfterUserLogin")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) AfterFilePut(ctx *bftp.Context, dstPath string, size int64, err error) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "AfterFilePut")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) AfterFileDeleted(ctx *bftp.Context, dstPath string, err error) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "AfterFileDeleted")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) AfterCurDirChanged(ctx *bftp.Context, oldCurDir, newCurDir string, err error) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "AfterCurDirChanged")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) AfterDirCreated(ctx *bftp.Context, dstPath string, err error) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "AfterDirCreated")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) AfterDirDeleted(ctx *bftp.Context, dstPath string, err error) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "AfterDirDeleted")
	m.lock.Unlock()
	return true
}
func (m *mockNotifier) AfterFileDownloaded(ctx *bftp.Context, dstPath string, size int64, err error) bool {
	m.lock.Lock()
	m.actions = append(m.actions, "AfterFileDownloaded")
	m.lock.Unlock()
	return true
}

func assetMockNotifier(t *testing.T, mock *mockNotifier, lastActions []string) {
	if len(lastActions) == 0 {
		return
	}
	mock.lock.Lock()
	assert.EqualValues(t, lastActions, mock.actions[len(mock.actions)-len(lastActions):])
	mock.lock.Unlock()
}

func TestNotification(t *testing.T) {
	err := os.MkdirAll("./testdata", os.ModePerm)
	assert.NoError(t, err)

	var perm = bftp.NewSimplePerm("test", "test")
	driver, err := file.NewDriver("./testdata")
	assert.NoError(t, err)

	opt := &bftp.Options{
		Name:   "test ftpd",
		Driver: driver,
		Port:   2121,
		Auth: &bftp.SimpleAuth{
			Name:     "admin",
			Password: "admin",
		},
		Perm:   perm,
		Logger: new(bftp.DiscardLogger),
	}

	mock := &mockNotifier{}

	runServer(t, opt, []bftp.Notifier{mock}, func() {
		// Give server 0.5 seconds to get to the listening state
		timeout := time.NewTimer(time.Millisecond * 500)

		for {
			f, err := ftp.Connect("localhost:2121")
			if err != nil && len(timeout.C) == 0 { // Retry errors
				continue
			}
			assert.NoError(t, err)

			assert.NoError(t, f.Login("admin", "admin"))
			assetMockNotifier(t, mock, []string{"BeforeLoginUser", "AfterUserLogin"})

			assert.Error(t, f.Login("admin", "1111"))
			assetMockNotifier(t, mock, []string{"BeforeLoginUser", "AfterUserLogin"})

			var content = `test`
			assert.NoError(t, f.Stor("server_test.go", strings.NewReader(content)))
			assetMockNotifier(t, mock, []string{"BeforePutFile", "AfterFilePut"})

			r, err := f.RetrFrom("/server_test.go", 2)
			assert.NoError(t, err)

			buf, err := ioutil.ReadAll(r)
			r.Close()
			assert.NoError(t, err)
			assert.EqualValues(t, "st", string(buf))
			assetMockNotifier(t, mock, []string{"BeforeDownloadFile", "AfterFileDownloaded"})

			err = f.Rename("/server_test.go", "/test.go")
			assert.NoError(t, err)

			err = f.MakeDir("/src")
			assert.NoError(t, err)
			assetMockNotifier(t, mock, []string{"BeforeCreateDir", "AfterDirCreated"})

			err = f.Delete("/test.go")
			assert.NoError(t, err)
			assetMockNotifier(t, mock, []string{"BeforeDeleteFile", "AfterFileDeleted"})

			err = f.ChangeDir("/src")
			assert.NoError(t, err)
			assetMockNotifier(t, mock, []string{"BeforeChangeCurDir", "AfterCurDirChanged"})

			err = f.RemoveDir("/src")
			assert.NoError(t, err)
			assetMockNotifier(t, mock, []string{"BeforeDeleteDir", "AfterDirDeleted"})

			err = f.Quit()
			assert.NoError(t, err)

			break
		}
	})
}
