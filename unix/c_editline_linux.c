#include <features.h>

#ifndef __GLIBC__

// secure_getenv is missing in musl, but it can be trivially implemented in
// terms of issetugid.
//
// Note that we must assume we're compiling against musl if we're not compiling
// against glibc because musl doesn't define a __MUSL__ macro [0]. Ideally we'd
// just check for the presence of secure_getenv, but cgo provides no facility
// for doing so.
//
// [0]: https://wiki.musl-libc.org/faq.html#Q:-Why-is-there-no-%3Ccode%3E__MUSL__%3C/code%3E-macro?
#include <stdlib.h>
#include <unistd.h>
char *secure_getenv(char const *name) {
	if (issetugid())
		return 0;
	return getenv(name);
}

#endif /* !__GLIBC__ */
