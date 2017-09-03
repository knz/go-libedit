#!/bin/bash
set -euxo pipefail

rm -f ../shim/*.c
rm -f ../wrap-*.c
rm -rf c-libedit
git add -u ..

rm -rf libbsd[-_]* libedit[-_]* build

apt-get source libbsd0
apt-get source libedit

mkdir build
(cd build && ../libedit-*/configure && make SUBDIRS=src)

mkdir -p c-libedit/linux-build c-libedit/editline

cp -a libedit-*/src/editline c-libedit/
cp -a libedit-*/src/*.[ch] c-libedit/
cp -a build/config.h build/src/*.h c-libedit/linux-build/

(cd c-libedit &&
     for i in *.c; do
	 echo "#include \"$i\"">../libedit-$i
	 echo "// Nothing to see here.">../../shim/libedit-$i
	 echo "#include \"libedit-$i\"">../../wrap-$i
     done)

git add ../shim/*.c
git add ../wrap-*.c
git add libedit-*.c
git add c-libedit
