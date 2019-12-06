package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v4"
	"log"
	"math/rand"
)

func PopulateUsers(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("populating users")
	stmt, err := tx.Prepare("INSERT INTO users (username) VALUES ($1)")
	if err != nil {
		return err
	}

	checkpoint := int64(0)
	for i := int64(0); i < config.NumUsers; i++ {
		fmt.Printf("\r[%d/%d]", i, config.NumUsers)
		_, err := stmt.ExecContext(ctx, gofakeit.Username())
		if err != nil {
			return  err
		}
		if i - checkpoint == 100000 {
			checkpoint = i
		}
	}
	fmt.Print("\n")
	return nil
}

func PopulateGroups(ctx context.Context, config BenchConfig, tx *sql.Tx) error  {
	log.Print("populating groups")
	stmt, err := tx.Prepare("INSERT INTO groups (name) VALUES ($1)")
	if err != nil {
		return err
	}

	for i := int64(0); i < config.NumGroups; i++ {
		fmt.Printf("\r[%d/%d]", i, config.NumGroups)
		_, err := stmt.ExecContext(ctx, gofakeit.HipsterWord())
		if err != nil {
			return  err
		}
	}
	fmt.Print("\n")
	return nil
}

func PopulateObjects(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("populating objects")
	stmt, err := tx.Prepare("INSERT INTO objects (name) VALUES ($1)")
	if err != nil {
		return err
	}

	for i := int64(0); i < config.NumObjects; i++ {
		fmt.Printf("\r[%d/%d]", i, config.NumObjects)
		_, err := stmt.ExecContext(ctx, gofakeit.BeerName())
		if err != nil {
			return  err
		}
	}
	fmt.Print("\n")
	return nil
}

func GenerateRandomUserPermissions(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("starting user permission generation (randomized)")
	enums := []string{"read", "alter", "owner", "admin"}

	selectEnum := func() string {
		return enums[rand.Intn(len(enums))]
	}

	stmt, err := tx.Prepare("INSERT INTO object_user_permissions (user_id, object_id, level) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}

	for objectID := int64(1); objectID < config.NumObjects; objectID++ {
		maxAssociated := rand.Intn(10)

		for i := 0; i < maxAssociated; i++ {
			userID := rand.Int63n(config.NumUsers-1)
			if userID == 0 {userID++}
			fmt.Printf("\r[%d/%d/%d]", userID, objectID, config.NumObjects)
			_, err := stmt.ExecContext(ctx, userID, objectID, selectEnum())
			if err != nil {
				return fmt.Errorf("failed on userID: %d, objectID: %d. %w", userID, objectID, err)
			}
		}
	}
	fmt.Print("\n")
	return  nil
}

func GenerateRandomGroupPermissions(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("starting group permission generation (randomized)")
	enums := []string{"read", "alter", "owner", "admin"}

	selectEnum := func() string {
		return enums[rand.Intn(len(enums))]
	}

	stmt, err := tx.Prepare("INSERT INTO object_group_permissions (group_id, object_id, level) VALUES ($1, $2, $3) ON CONFLICT  DO NOTHING")
	if err != nil {
		return err
	}

	for objectID := int64(1); objectID < config.NumObjects; objectID ++ {
		maxAssociated := rand.Intn(10)

		for i := 0; i < maxAssociated; i++ {
			groupID := rand.Int63n(config.NumGroups)
			if groupID == 0 {groupID++}
			fmt.Printf("\r[%d/%d/%d]", groupID, objectID, config.NumObjects)
			_, err := stmt.ExecContext(ctx, groupID, objectID, selectEnum())
			if err != nil {
				return fmt.Errorf("failed on groupID: %d, objectID: %d. %w", groupID, objectID, err)
			}
		}
	}
	fmt.Print("\n")
	return  nil
}

func StitchUsersToGroups(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("starting user to group stitching")
	stmt, err := tx.Prepare("INSERT INTO users_groups (user_id, group_id) VALUES ($1, $2)")
	if err != nil {
		return err
	}

	groupID := int64(1)
	currentGroupSize := int64(0)
	for userID := int64(1); userID < (config.NumUsers - 1); userID++ {
		fmt.Printf("\r[%d/%d]", userID, config.NumUsers)
		_, err := stmt.ExecContext(ctx, userID, groupID)
		if err != nil {
			return fmt.Errorf("StitchUsersToGroups failed on userID: %d, groupID: %d. %w", userID, groupID, err)
		}
		currentGroupSize++
		if currentGroupSize > config.AvgGroupSize {
			groupID++
			currentGroupSize = 0
		}
	}
	fmt.Print("\n")
	return nil
}