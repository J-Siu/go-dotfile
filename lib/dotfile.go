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

	"github.com/J-Siu/go-helper/v2/basestruct"
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
	*basestruct.Base

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

func (t *TypeDotfile) New(dirSrc string, dirDest string, mode FileProcMode, dirSkip, fileSkip *[]string, verbose bool) {
	t.Base = new(basestruct.Base)
	t.Initialized = true
	t.MyType = "TypeDotfile"
	prefix := t.MyType + ".New"

	if !(mode == Append || mode == Copy) {
		ezlog.Err().N(prefix).N("Mode error").M(mode).Out()
		return
	}

	t.DirDest = dirDest
	t.DirSrc = dirSrc
	t.Mode = mode

	t.DirSkip = dirSkip
	t.FileSkip = fileSkip
	t.Verbose = verbose

	// df.Conf = conf
	// df.Flag = flag
	// df.FlagUpdate = flagUpdate

	// cd to simplify path handling
	err := os.Chdir(t.DirSrc)
	if err == nil {
		t.Dirs, t.Files = t.dirFileGet(".")
	}
	ezlog.Debug().N(prefix).M(t).Out()
}

func (t *TypeDotfile) Run() {
	prefix := t.MyType + ".Run"
	var e error
	for _, fileDir := range t.Dirs {
		e = dirCreateHidden(fileDir, t.DirDest)
		if e != nil {
			os.Exit(1)
		}
	}
	// Append/Copy files
	for _, filepathSrc := range t.Files {
		filepathDest := path.Join(t.DirDest, pathHide(filepathSrc))
		e = t.processFile(path.Join(t.DirSrc, filepathSrc), filepathDest)
		errs.Queue(prefix, e)
	}
}

// Process file base on Mode(append|copy)
//
// Not using TypeDotfile.Err
func (t *TypeDotfile) processFile(src string, dest string) (err error) {
	prefix := t.MyType + ".ProcessFile"

	fileProcMode := t.Mode
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
		ezlog.Debug().N(prefix).M(str).Out()
	} else if t.Verbose {
		ezlog.Log().M(str).Out()
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
				ezlog.Debug().N(prefix).N("created").M(dirDest).Out()
			} else {
				ezlog.Err().N(prefix).N("ERR").M(e).Out()
			}
		}
	}
	return e
}

// Get list of directory and list of file
func (t *TypeDotfile) dirFileGet(dir string) (dirs []string, files []string) {
	err := symwalk.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if p != "." && !str.ArrayContainsSubString(t.DirSkip, "/"+p+"/") {
				dirs = append(dirs, p)
			}
		} else {
			if !contains(*t.FileSkip, path.Base(p)) && !str.ArrayContainsSubString(t.DirSkip, "/"+p) {
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
