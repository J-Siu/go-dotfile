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
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/J-Siu/go-dotfile/global"
	"github.com/J-Siu/go-dotfile/lib"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "up"},
	Short:   "Update dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
		var result []*[]string
		//
		property := lib.TypeDotfileProperty{
			DirDest:  &global.Conf.DirDest,
			DirSkip:  &global.Conf.DirSkip,
			Save:     global.FlagUpdate.Save,
			FileSkip: &global.Conf.FileSkip,
		}
		// Process copy
		property.Mode = lib.COPY
		for _, dir := range global.Conf.DirCP {
			var df lib.TypeDotfile
			property.DirSrc = &dir
			df.New(&property).Run()
			result = append(result, df.Result...)
		}
		// Process append
		property.Mode = lib.APPEND
		for _, dir := range global.Conf.DirAP {
			var df lib.TypeDotfile
			property.DirSrc = &dir
			df.New(&property)
			if !global.FlagUpdate.Save {
				// TODO: This is not complete. It should still run.
				df.Run()
				result = append(result, df.Result...)
			}
		}
		output(&result)
	},
}

func init() {
	cmd := updateCmd
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolVarP(&global.FlagUpdate.NonSkip, "non-skip", "n", false, "Show non-skip file only")
	cmd.Flags().BoolVarP(&global.FlagUpdate.Save, "save", "s", false, "Save changes")
}

func output(result *[]*[]string) {
	var (
		w   = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
		str string
	)
	for _, r := range *result {
		if ezlog.GetLogLevel() >= ezlog.DEBUG || global.Flag.Verbose || global.FlagUpdate.NonSkip && (*r)[0] != lib.SKIP.String() {
			str = strings.Join(*r, "\t")
			if !global.FlagUpdate.Save {
				str = "DryRun:\t" + str
			}
			fmt.Fprintln(w, str)
		}
	}
	w.Flush()
}
