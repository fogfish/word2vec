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
	"strings"

	"github.com/fogfish/word2vec"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(lookupCmd)
	lookupCmd.Flags().StringVarP(&lookupModel, "model", "m", "", "path to trained word2vec model")
	lookupCmd.Flags().IntVarP(&lookupVecSize, "vector", "v", 300, "vector size")
	lookupCmd.Flags().IntVarP(&lookupK, "size", "k", 30, "number of nearest elements")
}

var (
	lookupModel   string
	lookupVecSize int
	lookupK       int
)

var lookupCmd = &cobra.Command{
	Use:   "lookup",
	Short: "Query (Lookup) word2vec model for neighbor tokens (words).",
	Long: `
Query (Lookup) word2vec model for neighbor tokens (words).
	`,
	Example: `
  w2v lookup -m wap-v300w5e5s1h005-en.bin alexander
  w2v lookup -m wap-v300w5e5s1h005-en.bin emperor
	`,
	RunE: lookup,
	Args: cobra.MinimumNArgs(1),
}

func lookup(cmd *cobra.Command, args []string) error {
	w2v, err := word2vec.Load(lookupModel, embeddingVecSize)
	if err != nil {
		return err
	}

	seq := make([]word2vec.Nearest, lookupK)

	w2v.Lookup(strings.Join(args, " "), seq)
	for _, n := range seq {
		fmt.Printf("%15s : %.6f\n", n.Word, n.Distance)
	}

	return nil
}
