// Copyright 2018 The goftp Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package bftp

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

// Command represents a Command interface to a ftp Command
type Command interface {
	IsExtend() bool
	RequireParam() bool
	RequireAuth() bool
	Execute(*Session, string)
}

var (
	defaultCommands = map[string]Command{
		"ADAT": CommandAdat{},
		"ALLO": CommandAllo{},
		"APPE": CommandAppe{},
		"AUTH": CommandAuth{},
		"CDUP": CommandCdup{},
		"CWD":  CommandCwd{},
		"CCC":  CommandCcc{},
		"CONF": CommandConf{},
		"CLNT": CommandCLNT{},
		"DELE": CommandDele{},
		"ENC":  CommandEnc{},
		"EPRT": CommandEprt{},
		"EPSV": CommandEpsv{},
		"FEAT": CommandFeat{},
		"LIST": CommandList{},
		"LPRT": CommandLprt{},
		"NLST": CommandNlst{},
		"MDTM": CommandMdtm{},
		"MIC":  CommandMic{},
		"MLSD": CommandMLSD{},
		"MKD":  CommandMkd{},
		"MODE": CommandMode{},
		"NOOP": CommandNoop{},
		"OPTS": CommandOpts{},
		"PASS": CommandPass{},
		"PASV": CommandPasv{},
		"PBSZ": CommandPbsz{},
		"PORT": CommandPort{},
		"PROT": CommandProt{},
		"PWD":  CommandPwd{},
		"QUIT": CommandQuit{},
		"RETR": CommandRetr{},
		"REST": CommandRest{},
		"RNFR": CommandRnfr{},
		"RNTO": CommandRnto{},
		"RMD":  CommandRmd{},
		"SIZE": CommandSize{},
		"STAT": CommandStat{},
		"STOR": CommandStor{},
		"STRU": CommandStru{},
		"SYST": CommandSyst{},
		"TYPE": CommandType{},
		"USER": CommandUser{},
		"XCUP": CommandCdup{},
		"XCWD": CommandCwd{},
		"XMKD": CommandMkd{},
		"XPWD": CommandPwd{},
		"XRMD": CommandXRmd{},
	}
)

// DefaultCommands returns the default Commands
func DefaultCommands() map[string]Command {
	return defaultCommands
}

type CommandAdat struct{}

func (cmd CommandAdat) IsExtend() bool {
	return false
}

func (cmd CommandAdat) RequireParam() bool {
	return true
}

func (cmd CommandAdat) RequireAuth() bool {
	return true
}

func (cmd CommandAdat) Execute(sess *Session, param string) {
	sess.writeMessage(CODE_ACTION_NOTAKEN, "Action not taken")
}

// commandAllo responds to the ALLO FTP command.
//
// This is essentially a ping from the client so we just respond with an
// basic OK message.
type CommandAllo struct{}

func (cmd CommandAllo) IsExtend() bool {
	return false
}

func (cmd CommandAllo) RequireParam() bool {
	return false
}

func (cmd CommandAllo) RequireAuth() bool {
	return false
}

func (cmd CommandAllo) Execute(sess *Session, param string) {
	sess.writeMessage(CODE_COMMAND_IMPLEMENTED, "Obsolete")
}

// CommandAppe responds to the APPE FTP Command. It allows the user to upload a
// new file but always append if file exists otherwise create one.
type CommandAppe struct{}

func (cmd CommandAppe) IsExtend() bool {
	return false
}

func (cmd CommandAppe) RequireParam() bool {
	return true
}

func (cmd CommandAppe) RequireAuth() bool {
	return true
}

func (cmd CommandAppe) Execute(sess *Session, param string) {
	targetPath := sess.buildPath(param)
	sess.writeMessage(CODE_FILE_STATUS_OK, "Data transfer starting")

	if sess.preCommand != "REST" {
		sess.lastFilePos = -1
	}
	defer func() {
		sess.lastFilePos = -1
	}()

	var ctx = Context{
		Sess:  sess,
		Cmd:   "APPE",
		Param: param,
		Data:  make(map[string]interface{}),
	}
	if !sess.server.Options.Auth.IsPutable(&ctx, sess.user) {
		sess.writeMessage(CODE_FILE_ACTION_NOTAKEN, fmt.Sprint("No auth for put file."))
		return
	}
	if !sess.server.filters.BeforePutFile(&ctx, targetPath) {
		sess.writeMessage(CODE_FILE_ACTION_NOTAKEN, fmt.Sprint("This Command refused"))
		return
	}
	size, err := sess.server.Driver.PutFile(&ctx, targetPath, sess.dataConn, sess.lastFilePos)
	sess.server.filters.AfterFilePut(&ctx, targetPath, size, err)
	if err == nil {
		msg := fmt.Sprintf("OK, received %d bytes", size)
		sess.writeMessage(CODE_DATA_CONN_CLOSE, msg)
	} else {
		sess.writeMessage(CODE_FILE_ACTION_NOTAKEN, fmt.Sprint("error during transfer: ", err))
	}
}

type CommandAuth struct{}

func (cmd CommandAuth) IsExtend() bool {
	return false
}

func (cmd CommandAuth) RequireParam() bool {
	return true
}

func (cmd CommandAuth) RequireAuth() bool {
	return false
}

func (cmd CommandAuth) Execute(sess *Session, param string) {
	if param == "TLS" && sess.server.tlsConfig != nil {
		sess.writeMessage(CODE_TLS_AUTH_OK, "AUTH Command OK")
		err := sess.upgradeToTLS()
		if err != nil {
			sess.logf("Error upgrading connection to TLS %v", err.Error())
		}
	} else {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "Action not taken")
	}
}

type CommandCcc struct{}

func (cmd CommandCcc) IsExtend() bool {
	return false
}

func (cmd CommandCcc) RequireParam() bool {
	return true
}

func (cmd CommandCcc) RequireAuth() bool {
	return true
}

func (cmd CommandCcc) Execute(sess *Session, param string) {
	sess.writeMessage(CODE_ACTION_NOTAKEN, "Action not taken")
}

// cmdCdup responds to the CDUP FTP Command.
//
// Allows the client change their current directory to the parent.
type CommandCdup struct{}

func (cmd CommandCdup) IsExtend() bool {
	return false
}

func (cmd CommandCdup) RequireParam() bool {
	return false
}

func (cmd CommandCdup) RequireAuth() bool {
	return true
}

func (cmd CommandCdup) Execute(sess *Session, param string) {
	otherCmd := &CommandCwd{}
	otherCmd.Execute(sess, "..")
}

// CommandCwd responds to the CWD FTP Command. It allows the client to change the
// current working directory.
type CommandCwd struct{}

func (cmd CommandCwd) IsExtend() bool {
	return false
}

func (cmd CommandCwd) RequireParam() bool {
	return true
}

func (cmd CommandCwd) RequireAuth() bool {
	return true
}

func (cmd CommandCwd) Execute(sess *Session, param string) {
	path := sess.buildPath(param)
	var ctx = Context{
		Sess:  sess,
		Cmd:   "CWD",
		Param: param,
		Data:  make(map[string]interface{}),
	}
	info, err := sess.server.Driver.Stat(&ctx, path)
	if err != nil {
		sess.logf("%v", err)
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("Directory change to ", path, " failed."))
		return
	}
	if !info.IsDir() {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("Directory change to ", path, " is a file"))
		return
	}

	if !sess.server.filters.BeforeChangeCurDir(&ctx, sess.curDir, path) {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("This command refused by filter"))
	}
	err = sess.changeCurDir(path)
	sess.server.filters.AfterCurDirChanged(&ctx, sess.curDir, path, err)
	if err == nil {
		sess.writeMessage(CODE_FILE_COMMANG_OK, "Directory changed to "+path)
	} else {
		sess.logf("%v", err)
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("Directory change to ", path, " failed."))
	}
}

type CommandCLNT struct{}

func (cmd CommandCLNT) IsExtend() bool {
	return true
}

func (cmd CommandCLNT) RequireParam() bool {
	return false
}

func (cmd CommandCLNT) RequireAuth() bool {
	return false
}

func (cmd CommandCLNT) Execute(sess *Session, param string) {
	sess.clientSoft = param
	sess.writeMessage(CODE_COMMAND_OK, "OK")
}

// CommandDele responds to the DELE FTP Command. It allows the client to delete
// a file
type CommandDele struct{}

func (cmd CommandDele) IsExtend() bool {
	return false
}

func (cmd CommandDele) RequireParam() bool {
	return true
}

func (cmd CommandDele) RequireAuth() bool {
	return true
}

func (cmd CommandDele) Execute(sess *Session, param string) {
	path := sess.buildPath(param)
	var ctx = Context{
		Sess:  sess,
		Cmd:   "DELE",
		Param: param,
		Data:  make(map[string]interface{}),
	}
	if !sess.server.Options.Auth.IsDeleable(&ctx, sess.user) {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("No auth for delte file."))
		return
	}
	if !sess.server.filters.BeforeDeleteFile(&ctx, path) {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "This command refused. ")
		return
	}
	err := sess.server.Driver.DeleteFile(&ctx, path)
	sess.server.filters.AfterFileDeleted(&ctx, path, err)
	if err == nil {
		sess.writeMessage(CODE_FILE_COMMANG_OK, "File deleted")
	} else {
		sess.logf("%v", err)
		sess.writeMessage(CODE_ACTION_NOTAKEN, "File delete failed. ")
	}
}

type CommandEnc struct{}

func (cmd CommandEnc) IsExtend() bool {
	return false
}

func (cmd CommandEnc) RequireParam() bool {
	return true
}

func (cmd CommandEnc) RequireAuth() bool {
	return true
}

func (cmd CommandEnc) Execute(sess *Session, param string) {
	sess.writeMessage(CODE_ACTION_NOTAKEN, "Action not taken")
}

// CommandEprt responds to the EPRT FTP Command. It allows the client to
// request an active data socket with more options than the original PORT
// Command. It mainly adds ipv6 support.
type CommandEprt struct{}

func (cmd CommandEprt) IsExtend() bool {
	return true
}

func (cmd CommandEprt) RequireParam() bool {
	return true
}

func (cmd CommandEprt) RequireAuth() bool {
	return true
}

func (cmd CommandEprt) Execute(sess *Session, param string) {
	delim := string(param[0:1])
	parts := strings.Split(param, delim)
	addressFamily, err := strconv.Atoi(parts[1])
	if err != nil {
		sess.writeMessage(522, "Network protocol not supported, use (1,2)")
		return
	}
	if addressFamily != 1 && addressFamily != 2 {
		sess.writeMessage(522, "Network protocol not supported, use (1,2)")
		return
	}

	host := parts[2]
	port, err := strconv.Atoi(parts[3])
	if err != nil {
		sess.writeMessage(522, "Network protocol not supported, use (1,2)")
		return
	}
	socket, err := newActiveSocket(sess, host, port)
	if err != nil {
		sess.writeMessage(425, "Data connection failed")
		return
	}
	sess.dataConn = socket
	sess.writeMessage(200, "Connection established ("+strconv.Itoa(port)+")")
}

// CommandEpsv responds to the EPSV FTP Command. It allows the client to
// request a passive data socket with more options than the original PASV
// Command. It mainly adds ipv6 support, although we don't support that yet.
type CommandEpsv struct{}

func (cmd CommandEpsv) IsExtend() bool {
	return true
}

func (cmd CommandEpsv) RequireParam() bool {
	return false
}

func (cmd CommandEpsv) RequireAuth() bool {
	return true
}

func (cmd CommandEpsv) Execute(sess *Session, param string) {
	socket, err := sess.newPassiveSocket()
	if err != nil {
		sess.log(err)
		sess.writeMessage(425, "Data connection failed")
		return
	}

	msg := fmt.Sprintf("Entering Extended Passive Mode (|||%d|)", socket.Port())
	sess.writeMessage(229, msg)
}

type CommandFeat struct{}

func (cmd CommandFeat) IsExtend() bool {
	return false
}

func (cmd CommandFeat) RequireParam() bool {
	return false
}

func (cmd CommandFeat) RequireAuth() bool {
	return false
}

func (cmd CommandFeat) Execute(sess *Session, param string) {
	sess.writeMessageMultiline(CODE_SYSTEM_STATUS, sess.server.feats)
}

// CommandList responds to the LIST FTP Command. It allows the client to retrieve
// a detailed listing of the contents of a directory.
type CommandList struct{}

func (cmd CommandList) IsExtend() bool {
	return false
}

func (cmd CommandList) RequireParam() bool {
	return false
}

func (cmd CommandList) RequireAuth() bool {
	return true
}

func convertFileInfo(sess *Session, f os.FileInfo, p string) (FileInfo, error) {
	mode, err := sess.server.Perm.GetMode(p)
	if err != nil {
		return nil, err
	}
	if f.IsDir() {
		mode |= os.ModeDir
	}
	owner, err := sess.server.Perm.GetOwner(p)
	if err != nil {
		return nil, err
	}
	group, err := sess.server.Perm.GetGroup(p)
	if err != nil {
		return nil, err
	}
	return &fileInfo{
		FileInfo: f,
		mode:     mode,
		owner:    owner,
		group:    group,
	}, nil
}

func list(sess *Session, cmd, p, param string) ([]FileInfo, error) {
	var ctx = &Context{
		Sess:  sess,
		Cmd:   cmd,
		Param: param,
		Data:  make(map[string]interface{}),
	}
	info, err := sess.server.Driver.Stat(ctx, p)
	if err != nil {
		return nil, err
	}

	if info == nil {
		sess.logf("%s: no such file or directory.\n", p)
		return []FileInfo{}, nil
	}

	var files []FileInfo
	if info.IsDir() {
		// 如果是目录，则调用listdir，获取目录下所有文件信息
		err = sess.server.Driver.ListDir(ctx, p, func(f os.FileInfo) error {
			info, err := convertFileInfo(sess, f, path.Join(p, f.Name()))
			if err != nil {
				return err
			}
			files = append(files, info)
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		newInfo, err := convertFileInfo(sess, info, p)
		if err != nil {
			return nil, err
		}
		files = append(files, newInfo)
	}
	return files, nil
}
func parseListParam(param string) (path string) {
	if len(param) == 0 {
		path = param
	} else {
		fields := strings.Fields(param)
		i := 0
		for _, field := range fields {
			if !strings.HasPrefix(field, "-") {
				break
			}
			i = strings.LastIndex(param, " "+field) + len(field) + 1
		}
		path = strings.TrimLeft(param[i:], " ") //Get all the path even with space inside
	}
	return path
}

func (cmd CommandList) Execute(sess *Session, param string) {
	p := sess.buildPath(parseListParam(param))

	files, err := list(sess, "LIST", p, param)
	if err != nil {
		sess.writeMessage(CODE_ACTION_NOTAKEN, err.Error())
		return
	}

	sess.writeMessage(CODE_FILE_STATUS_OK, "Opening ASCII mode data connection for file list")
	sess.sendOutofbandData(listFormatter(files).Detailed())
}

// CommandLprt responds to the LPRT FTP Command. It allows the client to
// request an active data socket with more options than the original PORT
// Command.  FTP Operation Over Big Address Records.
type CommandLprt struct{}

func (cmd CommandLprt) IsExtend() bool {
	return true
}

func (cmd CommandLprt) RequireParam() bool {
	return true
}

func (cmd CommandLprt) RequireAuth() bool {
	return true
}

func (cmd CommandLprt) Execute(sess *Session, param string) {
	// No tests for this code yet

	parts := strings.Split(param, ",")

	addressFamily, err := strconv.Atoi(parts[0])
	if err != nil {
		sess.writeMessage(522, "Network protocol not supported, use 4")
		return
	}
	if addressFamily != 4 {
		sess.writeMessage(522, "Network protocol not supported, use 4")
		return
	}

	addressLength, err := strconv.Atoi(parts[1])
	if err != nil {
		sess.writeMessage(522, "Network protocol not supported, use 4")
		return
	}
	if addressLength != 4 {
		sess.writeMessage(522, "Network IP length not supported, use 4")
		return
	}

	host := strings.Join(parts[2:2+addressLength], ".")

	portLength, err := strconv.Atoi(parts[2+addressLength])
	if err != nil {
		sess.writeMessage(522, "Network protocol not supported, use 4")
		return
	}
	portAddress := parts[3+addressLength : 3+addressLength+portLength]

	// Convert string[] to byte[]
	portBytes := make([]byte, portLength)
	for i := range portAddress {
		p, _ := strconv.Atoi(portAddress[i])
		portBytes[i] = byte(p)
	}

	// convert the bytes to an int
	port := int(binary.BigEndian.Uint16(portBytes))

	// if the existing connection is on the same host/port don't reconnect
	if sess.dataConn.Host() == host && sess.dataConn.Port() == port {
		return
	}

	socket, err := newActiveSocket(sess, host, port)
	if err != nil {
		sess.writeMessage(CODE_FAILED_OPEN_DATA_CONN, "Data connection failed")
		return
	}
	sess.dataConn = socket
	sess.writeMessage(CODE_COMMAND_OK, "Connection established ("+strconv.Itoa(port)+")")
}

// CommandMdtm responds to the MDTM FTP Command. It allows the client to
// retreive the last modified time of a file.
type CommandMdtm struct{}

func (cmd CommandMdtm) IsExtend() bool {
	return false
}

func (cmd CommandMdtm) RequireParam() bool {
	return true
}

func (cmd CommandMdtm) RequireAuth() bool {
	return true
}

func (cmd CommandMdtm) Execute(sess *Session, param string) {
	path := sess.buildPath(param)
	stat, err := sess.server.Driver.Stat(&Context{
		Sess:  sess,
		Cmd:   "MDTM",
		Param: param,
		Data:  make(map[string]interface{}),
	}, path)
	if err == nil {
		sess.writeMessage(CODE_FILE_STATUS, stat.ModTime().Format("20060102150405"))
	} else {
		sess.writeMessage(CODE_FILE_ACTION_NOTAKEN, "File not available")
	}
}

type CommandMic struct{}

func (cmd CommandMic) IsExtend() bool {
	return false
}

func (cmd CommandMic) RequireParam() bool {
	return true
}

func (cmd CommandMic) RequireAuth() bool {
	return true
}

func (cmd CommandMic) Execute(sess *Session, param string) {
	sess.writeMessage(CODE_ACTION_NOTAKEN, "Action not taken")
}

// CommandMkd responds to the MKD FTP Command. It allows the client to create
// a new directory
type CommandMkd struct{}

func (cmd CommandMkd) IsExtend() bool {
	return false
}

func (cmd CommandMkd) RequireParam() bool {
	return true
}

func (cmd CommandMkd) RequireAuth() bool {
	return true
}

func (cmd CommandMkd) Execute(sess *Session, param string) {
	path := sess.buildPath(param)
	var ctx = Context{
		Sess:  sess,
		Cmd:   "MKD",
		Param: param,
		Data:  make(map[string]interface{}),
	}
	if !sess.server.Options.Auth.IsPutable(&ctx, sess.user) {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("No auth for create dir."))
		return
	}
	if !sess.server.filters.BeforeCreateDir(&ctx, path) {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("This Command refused."))
		return
	}
	err := sess.server.Driver.MakeDir(&ctx, path)
	sess.server.filters.AfterDirCreated(&ctx, path, err)
	if err == nil {
		sess.writeMessage(CODE_PATHNAME_CREATED, "Directory created")
	} else {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("Action not taken: ", err))
	}
}

type CommandMLSD struct{}

func (cmd CommandMLSD) IsExtend() bool {
	return true
}

func (cmd CommandMLSD) RequireParam() bool {
	return false
}

func (cmd CommandMLSD) RequireAuth() bool {
	return true
}

func toMLSDFormat(files []FileInfo) []byte {
	var buf bytes.Buffer
	for _, file := range files {
		var fileType = "file"
		if file.IsDir() {
			fileType = "dir"
		}
		/*Possible facts "Size" / "Modify" / "Create" /
				  "Type" / "Unique" / "Perm" /
				  "Lang" / "Media-Type" / "CharSet"
				  TODO: Perm pvals        = "a" / "c" / "d" / "e" / "f" /
		                     "l" / "m" / "p" / "r" / "w"
		*/
		fmt.Fprintf(&buf,
			"Type=%s;Modify=%s;Size=%d; %s\n",
			fileType,
			file.ModTime().Format("20060102150405"),
			file.Size(),
			file.Name(),
		)
	}
	return buf.Bytes()
}

func (cmd CommandMLSD) Execute(sess *Session, param string) {
	if param == "" {
		param = sess.curDir
	}
	p := sess.buildPath(param)

	files, err := list(sess, "MLSD", p, param)
	if err != nil {
		sess.writeMessage(CODE_ACTION_NOTAKEN, err.Error())
		return
	}

	sess.writeMessage(CODE_FILE_STATUS_OK, "Opening ASCII mode data connection for file list")
	sess.sendOutofbandData(toMLSDFormat(files))
}

// cmdMode responds to the MODE FTP Command.
//
// the original FTP spec had various options for hosts to negotiate how data
// would be sent over the data socket, In reality these days (S)tream mode
// is all that is used for the mode - data is just streamed down the data
// socket unchanged.
type CommandMode struct{}

func (cmd CommandMode) IsExtend() bool {
	return false
}

func (cmd CommandMode) RequireParam() bool {
	return true
}

func (cmd CommandMode) RequireAuth() bool {
	return true
}

func (cmd CommandMode) Execute(sess *Session, param string) {
	if strings.ToUpper(param) == "S" {
		sess.writeMessage(CODE_COMMAND_OK, "OK")
	} else {
		sess.writeMessage(CODE_CMD_IMPLEMENTED_WITH_PARAM, "MODE is an obsolete Command")
	}
}

// CommandNlst responds to the NLST FTP Command. It allows the client to
// retrieve a list of filenames in the current directory.
type CommandNlst struct{}

func (cmd CommandNlst) IsExtend() bool {
	return false
}

func (cmd CommandNlst) RequireParam() bool {
	return false
}

func (cmd CommandNlst) RequireAuth() bool {
	return true
}

func (cmd CommandNlst) Execute(sess *Session, param string) {
	var ctx = &Context{
		Sess:  sess,
		Cmd:   "NLST",
		Param: param,
		Data:  make(map[string]interface{}),
	}
	path := sess.buildPath(parseListParam(param))
	info, err := sess.server.Driver.Stat(ctx, path)
	if err != nil {
		sess.writeMessage(CODE_ACTION_NOTAKEN, err.Error())
		return
	}
	if !info.IsDir() {
		sess.writeMessage(CODE_ACTION_NOTAKEN, param+" is not a directory")
		return
	}

	var files []FileInfo
	err = sess.server.Driver.ListDir(ctx, path, func(f os.FileInfo) error {
		mode, err := sess.server.Perm.GetMode(path)
		if err != nil {
			return err
		}
		if info.IsDir() {
			mode |= os.ModeDir
		}
		owner, err := sess.server.Perm.GetOwner(path)
		if err != nil {
			return err
		}
		group, err := sess.server.Perm.GetGroup(path)
		if err != nil {
			return err
		}
		files = append(files, &fileInfo{
			FileInfo: f,
			mode:     mode,
			owner:    owner,
			group:    group,
		})
		return nil
	})
	if err != nil {
		sess.writeMessage(CODE_ACTION_NOTAKEN, err.Error())
		return
	}
	sess.writeMessage(CODE_FILE_STATUS_OK, "Opening ASCII mode data connection for file list")
	sess.sendOutofbandData(listFormatter(files).Short())
}

// cmdNoop responds to the NOOP FTP Command.
//
// This is essentially a ping from the client so we just respond with an
// basic 200 message.
type CommandNoop struct{}

func (cmd CommandNoop) IsExtend() bool {
	return false
}

func (cmd CommandNoop) RequireParam() bool {
	return false
}

func (cmd CommandNoop) RequireAuth() bool {
	return false
}

func (cmd CommandNoop) Execute(sess *Session, param string) {
	sess.writeMessage(200, "OK")
}

type CommandOpts struct{}

func (cmd CommandOpts) IsExtend() bool {
	return false
}

func (cmd CommandOpts) RequireParam() bool {
	return false
}

func (cmd CommandOpts) RequireAuth() bool {
	return false
}

func (cmd CommandOpts) Execute(sess *Session, param string) {
	parts := strings.Fields(param)
	if len(parts) != 2 {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "Unknow params")
		return
	}
	if strings.ToUpper(parts[0]) != "UTF8" {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "Unknow params")
		return
	}

	if strings.ToUpper(parts[1]) == "ON" {
		sess.writeMessage(CODE_COMMAND_OK, "UTF8 mode enabled")
	} else {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "Unsupported non-utf8 mode")
	}
}

// CommandPass respond to the PASS FTP Command by asking the driver if the
// supplied username and password are valid
type CommandPass struct{}

func (cmd CommandPass) IsExtend() bool {
	return false
}

func (cmd CommandPass) RequireParam() bool {
	return true
}

func (cmd CommandPass) RequireAuth() bool {
	return false
}

func (cmd CommandPass) Execute(sess *Session, param string) {
	auth := sess.server.Auth
	// If Driver implements Auth then call that instead of the Server version
	if driverAuth, found := sess.server.Driver.(Auth); found {
		auth = driverAuth
	}
	var ctx = Context{
		Sess:  sess,
		Cmd:   "PASS",
		Param: param,
		Data:  make(map[string]interface{}),
	}
	ok, err := auth.CheckPasswd(&ctx, sess.reqUser, param)
	sess.server.filters.AfterUserLogin(&ctx, sess.reqUser, param, ok, err)
	if err != nil {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "Checking password error")
		return
	}

	if ok {
		sess.user = sess.reqUser
		sess.reqUser = ""
		sess.writeMessage(CODE_USER_LOGINED, "Password ok, continue")
	} else {
		sess.writeMessage(CODE_UNLOGIN, "Incorrect password, not logged in")
	}
}

// CommandPasv responds to the PASV FTP Command.
//
// The client is requesting us to open a new TCP listing socket and wait for them
// to connect to it.
type CommandPasv struct{}

func (cmd CommandPasv) IsExtend() bool {
	return false
}

func (cmd CommandPasv) RequireParam() bool {
	return false
}

func (cmd CommandPasv) RequireAuth() bool {
	return true
}

func (cmd CommandPasv) Execute(sess *Session, param string) {
	listenIP := sess.passiveListenIP()
	// TODO: IPv6 for this Command is not implemented
	if strings.HasPrefix(listenIP, "::") {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "Action not taken")
		return
	}

	socket, err := sess.newPassiveSocket()
	if err != nil {
		sess.writeMessage(CODE_FAILED_OPEN_DATA_CONN, "Data connection failed")
		return
	}

	p1 := socket.Port() / 256
	p2 := socket.Port() - (p1 * 256)

	quads := strings.Split(listenIP, ".")
	target := fmt.Sprintf("(%s,%s,%s,%s,%d,%d)", quads[0], quads[1], quads[2], quads[3], p1, p2)
	msg := "Entering Passive Mode " + target
	sess.writeMessage(CODE_ENTER_PASV, msg)
}

type CommandPbsz struct{}

func (cmd CommandPbsz) IsExtend() bool {
	return false
}

func (cmd CommandPbsz) RequireParam() bool {
	return true
}

func (cmd CommandPbsz) RequireAuth() bool {
	return false
}

func (cmd CommandPbsz) Execute(sess *Session, param string) {
	if sess.tls && param == "0" {
		sess.writeMessage(CODE_COMMAND_OK, "OK")
	} else {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "Action not taken")
	}
}

// CommandPort responds to the PORT FTP Command.
//
// The client has opened a listening socket for sending out of band data and
// is requesting that we connect to it
type CommandPort struct{}

func (cmd CommandPort) IsExtend() bool {
	return false
}

func (cmd CommandPort) RequireParam() bool {
	return true
}

func (cmd CommandPort) RequireAuth() bool {
	return true
}

func (cmd CommandPort) Execute(sess *Session, param string) {
	nums := strings.Split(param, ",")
	portOne, _ := strconv.Atoi(nums[4])
	portTwo, _ := strconv.Atoi(nums[5])
	port := (portOne * 256) + portTwo
	host := nums[0] + "." + nums[1] + "." + nums[2] + "." + nums[3]
	socket, err := newActiveSocket(sess, host, port)
	if err != nil {
		sess.writeMessage(CODE_FAILED_OPEN_DATA_CONN, "Data connection failed")
		return
	}
	sess.dataConn = socket
	sess.writeMessage(CODE_COMMAND_OK, "Connection established ("+strconv.Itoa(port)+")")
}

type CommandProt struct{}

func (cmd CommandProt) IsExtend() bool {
	return false
}

func (cmd CommandProt) RequireParam() bool {
	return true
}

func (cmd CommandProt) RequireAuth() bool {
	return false
}

func (cmd CommandProt) Execute(sess *Session, param string) {
	if sess.tls && param == "P" {
		sess.writeMessage(CODE_COMMAND_OK, "OK")
	} else if sess.tls {
		sess.writeMessage(536, "Only P level is supported")
	} else {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "Action not taken")
	}
}

// CommandPwd responds to the PWD FTP Command.
//
// Tells the client what the current working directory is.
type CommandPwd struct{}

func (cmd CommandPwd) IsExtend() bool {
	return false
}

func (cmd CommandPwd) RequireParam() bool {
	return false
}

func (cmd CommandPwd) RequireAuth() bool {
	return true
}

func (cmd CommandPwd) Execute(sess *Session, param string) {
	sess.writeMessage(CODE_PATHNAME_CREATED, "\""+sess.curDir+"\" is the current directory")
}

// CommandQuit responds to the QUIT FTP Command. The client has requested the
// connection be closed.
type CommandQuit struct{}

func (cmd CommandQuit) IsExtend() bool {
	return false
}

func (cmd CommandQuit) RequireParam() bool {
	return false
}

func (cmd CommandQuit) RequireAuth() bool {
	return false
}

func (cmd CommandQuit) Execute(sess *Session, param string) {
	sess.writeMessage(CODE_CLOSE_CONN, "Goodbye")
	sess.Close()
}

// CommandRetr responds to the RETR FTP Command. It allows the client to
// download a file.
// REST can be followed by APPE, STOR, or RETR
type CommandRetr struct{}

func (cmd CommandRetr) IsExtend() bool {
	return false
}

func (cmd CommandRetr) RequireParam() bool {
	return true
}

func (cmd CommandRetr) RequireAuth() bool {
	return true
}

func (cmd CommandRetr) Execute(sess *Session, param string) {
	path := sess.buildPath(param)
	if sess.preCommand != "REST" {
		sess.lastFilePos = -1
	}
	defer func() {
		sess.lastFilePos = -1
	}()
	var ctx = Context{
		Sess:  sess,
		Cmd:   "RETR",
		Param: param,
		Data:  make(map[string]interface{}),
	}
	if !sess.server.Options.Auth.IsReadable(&ctx, sess.user) {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("No auth for read file."))
		return
	}
	if !sess.server.filters.BeforeDownloadFile(&ctx, path) {
		sess.writeMessage(CODE_ACTION_ABORTED, "This command refused.")
		return
	}
	var readPos = sess.lastFilePos
	if readPos < 0 {
		readPos = 0
	}
	size, data, err := sess.server.Driver.GetFile(&ctx, path, readPos)
	if err == nil {
		defer data.Close()
		sess.writeMessage(CODE_FILE_STATUS_OK, fmt.Sprintf("Data transfer starting %d bytes", size))
		err = sess.sendOutofBandDataWriter(data)
		sess.server.filters.AfterFileDownloaded(&ctx, path, size, err)
		if err != nil {
			sess.writeMessage(CODE_ACTION_ABORTED, "Error reading file")
		}
	} else {
		sess.server.filters.AfterFileDownloaded(&ctx, path, size, err)
		sess.writeMessage(CODE_ACTION_ABORTED, "File not available")
	}
}

type CommandRest struct{}

func (cmd CommandRest) IsExtend() bool {
	return false
}

func (cmd CommandRest) RequireParam() bool {
	return true
}

func (cmd CommandRest) RequireAuth() bool {
	return true
}

func (cmd CommandRest) Execute(sess *Session, param string) {
	var err error
	sess.lastFilePos, err = strconv.ParseInt(param, 10, 64)
	if err != nil {
		sess.writeMessage(CODE_ACTION_ABORTED, "File not available")
		return
	}

	sess.writeMessage(350, fmt.Sprint("Start transfer from ", sess.lastFilePos))
}

// CommandRnfr responds to the RNFR FTP Command. It's the first of two Commands
// required for a client to rename a file.
type CommandRnfr struct{}

func (cmd CommandRnfr) IsExtend() bool {
	return false
}

func (cmd CommandRnfr) RequireParam() bool {
	return true
}

func (cmd CommandRnfr) RequireAuth() bool {
	return true
}

func (cmd CommandRnfr) Execute(sess *Session, param string) {
	sess.renameFrom = ""
	p := sess.buildPath(param)
	if _, err := sess.server.Driver.Stat(&Context{
		Sess:  sess,
		Cmd:   "RNFR",
		Param: param,
		Data:  make(map[string]interface{}),
	}, p); err != nil {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("Action not taken: ", err))
		return
	}
	sess.renameFrom = p
	sess.writeMessage(350, "Requested file action pending further information.")
}

// cmdRnto responds to the RNTO FTP Command. It's the second of two Commands
// required for a client to rename a file.
type CommandRnto struct{}

func (cmd CommandRnto) IsExtend() bool {
	return false
}

func (cmd CommandRnto) RequireParam() bool {
	return true
}

func (cmd CommandRnto) RequireAuth() bool {
	return true
}

func (cmd CommandRnto) Execute(sess *Session, param string) {
	toPath := sess.buildPath(param)
	err := sess.server.Driver.Rename(&Context{
		Sess:  sess,
		Cmd:   "RNTO",
		Param: param,
		Data:  make(map[string]interface{}),
	}, sess.renameFrom, toPath)
	defer func() {
		sess.renameFrom = ""
	}()

	if err == nil {
		sess.writeMessage(250, "File renamed")
	} else {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("Action not taken: ", err))
	}
}

// cmdRmd responds to the RMD FTP Command. It allows the client to delete a
// directory.
type CommandRmd struct{}

func (cmd CommandRmd) IsExtend() bool {
	return false
}

func (cmd CommandRmd) RequireParam() bool {
	return true
}

func (cmd CommandRmd) RequireAuth() bool {
	return true
}

func (cmd CommandRmd) Execute(sess *Session, param string) {
	executeRmd("RMD", sess, param)
}

// cmdXRmd responds to the RMD FTP Command. It allows the client to delete a
// directory.
type CommandXRmd struct{}

func (cmd CommandXRmd) IsExtend() bool {
	return false
}

func (cmd CommandXRmd) RequireParam() bool {
	return true
}

func (cmd CommandXRmd) RequireAuth() bool {
	return true
}

func (cmd CommandXRmd) Execute(sess *Session, param string) {
	executeRmd("XRMD", sess, param)
}

func executeRmd(cmd string, sess *Session, param string) {
	p := sess.buildPath(param)
	var ctx = Context{
		Sess:  sess,
		Cmd:   cmd,
		Param: param,
		Data:  make(map[string]interface{}),
	}
	if param == "/" || param == "" {
		sess.writeMessage(CODE_ACTION_NOTAKEN, "Directory / cannot be deleted")
		return
	}

	var needChangeCurDir = strings.HasPrefix(param, sess.curDir)
	if !sess.server.Options.Auth.IsDeleable(&ctx, sess.user) {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("No auth for delte dir."))
		return
	}
	if !sess.server.filters.BeforeDeleteDir(&ctx, p) {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("This command refused."))
		return
	}
	err := sess.server.Driver.DeleteDir(&ctx, p)
	if needChangeCurDir {
		sess.curDir = path.Dir(param)
	}
	sess.server.filters.AfterDirDeleted(&ctx, p, err)
	if err == nil {
		sess.writeMessage(250, "Directory deleted")
	} else {
		sess.writeMessage(CODE_ACTION_NOTAKEN, fmt.Sprint("Directory delete failed: ", err))
	}
}

type CommandConf struct{}

func (cmd CommandConf) IsExtend() bool {
	return false
}

func (cmd CommandConf) RequireParam() bool {
	return true
}

func (cmd CommandConf) RequireAuth() bool {
	return true
}

func (cmd CommandConf) Execute(sess *Session, param string) {
	sess.writeMessage(CODE_ACTION_NOTAKEN, "Action not taken")
}

// CommandSize responds to the SIZE FTP Command. It returns the size of the
// requested path in bytes.
type CommandSize struct{}

func (cmd CommandSize) IsExtend() bool {
	return false
}

func (cmd CommandSize) RequireParam() bool {
	return true
}

func (cmd CommandSize) RequireAuth() bool {
	return true
}

func (cmd CommandSize) Execute(sess *Session, param string) {
	path := sess.buildPath(param)
	stat, err := sess.server.Driver.Stat(&Context{
		Sess:  sess,
		Cmd:   "SIZE",
		Param: param,
		Data:  make(map[string]interface{}),
	}, path)
	if err != nil {
		log.Printf("Size: error(%s)", err)
		sess.writeMessage(450, fmt.Sprintf("path %s not found", param))
	} else {
		sess.writeMessage(213, strconv.Itoa(int(stat.Size())))
	}
}

// CommandStat responds to the STAT FTP Command. It returns the stat of the
// requested path.
type CommandStat struct{}

func (cmd CommandStat) IsExtend() bool {
	return false
}

func (cmd CommandStat) RequireParam() bool {
	return false
}

func (cmd CommandStat) RequireAuth() bool {
	return true
}

func (cmd CommandStat) Execute(sess *Session, param string) {
	// system stat

	if param == "" {
		sess.writeMessage(211, fmt.Sprintf("%s FTP server status:\nVersion %s"+
			"Connected to %s (%s)\n"+
			"Logged in %s\n"+
			"TYPE: ASCII, FORM: Nonprint; STRUcture: File; transfer MODE: Stream\n"+
			"No data connection", sess.PublicIP(), version, sess.PublicIP(),
			version, sess.LoginUser()))
		sess.writeMessage(211, "End of status")
		return
	}

	var ctx = Context{
		Sess:  sess,
		Cmd:   "STAT",
		Param: param,
		Data:  make(map[string]interface{}),
	}

	// file or directory stat
	path := sess.buildPath(param)
	stat, err := sess.server.Driver.Stat(&ctx, path)
	if err != nil {
		log.Printf("Size: error(%s)", err)
		sess.writeMessage(450, fmt.Sprintf("path %s not found", path))
	} else {
		var files []FileInfo
		if stat.IsDir() {
			err = sess.server.Driver.ListDir(&ctx, path, func(f os.FileInfo) error {
				info, err := convertFileInfo(sess, f, filepath.Join(path, f.Name()))
				if err != nil {
					return err
				}
				files = append(files, info)
				return nil
			})
			if err != nil {
				sess.writeMessage(CODE_ACTION_NOTAKEN, err.Error())
				return
			}
			sess.writeMessage(213, "Opening ASCII mode data connection for file list")
		} else {
			info, err := convertFileInfo(sess, stat, path)
			if err != nil {
				sess.writeMessage(CODE_ACTION_NOTAKEN, err.Error())
				return
			}
			files = append(files, info)
			sess.writeMessage(212, "Opening ASCII mode data connection for file list")
		}
		sess.sendOutofbandData(listFormatter(files).Detailed())
	}
}

// CommandStor responds to the STOR FTP Command. It allows the user to upload a
// new file.
type CommandStor struct{}

func (cmd CommandStor) IsExtend() bool {
	return false
}

func (cmd CommandStor) RequireParam() bool {
	return true
}

func (cmd CommandStor) RequireAuth() bool {
	return true
}

func (cmd CommandStor) Execute(sess *Session, param string) {
	targetPath := sess.buildPath(param)
	sess.writeMessage(CODE_FILE_STATUS_OK, "Data transfer starting")

	if sess.preCommand != "REST" {
		sess.lastFilePos = -1
	}

	defer func() {
		sess.lastFilePos = -1
	}()

	var ctx = Context{
		Sess:  sess,
		Cmd:   "STOR",
		Param: param,
		Data:  make(map[string]interface{}),
	}
	if !sess.server.Options.Auth.IsPutable(&ctx, sess.user) {
		sess.writeMessage(450, fmt.Sprint("No auth for put file."))
		return
	}
	if !sess.server.filters.BeforePutFile(&ctx, targetPath) {
		sess.writeMessage(450, fmt.Sprint("This command refused."))
		return
	}
	size, err := sess.server.Driver.PutFile(&ctx, targetPath, sess.dataConn, sess.lastFilePos)
	sess.server.filters.AfterFilePut(&ctx, targetPath, size, err)
	if err == nil {
		msg := fmt.Sprintf("OK, received %d bytes", size)
		sess.writeMessage(226, msg)
	} else {
		sess.writeMessage(450, fmt.Sprint("error during transfer: ", err))
	}
}

// CommandStru responds to the STRU FTP Command.
//
// like the MODE and TYPE Commands, stru[cture] dates back to a time when the
// FTP protocol was more aware of the content of the files it was transferring,
// and would sometimes be expected to translate things like EOL markers on the
// fly.
//
// These days files are sent unmodified, and F(ile) mode is the only one we
// really need to support.
type CommandStru struct{}

func (cmd CommandStru) IsExtend() bool {
	return false
}

func (cmd CommandStru) RequireParam() bool {
	return true
}

func (cmd CommandStru) RequireAuth() bool {
	return true
}

func (cmd CommandStru) Execute(sess *Session, param string) {

	if strings.ToUpper(param) == "F" {
		sess.writeMessage(200, "OK")
	} else {
		sess.writeMessage(504, "STRU is an obsolete Command")
	}
}

// CommandSyst responds to the SYST FTP Command by providing a canned response.
type CommandSyst struct{}

func (cmd CommandSyst) IsExtend() bool {
	return false
}

func (cmd CommandSyst) RequireParam() bool {
	return false
}

func (cmd CommandSyst) RequireAuth() bool {
	return true
}

func (cmd CommandSyst) Execute(sess *Session, param string) {
	sess.writeMessage(215, "UNIX Type: L8")
}

// CommandType responds to the TYPE FTP Command.
//
//  like the MODE and STRU Commands, TYPE dates back to a time when the FTP
//  protocol was more aware of the content of the files it was transferring, and
//  would sometimes be expected to translate things like EOL markers on the fly.
//
//  Valid options were A(SCII), I(mage), E(BCDIC) or LN (for local type). Since
//  we plan to just accept bytes from the client unchanged, I think Image mode is
//  adequate. The RFC requires we accept ASCII mode however, so accept it, but
//  ignore it.
type CommandType struct{}

func (cmd CommandType) IsExtend() bool {
	return false
}

func (cmd CommandType) RequireParam() bool {
	return false
}

func (cmd CommandType) RequireAuth() bool {
	return true
}

func (cmd CommandType) Execute(sess *Session, param string) {
	if strings.ToUpper(param) == "A" {
		sess.writeMessage(200, "Type set to ASCII")
	} else if strings.ToUpper(param) == "I" {
		sess.writeMessage(200, "Type set to binary")
	} else {
		sess.writeMessage(500, "Invalid type")
	}
}

// CommandUser responds to the USER FTP Command by asking for the password
type CommandUser struct{}

func (cmd CommandUser) IsExtend() bool {
	return false
}

func (cmd CommandUser) RequireParam() bool {
	return true
}

func (cmd CommandUser) RequireAuth() bool {
	return false
}

func (cmd CommandUser) Execute(sess *Session, param string) {
	sess.reqUser = param
	sess.server.filters.BeforeLoginUser(&Context{
		Sess:  sess,
		Cmd:   "USER",
		Param: param,
		Data:  make(map[string]interface{}),
	}, sess.reqUser)
	sess.writeMessage(331, "User name ok, password required")
}
