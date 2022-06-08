// +ignore

package main

import (
	"log"

	"github.com/bin-work/go-example/pkg/bftp"
	"github.com/bin-work/go-example/pkg/bftp/driver/file"
)

func main() {
	driver, err := file.NewDriver("./")
	if err != nil {
		log.Fatal(err)
	}

	s, err := bftp.NewServer(&bftp.Options{
		Driver: driver,
		Auth: &bftp.MultiAuth{
			Accounts: make(map[string]bftp.Account),
		},
		Port: 21,
		Perm: bftp.NewSimplePerm("root", "root"),
	})
	s.Auth.Register(bftp.Account{
		Name:     "admin",
		Password: "admin",
		Putable:  true,
		Readable: true,
		Deleable: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
