package lib

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/strany"
)

// Record struct to store processed dotfile information
type TypeDotfileRecord struct {
	DesInfo      *os.FileInfo `json:"DesInfo"`
	DesPath      string       `json:"DesPath"`
	FileProcMode FileProcMode `json:"FileProcMode"`
	SrcInfo      *os.FileInfo `json:"SrcInfo"`
	SrcPath      string       `json:"SrcPath"`
}

type TypeDotfileRecords []*TypeDotfileRecord

func (t *TypeDotfileRecords) Output(noInfo, quiet, verbose, save bool) {
	const (
		STR_NO_MODTIME  = "---------- --:--:--"
		STR_TIME_FORMAT = "2006-01-02 15:04:05"
	)
	var (
		dupCopy      bool                        // true: duplicate copy to same file
		dupList      = make(map[string][]string) // map desPath to srcPath array
		recordStrArr []string
		tab_Writer   = tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)
	)
	for _, r := range *t {
		var (
			desModTimeStr string = STR_NO_MODTIME
			desSize       int64
		)
		// populate duplicate copy list
		if r.FileProcMode == COPY {
			dupList[r.DesPath] = append(dupList[r.DesPath], r.SrcPath)
			// It is duplicate copy if more than 1 src path
			if len(dupList[r.DesPath]) > 1 {
				dupCopy = true
			}
		}
		// output
		if ezlog.GetLogLevel() >= ezlog.DEBUG ||
			verbose ||
			!quiet && r.FileProcMode != SKIP {
			if !save { // Dry run prefix?
				recordStrArr = []string{"DryRun:"}
			} else {
				recordStrArr = nil
			}
			if noInfo { // file path only
				recordStrArr = append(recordStrArr,
					r.FileProcMode.String(),
					r.SrcPath,
					"->",
					r.DesPath,
				)
			} else { // full file info
				if r.DesInfo != nil {
					desModTimeStr = (*r.DesInfo).ModTime().Local().Format(STR_TIME_FORMAT)
					desSize = (*r.DesInfo).Size()
				}
				recordStrArr = append(recordStrArr,
					r.FileProcMode.String(),
					(*r.SrcInfo).Mode().String(),
					*strany.Any((*r.SrcInfo).Size()),
					(*r.SrcInfo).ModTime().Local().Format(STR_TIME_FORMAT),
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

	if dupCopy && !quiet {
		outputDupList(dupList)
	}
}

func outputDupList(dupList map[string][]string) {
	var (
		headerPrinted bool
	)
	for k, v := range dupList {
		if len(v) > 1 {
			if !headerPrinted {
				ezlog.Log().M("*** Duplicate Copy ***").Out()
				headerPrinted = true
			}
			ezlog.Log().N(k).Out()
			ezlog.Log().M(v).Out()
		}
	}
}
