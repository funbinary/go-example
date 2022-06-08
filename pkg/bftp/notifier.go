// Copyright 2020 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bftp

// Filter represents a notification operator interface
type Filter interface {
	BeforeLoginUser(ctx *Context, userName string) bool
	BeforePutFile(ctx *Context, dstPath string) bool
	BeforeDeleteFile(ctx *Context, dstPath string) bool
	BeforeChangeCurDir(ctx *Context, oldCurDir, newCurDir string) bool
	BeforeCreateDir(ctx *Context, dstPath string) bool
	BeforeDeleteDir(ctx *Context, dstPath string) bool
	BeforeDownloadFile(ctx *Context, dstPath string) bool
	AfterUserLogin(ctx *Context, userName, password string, passMatched bool, err error)
	AfterFilePut(ctx *Context, dstPath string, size int64, err error)
	AfterFileDeleted(ctx *Context, dstPath string, err error)
	AfterFileDownloaded(ctx *Context, dstPath string, size int64, err error)
	AfterCurDirChanged(ctx *Context, oldCurDir, newCurDir string, err error)
	AfterDirCreated(ctx *Context, dstPath string, err error)
	AfterDirDeleted(ctx *Context, dstPath string, err error)
}

type filterList []Filter

var (
	_ Filter = filterList{}
)

func (self filterList) BeforeLoginUser(ctx *Context, userName string) bool {
	for _, notifier := range self {
		notifier.BeforeLoginUser(ctx, userName)
	}
	return true
}

func (self filterList) BeforePutFile(ctx *Context, dstPath string) bool {
	for _, notifier := range self {
		notifier.BeforePutFile(ctx, dstPath)
	}
	return true
}

func (self filterList) BeforeDeleteFile(ctx *Context, dstPath string) bool {
	for _, notifier := range self {
		notifier.BeforeDeleteFile(ctx, dstPath)
	}
	return true
}

func (self filterList) BeforeChangeCurDir(ctx *Context, oldCurDir, newCurDir string) bool {
	for _, notifier := range self {
		notifier.BeforeChangeCurDir(ctx, oldCurDir, newCurDir)
	}
	return true
}

func (self filterList) BeforeCreateDir(ctx *Context, dstPath string) bool {
	for _, notifier := range self {
		notifier.BeforeCreateDir(ctx, dstPath)
	}
	return true
}

func (self filterList) BeforeDeleteDir(ctx *Context, dstPath string) bool {
	for _, notifier := range self {
		notifier.BeforeDeleteDir(ctx, dstPath)
	}
	return true
}

func (self filterList) BeforeDownloadFile(ctx *Context, dstPath string) bool {
	for _, notifier := range self {
		notifier.BeforeDownloadFile(ctx, dstPath)
	}
	return true
}

func (self filterList) AfterUserLogin(ctx *Context, userName, password string, passMatched bool, err error) {
	for _, notifier := range self {
		notifier.AfterUserLogin(ctx, userName, password, passMatched, err)
	}
}

func (self filterList) AfterFilePut(ctx *Context, dstPath string, size int64, err error) {
	for _, notifier := range self {
		notifier.AfterFilePut(ctx, dstPath, size, err)
	}
}

func (self filterList) AfterFileDeleted(ctx *Context, dstPath string, err error) {
	for _, notifier := range self {
		notifier.AfterFileDeleted(ctx, dstPath, err)
	}
}

func (self filterList) AfterFileDownloaded(ctx *Context, dstPath string, size int64, err error) {
	for _, notifier := range self {
		notifier.AfterFileDownloaded(ctx, dstPath, size, err)
	}
}

func (self filterList) AfterCurDirChanged(ctx *Context, oldCurDir, newCurDir string, err error) {
	for _, notifier := range self {
		notifier.AfterCurDirChanged(ctx, oldCurDir, newCurDir, err)
	}
}

func (self filterList) AfterDirCreated(ctx *Context, dstPath string, err error) {
	for _, notifier := range self {
		notifier.AfterDirCreated(ctx, dstPath, err)
	}
}

func (self filterList) AfterDirDeleted(ctx *Context, dstPath string, err error) {
	for _, notifier := range self {
		notifier.AfterDirDeleted(ctx, dstPath, err)
	}
}

// NullFilter implements Filter
type NullFilter struct{}

var (
	_ Filter = &NullFilter{}
)

// BeforeLoginUser implements Filter
func (NullFilter) BeforeLoginUser(ctx *Context, userName string) bool {
	return true
}

// BeforePutFile implements Filter
func (NullFilter) BeforePutFile(ctx *Context, dstPath string) bool {
	return true
}

// BeforeDeleteFile implements Filter
func (NullFilter) BeforeDeleteFile(ctx *Context, dstPath string) bool {
	return true
}

// BeforeChangeCurDir implements Filter
func (NullFilter) BeforeChangeCurDir(ctx *Context, oldCurDir, newCurDir string) bool {
	return true
}

// BeforeCreateDir implements Filter
func (NullFilter) BeforeCreateDir(ctx *Context, dstPath string) bool {
	return true
}

// BeforeDeleteDir implements Filter
func (NullFilter) BeforeDeleteDir(ctx *Context, dstPath string) bool {
	return true
}

// BeforeDownloadFile implements Filter
func (NullFilter) BeforeDownloadFile(ctx *Context, dstPath string) bool {
	return true
}

// AfterUserLogin implements Filter
func (NullFilter) AfterUserLogin(ctx *Context, userName, password string, passMatched bool, err error) {
}

// AfterFilePut implements Filter
func (NullFilter) AfterFilePut(ctx *Context, dstPath string, size int64, err error) {

}

// AfterFileDeleted implements Filter
func (NullFilter) AfterFileDeleted(ctx *Context, dstPath string, err error) {

}

// AfterFileDownloaded implements Filter
func (NullFilter) AfterFileDownloaded(ctx *Context, dstPath string, size int64, err error) {

}

// AfterCurDirChanged implements Filter
func (NullFilter) AfterCurDirChanged(ctx *Context, oldCurDir, newCurDir string, err error) {

}

// AfterDirCreated implements Filter
func (NullFilter) AfterDirCreated(ctx *Context, dstPath string, err error) {

}

// AfterDirDeleted implements Filter
func (NullFilter) AfterDirDeleted(ctx *Context, dstPath string, err error) {

}
