package handlers

import (
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func Checkerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
