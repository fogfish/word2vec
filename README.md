<p align="center">
  <img src="./doc/word2vec.png" height="120" />
  <h3 align="center">word2vec</h3>
  <p align="center"><strong>Golang "native" implementation of word2vec algorithm (port of word2vec++)</strong></p>

  <p align="center">
    <!-- Version -->
    <a href="https://github.com/fogfish/word2vec/releases">
      <img src="https://img.shields.io/github/v/tag/fogfish/word2vec?label=version" />
    </a>
    <!-- Documentation -->
    <a href="https://pkg.go.dev/github.com/fogfish/word2vec">
      <img src="https://pkg.go.dev/badge/github.com/fogfish/word2vec" />
    </a>
    <!-- Build Status  -->
    <a href="https://github.com/fogfish/word2vec/actions/">
      <img src="https://github.com/fogfish/word2vec/workflows/test/badge.svg" />
    </a>
    <!-- GitHub -->
    <a href="http://github.com/fogfish/word2vec">
      <img src="https://img.shields.io/github/last-commit/fogfish/word2vec.svg" />
    </a>
    <!-- Coverage -->
    <a href="https://coveralls.io/github/fogfish/word2vec?branch=main">
      <img src="https://coveralls.io/repos/github/fogfish/word2vec/badge.svg?branch=main" />
    </a>
    <!-- Go Card -->
    <a href="https://goreportcard.com/report/github.com/fogfish/word2vec">
      <img src="https://goreportcard.com/badge/github.com/fogfish/word2vec" />
    </a>
  </p>
</p>

--- 

The library enables [word2vec](https://en.wikipedia.org/wiki/Word2vec) algorithm for Golang using native runtime (no servers, no Python, etc). This Golang module implements CGO bridge towards Max Fomichev's [word2vec C++ library](https://github.com/maxoodf/word2vec).

## Getting started

### Building the C++ library

Use C++11 compatible compiler and cmake 3.1 to build the library. It is essential step before going further.

```bash
mkdir _build && cd _build
cmake -DCMAKE_BUILD_TYPE=Release ../libw2v
make
cp ../libw2v/lib/libw2v.dylib /usr/local/lib/libw2v.dylib
```

**Note**: the project does not distribute library binaries, it is upcoming feature. You have to build binaries by yourself for your target runtime or [raise an issue](https://github.com/fogfish/word2vec/issues) if any help is needed.

### Training the model

The trained model is required before moving on. Either use original Max Fomichev's [word2vec C++ utility](https://github.com/maxoodf/word2vec) or Golang's frond-end supplied by this project:

```
go install github.com/fogfish/word2vec/cmd@latest
```

In following examples, ["War and Peace" by Leo Tolstoy](./doc/leo-tolstoy-war-and-peace-en.txt) is used for training. We have also used [stop words](https://github.com/stopwords-iso/stopwords-en) to increase accuracy.

Let's start training with defining the config file:

```
cmd train config > wap-en.yaml

cmd train -C wap-en.yaml \
  -o wap-v300w5e10s1h010-en.bin \
  -f ../doc/leo-tolstoy-war-and-peace-en.txt
```

Name the output model after parameters used for training: `v` vector size, `w` nearby words window, `e` training epoch, architecture	skip-gram `s1` or CBoW `s0`, algorithm	H. softmax `h1`, N. Sampling `h0`.

The default arguments gives sufficient results, see the article [Word2Vec: Optimal hyperparameters and their impact on natural language processing downstream tasks](https://www.degruyter.com/document/doi/10.1515/comp-2022-0236/html?lang=en) for consideration about training options.


### Using word2vec

The latest version of the library is available at its `main` branch. All development, including new features and bug fixes, take place on the `main` branch using forking and pull requests as described in contribution guidelines. The stable version is available via Golang modules.

Use `go get` to retrieve the library and add it as dependency to your application.

```bash
go get -u github.com/fogfish/word2vec
```

The example below shows the usage patterns for the library

```go
import "github.com/fogfish/word2vec"

// 1. Load model
w2v, err := word2vec.Load(
	word2vec.WithModel("wap-v300w5e10s1h010-en.bin"),
	word2vec.WithVectosSize(300),
)

seq := make([]word2vec.Nearest, 30)
w2v.Lookup("alexander", seq)
```

See [the example](./cmd/opts/lookup.go) or try it our via command line

```bash
cmd lookup \
  -m wap-v300w5e10s1h010-en.bin \
  -k 30 \
  alexander
```


### Embeddings

Calculate embedding for document

```go
import "github.com/fogfish/word2vec"

// 1. Load model
w2v, err := word2vec.Load(
	word2vec.WithModel("wap-v300w5e10s1h010-en.bin"),
	word2vec.WithVectosSize(300),
)

// 2. Allocated the memory for vector
vec := make([]float32, 300)

// 3. Calculate embeddings for the document
doc := "braunau was the headquarters of the commander-in-chief"
err = w2v.Embedding(, vec)
```

See [the example](./cmd/opts/embedding.go) or try it our via command line

``` bash
cmd embedding \
  -m wap-v300w5e10s1h010-en.bin \
  ../doc/leo-tolstoy-war-and-peace-en.txt
```


## How To Contribute

The library is [MIT](LICENSE) licensed and accepts contributions via GitHub pull requests:

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Added some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

The build and testing process requires [Go](https://golang.org) version 1.21 or later.


### commit message

The commit message helps us to write a good release note, speed-up review process. The message should address two question what changed and why. The project follows the template defined by chapter [Contributing to a Project](http://git-scm.com/book/ch5-2.html) of Git book.

### bugs

If you experience any issues with the library, please let us know via [GitHub issues](https://github.com/fogfish/word2vec/issue). We appreciate detailed and accurate reports that help us to identity and replicate the issue. 


## License

[![See LICENSE](https://img.shields.io/github/license/fogfish/word2vec.svg?style=for-the-badge)](LICENSE)


## References

1. [Go Wiki: cgo](https://go.dev/wiki/cgo)
2. [Calling C code with cgo](https://spatocode.com/blog/calling-c-code-with-cgo)
3. [Pass struct and array of structs to C function from Go](https://stackoverflow.com/questions/19910647/pass-struct-and-array-of-structs-to-c-function-from-go)
4. [cgo - cast C struct to Go struct](https://groups.google.com/g/golang-nuts/c/JkvR4dQy9t4)
5. [Word2Vec (google code)](https://code.google.com/archive/p/word2vec/)
6. [word2vec patch for Mac OS X](https://github.com/William-Yeh/word2vec-mac)