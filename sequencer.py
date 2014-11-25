from Bio.SeqIO import parse, write, SeqRecord
from Bio.Seq import Seq
from random import randint, choice
import argparse
import sys


def split_reads(seqs, min_size, max_size, coverage, error, reversals):
    base_pair = {'A': 'C', 'C': 'A', 'G': 'T', 'T': 'G'}
    i, generated, target = 0, 0, sum([len(s.seq) for s in seqs]) * coverage
    while generated < target:
        seq = choice(seqs).seq
        size = randint(min_size, max_size)
        start = randint(0, len(seq)-size)
        i += 1
        generated += size
        read = seq[start:start+size].upper()
        if randint(0, 99) < reversals:
            read = Seq("".join([base_pair[b] for b in reversed(read)]))
        read = Seq("".join((b if randint(0, 99) < error
                        else choice(('A', 'C', 'G', 'T')) for b in read)))
        yield SeqRecord(read, 'fragment_%i' % i, '', '')


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Splits a RefSeq into reads')
    parser.add_argument('RefSeq', type=unicode,
                        help='Reference file or - for stdin')
    parser.add_argument('-o', dest='output', action='store', default='-',
                        help='Output fasta or - for stdin')
    parser.add_argument('--min', dest='min_size', action='store', default=800,
                        help='Minimum read length')
    parser.add_argument('--max', dest='max_size', action='store', default=900,
                        help='Maximum read length')
    parser.add_argument('--coverage', type=int, action='store', default=15,
                        help='Desired coverage of given RefSeq')
    parser.add_argument('--error', type=int, action='store', default=1,
                        help='Error rate, in %')
    parser.add_argument('--rev', type=int, action='store', default=50,
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

    seqs = list(parse(f, 'fasta'))
    for r in split_reads(seqs, args.min_size, args.max_size, args.coverage,
                            args.error, args.rev):
        write(r, g, "fasta")

    f.close()
    g.close()
