#include "debug.h"
#include "chunks.h"
#include "value.h"
#include <stdio.h>

static int simpleInstruction(const char *name, int offset) {
  printf("%s\n", name);
  return offset + 1;
}

static int constantInstruction(const char *name, Chunk *chunk, int offset) {
  uint8_t constant =
      chunk->code[offset + 1]; // offset + 1 = next chunk (const index)
  printf("%-16s %4d '", name, constant);
  printValue(chunk->consts.values[constant]);
  printf("'\n");
  return offset + 2;
}

static int loadConstantInstruction(const char *name, Chunk *chunk, int offset) {
  int constant = loadLongConst(chunk->code[offset + 1], chunk->code[offset + 2],
                               chunk->code[offset + 3]);
  printf("%-16s %4d '", name, constant);
  printValue(chunk->consts.values[constant]);
  printf("'\n");
  return offset + 4;
}

void disassembleChunk(Chunk *chunk, const char *name) {
  printf("== %s ==\n", name);
  for (int offset = 0; offset < chunk->count;) {
    offset = disassembleInstruction(chunk, offset);
  }
}

int disassembleInstruction(Chunk *chunk, int offset) {
  printf("%04d ", offset);
  // bool isNewLine = chunk->lines[offset].line == chunk->lines[offset -
  // 1].line;
  bool isNewLine = getLine(chunk, offset) == getLine(chunk, offset - 1);
  if (offset > 0 && isNewLine) {
    printf("   | ");
  } else {
    printf("%4d ", getLine(chunk, offset));
  }
  uint8_t instruction = chunk->code[offset];
  switch (instruction) {
  case OP_RETURN:
    return simpleInstruction("OP_RETURN", offset);
  case OP_CONST:
    return constantInstruction("OP_CONST", chunk, offset);
  case OP_CONST_LONG:
    return loadConstantInstruction("OP_LONG_CONST", chunk, offset);
  default:
    printf("Unknown opcode %d\n", instruction);
    return offset + 1;
  }
}
