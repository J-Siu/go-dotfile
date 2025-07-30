/*
Copyright © 2025 John, Sing Dao, Siu <john.sd.siu@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"os"
	"runtime"

	"github.com/J-Siu/go-dotfile/lib"
	"github.com/J-Siu/go-helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go-dotfile",
	Short: "A dotfile manager",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		helper.Debug = lib.Flag.Debug
		lib.Conf.Init()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	switch runtime.GOOS {
	case "darwin":
	case "linux":
	default:
		helper.Report("Only support Linux and MacOS.", "", false, true)
		os.Exit(1)
	}
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().BoolVarP(&lib.Flag.Debug, "debug", "d", false, "Enable debug")
	rootCmd.PersistentFlags().BoolVarP(&lib.Flag.Dryrun, "dryrun", "", false, "Dryrun")
	rootCmd.PersistentFlags().BoolVarP(&lib.Flag.Verbose, "verbose", "v", false, "Verbose")
	rootCmd.PersistentFlags().StringVarP(&lib.Conf.File, "config", "c", lib.DefaultConfFile, "Config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigType("json")
	viper.SetConfigFile(os.ExpandEnv(lib.Conf.File))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		helper.Report(err.Error(), "", true, true)
		os.Exit(1)
	}
}
