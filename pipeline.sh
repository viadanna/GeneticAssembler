#!/bin/bash
cat genome.fasta | python sequencer.py - -o reads.fasta
cat reads.fasta | julia grapher.jl > reads.edges
