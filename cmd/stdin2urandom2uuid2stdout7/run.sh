#!/bin/sh

randomdata=/dev/zero
randomdata=/dev/urandom

uuidcount=512
bufsize=16

dd \
	if="${randomdata}" \
	of=/dev/stdout \
	bs=${bufsize} \
	count=${uuidcount} \
	status=none |
	./stdin2urandom2uuid2stdout7 |
	cat > sample.d/sample.uuids.txt

cat sample.d/sample.uuids.txt |
	rev |
	uniq --skip-chars=20 |
	rev |
	cat -n |
	head -32
