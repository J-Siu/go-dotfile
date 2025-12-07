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
	"time"

	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/errs"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/J-Siu/go-helper/v2/str"
	"github.com/edwardrf/symwalk"
)

type FileProcMode int8

const (
	APPEND FileProcMode = iota
	COPY
	SKIP
)

type TypeDotfileProperty struct {
	DirDest  *string      `json:"DirDest"`
	DirSkip  *[]string    `json:"DirSkip"`
	DirSrc   *string      `json:"DirSrc"`
	FileSkip *[]string    `json:"FileSkip"`
	Mode     FileProcMode `json:"Mode"`
	Save     bool         `json:"Save"`
}

type TypeDotfile struct {
	*basestruct.Base
	*TypeDotfileProperty
	// --- calculate in Run()
	Dirs   *[]string   `json:"Dirs"`
	Files  *[]string   `json:"Files"`
	Result []*[]string `json:"Result"`
}

func (t *TypeDotfile) New(property *TypeDotfileProperty) *TypeDotfile {
	t.Base = new(basestruct.Base)
	t.Initialized = true
	t.MyType = "TypeDotfile"
	prefix := t.MyType + ".New"

	t.TypeDotfileProperty = property

	ezlog.Debug().N(prefix).M(t).Out()

	return t
}

func (t *TypeDotfile) Run() {
	prefix := t.MyType + ".Run"
	var e error
	// cd to simplify path handling
	t.Err = os.Chdir(*t.DirSrc)
	if t.Err == nil {
		t.Dirs, t.Files = t.getDirFile(".")
		ezlog.Debug().N(prefix).N("Dirs").Lm(t.Dirs).Out()
		ezlog.Debug().N(prefix).N("Files").Lm(t.Files).Out()
	}
	// create dirs
	if t.Err == nil && t.Dirs != nil {
		for _, fileDir := range *t.Dirs {
			t.Err = dirCreateHidden(fileDir, *t.DirDest)
			if t.Err != nil {
				errs.Queue(prefix, t.Err)
				break
			}
		}
	}
	// Append/Copy files
	if t.Err == nil && t.Files != nil {
		for _, filepathSrc := range *t.Files {
			filepathDest := path.Join(*t.DirDest, pathHide(filepathSrc))
			e = t.processFile(path.Join(*t.DirSrc, filepathSrc), filepathDest)
			errs.Queue(prefix, e)
		}
	}
}

// Process file base on Mode(append|copy)
//
// Not using TypeDotfile.Err
func (t *TypeDotfile) processFile(src, dest string) (err error) {
	// prefix := t.MyType + ".processFile"

	var (
		data          []byte
		desInfo       os.FileInfo
		desModTime    time.Time
		desModTimeStr string       = "---------- --:--:--"
		fileProcMode  FileProcMode = t.Mode
		srcInfo       os.FileInfo
		srcModTime    time.Time
		srcModTimeStr string = "---------- --:--:--"
		srcPermission os.FileMode
		timeFormat    string = "2006-01-02 15:04:05"
	)

	// Get File info before change

	desInfo, err = os.Stat(dest)
	if err == nil {
		desModTime = desInfo.ModTime()
		desModTimeStr = desModTime.Local().Format(timeFormat)
	}
	err = nil // Resetting err, as dest may not exist.

	srcInfo, err = os.Stat(src)
	if err == nil {
		srcPermission = srcInfo.Mode()
		srcModTime = srcInfo.ModTime()
		srcModTimeStr = srcModTime.Local().Format(timeFormat)
	}

	if fileProcMode == COPY && file.FileSame(src, dest) {
		fileProcMode = SKIP
	}

	if fileProcMode != SKIP && t.Save {

		// Read source file
		if err == nil {
			data, err = os.ReadFile(src)
		}

		// Append: add newline to destination file
		if err == nil {
			if fileProcMode == APPEND {
				b := []byte("\n")
				err = file.AppendByte(dest, &b)
				if err == nil {
					err = file.AppendByte(dest, &data)
				}
			} else {
				err = file.WriteByte(dest, &data, srcPermission)
			}
		}

		// Set dest permission
		if err == nil && fileProcMode == COPY {
			os.Chtimes(dest, srcModTime, srcModTime)
		}
	}

	if err == nil {
		var result []string
		result = append(result,
			fileProcMode.String(),
			srcPermission.String(),
			srcModTimeStr,
			src,
			"->",
			desModTimeStr,
			dest)

		t.Result = append(t.Result, &result)
	}

	return err
}

// Get list of directory and list of file
func (t *TypeDotfile) getDirFile(dir string) (*[]string, *[]string) {
	var (
		dirs  []string
		files []string
	)
	symwalk.Walk(dir, func(p string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if p != "." && !str.ArrayContainsSubString(t.DirSkip, "/"+p+"/", false) {
				dirs = append(dirs, p)
			}
		} else {
			base := path.Base(p)
			if !str.ArrayContains(t.FileSkip, &base, false) && !str.ArrayContainsSubString(t.DirSkip, "/"+p, false) {
				files = append(files, p)
			}
		}
		return nil
	})
	return &dirs, &files
}

// Create dotted/hidden directory
func dirCreateHidden(dir, dirBase string) (e error) {
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

// Add "."" in front of path if there is none
func pathHide(p string) string {
	if strings.HasPrefix(p, ".") {
		return p
	}
	return "." + p
}
