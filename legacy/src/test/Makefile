LIBS =
CC = gcc
CFLAGS = -O3 -Wall

.PHONY: default all clean

default: test-editdistance

test-editdistance.o: test-editdistance.c
	$(CC) $(CFLAGS) -c $< -o $@

test-editdistance: test-editdistance.o ../editdistance.o
	$(CC) -o $@ $^ $(LIBS)

clean:
	-rm -f *.o
	-rm -f test-editdistance
