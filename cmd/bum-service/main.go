package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
	_ "net/http/pprof"

	"bum-service/config"
)

//nolint:gochecknoglobals // in order to use preRun command we need this global variable
var cfg *config.Config

func main() {

	//// TODO: убрать т.к это временная штука для профайлинга
	//go func() {
	//	log.Fatal(http.ListenAndServe(":5555", nil))
	//}()

	if err := rootCmd().Execute(); err != nil {
		log.Fatalf("rootCmd().Execute() error=%v", err)
	}
}
