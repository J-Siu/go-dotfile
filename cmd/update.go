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
	"github.com/J-Siu/go-helper/v2/strany"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "up"},
	Short:   "Update dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
		var (
			df       lib.TypeDotfile
			property = lib.TypeDotfileProperty{
				DirDest:  &global.Conf.DirDest,
				DirSkip:  &global.Conf.DirSkip,
				Save:     global.FlagUpdate.Save,
				FileSkip: &global.Conf.FileSkip,
			}
			records lib.TypeDotfileRecords
		)
		// Process copy
		property.Mode = lib.COPY
		for _, dir := range global.Conf.DirCP {
			property.DirSrc = &dir
			df.New(&property).Run()
			records = append(records, df.Records...)
		}
		// Process append
		property.Mode = lib.APPEND
		for _, dir := range global.Conf.DirAP {
			property.DirSrc = &dir
			df.New(&property)
			df.Run()
			records = append(records, df.Records...)
		}
		// Output
		output(&records)
	},
}

func init() {
	cmd := updateCmd
	rootCmd.AddCommand(cmd)
	cmd.Flags().BoolVarP(&global.FlagUpdate.NoInfo, "noinfo", "n", false, "Do not print file info")
	cmd.Flags().BoolVarP(&global.FlagUpdate.Quiet, "quiet", "q", false, "Show non-skip file only")
	cmd.Flags().BoolVarP(&global.FlagUpdate.Save, "save", "s", false, "Save changes")
}

func output(records *lib.TypeDotfileRecords) {
	const (
		noModTimeStr = "---------- --:--:--"
		timeFormat   = "2006-01-02 15:04:05"
	)
	var (
		dupCopy      bool                        // true: duplicate copy to same file
		dupList      = make(map[string][]string) // map desPath to srcPath array
		recordStrArr []string
		tab_Writer   = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	)
	for _, r := range *records {
		var (
			desModTimeStr string = noModTimeStr
			desSize       int64
		)
		// populate duplicate copy list
		if r.FileProcMode == lib.COPY {
			dupList[r.DesPath] = append(dupList[r.DesPath], r.SrcPath)
			// It is duplicate copy if more than 1 src path
			if len(dupList[r.DesPath]) > 1 {
				dupCopy = true
			}
		}
		// output
		if ezlog.GetLogLevel() >= ezlog.DEBUG ||
			global.Flag.Verbose ||
			!global.FlagUpdate.Quiet && r.FileProcMode != lib.SKIP {
			// Dry run prefix?
			if !global.FlagUpdate.Save {
				recordStrArr = []string{"DryRun:"}
			} else {
				recordStrArr = nil
			}
			if global.FlagUpdate.NoInfo {
				// file path only
				recordStrArr = append(recordStrArr,
					r.FileProcMode.String(),
					r.SrcPath,
					"->",
					r.DesPath,
				)
			} else {
				// full file info
				if r.DesInfo != nil {
					desModTimeStr = (*r.DesInfo).ModTime().Local().Format(timeFormat)
					desSize = (*r.DesInfo).Size()
				}
				recordStrArr = append(recordStrArr,
					r.FileProcMode.String(),
					(*r.SrcInfo).Mode().String(),
					*strany.Any((*r.SrcInfo).Size()),
					(*r.SrcInfo).ModTime().Local().Format(timeFormat),
					r.SrcPath,
					"->",
					*strany.Any(desSize),
					desModTimeStr,
					r.DesPath,
				)
			}
			// send to tabwriter
			fmt.Fprintln(tab_Writer, strings.Join(recordStrArr, "\t"))
		}
	}
	tab_Writer.Flush()

	if dupCopy && !global.FlagUpdate.Quiet {
		outputDupList(dupList)
	}
}

func outputDupList(dupList map[string][]string) {
	ezlog.Log().M("*** Duplicate Copy ***").Out()
	if len(dupList) > 0 {
		ezlog.Log().M(dupList).Out()
	}
}
