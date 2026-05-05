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

type ModeDirPair struct {
	Mode lib.FileProcMode
	Dirs *[]string
}

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "up"},
	Short:   "Update dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			df            lib.TypeDotfile
			mode_dir_pair = []ModeDirPair{
				{lib.COPY, &global.Conf.DirCP},
				{lib.APPEND, &global.Conf.DirAP},
			}
			property = lib.TypeDotfileProperty{
				DirDest:  &global.Conf.DirDest,
				DirSkip:  &global.Conf.DirSkip,
				Save:     global.FlagUpdate.Save,
				FileSkip: &global.Conf.FileSkip,
			}
			records lib.TypeDotfileRecords
		)
		for _, m := range mode_dir_pair {
			property.Mode = m.Mode
			for _, dir := range *m.Dirs {
				property.DirSrc = &dir
				df.New(&property).Run()
				records = append(records, df.Records...)
			}
		}
		// Output
		records.Output(global.FlagUpdate.NoInfo, global.FlagUpdate.Quiet, global.Flag.Verbose, global.FlagUpdate.Save)
	},
}

func init() {
	cmd := updateCmd
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolVarP(&global.FlagUpdate.NoInfo, "noinfo", "n", false, "Do not print file info")
	cmd.Flags().BoolVarP(&global.FlagUpdate.Quiet, "quiet", "q", false, "Show non-skip file only")
	cmd.Flags().BoolVarP(&global.FlagUpdate.Save, "save", "s", false, "Save changes")
}
