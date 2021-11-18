package cmd

import (
	"fmt"
	"os"

	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mailbadger/app/s3"
	"github.com/mailbadger/app/storage"
)

const version = "v1.0.0"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fixtures",
	Short: "fixtures is a cli for generating test data for mailbadger",
	Long: `Fixturenamenames can generate testing data a user with a few campaigns alongside with a few templates. Also about 
hundreds of subscribers in a few segments`,
	Version: version,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var (
	// db represents the connection to the database
	db storage.Storage
	// s3Client represents the s3 client
	s3Client *awss3.S3
	// username for the user with fixtures
	username string
	// password for the user with fixtures
	password string
	// secret for the user with fixtures
	secret string
)

func init() {
	// viper reads conf file app.env located in fixtures
	initConfig()

	var err error
	logrus.SetLevel(logrus.InfoLevel)
	logrus.SetOutput(os.Stdout)

	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "", "Username for the user with fixtures")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "p", "", "Password for the user with fixtures")
	rootCmd.PersistentFlags().StringVarP(&password, "secret", "s", "", "Secret for api key for the user with fixtures")

	rootCmd.InitDefaultVersionFlag()

	// Connecting to database
	driver := viper.GetString("DATABASE_DRIVER")
	conf := storage.MakeConfigFromEnv(driver)
	db = storage.New(driver, conf)

	// Creating s3 client
	s3Client, err = s3.NewS3Client(
		viper.GetString("AWS_ACCESS_KEY_ID"),
		viper.GetString("AWS_SECRET_ACCESS_KEY"),
		viper.GetString("AWS_REGION"),
	)
	if err != nil {
		fmt.Printf("[ERROR %s] failed to create s3 client", err.Error())
		os.Exit(1)
	}
}

// initConfig reads configuration file
func initConfig() {
	viper.SetConfigFile(".env.local")
	viper.SetConfigType("env")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("[ERROR %s] failed to read config file", err.Error())
		os.Exit(1)
	}
}
