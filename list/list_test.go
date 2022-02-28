package list

import (
	"testing"
)

func TestList(t *testing.T) {
	t.Log("log1")
	t.Log(t.Failed())
	t.Error("error")
	t.Error("error2")
	t.Log(t.Failed())
	t.Fail()
	t.Log(t.Failed())
	t.Fatal("fatal")
	t.Fatal("fatal2")
}
