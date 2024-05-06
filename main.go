package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"

	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"github.com/abhinavrobinson/Go-Microservice/account"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/log/level"
)

const dbsource = "user='koyeb-adm' password=3JihHjqwX6xT host=ep-summer-dawn-a27pf9in.eu-central-1.pg.koyeb.app dbname='koyebdb'"

func main() {
	var httpAddr = flag.String("http", ":8080", "http listen address")
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "account",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")

	defer level.Info(logger).Log("msg", "service ended")

	var db *sql.DB
	{
		var err error
		db, err = sql.Open("postgres", dbsource)
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	flag.Parse()
	ctx := context.Background()

	var srv account.Service
	{
		repository := account.NewRepo(db, logger)
		srv = account.NewService(repository, logger)
	}

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	endpoints := account.MakeEndpoints(srv)

	go func() {
		fmt.Println("listening on port", *httpAddr)
		handler := account.NewHTTPServer(ctx, endpoints)
		errs <- http.ListenAndServe(*httpAddr, handler)
	}()

	level.Error(logger).Log("exit", <-errs)
}
