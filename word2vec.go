//
// Copyright (C) 2024 Dmitry Kolesnikov
//
// This file may be modified and distributed under the terms
// of the MIT license.  See the LICENSE file for details.
// https://github.com/fogfish/word2vec
//

package word2vec

/*
#cgo CFLAGS: -Ilibw2v/include
#cgo LDFLAGS: -Llibw2v/lib -lw2v -lstdc++
#include <stdlib.h>
#include "w2v.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

// word2vec model
type Model struct {
	fileModel  string
	vectorSize int

	h unsafe.Pointer
}

// Loads pre-trained model
func Load(model string, vector int) (w2v Model, err error) {
	w2v.fileModel = model
	w2v.vectorSize = 300
	if vector != 0 {
		w2v.vectorSize = vector
	}

	name := C.CString(w2v.fileModel)
	defer C.free(unsafe.Pointer(name))

	w2v.h = C.Load(name)
	if uintptr(w2v.h) == 0 {
		return w2v, fmt.Errorf("unable to load model")
	}

	return w2v, nil
}

//
//
//

// Calculates embedding vector for input term (word)
func (w2v Model) VectorOf(word string, vector []float32) error {
	cword := C.CString(word)
	defer C.free(unsafe.Pointer(cword))

	ptr := C.VectorOf(w2v.h, cword)
	if ptr == nil {
		return errors.New("unknown tokens")
	}

	array := unsafe.Slice((*float32)(ptr), w2v.vectorSize)

	copy(vector, array)

	C.free(unsafe.Pointer(ptr))

	return nil
}

// Calculates embedding for document
func (w2v Model) Embedding(doc string, vector []float32) error {
	cdoc := C.CString(doc)
	defer C.free(unsafe.Pointer(cdoc))

	ptr := C.Embedding(w2v.h, cdoc)
	if ptr == nil {
		return errors.New("unknown tokens")
	}

	array := unsafe.Slice((*float32)(ptr), w2v.vectorSize)

	copy(vector, array)

	C.free(unsafe.Pointer(ptr))

	return nil
}

//
//
//

type Nearest struct {
	Word     string
	Distance float32
}

type nearest_t struct {
	seq *C.float
	len C.ulong
	buf *C.char
}

// Lookup nearest words from the model
func (w2v Model) Lookup(query string, seq []Nearest) error {
	cq := C.CString(query)
	defer C.free(unsafe.Pointer(cq))

	k := len(seq)
	bag := (nearest_t)(C.Lookup(w2v.h, cq, C.ulong(k)))

	if bag.seq == nil || bag.buf == nil {
		return errors.New("unknown tokens")
	}

	seqd := unsafe.Slice((*float32)(bag.seq), k)
	seqw := unsafe.Slice((*C.char)(bag.buf), bag.len)

	p := 0
	for i := 0; i < k; i++ {
		seq[i].Distance = seqd[i]
		seq[i].Word = C.GoString(&seqw[p])
		p += len(seq[i].Word) + 1
	}

	C.free(unsafe.Pointer(bag.seq))
	C.free(unsafe.Pointer(bag.buf))

	return nil
}
