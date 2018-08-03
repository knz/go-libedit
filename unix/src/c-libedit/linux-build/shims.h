#include <features.h>

#ifdef __GLIBC__
// Only __secure_getenv exists prior to glibc 2.17.
#if __GLIBC__ <= 2 && __GLIBC_MINOR__ < 17
#define secure_getenv __secure_getenv
#endif
#else
// We provide a shim implementation for other C libraries.
char *secure_getenv(char const *name);
#endif
