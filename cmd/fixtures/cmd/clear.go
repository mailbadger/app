package cmd

import (
	"context"
	"errors"
	"fmt"

	"github.com/mailbadger/app/entities"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/mailbadger/app/services/templates"
)

var u *entities.User

// clearCmd represents the clear command
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear deletes all data for the user",
	RunE: func(cmd *cobra.Command, args []string) error {
		return clear()
	},
	PreRunE: func(cmd *cobra.Command, args []string) (err error) {
		u, err = db.GetUserByUsername(username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return fmt.Errorf("failed to fetch %s user: %w", username, err)
		}

		if u.Source != viper.GetString("USER_SOURCE") || !u.Password.Valid {
			return fmt.Errorf("user: %s can not be cleared, it is not created by %s", username, viper.GetString("USER_SOURCE"))
		}

		pwString := ""
		fmt.Printf("password for %s: ", username)
		_, err = fmt.Scanf("%s", &pwString)
		if err != nil {
			return fmt.Errorf("failed to scan password %w", err)
		}

		fmt.Println()

		err = bcrypt.CompareHashAndPassword([]byte(u.Password.String), []byte(pwString))
		if err != nil {
			return fmt.Errorf("invalid credentials %w", err)
		}

		return nil
	},
}

func clear() error {
	apiKeys, err := db.GetAPIKeys(u.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch api keys for user: %w", err)
	}

	fmt.Println("deleted api keys")
	for _, apiKey := range apiKeys {
		err = db.DeleteAPIKey(apiKey.ID, u.ID)
		if err != nil {
			return fmt.Errorf("failed to delete api key with id %d: %w", apiKey.ID, err)
		}

		fmt.Printf(" secret: %s\n", apiKey.SecretKey)
	}

	fmt.Println()

	subscribers, err := db.GetAllSubscribersForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch subscribers for user: %w", err)
	}

	fmt.Println("deleted subscribers")
	for _, s := range subscribers {
		err = db.DeleteSubscriber(s.ID, u.ID)
		if err != nil {
			return fmt.Errorf("failed to delete subscriber: %w", err)
		}

		fmt.Printf(" name: %s, email: %s\n", s.Name, s.Email)
	}

	fmt.Println()

	err = db.DeleteAllSegmentsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all segments for user: %w", err)
	}

	fmt.Printf("deleted all segments\n\n")

	err = db.DeleteAllEventsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete subscriber events for 'badger' user: %w", err)
	}

	fmt.Printf("deleted all events\n\n")

	err = db.DeleteAllBouncesForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all bounces for user: %w", err)
	}

	fmt.Printf("deleted all bounces\n\n")

	err = db.DeleteAllCampaignFailedLogsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all campaign failed logs for user: %w", err)
	}

	fmt.Printf("deleted all campaigns failed logs\n\n")

	err = db.DeleteAllClicksForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all clicks for user: %w", err)
	}

	fmt.Printf("deleted all clicks\n\n")

	err = db.DeleteAllComplaintsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all complaints for user: %w", err)
	}

	fmt.Printf("deleted all complaints\n\n")

	err = db.DeleteAllDeliveriesForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all deliveries for user: %w", err)
	}

	fmt.Printf("deleted all deliveries\n\n")

	err = db.DeleteAllOpensForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all opens for user: %w", err)
	}

	fmt.Printf("deleted all opens\n\n")

	err = db.DeleteAllSendsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all sends for user: %w", err)
	}

	fmt.Printf("deleted all sends\n\n")

	err = db.DeleteAllCampaignsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all campaigns for user: %w", err)
	}

	fmt.Printf("deleted all campaigns\n\n")

	allTemplates, err := db.GetAllTemplatesForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch all tmeplates for user: %w", err)
	}

	for _, t := range allTemplates {
		err = templates.New(db, s3Client, templates.TemplateBucket(viper.GetString("TEMPLATES_BUCKET"))).DeleteTemplate(context.Background(), t.ID, u.ID)
		if err != nil {
			return fmt.Errorf("failed to delete template	: %w", err)
		}
	}

	err = db.DeleteAllReportsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all reports for user: %w", err)
	}

	fmt.Printf("deleted all reports\n\n")

	err = db.DeleteSesKeys(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete ses keys: %w", err)
	}

	fmt.Printf("deleted ses keys\n\n")

	err = db.DeleteAllSessionsForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all sessions for user: %w", err)
	}

	fmt.Printf("deleted all sessions\n\n")

	err = db.DeleteAllTokensForUser(u.ID)
	if err != nil {
		return fmt.Errorf("failed to delete all tokens for user: %w", err)
	}

	fmt.Printf("deleted all tokens\n\n")

	err = db.DeleteUser(u)
	if err != nil {
		return fmt.Errorf("failed to delete 'badger' user: %w", err)
	}

	fmt.Printf("deleted user\n username: %s, uuid %s\n", u.Username, u.UUID)

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
