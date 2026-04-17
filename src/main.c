#include "chunks.h"
#include "common.h"
#include "debug.h"
#include <stdint.h>

int main(int argc, const char *argv[]) {
  Chunk chunk;
  initChunk(&chunk);
  writeChunk(&chunk, OP_CONST, 127);
  writeConstant(&chunk, 0.6, 127);
  writeChunk(&chunk, OP_RETURN, 127);
  writeChunk(&chunk, OP_CONST_LONG, 128);
  writeConstant(&chunk, 0.7, 128);
  writeChunk(&chunk, OP_RETURN, 128);
  writeChunk(&chunk, OP_CONST, 129);
  writeConstant(&chunk, 0.8, 129);
  writeChunk(&chunk, OP_RETURN, 129);
  disassembleChunk(&chunk, "Chunk Block");
  freeChunk(&chunk);
  return 0;
}
