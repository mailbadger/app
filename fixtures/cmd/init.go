package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"github.com/mailbadger/app/entities"
	"github.com/mailbadger/app/services/templates"
)

// initCmd represents the fullInit command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init creates testing data",
	Long: `init will create a user with few campaigns and templates, also it will create hundred of subscribers
into few different segments.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return fullInit()
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if username == "" {
			return errors.New("username flag is required")
		}

		if password == "" {
			password = uuid.NewString()
		}

		if secret == "" {
			secret = uuid.NewString()
		}

		return nil
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

	err = createCampaignsAndTemplates(u.ID)
	if err != nil {
		return err
	}

	return nil
}

// createUser creates a user - badger
func createUser() (*entities.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password")
	}

	uuid := uuid.New().String()

	u := &entities.User{
		UUID:     uuid,
		Username: username,
		Password: sql.NullString{
			String: string(hashedPassword),
			Valid:  true,
		},
		Active:     true,
		Verified:   true,
		BoundaryID: 1,
		Source:     viper.GetString("USER_SOURCE"),
	}

	err = db.CreateUser(u)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	fmt.Printf("created user:\n username: %s\n password: %s\n uuid: %s\n\n", username, password, uuid)

	err = db.CreateAPIKey(&entities.APIKey{
		UserID:    u.ID,
		User:      *u,
		SecretKey: secret,
		Active:    true,
	})
	if err != nil {
		fmt.Printf("[ERROR %s] failed to create api key with secret %s\n\n", err.Error(), secret)
	} else {
		fmt.Printf("created api key\n secret: %s\n\n", secret)
	}

	return u, nil
}

// createSubscribersAndSegments creates 100 subscribers in one segment
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

	fmt.Printf("created segment:\n name: %s\n\n", fullSegment.Name)

	fmt.Println("created subscribers:")
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
			fmt.Printf("[ERROR %s] failed to create subscriber\n name: %s, email: %s, segment_id: %d\n", err.Error(), name, email, fullSegment.ID)
			continue
		}

		fmt.Printf(" name: %s email: %s\n", name, email)
	}

	fmt.Println()

	return nil
}

// createCampaignsAndTemplates creates 5 campaigns with one template
func createCampaignsAndTemplates(userID int64) error {
	template := &entities.Template{
		BaseTemplate: entities.BaseTemplate{
			UserID:      userID,
			Name:        "init-template",
			SubjectPart: "Welcome {{name}}",
		},
		HTMLPart: `<h1>{{header}}</h1>
	Dear {{name}},
	You have subscribed to our newsletter.
	<ul>
    	<li>Facebook <a href="{{fb_link}}">fb link</a><li>
    	<li>Instagram <a href="{{i_link}}">ig link</a><li>
    	<li>Twitter <a href="{{t_link}}">tw link</a><li>
	<ul>
	<a href="{{unsubscribe_url}}"><button>Unsubscribe</buntton></a>`,
		TextPart: `{{header}}
	Dear {{name}},
	You have subscribed to our newsletter.
	Facebook {{fb_link}}
	Instagram {{i_link}}
	Twitter {{t_link}}
	Unsubscribe => {{unsubscribe_url}}`,
	}

	err := templates.New(db, s3Client, templates.TemplateBucket(viper.GetString("TEMPLATES_BUCKET"))).AddTemplate(context.Background(), template)
	if err != nil {
		return fmt.Errorf("failed to create 'init-template' template: %w", err)
	}

	fmt.Printf("created template:\n name: %s\n\n", template.Name)

	fmt.Println("created campaigns:")
	for i := 1; i <= 5; i++ {
		name := "campaign" + strconv.Itoa(i)
		campaign := &entities.Campaign{
			UserID:     userID,
			Name:       name,
			TemplateID: template.ID,
		}

		err := db.CreateCampaign(campaign)
		if err != nil {
			fmt.Printf("[ERROR %s] failed to create campaign\n name: %s, template_id: %d", err.Error(), name, template.ID)
			continue
		}

		fmt.Printf(" name: %s, template_id: %d\n", name, template.ID)
	}

	fmt.Println()

	return nil
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// fullInitCmd.PersistentFlags().String("foo", "", "A help for foo")
}
