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
)

type ConfigCorpus struct {
	// filename of a train text corpus
	Dataset string

	// filename of the stop-words set [optional]
	StopWords string

	// word tokenizer
	Tokenizer string

	// sequence tokenizer
	Sequencer string
}

func NewCorpusDefault() ConfigCorpus {
	return ConfigCorpus{
		Tokenizer: " \n,.-!?:;/\"#$%&'()*+<=>@[]\\^_`{|}~\t\v\f\r",
		Sequencer: ".\n?!",
	}
}

type ConfigWordVector struct {
	// words vector dimension
	Vector int

	// nearby words frame or window
	Window int

	// threshold for occurrence of words
	Threshold float64

	// exclude words that appear less than [value] times from vocabulary
	Frequency int
}

func NewWordVectorDefault() ConfigWordVector {
	return ConfigWordVector{
		Vector:    300,
		Window:    5,
		Threshold: 1e-3,
		Frequency: 5,
	}
}

type ConfigLearning struct {
	// number of training iterations, epoch
	Epoch int

	// starting learning rate
	Rate float64
}

func NewLearningDefault() ConfigLearning {
	return ConfigLearning{
		Epoch: 5,
		Rate:  0.05,
	}
}

type Config struct {
	Corpus   ConfigCorpus
	Vector   ConfigWordVector
	Learning ConfigLearning

	// choose of the learning model:
	//  - Continuous Bag of Words (CBOW)
	//  - Skip-Gram
	UseSkipGram bool
	UseCBOW     bool

	// the computationally efficient approximation
	//  - Negative Sampling (NS)
	//  - Hierarchical Softmax (HS)
	UseNegativeSampling    bool
	UseHierarchicalSoftMax bool

	// number of negative examples (NS option)
	SizeNegativeSampling int

	Output  string
	Threads int
	Verbose bool
}

func NewConfigDefault() Config {
	return Config{
		Corpus:      NewCorpusDefault(),
		Vector:      NewWordVectorDefault(),
		Learning:    NewLearningDefault(),
		UseSkipGram: true,
		UseCBOW:     false,

		UseNegativeSampling:    true,
		UseHierarchicalSoftMax: false,
		SizeNegativeSampling:   5,

		Threads: 12,
		Verbose: true,
	}
}

type Trainer struct {
	config Config
	h      unsafe.Pointer
}

// type Option func(*ConfigTrainer)

// func WithDefault() Option {
// 	return func(t *ConfigTrainer) {
// 		t.config = NewTrainerDefault()
// 	}
// }

// func WithCorpusDataset(dataset string) Option {
// 	return func(t *Trainer) {
// 		t.corpus.dataset = dataset
// 	}
// }

// func WithOutput(output string) Option {
// 	return func(t *Trainer) {
// 		t.output = output
// 	}
// }

// func WithVerbose(verbose bool) Option {
// 	return func(t *Trainer) {
// 		t.verbose = verbose
// 	}
// }

// func WithThreads(threads int) Option {
// 	return func(t *Trainer) {
// 		t.threads = threads
// 	}
// }

//
//

func Train(config Config) error {
	var w2v Trainer
	w2v.config = config

	dataset := C.CString(w2v.config.Corpus.Dataset)
	defer C.free(unsafe.Pointer(dataset))

	fileStopWords := C.CString(w2v.config.Corpus.StopWords)
	defer C.free(unsafe.Pointer(fileStopWords))

	fileModel := C.CString(w2v.config.Output)
	defer C.free(unsafe.Pointer(fileModel))

	withHS := C.uchar(0)
	if w2v.config.UseHierarchicalSoftMax {
		withHS = C.uchar(1)
	}

	withSG := C.uchar(0)
	if w2v.config.UseSkipGram {
		withSG = C.uchar(1)
	}

	tokenizer := C.CString(w2v.config.Corpus.Tokenizer)
	defer C.free(unsafe.Pointer(tokenizer))

	sequencer := C.CString(w2v.config.Corpus.Sequencer)
	defer C.free(unsafe.Pointer(sequencer))

	verbose := C.uchar(0)
	if w2v.config.Verbose {
		verbose = C.uchar(1)
	}

	w2v.h = C.Train(
		dataset,
		fileStopWords,
		fileModel,
		C.ushort(w2v.config.Vector.Frequency),
		C.ushort(w2v.config.Vector.Vector),
		C.uchar(w2v.config.Vector.Window),
		C.float(w2v.config.Vector.Threshold),
		withHS,
		C.uint8_t(w2v.config.SizeNegativeSampling),
		C.uint8_t(w2v.config.Threads),
		C.uint8_t(w2v.config.Learning.Epoch),
		C.float(w2v.config.Learning.Rate),
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
