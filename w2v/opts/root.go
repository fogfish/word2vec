//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/word2vec
//

package opts

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Execute is entry point for cobra cli application
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		e := err.Error()
		fmt.Println(strings.ToUpper(e[:1]) + e[1:])
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "w2v",
	Short: "CLI for word2vec model training and execution of simple use-cases.",
	Long: `
CLI for word2vec model training and execution of simple use-cases.
	`,
	Run: root,
}

func root(cmd *cobra.Command, args []string) {
	cmd.Help()
}
