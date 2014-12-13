#!/bin/bash
if [ ! -f "grapher" ]; then
	[ -z "$(which gcc)" ] && echo -e "Missing gcc in PATH" && exit 1
	gcc -O2 src/grapher.c -o grapher
fi

SIZE="$1"
[ -z "$SIZE" ] && SIZE=1000 
>&2 echo "Generating random $SIZEbp genome"
python src/generator.py $SIZE > generated.fasta
>&2 echo "Splitting Reference Genome into Reads"
python src/sequencer.py --min=25 --max=50 < generated.fasta > reads.fasta
>&2 echo "Building Overlap Graph edges from Reads"
python src/grapher.py < reads.fasta > reads.edges
>&2 echo "Building Hamiltonian path using greedy algorithm"
python src/greedy.py < reads.edges
>&2 echo "Building Hamiltonian path using genetic algorithm"
python src/genetic.py < reads.edges
