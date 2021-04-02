/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"database/sql"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/bcrypt"

	"github.com/mailbadger/app/entities"
)

// fullInitCmd represents the fullInit command
var fullInitCmd = &cobra.Command{
	Use:   "full-init",
	Short: "full-init creates testing data",
	Long: `full-init will create a user with few campaigns and templates, also it will create hundred of subscribers
into few different segments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fullInit()
	},
}

func fullInit() error {
	u, err := createUser()
	if err != nil {
		return err
	}

	err = createSubscribersAndSegments(u.ID)
	if err != nil {
		return err
	}

	return nil
}

// createUser creates a user - badger
func createUser() (*entities.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("changeme"), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	uuid := uuid.New()
	username := "badger"

	u := &entities.User{
		UUID:     uuid.String(),
		Username: username,
		Password: sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		},
		Active:     true,
		Verified:   true,
		BoundaryID: 1,
		Source:     "fixtures",
	}

	err = db.CreateUser(u)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("created new user:\n username: %s\n uuid: %s\n\n", username, uuid.String())

	return u, nil
}

func createSubscribersAndSegments(userID int64) error {
	fullSegment := &entities.Segment{
		Model:  entities.Model{},
		Name:   "full segment",
		UserID: userID,
	}

	err := db.CreateSegment(fullSegment)
	if err != nil {
		return fmt.Errorf("failed to create segment: %w", err)
	}

	fmt.Printf("created new segment:\n user_id: %d\n name: %s\n\n", userID, fullSegment.Name)

	for i := 1; i <= 100; i++ {
		name := "subscriber" + strconv.Itoa(i)
		email := name + "@mail.com"
		subscriber := &entities.Subscriber{
			UserID:   userID,
			Name:     name,
			Email:    email,
			MetaJSON: nil,
			Segments: []entities.Segment{
				*fullSegment,
			},
			Blacklisted: false,
			Active:      true,
		}

		err = db.CreateSubscriber(subscriber)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"name":       name,
				"segment_id": fullSegment.ID,
			}).Errorf("failed to create subscriber")
		}

		fmt.Printf("created new subscriber:\n name: %s\n user_id: %d\n email: %s\n\n", name, userID, email)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(fullInitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fullInitCmd.PersistentFlags().String("foo", "", "A help for foo")
}
