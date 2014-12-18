#!/usr/bin/python
"""
    Generator: Generates random genomes for experiments
    Author: Paulo Viadanna
"""
from fasta import write
from random import choice
import argparse
import sys


def generate_sequence(size, alphabet='ACGT'):
    return "".join([choice(alphabet) for _ in xrange(size)])


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Generate random sequence')
    parser.add_argument('Size', type=int, nargs='?', default=1000,
                        help='Size of generated sequence')
    parser.add_argument('Alphabet', type=str, nargs='?', default='ACGT',
                        help='Alphabet used to generate sequence')
    args = parser.parse_args()

    g = sys.stdout
    write(generate_sequence(args.Size, args.Alphabet),
          'random_genome_size_%d' % args.Size, g)
