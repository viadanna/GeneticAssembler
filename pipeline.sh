#!/bin/bash
if [ ! -f "grapher" ]; then
	[ -z "$(which gcc)" ] && echo -e "Missing gcc in PATH" && exit 1
	gcc -O2 src/grapher.c -o grapher
fi
>&2 echo "Splitting Reference Genome into Reads"
python src/sequencer.py -o reads.fasta < data/genome.fasta
>&2 echo "Building Overlap Graph edges from Reads"
python src/grapher.py < reads.fasta > reads.edges
>&2 echo "Building Hamiltonian path using greedy algorithm"
python src/greedy.py < reads.edges
>&2 echo "Building Hamiltonian path using genetic algorithm"
python src/genetic.py -v < reads.edges
