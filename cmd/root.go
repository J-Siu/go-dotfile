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

	"github.com/J-Siu/go-dotfile/global"
	"github.com/J-Siu/go-dotfile/lib"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "go-dotfile",
	Short:   "A dotfile manager",
	Version: global.Version,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		ezlog.SetLogLevel(ezlog.ERR)
		if global.Flag.Debug {
			ezlog.SetLogLevel(ezlog.DEBUG)
		}
		if global.Flag.Trace {
			ezlog.SetLogLevel(ezlog.TRACE)
		}
		ezlog.Debug().N("Version").Mn(global.Version).Nn("Flag").M(&global.Flag).Out()
		global.Conf.New()
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		if errs.NotEmpty() {
			ezlog.Err().M(errs.Errs).Out()
		}
	},
}

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
		ezlog.Log().M("Only support Linux and MacOS.").Out()
		os.Exit(1)
	}
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.Debug, "debug", "d", false, "Enable debug")
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.Trace, "trace", "t", false, "Enable trace")
	rootCmd.PersistentFlags().BoolVarP(&global.Flag.Verbose, "verbose", "v", false, "Verbose")
	rootCmd.PersistentFlags().StringVarP(&global.Conf.FileConf, "config", "c", lib.Default.FileConf, "Config file")
}
