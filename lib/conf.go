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

	"github.com/J-Siu/go-basestruct"
	"github.com/J-Siu/go-helper/v2/ezlog"
	"github.com/J-Siu/go-helper/v2/file"
	"github.com/spf13/viper"
)

var Default = TypeConf{
	FileConf: "$HOME/.config/go-dotfile.json",
}

const (
	ConfDirAP   = "DirAP"
	ConfDirCP   = "DirCP"
	ConfDirDest = "DirDest"
)

type TypeConf struct {
	basestruct.Base

	DirAP    []string `json:"DirAP"`
	DirCP    []string `json:"DirCP"`
	DirDest  string   `json:"DirDest"`
	DirSkip  []string `json:"DirSkip"`
	FileConf string   `json:"FileConf"`
	FileSkip []string `json:"FileSkip"`
}

func (c *TypeConf) New() {
	c.Initialized = true
	c.MyType = "TypeConf"
	prefix := c.MyType + ".New"

	c.setDefault()
	ezlog.Debug().N(prefix).Nn("Default").M(c).Out()

	c.readFileConf()
	ezlog.Debug().N(prefix).Nn("Raw").M(c).Out()

	// TODO: add flag

	c.expand()
	ezlog.Debug().N(prefix).Nn("Expand").M(c).Out()

	// Check DirDest
	if !file.IsDir(c.DirDest) {
		ezlog.Err().N(prefix).N("DirDest does not exist").M(c.DirDest).Out()
		os.Exit(1)
	}

}

func (c *TypeConf) readFileConf() {
	c.MyType = "TypeConf"
	prefix := c.MyType + ".readFileConf"

	viper.SetConfigType("json")
	viper.SetConfigFile(file.TildeEnvExpand(c.FileConf))
	viper.AutomaticEnv()
	c.Err = viper.ReadInConfig()

	if c.Err == nil {
		c.Err = viper.Unmarshal(&c)
	} else {
		ezlog.Debug().N(prefix).M(c.Err).Out()
		os.Exit(1)
	}
}

// Should be called before reading config file
func (c *TypeConf) setDefault() {
	if c.FileConf == "" {
		c.FileConf = Default.FileConf
	}
	c.DirDest, _ = os.UserHomeDir()
}

func (c *TypeConf) expand() {
	c.DirDest = file.TildeEnvExpand(c.DirDest)
	c.FileConf = file.TildeEnvExpand(c.FileConf)
	for i := range c.DirAP {
		c.DirAP[i] = file.TildeEnvExpand(c.DirAP[i])
	}
	for i := range c.DirCP {
		c.DirCP[i] = file.TildeEnvExpand(c.DirCP[i])
	}
	for i := range c.DirSkip {
		c.DirSkip[i] = file.TildeEnvExpand(c.DirSkip[i])
	}
	for i := range c.FileSkip {
		c.FileSkip[i] = file.TildeEnvExpand(c.FileSkip[i])
	}
}
