LIBS =
CC = gcc
CFLAGS = -O3 -Wall -fPIC

.PHONY: default all clean

default: pyeditdistance.so 

pyeditdistance.o: pyeditdistance.c
	$(CC) $(CFLAGS) -shared -c $< -o $@ -I/usr/include/python2.7

pyeditdistance.so: pyeditdistance.o editdistance.o
	$(CC) -o $@ -shared $^ $(LIBS)

clean:
	-rm -f *.o
	-rm -f pyeditdistance.so

