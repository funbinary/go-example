package testing

import (
	"fmt"
	"testing"
	"time"
)

func TestTesting(t *testing.T) {
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

func TestFail(t *testing.T) {
	go func() {
		for {
			t.FailNow()
			fmt.Println("sub goroutine running")
		}
	}()

	for {
		t.Log("main goroutine running")
		time.Sleep(2 * time.Second)
		break
	}
}

func TestSkip(t *testing.T) {
	t.Log("1")
	t.Skip("skip")
	t.Log("2")
}

func TestSkiped(t *testing.T) {
	t.Log("1")
	t.Skipped()
	t.Log("2")
}

func TestSkipNow(t *testing.T) {
	t.Log("1")
	t.SkipNow()
	t.Log("2")
}

func TestName(t *testing.T) {
	t.Log("1")
	t.Log(t.Name())
	t.Log("2")
}

func Test(t *testing.T) {
	t.Log("1")

}
