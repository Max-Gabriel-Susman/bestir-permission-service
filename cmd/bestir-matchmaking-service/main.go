package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/foundation/database"
	"github.com/Max-Gabriel-Susman/bestir-permissionmaking-service/internal/handler"
	env "github.com/caarlos0/env/v6"
	"github.com/pkg/errors"

	"go.uber.org/zap"
)

var (
	GitSHA = "~git~" // populated with ldflags at build time(wtf is an ldflag?)

	// File ssm-params.yml contains precomputed data
	// procedure is documented @ https://blog.carlmjohnson.net/post/2021/how-to-use-go-embed/
	//--+-go:embed ssm-params.yml(not implemented yet)
	ssmParams []byte
)

const (
	exitCodeErr       = 1
	exitCodeInterrupt = 2
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	defer func() {
		signal.Stop(signalChan)
		cancel()
	}()
	go func() {
		select {
		case <-signalChan: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-signalChan // second signal, hard exit
		os.Exit(exitCodeInterrupt)
	}()
	if err := run(ctx, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(exitCodeErr)
	}
}

func run(ctx context.Context, _ []string) error {
	// aws shit, currently unsupported, but soon
	//wsCfg, err := aws.NewConfig(ctx)
	//f err != nil {
	//	return errors.Wrap(err, "could not create aws sdk config")
	//

	//f _, ok := os.LookupEnv("SSM_DISABLE"); !ok {
	//	if err := awsParseSSMParams(ctx, awsCfg, ssmParams); err != nil {
	//		return err
	//	}
	//

	// open api shit for documentation

	// cfg and setup shit right hurr, we gotta alter it for my database setup
	var cfg struct {
		ServiceName string `env:"SERVICE_NAME" envDefault:"bestir-permissionmaking-service"`
		Env         string `env:"ENV" envDefault:"local"`
		Database    struct {
			User   string `env:"permission_DB_USER,required"`
			Pass   string `env:"permission_DB_PASSWORD,required"`
			Host   string `env:"permission_DB_HOST"`
			Port   string `env:"permission_DB_PORT" envDefault:"3306"`
			DBName string `env:"permission_DB_Name" envDefault:"permission"`
			Params string `env:"permission_DB_Param_Overrides" envDefault:"parseTime=true"`
		}
		Datadog struct {
			Disable bool `env:"DD_DISABLE"`
		}
		Migration struct {
			Enable bool `env:"ENABLE_MIGRATE"`
		}
	}
	if err := env.Parse(&cfg); err != nil {
		return errors.Wrap(err, "parsing configuration")
	}
	// cfg.Datadog.Disable = true

	// Create base logger
	z, err := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:      false,
		Encoding:         "json",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}.Build(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		return errors.Wrap(err, "initializing zap logger")
	}

	z = z.With(
		zap.String("service", cfg.ServiceName),
		zap.String("version.git_sha", GitSHA),
		zap.String("env", cfg.Env),
	)
	// zl := bestirlog.WrapZap(z)

	// Intitialize tracing

	/*
		if !cfg.Datadog.Disable {
			// Configure tracing
		}
	*/

	// Migrate - issa broken Luigi
	// if err := goose.EnsureMigrations(ctx, zl, goose.Config{
	// 	User:     cfg.Database.User,
	// 	Password: cfg.Database.Pass,
	// 	Port:     cfg.Database.Port,
	// 	Host:     cfg.Database.Host,
	// 	DryRun:   !cfg.Migration.Enable,
	// 	Name:     cfg.Database.DBName,
	// }); err != nil {
	// 	return err
	// }

	// dsn: usr:identity@tcp(127.0.0.1:3306)/identity
	db, err := database.Open(database.Config{
		User:     cfg.Database.User,
		Password: cfg.Database.Pass,
		Host:     cfg.Database.Host,
		Name:     cfg.Database.DBName,
		Params:   cfg.Database.Params,
	}, cfg.ServiceName)
	if err != nil {
		return errors.Wrap(err, "connecting to db")
	}
	defer func() {
		// zl.Info(ctx, "stopping database")
		db.Close()
	}()

	// func dsn(dbName string) string {
	// 	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
	// }

	// db, err := sql.Open("mysql", dsn(""))
	/*
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.Database.User, cfg.Database.Pass, cfg.Database.Host, ""))
		if err != nil {
			log.Printf("Error %s when opening DB\n", err)
			return err
		}
		defer db.Close()
	*/
	// needs cleaner implementation, logic should be moved elsewhere
	// query := `CREATE TABLE IF NOT EXISTS permission (
	// 		id CHAR(36) NOT NULL,
	// 		[name] VARCHAR(255) NOT NULL,
	// 		PRIMARY KEY (id)
	// 	) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;`
	//
	// ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancelfunc()
	// // _, err = db.ExecContext(ctx, query)
	// _, err = db.Exec(query)
	// if err != nil {
	// 	log.Printf("Error %s when creating product table", err)
	// 	return err
	// }
	// log.Println("successfuly created permission table")
	//*/

	// If DD is enabled, configure db to send stats info
	//if !cfg.Datadog.Disable {
	// statsdAddress := ddAgentAddress(8125)
	// statsd, err := statsd.New(statsdAddress, statsd.WithMaxBytesPerPayload(4096))
	// if err != nil {
	// 	return errors.Wrap(err, "could not start statsd client")
	// }
	// // Start Stats Reporting for db
	// go func() {
	// 	zl.Info(ctx, "Starting reporting DB metrics", zap.String("statsd.address", statsdAddress))
	// 	defer zl.Info(ctx, "Stopped reporting DB metrics")
	// 	sr := database.NewStatsReporter(db, statsd, zl)

	// 	sr.ReportDBStats(ctx, []string{
	// 		fmt.Sprintf("service:%s", cfg.ServiceName),
	// 		fmt.Sprintf("version:%s", GitSHA),
	// 		fmt.Sprintf("env:%s", cfg.Env),
	// 		"collabs.squad:red",
	// 	}, 1)
	// }()
	//}

	// CLIENTS N SHIT

	// event bridge shit

	// we gott reconfigure the service to use pgx now
	h := handler.API(handler.Deps{DB: db})

	// Start API Service
	api := http.Server{
		Handler: h,
		// Addr:              "127.0.0.1:80",
		Addr:              "0.0.0.0:80",
		ReadHeaderTimeout: 2 * time.Second,
	}

	// Make a channel to listen for errors coming from the listener
	serverErrors := make(chan error, 1)

	// Start listening for requests
	go func() {
		// log info about this
		serverErrors <- api.ListenAndServe()
	}()
	// Shutdown

	// logic for handling shutdown gracefully
	select {
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	case <-ctx.Done():
		// log something

		// request a deadline for completion
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()

		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
