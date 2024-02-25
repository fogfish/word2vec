//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/word2vec
//

package opts

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/fogfish/word2vec/trainer"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(trainCmd)
	trainCmd.Flags().StringVarP(&trainConfig, "config", "C", "", "configure word2vec model")
	trainCmd.Flags().StringVarP(&trainCorpus, "corpus", "f", "", "training corpus")
	trainCmd.Flags().StringVarP(&trainOutput, "output", "o", "", "output model")
	trainCmd.Flags().IntVarP(&trainThreads, "threads", "t", 12, "number of training threads")
	trainCmd.Flags().BoolVar(&trainSilent, "silent", false, "silent training")

	trainCmd.AddCommand(configCmd)
}

var (
	trainConfig  string
	trainCorpus  string
	trainOutput  string
	trainThreads int
	trainSilent  bool
)

var trainCmd = &cobra.Command{
	Use:   "train",
	Short: "train word2vec model",
	Long: `
Train word2vec model from own corpus. The training process requires:
* text corpus on target language
* stop words, obtain from https://github.com/stopwords-iso

Configure training process through config file:

  w2v train config > wap-en.yaml

The default params gives sufficient results but feel free to tune them.

Consider naming of the model after parameters used for training:
* "v" vector size
* "w" nearby words window 
* "e" training epoch
* architecture skip-gram "s1" or CBoW "s0"
* algorithm H. softmax "h1", N. Sampling "h0"
	`,
	Example: `
cmd train -C wap-en.yaml \
  -o wap-v300w5e10s1h010-en.bin \
  -f ../doc/leo-tolstoy-war-and-peace-en.txt
	`,
	RunE: train,
}

func train(cmd *cobra.Command, args []string) error {
	if err := configure(trainConfig); err != nil {
		return err
	}

	cfg := NewConfig()
	cfg.Corpus.Dataset = trainCorpus
	cfg.Output = trainOutput
	cfg.Threads = trainThreads
	cfg.Verbose = !trainSilent

	err := trainer.Train(cfg)
	if err != nil {
		return err
	}

	return nil
}

//
//
//

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Generate config file with default params for training",
	Long: `
Generate config file with default params for training
	`,
	RunE: configFile,
}

func configFile(cmd *cobra.Command, args []string) error {
	s := `##
## word2vec configuration

##
## Corpus specification
corpus:
  ## These words will be excluded from training vocabulary.
  ## Stop-words are separated by any of word delimiter char (see below).
  ## stopwords: ./path-to/stopwords.txt

  ##
  ## Words delimiter chars
  tokenizer: " \n,.-!?:;/\"#$%&'()*+<=>@[]\\^_` + "`" + `{|}~\t\v\f\r"

  ##
  ## End of sentence chars
  sequencer: ".\n?!"

##
## "Word" / "Token" specification
word:
  ## Words vector dimension. 
  ## Large vectors are usually better, but it requires more training data.
  vector: 300

  ##
  ## Nearby words window.
  ## It defines how many words we will include in training of the word
  ## inside of corpus - [value] words behind and [value] words ahead.
  window: 5

  ##
  ## threshold for occurrence of words. 
  ## The value used for down-sampling the dataset, words with higher frequency
  ## in the training data will be randomly down-sampled. This parameter controls
  ## the down-sampling algorithm
  threshold: 1e-3

  ##
  ## exclude words that appear less than [value] times from vocabulary,
  frequency: 0

##
## Learning parameters
learning:
  ## Number of training iterations.
  ## More iterations makes a more precise model, but computational cost is
  ## linearly proportional to iterations.
  epoch: 5

  ##
  ## Starting learning rate
  rate: 0.05

##
## Use Skip-Gram algorithms
skip-gram:
  enabled: true

##
## Use Continuous Bag of Words (CBOW)
cbow:
  enabled: false

##
## Use Negative Sampling the computationally efficient approximation
## umber of negative examples (NS option), default value is 5. Values in the range 5–20 are useful for small training datasets, while for large datasets the value can be as small as 2–5. 
negative-sampling:
  enabled: true

  ## size of negative samples.
  ## Values in the range [5, 20] are useful for small training datasets.
  ## Large datasets the value can be as small as [2, 5]. 
  size: 5

##
## Use Hierarchical Softmax the computationally efficient approximation
hierarchical-softmax:
  enabled: false
`

	_, err := os.Stdout.WriteString(s)
	return err
}

// configure word2vec training process
func configure(config string) error {
	if config == "" {
		return errors.New("config file is not defined")
	}

	path := filepath.Dir(config)
	viper.AddConfigPath(path)
	viper.AddConfigPath("$HOME/.config/")
	viper.AddConfigPath(".")

	file := strings.TrimSuffix(
		filepath.Base(config),
		filepath.Ext(config),
	)

	viper.SetConfigName(file)
	viper.SetConfigType("yaml")
	return viper.ReadInConfig()
}
