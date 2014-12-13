LIBS =
CC = gcc
CFLAGS = -O3 -Wall -fPIC

.PHONY: default all clean

default: pyeditdistance.so

src/pyeditdistance.o: src/pyeditdistance.c
	$(CC) $(CFLAGS) -shared -c $< -o $@ -I/usr/include/python2.7

pyeditdistance.so: src/pyeditdistance.o src/editdistance.o
	$(CC) -o $@ -shared $^ $(LIBS)

clean:
	-rm -f src/*.o
	-rm -f pyeditdistance.so

