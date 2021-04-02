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
	"context"
	"errors"
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/spf13/cobra"

	"github.com/mailbadger/app/services/templates"
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

	apiKeys, err := db.GetAPIKeys(u.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch api keys for user: %w", err)
	}

	for _, apiKey := range apiKeys {
		err = db.DeleteAPIKey(apiKey.ID, u.ID)
		if err != nil {
			return fmt.Errorf("failed to delete ai key: %w", err)
		}
	}

	subscribers, err := db.GetAllSubscribersForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch subscribers for user: %w", err)
	}

	for _, s := range subscribers {
		err = db.DeleteSubscriber(s.ID, u.ID)
		if err != nil {
			return fmt.Errorf("failed to delete subscriber: %w", err)
		}
	}

	err = db.DeleteAllSegmentsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all segments for user: %w", err)
	}

	err = db.DeleteAllEventsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete subscriber events for 'badger' user: %w", err)
	}

	err = db.DeleteAllBouncesForUSer(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all bounces for user: %w", err)
	}

	err = db.DeleteAllCampaignFailedLogsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all campaign failed logs for user: %w", err)
	}

	err = db.DeleteAllClicksForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all clicks for user: %w", err)
	}

	err = db.DeleteAllComplaintsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all complaints for user: %w", err)
	}

	err = db.DeleteAllDeliveriesForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all deliveries for user: %w", err)
	}

	err = db.DeleteAllOpensForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all opens for user: %w", err)
	}

	err = db.DeleteAllSendsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all sends for user: %w", err)
	}

	err = db.DeleteAllCampaignsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all campaigns for user: %w", err)
	}

	allTemplates, err := db.GetAllTemplatesForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch all tmeplates for user: %w", err)
	}

	for _, t := range allTemplates {
		err = templates.New(db, s3Client).DeleteTemplate(context.Background(), t.ID, u.ID)
		if err != nil {
			return fmt.Errorf("failed to delete template	: %w", err)
		}
	}

	err = db.DeleteAllReportsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all reports for user: %w", err)
	}

	err = db.DeleteSesKeys(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete ses keys: %w", err)
	}

	err = db.DeleteAllSessionsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all sessions for user: %w", err)
	}

	err = db.DeleteAllTokensForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all tokens for user: %w", err)
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
