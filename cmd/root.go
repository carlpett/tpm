package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
)

var rootCmd = &cobra.Command{
	Use:     "tpm",
	Version: "0.3.0",
	Short:   "Terraform Provider Manager",
	Long:    "Terraform Provider Manager is a simple CLI to manage Terraform providers in the Terraform plugin cache directory",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
		os.Exit(1)
	},
}

func Execute() {
	rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global Flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file for tpm")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "enable debug mode")
	rootCmd.PersistentFlags().StringP("terraform-plugin-cache-dir", "p", filepath.Join(os.ExpandEnv("$HOME"), "/.terraform.d/plugin-cache"), "the location of the Terraform plugin cache directory")
	rootCmd.PersistentFlags().StringP("terraform-registry", "r", "registry.terraform.io", "the Terraform registry provider hostname")
	rootCmd.PersistentFlags().VisitAll(bindCustomFlag)
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		// Config name
		viper.SetConfigName("config")
		viper.SetConfigType("json")

		// Config location
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME/.tpm")
		viper.AddConfigPath("/usr/local/etc/tpm")
	}

	// Environment variables
	viper.SetEnvPrefix("TPM")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("debug") {
			log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)
			log.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		if viper.GetBool("debug") {
			log.Println("Error reading config file,", err)
		}
	}
}
