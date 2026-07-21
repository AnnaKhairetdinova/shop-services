package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/AnnaKhairetdinova/user-service/config"
	pb "github.com/AnnaKhairetdinova/user-service/docs/proto/user"
	"github.com/AnnaKhairetdinova/user-service/internal/app/service"
	"github.com/AnnaKhairetdinova/user-service/internal/infrastructure/repository"
	"github.com/AnnaKhairetdinova/user-service/internal/presentation/handler"
	"github.com/AnnaKhairetdinova/user-service/internal/presentation/interceptor"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := dbConnect(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("БД: %v", err)
	}
	defer pool.Close()

	if err := runMigrations(cfg.DBUrl); err != nil {
		log.Fatalf("миграции: %v", err)
	}

	userRepository := repository.NewUserRepository(pool)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewGRPCHandler(userService)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptor.Logging,
			interceptor.Recovery,
		),
	)

	pb.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	lis, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	go grpcServer.Serve(lis)

	<-ctx.Done()
	grpcServer.GracefulStop()
}

func dbConnect(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, dsn)
	pool.Ping(ctx)
	if err != nil {
		pool.Close()
		fmt.Errorf("ошибка: %w", err)
	}

	return pool, nil
}

func runMigrations(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	migrationsPath := "file://migrations"

	m, err := migrate.NewWithDatabaseInstance(
		migrationsPath,
		"postgres",
		driver,
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("миграции уже накатаны")
		} else {
			return err
		}
	}

	return err
}
