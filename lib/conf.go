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

	"github.com/J-Siu/go-helper"
	"github.com/spf13/viper"
)

const (
	ConfDirAP   = "DirAP"
	ConfDirCP   = "DirCP"
	ConfDirDest = "DirDest"
)

type TypeConf struct {
	Err    error
	init   bool
	myType string

	DirAP    []string `json:"DirAP"`
	DirCP    []string `json:"DirCP"`
	DirDest  string   `json:"DirDest"`
	DirSkip  []string `json:"DirSkip"`
	FileConf string   `json:"FileConf"`
	FileSkip []string `json:"FileSkip"`
}

func (c *TypeConf) Init() {
	c.init = true
	c.myType = "TypeConf"
	prefix := c.myType + ".Init"

	c.setDefault()
	helper.ReportDebug(c.FileConf, prefix+": Config file", false, true)

	c.readFileConf()
	helper.ReportDebug(c, prefix+": Raw", false, true)

	c.expand()
	helper.ReportDebug(c, prefix+": Expand", false, true)
}

func (c *TypeConf) readFileConf() {
	c.myType = "TypeConf"
	prefix := c.myType + ".readFileConf"

	viper.SetConfigType("json")
	viper.SetConfigFile(helper.TildeEnvExpand(Conf.FileConf))
	viper.AutomaticEnv()
	c.Err = viper.ReadInConfig()

	if c.Err == nil {
		c.Err = viper.Unmarshal(&c)
	} else {
		helper.Report(c.Err.Error(), prefix, true, true)
		os.Exit(1)
	}
}

// Should be called before reading config file
func (c *TypeConf) setDefault() {
	c.DirDest, _ = os.UserHomeDir()
	c.FileConf = Default.FileConf
}

func (c *TypeConf) expand() {
	c.DirDest = helper.TildeEnvExpand(c.DirDest)
	c.FileConf = helper.TildeEnvExpand(c.FileConf)
	for i := range c.DirAP {
		c.DirAP[i] = helper.TildeEnvExpand(c.DirAP[i])
	}
	for i := range c.DirCP {
		c.DirCP[i] = helper.TildeEnvExpand(c.DirCP[i])
	}
	for i := range c.DirSkip {
		c.DirSkip[i] = helper.TildeEnvExpand(c.DirSkip[i])
	}
	for i := range c.FileSkip {
		c.FileSkip[i] = helper.TildeEnvExpand(c.FileSkip[i])
	}
}
