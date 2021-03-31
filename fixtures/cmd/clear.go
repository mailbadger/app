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
	"errors"
	"fmt"
	"strconv"

	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear deletes all the stuff created with full-init",
	RunE: func(cmd *cobra.Command, args []string) error {
		return clear()
	},
}

func clear() error {
	u, err := db.GetUserByUsername("badger")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("failed to fetch 'badger' user: %w", err)
	}

	for i := 1; i <= 100; i++ {
		email := "subscriber" + strconv.Itoa(i) + "@mail.com"
		err = db.DeleteSubscriberByEmail(email, u.ID)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"user_id":          u.ID,
				"subscriber_email": email,
			}).Error("failed to delete subscriber")
			return fmt.Errorf("failed to delete subscriber: %w", err)
		}
	}

	fullSegment, err := db.GetSegmentByName("full segment", u.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return fmt.Errorf("failed to fetch 'full segment' segment: %w", err)
	}

	err = db.DeleteSegment(fullSegment.ID, u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete 'full segment' segment: %w", err)
	}

	err = db.DeleteAllEventsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete subscriber events for 'badger' user: %w", err)
	}

	err = db.DeleteUser(u)
	if err != nil {
		return fmt.Errorf("failed to delete 'badger' user: %w", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(clearCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clearCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clearCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
