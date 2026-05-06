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

	"github.com/J-Siu/go-helper/v2/basestruct"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/spf13/viper"
)

var Default = TypeConf{
	FileConf: "$HOME/.config/go-dotfile.json",
}

type TypeConf struct {
	*basestruct.Base

	DirAP    []string `json:"DirAP,omitempty"`
	DirCP    []string `json:"DirCP,omitempty"`
	DirDest  string   `json:"DirDest,omitempty"`
	DirSkip  []string `json:"DirSkip,omitempty"`
	FileConf string   `json:"FileConf,omitempty"`
	FileSkip []string `json:"FileSkip,omitempty"`
}

func (t *TypeConf) New() {
	t.Base = new(basestruct.Base)
	t.Initialized = true
	t.MyType = "TypeConf"
	prefix := t.MyType + ".New"

	t.setDefault()
	ezlog.Debug().N(prefix).N("Default").Lm(t).Out()

	t.readFileConf()
	ezlog.Debug().N(prefix).N("Raw").Lm(t).Out()

	t.expand()
	ezlog.Debug().N(prefix).N("Expand").Lm(t).Out()

	// Check DirDest
	if !file.IsDir(t.DirDest) {
		ezlog.Err().N(prefix).N("DirDest does not exist").M(t.DirDest).Out()
		os.Exit(1)
	}
}

func (t *TypeConf) readFileConf() {
	t.MyType = "TypeConf"
	prefix := t.MyType + ".readFileConf"

	viper.SetConfigType("json")
	viper.SetConfigFile(file.TildeEnvExpand(t.FileConf))
	viper.AutomaticEnv()
	if t.Err = viper.ReadInConfig(); t.Err == nil {
		t.Err = viper.Unmarshal(&t)
	} else {
		ezlog.Debug().N(prefix).M(t.Err).Out()
		os.Exit(1)
	}
}

// Should be called before reading config file
func (t *TypeConf) setDefault() {
	if t.FileConf == "" {
		t.FileConf = Default.FileConf
	}
	t.DirDest, _ = os.UserHomeDir()
}

func (t *TypeConf) expand() {
	t.DirDest = file.TildeEnvExpand(t.DirDest)
	t.FileConf = file.TildeEnvExpand(t.FileConf)

	strArrays := [][]string{t.DirAP, t.DirCP, t.DirSkip, t.FileSkip}
	for _, arr := range strArrays {
		for i := range arr {
			arr[i] = file.TildeEnvExpand(arr[i])
		}
	}
}
