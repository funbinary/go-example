// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bftp

import (
	"crypto/subtle"
	"sync"

	"github.com/bin-work/go-example/pkg/errors"
)

// Auth is an interface to auth your ftp user login.
type Auth interface {
	CheckPasswd(*Context, string, string) (bool, error)
	Register(account Account)
	Update(account Account)
	Delete(name string)
	IsPutable(ctx *Context, name string) bool
	IsReadable(ctx *Context, name string) bool
	IsDeleable(ctx *Context, name string) bool
}

//
//var (
//	_ Auth = &SimpleAuth{}
//)

// SimpleAuth implements Auth interface to provide a memory user login auth
type SimpleAuth struct {
	Name     string
	Password string
}

// CheckPasswd will check user's password
func (a *SimpleAuth) CheckPasswd(ctx *Context, name, pass string) (bool, error) {
	return constantTimeEquals(name, a.Name) && constantTimeEquals(pass, a.Password), nil
}

func (self *SimpleAuth) Register(account Account) {}
func (self *SimpleAuth) Update(account Account)   {}
func (self *SimpleAuth) Delete(name string)       {}

func (self *SimpleAuth) IsPutable(ctx *Context, name string) bool  { return true }
func (self *SimpleAuth) IsReadable(ctx *Context, name string) bool { return true }
func (self *SimpleAuth) IsDeleable(ctx *Context, name string) bool { return true }

func constantTimeEquals(a, b string) bool {
	return len(a) == len(b) && subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

type Account struct {
	Name     string
	Password string
	Putable  bool // 是否具有上传权限
	Readable bool // 是否具有下载权限
	Deleable bool // 是否具有下载权限
}

type MultiAuth struct {
	Accounts map[string]Account
	mut      sync.RWMutex
}

func NewMultiAuth() *MultiAuth {
	return &MultiAuth{
		Accounts: make(map[string]Account),
	}
}

// Register
//  @Description: 注册账户，如果已经存在将会替换
//  @param account
//
func (self *MultiAuth) Register(account Account) {
	self.mut.Lock()
	self.Accounts[account.Name] = account
	self.mut.Unlock()
}

// CheckPasswd
//  @Description: 校验密码
//  @param ctx
//  @param name 用户名
//  @param pass 密码
//  @return bool 是否匹配
//  @return error 错误说明
//
func (self *MultiAuth) CheckPasswd(ctx *Context, name, pass string) (bool, error) {
	self.mut.RLock()
	defer self.mut.RUnlock()
	if v, ok := self.Accounts[name]; ok {
		if v.Password == pass {
			return true, nil
		} else {
			return false, errors.Errorf("Match password error.")
		}
	}
	return false, errors.Errorf("No found User %s.", name)
}

// IsPutable
//  @Description: name是否具有上传权限
//  @param ctx
//  @param name
//  @return bool true - 具有上传权限
//
func (self *MultiAuth) IsPutable(ctx *Context, name string) bool {
	self.mut.RLock()
	defer self.mut.RUnlock()
	if v, ok := self.Accounts[name]; ok {
		return v.Putable
	}
	return false
}

// IsReadable
//  @Description:name是否具有读取权限
//  @param ctx
//  @param name
//  @return bool true - 具有读取权限
//
func (self *MultiAuth) IsReadable(ctx *Context, name string) bool {
	self.mut.RLock()
	defer self.mut.RUnlock()
	if v, ok := self.Accounts[name]; ok {
		return v.Readable
	}
	return false
}

// IsDeleable
//  @Description: name是否具有删除权限
//  @param ctx
//  @param name
//  @return bool true - 具有删除权限
//
func (self *MultiAuth) IsDeleable(ctx *Context, name string) bool {
	self.mut.RLock()
	defer self.mut.RUnlock()
	if v, ok := self.Accounts[name]; ok {
		return v.Deleable
	}
	return false
}

// Update
//  @Description: 更新用户信息
//  @param account
//
func (self *MultiAuth) Update(account Account) {
	self.mut.Lock()
	self.Accounts[account.Name] = account
	self.mut.Unlock()
}

// Delete
//  @Description: 删除name
//  @param name
//
func (self *MultiAuth) Delete(name string) {
	self.mut.Lock()
	delete(self.Accounts, name)
	self.mut.Unlock()
}
