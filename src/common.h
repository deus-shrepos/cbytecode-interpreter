#ifndef clox_common_h
#define clox_common_h

#include <stdbool.h>
#include <stddef.h>
#include <stdint.h>
#include <stdio.h>

#define PRINT_ARRAY(arr, len, fmt, ...) \
    do { \
        for (size_t _i = 0; _i < (len); _i++) { \
            __typeof__((arr)[0]) _elm = (arr)[_i]; \
            printf(fmt, __VA_ARGS__); \
        } \
    } while (0)


#endif
