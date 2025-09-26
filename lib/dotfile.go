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
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/J-Siu/go-basestruct"
	"github.com/J-Siu/go-helper"
	"github.com/edwardrf/symwalk"
)

type FileProcMode int8

const (
	Append FileProcMode = iota // Processing mode for TypeDotfile.Mode
	Copy                       // Processing mode for TypeDotfile.Mode
	Skip
)

type TypeDotfile struct {
	basestruct.Base

	Dirs    []string     `json:"dirs,omitempty"`
	Files   []string     `json:"files,omitempty"`
	DirDest string       `json:"dir_dest,omitempty"`
	DirSrc  string       `json:"dir_src,omitempty"`
	Mode    FileProcMode `json:"mode,omitempty"`
}

func (df *TypeDotfile) New(dirSrc string, dirDest string, mode FileProcMode) {
	df.Initialized = true
	df.MyType = "TypeDotfile"
	prefix := df.MyType + ".Init"

	if !(mode == Append || mode == Copy) {
		helper.Report("mode error", prefix, false, true)
		return
	}

	df.DirDest = dirDest
	df.DirSrc = dirSrc
	df.Mode = mode

	// cd to simplify path handling
	err := os.Chdir(df.DirSrc)
	if err == nil {
		df.Dirs, df.Files = dirFileGet(".")
	}
	helper.ReportDebug(df, prefix+" Dotfile", false, false)
}

func (df *TypeDotfile) Run() {
	prefix := df.MyType + ".Process"
	var err error
	for _, fileDir := range df.Dirs {
		err = dirCreateHidden(fileDir, df.DirDest)
		if err != nil {
			os.Exit(1)
		}
	}
	// Append/Copy files
	for _, filepathSrc := range df.Files {
		filepathDest := path.Join(df.DirDest, pathHide(filepathSrc))
		err = df.processFile(path.Join(df.DirSrc, filepathSrc), filepathDest)
		helper.ErrsQueue(err, prefix)
	}
}

// Process file base on Mode(append|copy)
//
// Not using TypeDotfile.Err
func (df *TypeDotfile) processFile(src string, dest string) (err error) {
	prefix := df.MyType + ".ProcessFile"

	fileProcMode := df.Mode
	filePermStr := ".........."

	// Destination FileMode
	destFlag := os.O_CREATE | os.O_WRONLY
	if fileProcMode == Append {
		destFlag |= os.O_APPEND
	}
	if fileProcMode == Copy {
		destFlag |= os.O_TRUNC
		changed, err := fileChanged(src, dest)
		if err != nil {
			return err
		}
		if !changed {
			fileProcMode = Skip
		}
	}

	if fileProcMode != Skip {
		srcInfo, err := os.Stat(src)
		if err != nil {
			return err
		}
		srcPermission := srcInfo.Mode()
		srcModTime := srcInfo.ModTime()

		// Read source file
		data, err := os.ReadFile(src)
		if err != nil {
			return err
		}

		// Open destination file
		f, err := os.OpenFile(dest, destFlag, srcPermission)
		if err != nil {
			return err
		}

		// Append: add newline to destination file
		if fileProcMode == Append {
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
		if fileProcMode == Copy {
			os.Chtimes(dest, srcModTime, srcModTime)
			filePermStr = srcPermission.String()
		}
	}

	str := fmt.Sprintf("%-6s %s %s -> %s", fileProcMode.String(), filePermStr, src, dest)
	if Flag.Debug {
		helper.Report(str, prefix, false, true)
	} else if Flag.Verbose {
		fmt.Println(str)
	}

	return err
}

// Array contains
func contains[T comparable](arr []T, x T) bool {
	for _, v := range arr {
		if v == x {
			return true
		}
	}
	return false
}

// Array of string contains a substring
func containsArraySubString(strArray []string, str string) bool {
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

// Create dotted/hidden directory
func dirCreateHidden(dir string, dirBase string) (err error) {
	var prefix = "DirCreate"
	if !(dir == "." || dir == "") {
		dirDest := path.Join(dirBase, pathHide(dir))
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
func dirFileGet(dir string) (dirs []string, files []string) {
	err := symwalk.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if p != "." && !containsArraySubString(Conf.DirSkip, "/"+p+"/") {
				dirs = append(dirs, p)
			}
		} else {
			if !contains(Conf.FileSkip, path.Base(p)) && !containsArraySubString(Conf.DirSkip, "/"+p) {
				files = append(files, p)
			}
		}
		return nil
	})
	// prefix := "dirFileGet"
	// helper.Report(err, prefix, false, true)
	if err == nil {
		return dirs, files
	} else {
		return []string{}, []string{}
	}
}

// Add "."" in front of path if there is none
func pathHide(p string) string {
	if strings.HasPrefix(p, ".") {
		return p
	}
	return "." + p
}

// False: if (src and dst have same modification time and size)
//
// True: else
//
// Ignore destination file stat() err
func fileChanged(src string, dst string) (changed bool, err error) {
	var (
		infoDst os.FileInfo
		infoSrc os.FileInfo
	)

	changed = true

	infoSrc, err = os.Stat(src)
	if err == nil {
		infoDst, err = os.Stat(dst)

		if err == nil {
			if time.Time.Equal(infoDst.ModTime(), infoSrc.ModTime()) &&
				infoDst.Size() == infoSrc.Size() {
				changed = false
			}
		}
		// it is for destination file stat to fail
		err = nil
	}

	return changed, err
}
