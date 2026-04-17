CC := @gcc
CFLAGS := -std=c99 -Wall


SRC_DIR := src
BUILD_DIR := build
TARGET := bin/bcinter
VENDOR_DIR := vendor/unity
TEST_DIR := tests/unit

SRCS := $(wildcard $(SRC_DIR)/*.c)
OBJS := $(patsubst $(SRC_DIR)/%.c, $(BUILD_DIR)/%.o, $(SRCS))

LANG_SRC := $(filter-out $(SRC_DIR)/main.c, $(SRCS))
TEST_SRC := $(wildcard $(TEST_DIR)/test_*.c)
TEST_BINS := $(patsubst $(TEST_DIR)/test_%.c, $(BUILD_DIR)/test_%, $(TEST_SRC))


.PHONY: all

all: $(TARGET)

$(TARGET): $(OBJS)
	$(CC) $(LDFLAGS) $^ -o $@ $(LDLIBS)

$(BUILD_DIR)/%.o: $(SRC_DIR)/%.c
	$(CC) $(CFLAGS) -c $< -o $@

clean:
	$(RM) -r $(BUILD_DIR) bin

run: all
	./$(TARGET)

test-unit: $(TEST_BINS)
	@echo "Running unit tests..."
	@failed=0; \
	for t in $(TEST_BINS); do \
			$$t || $$failed=$$((failed+1)); \
	done; \
	[ $$failed -eq 0 ] && echo "All unit tests passed." \
				|| (echo "$$failed test suit(s) failed" && exit 1)

$(BUILD_DIR)/test_%: $(TEST_DIR)/test_%.c $(LANG_SRC) $(VENDOR_DIR)/unity.c
	$(CC) $(CFLAGS) -I$(SRC_DIR) -I$(VENDOR_DIR) $^ -o $@


$(BUILD_DIR):
		mkdir -p $@