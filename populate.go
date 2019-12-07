package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/brianvoe/gofakeit/v4"
	"log"
	"math/rand"
	"github.com/lib/pq"
)

func PopulateUsers(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("populating users")
	stmt, err := tx.Prepare(pq.CopyIn("users", "username"))
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
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return  err
	}

	fmt.Print("\n")
	log.Print("finished populating users")
	return nil
}

func PopulateGroups(ctx context.Context, config BenchConfig, tx *sql.Tx) error  {
	log.Print("populating groups")
	stmt, err := tx.Prepare(pq.CopyIn("groups", "name"))
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
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return  err
	}

	fmt.Print("\n")
	log.Print("finished populating groups")
	return nil
}

func PopulateObjects(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("populating objects")
	stmt, err := tx.Prepare(pq.CopyIn("objects", "name"))
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
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return  err
	}

	fmt.Print("\n")
	log.Print("finished populating objects")
	return nil
}

func GenerateRandomUserPermissions(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("starting user permission generation (randomized)")
	enums := []string{"read", "alter", "owner", "admin"}

	selectEnum := func() string {
		return enums[rand.Intn(len(enums))]
	}

	_, err := tx.ExecContext(ctx, `CREATE TEMPORARY TABLE object_user_permissions_temp (user_id INTEGER, object_id INTEGER, level permission_level) ON COMMIT DROP`)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("object_user_permissions_temp", "user_id", "object_id", "level"))
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
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return  err
	}

	_, err = tx.Exec("INSERT INTO object_user_permissions (user_id, object_id, level) SELECT  * FROM object_user_permissions_temp ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}

	fmt.Print("\n")
	log.Print("finised user permission generation")
	return  nil
}

func GenerateRandomGroupPermissions(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("starting group permission generation (randomized)")
	enums := []string{"read", "alter", "owner", "admin"}

	selectEnum := func() string {
		return enums[rand.Intn(len(enums))]
	}

	_, err := tx.ExecContext(ctx, `CREATE TEMPORARY TABLE object_group_permissions_temp (group_id INTEGER, object_id INTEGER, level permission_level) ON COMMIT DROP;`)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(pq.CopyIn("object_group_permissions_temp", "group_id", "object_id", "level"))
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
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return  err
	}

	_, err = tx.Exec("INSERT INTO object_group_permissions (group_id, object_id, level) SELECT  * FROM object_group_permissions_temp ON CONFLICT DO NOTHING")
	if err != nil {
		return err
	}

	fmt.Print("\n")
	log.Print("starting group permission generation")
	return  nil
}

func StitchUsersToGroups(ctx context.Context, config BenchConfig, tx *sql.Tx) error {
	log.Print("starting user to group stitching")
	stmt, err := tx.Prepare(pq.CopyIn("users_groups", "user_id", "group_id"))
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

			if groupID > config.NumGroups {
				groupID = 1
			}
		}
	}
	_, err = stmt.ExecContext(ctx)
	if err != nil {
		return  err
	}

	fmt.Print("\n")
	log.Print("finished stitching users to groups")
	return nil
}