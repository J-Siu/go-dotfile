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

	"github.com/J-Siu/go-helper"
	"github.com/edwardrf/symwalk"
)

func FileCopy(fileSrc string, fileDest string) error {
	prefix := "FileCP"
	helper.ReportDebug(fileSrc+" -> "+fileDest, prefix, false, true)

	data, err := os.ReadFile(fileSrc)
	if err != nil {
		return err
	}

	return os.WriteFile(fileDest, data, 0666)
}

func FileAppend(fileSrc string, fileDest string) error {
	prefix := "FileAP"
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

func DirFileGet(dir string) (dirs []string, files []string) {
	// var files []string
	symwalk.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			dirs = append(dirs, p)
		} else {
			files = append(files, p)
		}
		return nil
	})
	return dirs, files
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

const (
	ProcModeAppend int = 0
	ProcModeCopy   int = 1
)

type Dotfile struct {
	Dirs    []string
	Files   []string
	DirDest string
	DirSrc  string
	Mode    int // modeAppend | modeCopy
}

func (df *Dotfile) Init(dirSrc string, dirDest string, mode int) {
	prefix := "Dotfiles.Init()"
	if !(mode == ProcModeAppend || mode == ProcModeCopy) {
		helper.Report("mode error", prefix, false, true)
		return
	}

	df.DirDest = dirDest
	df.DirSrc = dirSrc
	df.Mode = mode

	os.Chdir(df.DirSrc)
	df.Dirs, df.Files = DirFileGet(".")
	// helper.ReportDebug(df.Dirs, prefix+" Dirs", false, false)
	// helper.ReportDebug(df.Files, prefix+" Files", false, false)
}

func (df *Dotfile) Process() {
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
			helper.Report(df.Mode, prefix+" df.Mode error", false, true)
		}
	}
}
