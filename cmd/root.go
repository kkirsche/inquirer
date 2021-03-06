// Copyright © 2016 Kevin Kirsche <kev.kirsche@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"log"
	"log/syslog"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var oid string
var syslogger *syslog.Writer
var saveMethod string
var file *os.File

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "inquirer",
	Short: "Inquirer is used to poll a remote host for SNMP data",
	Long: `Inquirer is a tool to be leveraged from the command line to retrieve
SNMP v2c data from a remote device or from devices.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "settings", "s", "", "config file (default is $HOME/.inquirer.json)")
	RootCmd.PersistentFlags().StringP("ip", "i", "127.0.0.1", "remote host to query")
	RootCmd.PersistentFlags().StringP("community", "c", "Public", "remote host to query")
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".inquirer") // name of config file (without extension)
	viper.SetConfigType("json")      // extension of config file, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop"
	viper.AddConfigPath("$HOME")     // adding home directory as first search path
	viper.AutomaticEnv()             // read in environment variables that match

	viper.BindPFlag("ip", RootCmd.Flags().Lookup("ip"))
	viper.BindPFlag("community", RootCmd.Flags().Lookup("community"))

	viper.SetDefault("cron.save_via", "stdout")
	viper.SetDefault("cron.save_file_path", "")
	viper.SetDefault("cron.save_filename", "results")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		viper.WatchConfig()
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func writeToOutputMethod(v ...interface{}) {
	switch saveMethod {
	case "stdout":
		fmt.Print(v...)
	case "file":
		_, err := fmt.Fprint(file, v...)
		if err != nil {
			log.Fatal(err.Error())
		}
	case "syslog":
		compoundString := fmt.Sprintf("%v", v...)
		err := syslogger.Notice(compoundString)
		if err != nil {
			log.Fatal(err.Error())
		}
	default:
		fmt.Print(v...)
	}
}
