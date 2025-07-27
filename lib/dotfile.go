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

func Contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

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

func DirFileGet(dir string) (dirs []string, files []string) {
	// prefix := "DirFileGet"
	symwalk.Walk(dir, func(p string, info os.FileInfo, err error) error {
		// helper.ReportDebug(p, "\n"+prefix, false, true)
		if info.IsDir() {
			if !ContainsArraySubString(Conf.DirSkip, "/"+p+"/") {
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

func FileChmod(fileSrc string, fileDest string) (err error) {
	info, err := os.Stat(fileSrc)
	if err != nil {
		return err
	}

	prefix := "FileChmod"
	helper.ReportDebug(fileSrc+" -> "+fileDest+"("+info.Mode().String()+")", prefix, false, true)

	return os.Chmod(fileDest, info.Mode())
}

func FileCopy(fileSrc string, fileDest string) error {
	prefix := "FileCopy"
	helper.ReportDebug(fileSrc+" -> "+fileDest, prefix, false, true)

	data, err := os.ReadFile(fileSrc)
	if err != nil {
		return err
	}

	err = os.WriteFile(fileDest, data, 0666)
	if err != nil {
		return err
	}

	return FileChmod(fileSrc, fileDest)
}

func FileAppend(fileSrc string, fileDest string) error {
	prefix := "FileAppend"
	helper.ReportDebug(fileSrc+" -> "+fileDest, prefix, false, true)

	// Read source file
	data, err := os.ReadFile(fileSrc)
	if err != nil {
		return err
	}

	// Open destination file with append
	f, err := os.OpenFile(fileDest, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Append newline to destination file
	_, err = f.Write([]byte("\n"))
	if err != nil {
		return err
	}

	// Append source content to destination file
	_, err = f.Write(data)
	f.Close()
	return err
}

const (
	ProcModeAppend string = "append"
	ProcModeCopy   string = "copy"
)

type TypeDotfile struct {
	Dirs    []string
	Files   []string
	DirDest string
	DirSrc  string
	Mode    string // modeAppend | modeCopy
}

func (df *TypeDotfile) Init(dirSrc string, dirDest string, mode string) {
	prefix := "Dotfiles.Init"
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
}

func (df *TypeDotfile) Process() {
	prefix := "Dotfiles.Process"
	// Create directories
	for _, fileDir := range df.Dirs {
		err := DirCreate(fileDir, df.DirDest)
		if err != nil {
			os.Exit(1)
		}
	}
	// Append/Copy files
	for _, fileSrc := range df.Files {
		var fileDest = path.Join(df.DirDest, "."+fileSrc)
		switch df.Mode {
		case ProcModeAppend:
			FileAppend(path.Join(df.DirSrc, fileSrc), fileDest)
		case ProcModeCopy:
			FileCopy(path.Join(df.DirSrc, fileSrc), fileDest)
		default:
			helper.Report("df.Mode error: "+df.Mode, prefix, false, true)
		}
	}
}
