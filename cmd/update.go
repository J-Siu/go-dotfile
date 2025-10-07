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
	"github.com/J-Siu/go-dotfile/global"
	"github.com/J-Siu/go-dotfile/lib"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "up"},
	Short:   "Update dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
		property := lib.TypeDotfileProperty{
			DirDest:  &global.Conf.DirDest,
			DirSkip:  &global.Conf.DirSkip,
			FileSkip: &global.Conf.FileSkip,
			Verbose:  global.Flag.Verbose,
		}

		// Process copy
		property.Mode = lib.Copy
		for _, dir := range global.Conf.DirCP {
			var df lib.TypeDotfile
			property.DirSrc = &dir
			df.New(&property)
			if !global.FlagUpdate.Dryrun {
				df.Run()
			}
		}
		// Process append
		property.Mode = lib.Append
		for _, dir := range global.Conf.DirAP {
			var df lib.TypeDotfile
			property.DirSrc = &dir
			df.New(&property)
			if !global.FlagUpdate.Dryrun {
				df.Run()
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.PersistentFlags().BoolVarP(&global.FlagUpdate.Dryrun, "dryrun", "", false, "Dryrun")
}
