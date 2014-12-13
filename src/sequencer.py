#!/usr/bin/python
"""
    Sequencer: Splits given contigs into a number of fasta reads
    Author: Paulo Viadanna
"""
from fasta import parse, write
from random import randint, choice
import argparse
import sys

base_pair = {'A': 'C', 'C': 'A', 'G': 'T', 'T': 'G'}


def reverse_strip(x):
    return "".join([base_pair[b] for b in reversed(x)])


def split_reads(seqs, min_size, max_size, coverage, error, reversals):
    i, generated, target = 0, 0, sum([len(s) for s in seqs]) * coverage
    while generated < target:
        seq = choice(seqs)
        size = randint(min_size, max_size)
        start = randint(0, len(seq)-size)
        i += 1
        generated += size
        read = seq[start:start+size].upper()
        if reversals > 0 and randint(0, 99) < reversals:
            read = reverse_strip(read)
        if error > 0:
            read = "".join([b if randint(0, 99) < error
                else choice(('A', 'C', 'G', 'T', '')) for b in read])
        yield read


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Splits a RefSeq into reads')
    parser.add_argument('RefSeq', type=unicode, nargs='?', default='-',
                        help='Reference file or - for stdin')
    parser.add_argument('-o', dest='output', action='store', default='-',
                        help='Output fasta or - for stdout')
    parser.add_argument('--min', dest='min_size', action='store', default=800,
                        type=int, help='Minimum read length')
    parser.add_argument('--max', dest='max_size', action='store', default=900,
                        type=int, help='Maximum read length')
    parser.add_argument('--coverage', type=float, action='store', default=11,
                        help='Desired coverage of given RefSeq')
    parser.add_argument('--error', type=float, action='store', default=0,
                        help='Error rate, in %')
    parser.add_argument('--rev', type=float, action='store', default=0,
                        help='Reverse strip rate, in %')
    args = parser.parse_args()

    if args.RefSeq == '-':
        f = sys.stdin
    else:
        f = open(args.RefSeq, 'rb')

    if args.output == '-':
        g = sys.stdout
    else:
        g = open(args.output, 'wb')

    seqs = [s for _, s in parse(f)]

    for i, fragment in enumerate(split_reads(seqs, args.min_size,
                                             args.max_size, args.coverage,
                                             args.error, args.rev)):
        write(fragment, 'fragment_%d' % i, g)

    if args.RefSeq != '-':
        f.close()
    if args.output != '-':
        g.close()
