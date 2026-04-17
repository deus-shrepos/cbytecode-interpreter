#ifndef clox_chunk_h
#define clox_chunk_h

#include "common.h"
#include "value.h"
#include <stdint.h>

// ByteCode Enum
typedef enum {
  OP_RETURN,
  OP_CONST,
  OP_CONST_LONG,
} OpCode;


typedef struct {
  int line;
  int count;
} Lines;

typedef struct {
  int count;
  int capacity;
  int currConstIndex;
  uint8_t *code;
  Lines *lines;
  ValueArray consts;
} Chunk;

void initChunk(Chunk *chunk);
void writeChunk(Chunk *chunk, uint8_t bytes, int line);
void writeConstant(Chunk *chunk, Value value, int line);
void freeChunk(Chunk *chunk);
int addConst(Chunk *chunk, Value value);
int loadLongConst(uint8_t a, uint8_t b, uint8_t c);

// void freeChunk(Chunk *chunk);
int getLine(Chunk *chunk, int instructionIndex);

#endif
