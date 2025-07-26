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
	"regexp"

	"github.com/spf13/viper"
)

func tildeEnvExpand(strIn string) (strOut string) {
	re := regexp.MustCompile(`^~/`)
	strOut = re.ReplaceAllString(strIn, "$$HOME/")
	strOut = os.ExpandEnv(strOut)

	// prefix := "tildeEnvExpand"
	// helper.ReportDebug(strIn+" -> "+strOut, prefix, false, true)

	return strOut
}

const (
	ConfDirAP   = "DirAP"
	ConfDirCP   = "DirCP"
	ConfDirDest = "DirDest"
)

type TypeConf struct {
	DirAP    []string `json:"DirAP"`
	DirCP    []string `json:"DirCP"`
	DirDest  string   `json:"DirDest"`
	File     string   `json:"-"`
	FileSkip []string `json:"FileSkip"`
}

func (s *TypeConf) Init() {
	s.File = viper.ConfigFileUsed()
	viper.Unmarshal(&s)
	for i := range s.DirAP {
		s.DirAP[i] = tildeEnvExpand(s.DirAP[i])
	}
	for i := range s.DirCP {
		s.DirCP[i] = tildeEnvExpand(s.DirCP[i])
	}
	for i := range s.FileSkip {
		s.FileSkip[i] = tildeEnvExpand(s.FileSkip[i])
	}
	s.DirDest = tildeEnvExpand(s.DirDest)
	s.File = tildeEnvExpand(s.File)
}
