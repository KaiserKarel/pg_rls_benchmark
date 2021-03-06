package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"testing"
)

func fatal(msg string, args ...interface{})  {
	panic(fmt.Sprintf(msg, args...))
}

type Option func(ctx context.Context, config BenchConfig, tx *sql.Tx) error

func Benchmark(b testing.B, cfg BenchConfig,  options ...Option)  {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=karel "+
		"password=dwleml123 dbname=rls_bench sslmode=disable")
	if err != nil {
		log.Fatalf("unable to open db: %v", err)
	}

	err = Initialize(db)
	if err != nil {
		fatal("unable to initialize db: %v", err)
	}

	tx, err := db.BeginTx(context.Background(), nil)
	if err != nil {
		fatal("unable to initialize tx: %v", err)
	}

	ctx := context.Background()

	for _, opt := range options {
		err = opt(ctx, cfg, tx)
		if err != nil {
			fatal("unable to run option: %v", err)
		}
	}
	defer tx.Commit()
}

func main()  {

	// shortConf :=  BenchConfig{
	// 	NumUsers: 10,
	// 	NumGroups: 1,
	// 	NumObjects: 100,
	// 	AvgGroupSize: 1,
	// }

	bigConf :=  BenchConfig{
		NumUsers: 10000000,
		NumGroups: 200000,
		NumObjects: 100000000,
		AvgGroupSize: 100,
	}

	Benchmark(testing.B{}, bigConf,
		PopulateUsers,
		PopulateGroups,
		StitchUsersToGroups,
		PopulateObjects,
		GenerateRandomUserPermissions,
		GenerateRandomGroupPermissions,
		)
}
