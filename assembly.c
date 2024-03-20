#include <stdint.h>

// Fast but slightly inaccurate multiply for AVR.
// The result is the same on all architectures to ease testing (tested by
// comparing many random numbers).
//
// TODO: it would be nice if LLVM could be improved to emit optimized
// instructions for the standard C expressions, instead of whatever it does now.
uint16_t ledsgo_mul16(uint16_t x, uint16_t y) {
	uint16_t result = 0;
#if __AVR__
	__asm__ (
		// result = (x >> 8) * (y >> 8)
		"mul %B[x], %B[y]\n\t"
		"movw %A[result], r0\n\t"

		// result += ((x & 0xff) * (y >> 8)) >> 8
		"mul %A[x], %B[y]\n\t"
		"add %A[result], r1\n\t"
		"clr r1\n\t"
		"adc %B[result], r1\n\t"

		// result += ((x >> 8) * (y & 0xff)) >> 8;
		"mul %B[x], %A[y]\n\t"
		"add %A[result], r1\n\t"
		"clr r1\n\t"
		"adc %B[result], r1\n\t"

		// r1 has already been cleared

		: [result]"=&r"(result)
		: [x]"d"(x),
		  [y]"d"(y)
	);
#else
	result += (x >> 8) * (y >> 8);
	result += ((x & 0xff) * (y >> 8)) >> 8;
	result += ((x >> 8) * (y & 0xff)) >> 8;
#endif
	return result;
}
