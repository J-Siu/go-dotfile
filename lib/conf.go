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
	DirAP    []string `json:"DirAP"`
	DirCP    []string `json:"DirCP"`
	DirDest  string   `json:"DirDest"`
	DirSkip  []string `json:"DirSkip"`
	FileConf string   `json:"FileConf"`
	FileSkip []string `json:"FileSkip"`
}

func (conf *TypeConf) Init() {
	prefix := "TypeConf.Init"

	conf.setDefault()

	helper.ReportDebug(conf.FileConf, prefix+": Config file", false, true)

	conf.readFileConf()

	helper.ReportDebug(conf, prefix+": Raw", false, true)

	conf.expand()

	helper.ReportDebug(conf, prefix+": Expand", false, true)
}

func (conf *TypeConf) readFileConf() {
	prefix := "TypeConf.readFileConf"
	viper.SetConfigType("json")
	viper.SetConfigFile(helper.TildeEnvExpand(Conf.FileConf))
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		viper.Unmarshal(&conf)
	} else {
		helper.Report(err.Error(), prefix, true, true)
		os.Exit(1)
	}
}

func (conf *TypeConf) setDefault() {
	if Conf.FileConf == "" {
		Conf.FileConf = Default.FileConf
	}
}

func (conf *TypeConf) expand() {
	conf.DirDest = helper.TildeEnvExpand(conf.DirDest)
	conf.FileConf = helper.TildeEnvExpand(conf.FileConf)
	for i := range conf.DirAP {
		conf.DirAP[i] = helper.TildeEnvExpand(conf.DirAP[i])
	}
	for i := range conf.DirCP {
		conf.DirCP[i] = helper.TildeEnvExpand(conf.DirCP[i])
	}
	for i := range conf.DirSkip {
		conf.DirSkip[i] = helper.TildeEnvExpand(conf.DirSkip[i])
	}
	for i := range conf.FileSkip {
		conf.FileSkip[i] = helper.TildeEnvExpand(conf.FileSkip[i])
	}
}
