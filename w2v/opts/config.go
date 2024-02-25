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

	"github.com/fogfish/word2vec/trainer"
	"github.com/spf13/viper"
)

// Config Trainer from Config file
func NewConfig() trainer.Config {
	return trainer.Config{
		Corpus:   NewCorpusFromConfig(),
		Vector:   NewWordVectorFromConfig(),
		Learning: NewLearningFromConfig(),

		UseSkipGram: viper.GetBool("skip-gram.enabled"),
		UseCBOW:     viper.GetBool("cbow.enabled"),

		UseNegativeSampling:    viper.GetBool("negative-sampling.enabled"),
		SizeNegativeSampling:   viper.GetInt("negative-sampling.size"),
		UseHierarchicalSoftMax: viper.GetBool("hierarchical-softmax.enabled"),
	}
}

func NewCorpusFromConfig() trainer.ConfigCorpus {
	tokenizer := viper.GetString("corpus.tokenizer")
	fmt.Printf("%x", []byte(tokenizer))
	if tokenizer == "" {
		tokenizer = " \n,.-!?:;/\"#$%&'()*+<=>@[]\\^_`{|}~\t\v\f\r"
	}

	sequencer := viper.GetString("corpus.sequencer")
	if sequencer == "" {
		sequencer = ".\n?!"
	}

	return trainer.ConfigCorpus{
		StopWords: viper.GetString("corpus.stopwords"),
		Tokenizer: tokenizer,
		Sequencer: sequencer,
	}
}

func NewWordVectorFromConfig() trainer.ConfigWordVector {
	vector := viper.GetInt("word.vector")
	if vector == 0 {
		vector = 300
	}

	window := viper.GetInt("word.window")
	if window == 0 {
		window = 5
	}

	return trainer.ConfigWordVector{
		Vector:    vector,
		Window:    window,
		Threshold: viper.GetFloat64("word.threshold"),
		Frequency: viper.GetInt("word.frequency"),
	}
}

func NewLearningFromConfig() trainer.ConfigLearning {
	return trainer.ConfigLearning{
		Epoch: viper.GetInt("learning.epoch"),
		Rate:  viper.GetFloat64("learning.rate"),
	}
}
