package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	flags "github.com/jessevdk/go-flags"
	_ "github.com/lib/pq"

	"github.com/xanderflood/database/cmd/api/router"
	"github.com/xanderflood/database/lib/tools"
	"github.com/xanderflood/database/lib/web"
	"github.com/xanderflood/database/pkg/dbi"
)

var opts struct {
	LogLevel   string `long:"log-level" env:"LOG_LEVEL" default:"INFO"`
	DBString   string `long:"db-string" env:"DB_STRING" default:"database"`
	PublicRoot string `long:"public-root" env:"PUBLIC_ROOT" required:"true"`
}

func main() {
	fmt.Println("Starting database server")

	logger := tools.NewStdoutLogger(tools.LogLevelDebug, "database")

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	dynamoClient := dynamodb.New(session.New(&aws.Config{}))

	server := &router.Server{
		DB:   &dbi.Database{Client: dynamoClient},
		Vars: web.MuxVars{},
	}

	router := router.New(
		server,
		logger,
		opts.PublicRoot,
	)

	s := &http.Server{
		Addr:           ":8080",
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
