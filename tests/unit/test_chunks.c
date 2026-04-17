#include "chunks.h"
#include "unity.h"

static Chunk chunk;

// setup/teardown - run before/after every test function
void setUp(void) {
  printf("\ninitalising chunk\n");
  initChunk(&chunk);
}; // allocate shared fixtures if needed
void tearDown(void) {
  printf("\nfreeing chunk...\n");
  freeChunk(&chunk);
}; // free here

void test_chunks_allocate_chunk(void) {
  TEST_ASSERT_EQUAL_INT(0, chunk.count);
  TEST_ASSERT_EQUAL_INT(0, chunk.capacity);
  TEST_ASSERT_EQUAL_INT(0, chunk.currConstIndex);
  TEST_ASSERT_NULL(chunk.code);
  TEST_ASSERT_NULL(chunk.lines);
}

void test_chunks_writechunk_newChunk(void) {
  writeChunk(&chunk, OP_CONST, 1);
  TEST_ASSERT_EQUAL_UINT8(OP_CONST, chunk.code[0]);
  TEST_ASSERT_EQUAL_INT(1, chunk.count);
}

void test_chunks_writeChunk_same_line(void) {
  writeChunk(&chunk, OP_CONST, 1);
  writeChunk(&chunk, OP_RETURN, 1);
  TEST_ASSERT_EQUAL(2, chunk.count);
  TEST_ASSERT_EQUAL(1, chunk.currConstIndex);
  TEST_ASSERT_EQUAL(2, chunk.lines[0].count);
}

void test_chunks_writeChunk_different_lines(void) {
  writeChunk(&chunk, OP_CONST, 2);
  writeChunk(&chunk, OP_CONST, 3);
  TEST_ASSERT_EQUAL(2, chunk.count);
  TEST_ASSERT_EQUAL(2, chunk.currConstIndex);
  TEST_ASSERT_EQUAL(2, chunk.lines[1].line);
  TEST_ASSERT_EQUAL(1, chunk.lines[0].count);
  TEST_ASSERT_EQUAL(1, chunk.lines[1].count);
}

void test_chunks_writeChunk_capacity_growth(void) {
  // int initial_capacity = chunk.capacity;
  for (int i = 0; i < 100; i++) {
    writeChunk(&chunk, (uint8_t)i, i);
  }
  TEST_ASSERT(chunk.capacity >= 100);
  TEST_ASSERT_EQUAL(100, chunk.count);
}

int main(void) {
  UNITY_BEGIN();
  RUN_TEST(test_chunks_allocate_chunk);
  RUN_TEST(test_chunks_writechunk_newChunk);
  RUN_TEST(test_chunks_writeChunk_same_line);
  // RUN_TEST(test_chunks_writeChunk_different_lines);
  // RUN_TEST(test_chunks_writeChunk_capacity_growth);
  return UNITY_END();
}
