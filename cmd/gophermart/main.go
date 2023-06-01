package main

import (
	"flag"
	"github.com/riyce/gophermart/internal"
	"log"
	"os"
)

func main() {
	s := new(internal.Server)

	conf := internal.NewConfig()

	addr := flag.String("a",
		"127.0.0.1:8000",
		"Server address")

	accrualAddress := flag.String("r",
		"http://127.0.0.1:8080",
		"Accrual service address")

	dns := flag.String("d",
		"postgres://postgres:postgres@localhost:5432/gophersmart?sslmode=disable",
		"Database DNS")

	flag.Parse()

	if os.Getenv("RUN_ADDRESS") == "" {
		conf.Address = *addr
	}

	if os.Getenv("DATABASE_URI") == "" {
		conf.DBURI = *dns
	}

	if os.Getenv("ACCRUAL_SYSTEM_ADDRESS") == "" {
		conf.AccrualAddress = *accrualAddress
	}

	if err := s.Run(conf); err != nil {
		log.Fatal(err)
	}
}
