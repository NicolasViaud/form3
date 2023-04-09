package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	secret "github.com/form3tech/innsecure/pkg"

	"github.com/form3tech/innsecure"
	"github.com/form3tech/innsecure/jwtauth"
	"github.com/form3tech/innsecure/postgres"
	"github.com/go-kit/kit/log"
	_ "github.com/lib/pq"
)

func main() {
	var (
		httpAddr = flag.String("http.addr", ":8080", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewJSONLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	ctx := context.Background()

	vaultSecretProvider, _ := strconv.ParseBool(os.Getenv("VAULT_SECRET_PROVIDER"))
	vaultRoleName := os.Getenv("VAULT_ROLE_NAME") //Used only if VaultSecretProvider is true.
	secretProviderConfiguration := secret.Configuration{
		VaultSecretProvider: vaultSecretProvider,
		VaultRoleName:       vaultRoleName,
	}
	secretProvider, err := secret.NewSecretProvider(ctx, secretProviderConfiguration)
	if err != nil {
		panic(err)
	}

	var s innsecure.Service
	{
		host := os.Getenv("DB_HOST")
		databaseCredentials, err := secretProvider.GetDatabaseCredentials(ctx)
		if err != nil {
			panic(err)
		}
		db, err := postgres.NewConnection(host, databaseCredentials.Username, databaseCredentials.Password)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		r := postgres.NewRepo(db)
		s = innsecure.NewBookingService(r)
	}

	var h http.Handler
	{
		hs256Secret, err := secretProvider.GetHS256Secret(ctx)
		if err != nil {
			panic(err)
		}
		jwtmw := jwtauth.NewMiddleware(hs256Secret)
		e := innsecure.MakeServerEndpoints(s, jwtmw)
		h = innsecure.MakeHTTPHandler(e, log.With(logger, "component", "HTTP"))
	}

	errs := make(chan error)

	// Shutdown handler
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	// HTTP Transport
	go func() {
		logger.Log("transport", "HTTP", "addr", *httpAddr)
		errs <- http.ListenAndServe(*httpAddr, h)
	}()

	logger.Log("exit", <-errs)
}
