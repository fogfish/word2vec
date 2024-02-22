#ifndef WRAP_SUM_H
#define WRAP_SUM_H

#include <inttypes.h>
#include <stddef.h>

// __cplusplus tells the compiler that inside code is compiled with the c++ compiler
#ifdef __cplusplus
// extern "C" tells C++ compiler exports the symbols without a name manging.
extern "C"
{
#endif

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
      uint8_t verbose);
  void *Load(const char *file);
  void Free(void *fd);

  float *VectorOf(void *fd, const char *word);
  float *Embedding(void *fd, const char *doc);

  struct nearest_t
  {
    float *seq;
    size_t len;
    char *buf;
  };

  struct nearest_t Lookup(void *fd, const char *query, size_t k);
#ifdef __cplusplus
}
#endif
#endif