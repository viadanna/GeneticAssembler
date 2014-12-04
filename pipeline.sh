#!/bin/bash
echo "Splitting Reference Genome into Reads"
cat data/genome.fasta | python src/sequencer.py - -o reads.fasta
echo "Building Overlap Graph edges from Reads"
cat reads.fasta | julia src/grapher.jl > reads.edges
echo "Building Hamiltonian path from OG"
cat reads.edges | python src/greedy.py -
