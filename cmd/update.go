/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/J-Siu/go-dotfile/lib"
	"github.com/J-Siu/go-helper"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"u", "up"},
	Short:   "Update dotfiles",
	Run: func(cmd *cobra.Command, args []string) {
		helper.Debug = lib.Flag.Debug

		dirDest := lib.Conf.DirDest
		if dirDest == "" {
			dirDest, _ = os.UserHomeDir()
		}
		for _, dir := range lib.Conf.DirCP {
			var df lib.TypeDotfile
			df.Init(dir, dirDest, lib.ProcModeCopy)
			df.Process()
		}
		for _, dir := range lib.Conf.DirAP {
			var df lib.TypeDotfile
			df.Init(dir, dirDest, lib.ProcModeAppend)
			df.Process()
		}
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
