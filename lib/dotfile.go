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

	"github.com/J-Siu/go-basestruct"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/J-Siu/go-helper/v2/str"
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

	DirDest string       `json:"dir_dest,omitempty"`
	DirSrc  string       `json:"dir_src,omitempty"`
	Dirs    []string     `json:"dirs,omitempty"`
	Files   []string     `json:"files,omitempty"`
	Mode    FileProcMode `json:"mode,omitempty"`

	DirSkip  *[]string
	FileSkip *[]string
	Verbose  bool

	// Conf       *TypeConf
	// Flag       *TypeFlag
	// FlagUpdate *TypeFlagUpdate
}

func (df *TypeDotfile) New(dirSrc string, dirDest string, mode FileProcMode, dirSkip, fileSkip *[]string, verbose bool) {
	df.Initialized = true
	df.MyType = "TypeDotfile"
	prefix := df.MyType + ".New"

	if !(mode == Append || mode == Copy) {
		ezlog.Err().Name(prefix).Name("Mode error").Msg(mode).Out()
		return
	}

	df.DirDest = dirDest
	df.DirSrc = dirSrc
	df.Mode = mode

	df.DirSkip = dirSkip
	df.FileSkip = fileSkip
	df.Verbose = verbose

	// df.Conf = conf
	// df.Flag = flag
	// df.FlagUpdate = flagUpdate

	// cd to simplify path handling
	err := os.Chdir(df.DirSrc)
	if err == nil {
		df.Dirs, df.Files = df.dirFileGet(".")
	}
	ezlog.Debug().Name(prefix).Msg(df).Out()
}

func (df *TypeDotfile) Run() {
	prefix := df.MyType + ".Run"
	var e error
	for _, fileDir := range df.Dirs {
		e = dirCreateHidden(fileDir, df.DirDest)
		if e != nil {
			os.Exit(1)
		}
	}
	// Append/Copy files
	for _, filepathSrc := range df.Files {
		filepathDest := path.Join(df.DirDest, pathHide(filepathSrc))
		e = df.processFile(path.Join(df.DirSrc, filepathSrc), filepathDest)
		errs.Queue(prefix, e)
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
		same := file.FileSame(src, dest)
		if same {
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
	if ezlog.GetLogLevel() >= ezlog.DebugLevel {
		ezlog.Debug().Name(prefix).Msg(str).Out()
	} else if df.Verbose {
		ezlog.Log().Msg(str).Out()
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

// Create dotted/hidden directory
func dirCreateHidden(dir string, dirBase string) (e error) {
	var prefix = "DirCreate"
	if !(dir == "." || dir == "") {
		dirDest := path.Join(dirBase, pathHide(dir))
		if !file.IsDir(dirDest) {
			e = os.MkdirAll(dirDest, os.ModePerm)
			if e == nil {
				ezlog.Debug().Name(prefix).Name("created").Msg(dirDest).Out()
			} else {
				ezlog.Err().Name(prefix).Name("ERR").Msg(e).Out()
			}
		}
	}
	return e
}

// Get list of directory and list of file
func (df *TypeDotfile) dirFileGet(dir string) (dirs []string, files []string) {
	err := symwalk.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if p != "." && !str.ArrayContainsSubString(df.DirSkip, "/"+p+"/") {
				dirs = append(dirs, p)
			}
		} else {
			if !contains(*df.FileSkip, path.Base(p)) && !str.ArrayContainsSubString(df.DirSkip, "/"+p) {
				files = append(files, p)
			}
		}
		return nil
	})
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
