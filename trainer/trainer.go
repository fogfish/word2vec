//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/word2vec
//

package trainer

/*
#cgo CFLAGS: -I../libw2v/include
#cgo LDFLAGS: -L . -lw2v
#include <stdlib.h>
#include "w2v.h"
*/
import "C"
import (
	"fmt"
	"unsafe"

	"github.com/spf13/viper"
)

type ConfigCorpus struct {
	// filename of a train text corpus
	dataset string

	// filename of the stop-words set [optional]
	stopwords string

	// word tokenizer
	tokenizer string

	// sequence tokenizer
	sequencer string
}

func NewCorpusFromConfig() ConfigCorpus {
	tokenizer := viper.GetString("corpus.tokenizer")
	fmt.Printf("%x", []byte(tokenizer))
	if tokenizer == "" {
		tokenizer = " \n,.-!?:;/\"#$%&'()*+<=>@[]\\^_`{|}~\t\v\f\r"
	}

	sequencer := viper.GetString("corpus.sequencer")
	if sequencer == "" {
		sequencer = ".\n?!"
	}

	return ConfigCorpus{
		stopwords: viper.GetString("corpus.stopwords"),
		tokenizer: tokenizer,
		sequencer: sequencer,
	}
}

func NewCorpusDefault() ConfigCorpus {
	return ConfigCorpus{
		tokenizer: " \n,.-!?:;/\"#$%&'()*+<=>@[]\\^_`{|}~\t\v\f\r",
		sequencer: ".\n?!",
	}
}

type ConfigWordVector struct {
	// words vector dimension
	vector int

	// nearby words frame or window
	window int

	// threshold for occurrence of words
	threshold float64

	// exclude words that appear less than [value] times from vocabulary
	frequency int
}

func NewWordVectorFromConfig() ConfigWordVector {
	vector := viper.GetInt("word.vector")
	if vector == 0 {
		vector = 300
	}

	window := viper.GetInt("word.window")
	if window == 0 {
		window = 5
	}

	return ConfigWordVector{
		vector:    vector,
		window:    window,
		threshold: viper.GetFloat64("word.threshold"),
		frequency: viper.GetInt("word.frequency"),
	}
}

func NewWordVectorDefault() ConfigWordVector {
	return ConfigWordVector{
		vector:    300,
		window:    5,
		threshold: 1e-3,
		frequency: 5,
	}
}

type ConfigLearning struct {
	// number of training iterations, epoch
	epoch int

	// starting learning rate
	rate float64
}

func NewLearningFromConfig() ConfigLearning {
	return ConfigLearning{
		epoch: viper.GetInt("learning.epoch"),
		rate:  viper.GetFloat64("learning.rate"),
	}
}

func NewLearningDefault() ConfigLearning {
	return ConfigLearning{
		epoch: 5,
		rate:  0.05,
	}
}

type Trainer struct {
	output   string
	corpus   ConfigCorpus
	vector   ConfigWordVector
	learning ConfigLearning

	// choose of the learning model:
	//  - Continuous Bag of Words (CBOW)
	//  - Skip-Gram
	useSkipGram bool
	useCBOW     bool

	// the computationally efficient approximation
	//  - Negative Sampling (NS)
	//  - Hierarchical Softmax (HS)
	useNegativeSampling    bool
	useHierarchicalSoftMax bool

	// number of negative examples (NS option)
	sizeNegativeSampling int

	threads int
	verbose bool

	h unsafe.Pointer
}

type Option func(*Trainer)

func WithDefault() Option {
	return func(t *Trainer) {
		t.corpus = NewCorpusDefault()
		t.vector = NewWordVectorDefault()
		t.learning = NewLearningDefault()

		t.useSkipGram = true
		t.useCBOW = false

		t.useNegativeSampling = true
		t.useHierarchicalSoftMax = false
		t.sizeNegativeSampling = 5

		t.threads = 12
		t.verbose = true
	}
}

func WithConfigFile() Option {
	return func(t *Trainer) {
		t.corpus = NewCorpusFromConfig()
		t.vector = NewWordVectorFromConfig()
		t.learning = NewLearningFromConfig()

		t.useSkipGram = viper.GetBool("skip-gram.enabled")
		t.useCBOW = viper.GetBool("cbow.enabled")

		t.useNegativeSampling = viper.GetBool("negative-sampling.enabled")
		t.sizeNegativeSampling = viper.GetInt("negative-sampling.size")
		t.useHierarchicalSoftMax = viper.GetBool("hierarchical-softmax.enabled")
	}
}

func WithCorpusDataset(dataset string) Option {
	return func(t *Trainer) {
		t.corpus.dataset = dataset
	}
}

func WithOutput(output string) Option {
	return func(t *Trainer) {
		t.output = output
	}
}

func WithVerbose(verbose bool) Option {
	return func(t *Trainer) {
		t.verbose = verbose
	}
}

func WithThreads(threads int) Option {
	return func(t *Trainer) {
		t.threads = threads
	}
}

//
//

func Train(opts ...Option) error {
	var w2v Trainer
	WithDefault()(&w2v)
	for _, opt := range opts {
		opt(&w2v)
	}

	dataset := C.CString(w2v.corpus.dataset)
	defer C.free(unsafe.Pointer(dataset))

	fileStopWords := C.CString(w2v.corpus.stopwords)
	defer C.free(unsafe.Pointer(fileStopWords))

	fileModel := C.CString(w2v.output)
	defer C.free(unsafe.Pointer(fileModel))

	withHS := C.uchar(0)
	if w2v.useHierarchicalSoftMax {
		withHS = C.uchar(1)
	}

	withSG := C.uchar(0)
	if w2v.useSkipGram {
		withSG = C.uchar(1)
	}

	tokenizer := C.CString(w2v.corpus.tokenizer)
	defer C.free(unsafe.Pointer(tokenizer))

	sequencer := C.CString(w2v.corpus.sequencer)
	defer C.free(unsafe.Pointer(sequencer))

	verbose := C.uchar(0)
	if w2v.verbose {
		verbose = C.uchar(1)
	}

	w2v.h = C.Train(
		dataset,
		fileStopWords,
		fileModel,
		C.ushort(w2v.vector.frequency),
		C.ushort(w2v.vector.vector),
		C.uchar(w2v.vector.window),
		C.float(w2v.vector.threshold),
		withHS,
		C.uint8_t(w2v.sizeNegativeSampling),
		C.uint8_t(12),
		C.uint8_t(w2v.learning.epoch),
		C.float(w2v.learning.rate),
		withSG,
		tokenizer,
		sequencer,
		verbose,
	)
	if uintptr(w2v.h) == 0 {
		return fmt.Errorf("unable to train model")
	}

	return nil
}
