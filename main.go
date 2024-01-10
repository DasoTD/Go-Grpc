package main

import (
	"context"
	"net"
	"os"

	"github.com/dasotd/go_grpc/api"
	db "github.com/dasotd/go_grpc/db/sqlc"
	pb "github.com/dasotd/go_grpc/pb"
	"github.com/dasotd/go_grpc/util"
	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)
func main() {
    config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config")
	}
	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}

	runDBMigration(config.MigrationURL, config.DBSource)

	grpcDB := db.NewGrpc(connPool)

    lis, err := net.Listen("tcp", "[::1]:8080")
    if err != nil {
        log.Fatal().Err(err).Msg("failed to listen")
    }

    grpcServer := grpc.NewServer()
    api := &api.Server(grpcDB) 

    pb.RegisterAccountAPIServer(grpcServer, api)
    err = grpcServer.Serve(lis)
    reflection.Register(grpcServer)

    if err != nil {
        log.Fatal().Err(err).Msg("Error strating server:")
    }
}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}