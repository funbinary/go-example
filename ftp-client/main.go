package main

import (
	"bytes"
	"fmt"
	. "github.com/jlaffaye/ftp"
)

var conn *ServerConn

func main() {
	// 连接
	//conn, err := Dial("localhost:21", DialWithDebugOutput(os.Stdout))
	var err error

	conn, err = Dial("localhost:21")
	if err != nil {
		panic(err)
	}
	// 登录
	if err = conn.Login("admin", "admin"); err != nil {
		panic(err)
	}
	curDir, err := conn.CurrentDir()
	fmt.Println("curDir:", curDir)
	fmt.Println("gettime is support:", conn.IsGetTimeSupported())
	fmt.Println("settime is support:", conn.IsSetTimeSupported())
	fmt.Println("mlst is support:", conn.IsTimePreciseInList())
	// 读取文件
	//namelist(curDir)
	//list(curDir)
	walk(curDir)
	//propery(curDir)
	//retr(curDir + "/ui/config.toml")
	//stor(curDir + "config.txt")
	//time.Sleep(2 * time.Second)

	// 退出
	if err = conn.Quit(); err != nil {
		panic(err)
	}
}

func stor(dir string) {
	data := bytes.NewBufferString("Hello aaa")

	if err := conn.Stor(dir, data); err != nil {
		fmt.Println(err)
	}

}

func retr(dir string) {
	r, err := conn.Retr(dir)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	//buf, err := ioutil.ReadAll(r)
	//fmt.Println(string(buf))

}

func propery(dir string) {
	fmt.Println(conn.FileSize(dir))
	fmt.Println(conn.GetTime(dir))

}

func walk(dir string) {
	w := conn.Walk(dir)
	//fmt.Println(w.Path())
	for w.Next() {
		fmt.Println(w.Path(), w.Stat().Size)
	}
}

func namelist(dir string) {
	fmt.Println("列出所有文件信息：")
	entries, err := conn.NameList(dir)
	if err != nil {
		panic(err)
	}
	fmt.Println("============")
	for _, v := range entries {
		fmt.Println(v)
		fmt.Println("============")
	}

}

func list(dir string) {
	fmt.Println("列出所有文件信息：")
	entries, err := conn.List(dir)
	if err != nil {
		panic(err)
	}
	fmt.Println("============")
	for _, v := range entries {
		fmt.Println(v.Name)
		if v.Type == EntryTypeFile {
			fmt.Println("文件")
		} else if v.Type == EntryTypeLink {
			fmt.Println("连接类型")
		} else if v.Type == EntryTypeFolder {
			fmt.Println("目录")
		}
		fmt.Println(v.Time) //最后修改时间
		fmt.Println(v.Size) // 大小
		fmt.Println(v.Target)
		fmt.Println("============")
	}

}
