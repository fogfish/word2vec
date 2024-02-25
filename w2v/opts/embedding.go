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
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fogfish/hnsw/vector"
	"github.com/fogfish/word2vec"
	"github.com/kshard/fvecs"
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
Vectors are written to .fvecs file next to input text. The text used for
the assessment of embeddings is written to .bvecs file. 
	`,
	Example: `
  w2v embedding -m wap-v300w5e10s1h010-en.bin doc/leo-tolstoy-war-and-peace-en.txt
	`,
	RunE: embedding,
	Args: cobra.MinimumNArgs(1),
}

type Node struct {
	ID     int
	Vector vector.V32
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
	fdata := strings.TrimSuffix(text, filepath.Ext(text)) + ".fvecs"
	bdata := strings.TrimSuffix(text, filepath.Ext(text)) + ".bvecs"

	feg, err := os.Create(fdata)
	if err != nil {
		return err
	}
	defer feg.Close()
	fw := fvecs.NewEncoder[float32](feg)

	beg, err := os.Create(bdata)
	if err != nil {
		return err
	}
	defer beg.Close()
	bw := fvecs.NewEncoder[byte](beg)

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
			slog.Warn("skip", "text", t)
			continue
		}

		cnt++
		vol += len(txt)

		// if cnt%1000000 == 0 {
		// 	os.Stderr.WriteString(
		// 		fmt.Sprintf("%v\t%d\t%d\t%d\t%d\n", time.Since(t), cnt, vol, int(time.Since(t).Nanoseconds())/cnt, vol/cnt),
		// 	)
		// }

		if err := fw.Write(vec); err != nil {
			return err
		}
		if err := bw.Write([]byte(txt)); err != nil {
			return err
		}

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
