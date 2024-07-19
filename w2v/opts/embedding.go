//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/word2vec
//

package opts

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fogfish/word2vec"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(embeddingCmd)
	embeddingCmd.Flags().StringVarP(&embeddingModel, "model", "m", "", "path to trained word2vec model")
	embeddingCmd.Flags().IntVarP(&embeddingVecSize, "vector", "v", 300, "vector size")
}

var (
	embeddingModel   string
	embeddingVecSize int
)

var embeddingCmd = &cobra.Command{
	Use:   "embedding",
	Short: "Calculate embedding for input text",
	Long: `
Calculates embedding for input text. 

For each paragraph (split by \n) from input text, it calculates embeddings vector.
Vectors are written to standard output.
`,
	Example: `
  w2v embedding -m wap-v300_w5_e5_s1_h010-en.bin doc/leo-tolstoy-war-and-peace-en.txt
	`,
	RunE: embedding,
	Args: cobra.MinimumNArgs(1),
}

func embedding(cmd *cobra.Command, args []string) error {
	w2v, err := word2vec.Load(embeddingModel, embeddingVecSize)
	if err != nil {
		return err
	}

	for _, text := range args {
		if err := embeddingText(w2v, text); err != nil {
			return err
		}
	}

	return nil
}

func embeddingText(w2v word2vec.Model, text string) error {
	fd, err := os.Open(text)
	if err != nil {
		return err
	}
	defer fd.Close()

	t := time.Now()
	cnt := 1
	vol := 0
	vec := make([]float32, embeddingVecSize)

	scanner := bufio.NewScanner(fd)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	for scanner.Scan() {
		txt := strings.Trim(
			scanner.Text(),
			" \n,.-!?:;/\"#$%&'()*+<=>@[]\\^_`{|}~\t\v\f\r",
		)

		if len(txt) == 0 {
			continue
		}

		if err := w2v.Embedding(txt, vec); err != nil {
			continue
		}

		cnt++
		vol += len(txt)

		seq := make([]string, len(vec)+1)
		seq[0] = txt
		for i, x := range vec {
			seq[i+1] = fmt.Sprintf("%f", x)
		}

		os.Stdout.WriteString(strings.Join(seq, " "))
		os.Stdout.WriteString("\n")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	os.Stderr.WriteString(
		fmt.Sprintf("==> %s\n", text),
	)
	os.Stderr.WriteString(
		fmt.Sprintf("\tvectors: %v\n", cnt),
	)
	os.Stderr.WriteString(
		fmt.Sprintf("\t   time: %v\n", time.Since(t)),
	)
	os.Stderr.WriteString(
		fmt.Sprintf("\t  op/ns: %v\n", int(time.Since(t).Nanoseconds())/cnt),
	)

	return nil
}
