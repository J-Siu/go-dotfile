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
package lib

import (
	"os"
	"path"
	"strings"

	"github.com/J-Siu/go-helper"
	"github.com/edwardrf/symwalk"
)

// Array contains
func Contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

// Array of string contains a substring
func ContainsArraySubString(strArray []string, str string) bool {
	// prefix := "ContainsArraySubString"
	// helper.ReportDebug(str, prefix, false, true)
	for _, s := range strArray {
		if strings.Contains(str, s) {
			// helper.ReportDebug(str+" ~= "+s, prefix, false, true)
			return true
		}
	}
	// helper.ReportDebug(str+" not match", prefix, false, true)
	return false
}

func DirCreate(dir string, dirBase string) (err error) {
	var prefix = "DirCreate"
	if !(dir == "." || dir == "") {
		dirDest := path.Join(dirBase, "."+dir)
		if !DirExists(dirDest) {
			err = os.MkdirAll(dirDest, os.ModePerm)
			if err == nil {
				helper.ReportDebug(dirDest, prefix+" created", true, true)
			} else {
				helper.Report(err, prefix+" ERROR", false, true)
			}
		}
	}
	return err
}

func DirExists(p string) bool {
	if info, err := os.Stat(p); err == nil {
		return info.IsDir()
	}
	return false
}

// Get list of directory and list of file
func DirFileGet(dir string) (dirs []string, files []string) {
	// prefix := "DirFileGet"
	symwalk.Walk(dir, func(p string, info os.FileInfo, err error) error {
		// helper.ReportDebug(p, "\n"+prefix, false, true)
		if info.IsDir() {
			if p != "." && !ContainsArraySubString(Conf.DirSkip, "/"+p+"/") {
				dirs = append(dirs, p)
			}
		} else {
			if !Contains(Conf.FileSkip, path.Base(p)) && !ContainsArraySubString(Conf.DirSkip, "/"+p) {
				files = append(files, p)
			}
		}
		return nil
	})
	return dirs, files
}

// Add "."" in front of path if there is none
func PathHide(p string) string {
	if strings.HasPrefix(p, ".") {
		return p
	}
	return "." + p
}

const (
	ProcModeAppend string = "append" // Processing mode for TypeDotfile.Mode
	ProcModeCopy   string = "copy"   // Processing mode for TypeDotfile.Mode
)

type TypeDotfile struct {
	Err    error
	myType string
	init   bool

	Dirs    []string
	Files   []string
	DirDest string
	DirSrc  string
	Mode    string // modeAppend | modeCopy
}

func (df *TypeDotfile) Init(dirSrc string, dirDest string, mode string) {
	df.init = true
	df.myType = "TypeConf"
	prefix := df.myType + ".Init"

	if !(mode == ProcModeAppend || mode == ProcModeCopy) {
		helper.Report("mode error", prefix, false, true)
		return
	}

	df.DirDest = dirDest
	df.DirSrc = dirSrc
	df.Mode = mode

	os.Chdir(df.DirSrc)
	df.Dirs, df.Files = DirFileGet(".")
	helper.ReportDebug(df, prefix+" Dotfile", false, false)
	df.ProcessCheckMode()
}

func (df *TypeDotfile) Process() {
	df.ProcessCheckMode()
	// Create directories
	for _, fileDir := range df.Dirs {
		err := DirCreate(fileDir, df.DirDest)
		if err != nil {
			os.Exit(1)
		}
	}
	// Append/Copy files
	for _, filepathSrc := range df.Files {
		filepathDest := path.Join(df.DirDest, PathHide(filepathSrc))
		df.ProcessFile(path.Join(df.DirSrc, filepathSrc), filepathDest)
	}
}

// Mode must be set
//
// If Mode != ProcModeAppend|ProcModeCopy then force exit
func (df *TypeDotfile) ProcessCheckMode() {
	prefix := df.myType + ".CheckMode"
	switch df.Mode {
	case ProcModeAppend:
	case ProcModeCopy:
	default:
		helper.Report(df.myType+".Mode error: "+df.Mode, prefix, false, true)
		os.Exit(1)
	}
}

// Process file base on Mode(append|copy)
//
// Not using TypeDotfile.Err
func (df *TypeDotfile) ProcessFile(src string, dest string) error {
	prefix := df.myType + ".ProcessFile"
	// Read source file
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	fileMode := os.O_CREATE | os.O_WRONLY
	if df.Mode == ProcModeAppend {
		fileMode |= os.O_APPEND
	}
	if df.Mode == ProcModeCopy {
		fileMode |= os.O_TRUNC
	}
	info, err := os.Stat(src)
	if err != nil {
		return err
	}
	filePermission := info.Mode()

	// Open destination file with append
	f, err := os.OpenFile(dest, fileMode, filePermission)
	if err != nil {
		return err
	}

	// Append newline to destination file
	if df.Mode == ProcModeAppend {
		_, err = f.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	// Append source content to destination file
	_, err = f.Write(data)
	if err != nil {
		return err
	}

	f.Close()

	// Set dest permission
	fpStr := ".........."
	if df.Mode == ProcModeCopy {
		err = os.Chmod(dest, filePermission)
		fpStr = filePermission.String() + "  "
	}

	if Flag.Debug || Flag.Verbose {
		if Flag.Debug {
		}
		helper.Report(fpStr+" "+df.Mode+" "+src+" -> "+dest, prefix, false, true)
	}

	return err
}
