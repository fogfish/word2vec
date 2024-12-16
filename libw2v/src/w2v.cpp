#include "w2v.h"
#include "word2vec.hpp"

#include <inttypes.h>
#include <iostream>

#include <iostream>
#include <iomanip>
#include <cstring>
#include <stdexcept>

class H
{
public:
  std::unique_ptr<w2v::w2vModel_t> model;

  H();
  ~H();
};

H::H()
{
  model.reset(new w2v::w2vModel_t());
}

H::~H()
{
  model.reset();
}

void *Train(
    char *fileTrain,
    char *fileStopWords,
    char *fileModel,
    uint16_t minWordFreq,
    uint16_t vectorSize,
    uint8_t window,
    float sample,
    uint8_t withHS,
    uint8_t negative,
    uint8_t threads,
    uint8_t iterations,
    float alpha,
    uint8_t withSG,
    char *wordDelimiterChars,
    char *endOfSentenceChars,
    uint8_t verbose)
{
  w2v::trainSettings_t trainSettings;
  trainSettings.size = vectorSize;
  trainSettings.window = window;
  trainSettings.sample = sample;
  trainSettings.withHS = withHS;
  trainSettings.negative = negative;
  trainSettings.threads = threads;
  trainSettings.iterations = iterations;
  trainSettings.minWordFreq = minWordFreq;
  trainSettings.alpha = alpha;
  trainSettings.withSG = withSG;
  trainSettings.wordDelimiterChars = wordDelimiterChars;
  trainSettings.endOfSentenceChars = endOfSentenceChars;

  std::string trainFile;
  trainFile = fileTrain;

  std::string modelFile;
  modelFile = fileModel;

  std::string stopWordsFile;
  stopWordsFile = fileStopWords;

  if (verbose)
  {
    std::cout << "Train data file: " << trainFile << std::endl;
    std::cout << "Output model file: " << modelFile << std::endl;
    std::cout << "Stop-words file: " << stopWordsFile << std::endl;
    std::cout << "Training model: " << (trainSettings.withSG ? "Skip-Gram" : "CBOW") << std::endl;
    std::cout << "Sample approximation method: ";
    if (trainSettings.withHS)
    {
      std::cout << "Hierarchical softmax" << std::endl;
    }
    else
    {
      std::cout << "Negative sampling with number of negative examples = "
                << static_cast<int>(trainSettings.negative) << std::endl;
    }
    std::cout << "Number of training threads: " << static_cast<int>(trainSettings.threads) << std::endl;
    std::cout << "Number of training iterations: " << static_cast<int>(trainSettings.iterations) << std::endl;
    std::cout << "Min word frequency: " << static_cast<int>(trainSettings.minWordFreq) << std::endl;
    std::cout << "Vector size: " << static_cast<int>(trainSettings.size) << std::endl;
    std::cout << "Max skip length: " << static_cast<int>(trainSettings.window) << std::endl;
    std::cout << "Threshold for occurrence of words: " << trainSettings.sample << std::endl;
    std::cout << "Starting learning rate: " << trainSettings.alpha << std::endl;
    std::cout << std::endl
              << std::flush;
  }

  auto h = new H();
  bool trained;
  if (verbose)
  {
    trained = h->model->train(
        trainSettings, trainFile, stopWordsFile,
        [](float _percent)
        {
          std::cout << "\rParsing train data... "
                    << std::fixed << std::setprecision(2)
                    << _percent << "%" << std::flush;
        },
        [](std::size_t _vocWords, std::size_t _trainWords, std::size_t _totalWords)
        {
          std::cout << std::endl
                    << "Vocabulary size: " << _vocWords << std::endl
                    << "Train words: " << _trainWords << std::endl
                    << "Total words: " << _totalWords << std::endl
                    << std::endl;
        },
        [](float _alpha, float _percent)
        {
          std::cout << "\r                                                                  \r"
                    << "alpha: "
                    << std::fixed << std::setprecision(6)
                    << _alpha
                    << ", progress: "
                    << std::fixed << std::setprecision(2)
                    << _percent << "%"
                    << std::flush;
        });
    std::cout << std::endl;
  }
  else
  {
    trained = h->model->train(trainSettings, trainFile, stopWordsFile, nullptr, nullptr, nullptr);
  }
  if (!trained)
  {
    std::cerr << "Training failed: " << h->model->errMsg() << std::endl;
    return 0;
  }

  if (!h->model->save(modelFile))
  {
    std::cerr << "Model file saving failed: " << h->model->errMsg() << std::endl;
    return 0;
  }

  return h;
}

// Load model
void *Load(const char *file)
{
  auto h = new H();

  if (!h->model->load(file))
  {
    std::cerr << h->model->errMsg() << '\n';
    return 0;
  }

  return h;
}

void Free(void *fd)
{
  auto h = reinterpret_cast<H *>(fd);
  delete h;
}

float *VectorOf(void *fd, const char *word)
{
  try
  {
    auto h = reinterpret_cast<H *>(fd);
    w2v::word2vec_t vec(h->model, word);

    float *vector = (float *)malloc(sizeof(float) * vec.size());
    std::copy(vec.begin(), vec.end(), vector);

    return vector;
  }
  catch (const std::exception &e)
  {
    return 0;
  }
}

float *Embedding(void *fd, const char *doc)
{
  try
  {
    auto h = reinterpret_cast<H *>(fd);
    w2v::doc2vec_t vec(h->model, doc);

    float *vector = (float *)malloc(sizeof(float) * vec.size());
    std::copy(vec.begin(), vec.end(), vector);

    return vector;
  }
  catch (const std::exception &e)
  {
    return 0;
  }
}

struct nearest_t Lookup(void *fd, const char *query, size_t k)
{
  try
  {
    auto h = reinterpret_cast<H *>(fd);
    w2v::doc2vec_t vec(h->model, query);

    std::vector<std::pair<std::string, float>> nearests;
    h->model->nearest(vec, nearests, k);

    float *seqd = (float *)malloc(sizeof(float) * k);

    size_t len = 0;
    for (auto i = size_t(0); i < k; i++)
    {
      len += nearests[i].first.length() + 1;
    }
    char *seqw = (char *)malloc(len);

    size_t p = 0;

    for (auto i = size_t(0); i < k; i++)
    {
      *(seqd + i) = nearests[i].second;

      strcpy(seqw + p, nearests[i].first.c_str());

      p += nearests[i].first.length();
      *(seqw + p) = '\0';

      p++;
    }
    return nearest_t{seqd, len, seqw};
  }
  catch (const std::exception &e)
  {
    return nearest_t{0, 0, 0};
  }
}
