#include "chunks.h"
#include "common.h"
#include "memory.h"
#include "stdlib.h"
#include "value.h"
#include <stdio.h>
#include <string.h>

void initChunk(Chunk *chunk) {
  chunk->count = 0;
  chunk->capacity = 0;
  chunk->currConstIndex = 0;
  chunk->code = NULL;
  chunk->lines = NULL;
  initValueArray(&chunk->consts);
}

void writeChunk(Chunk *chunk, uint8_t byte, int line) {
  if (chunk->capacity < chunk->count + 1) {
    int oldCapacity = chunk->capacity;
    chunk->capacity = GROW_CAPACITY(oldCapacity);
    chunk->code =
        GROW_ARRAY(uint8_t, chunk->code, oldCapacity, chunk->capacity);
    chunk->lines =
        GROW_ARRAY(Lines, chunk->lines, oldCapacity, chunk->capacity);
  }
  chunk->code[chunk->count] = byte;
  if (chunk->count > 0 &&
      chunk->lines[chunk->currConstIndex - 1].line == line) {
    chunk->lines[chunk->currConstIndex - 1].count++;
  } else {
    chunk->lines[chunk->currConstIndex].count++;
    chunk->lines[chunk->currConstIndex].line = line;
    chunk->currConstIndex++;
  }
  // printf("Count=%d\n", chunk->lines[chunk->currConstIndex].count);
  chunk->count++;
}

void writeConstant(Chunk *chunk, Value value, int line) {
  int idx = addConst(chunk, value);
  switch (chunk->code[chunk->count - 1]) {
  case OP_CONST:
    writeChunk(chunk, idx & 0xFF, line);
    return;
  case OP_CONST_LONG:
    writeChunk(chunk, idx & 0xFF, line);
    writeChunk(chunk, (idx >> 8) & 0xFF, line);
    writeChunk(chunk, (idx >> 16), line);
    return;
  default:
    printf("invalid constant instruction");
    exit(0);
  }
  // int constIndex = addConst(chunk, value);
}

int addConst(Chunk *chunk, Value value) {
  writeValueArray(&chunk->consts, value);
  return chunk->consts.count - 1; // return where the constant was added
}

int getLine(Chunk *chunk, int offset) {
  int idx = 0;
  while (chunk->lines[idx].count <= offset) {
    offset = offset - chunk->lines[idx].count;
    idx++;
  }
  return chunk->lines[idx].line;
}

// just a helper function
int loadLongConst(uint8_t a, uint8_t b, uint8_t c) {
  return (a & 0xff) | (b << 8) | (c << 16);
}

void freeChunk(Chunk *chunk) {
  FREE_ARRAY(uint8_t, chunk->code, chunk->capacity);
  FREE_ARRAY(Lines, chunk->lines, chunk->capacity);
  freeValueArray(&chunk->consts);
  initChunk(chunk);
}
