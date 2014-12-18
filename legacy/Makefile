LIBS =
CC = gcc
CFLAGS = -O3 -Wall -fPIC

.PHONY: default all clean

default: pyeditdistance.so grapher.so

src/pyeditdistance.o: src/pyeditdistance.c
	$(CC) $(CFLAGS) -shared -c $< -o $@ -I/usr/include/python2.7

pyeditdistance.so: src/pyeditdistance.o src/editdistance.o
	$(CC) -o $@ -shared $^ $(LIBS)

src/grapher.o: src/grapher.c
	$(CC) $(CFLAGS) -shared -c $< -o $@

grapher.so: src/grapher.o
	$(CC) -o $@ -shared $^ $(LIBS)

clean:
	-rm -f src/*.o
	-rm -f pyeditdistance.so
	-rm -f grapher.so

