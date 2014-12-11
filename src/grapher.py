#!/usr/bin/python
"""
    Grapher: Generate Multigraph for SCS from given reads
    Author: Paulo Viadanna
"""
from ctypes import cdll, c_char_p, c_int
# from sequencer import reverse_strip
import sys
import argparse


c_sp_align = cdll.LoadLibrary('./grapher').sp_align
c_sp_align.argtypes = [c_char_p, c_char_p]
c_sp_align.restype = c_int


def sp_align(a, b):
    return c_sp_align(c_char_p(a), c_char_p(b))


def multi_graph(V):
    for i, f in enumerate(F):
        for j, f_linha in enumerate(F):
            if i == j:
                continue
            t = sp_align(f, f_linha)
            yield (i, j, t)


if __name__ == '__main__':
    from Bio.SeqIO import parse
    parser = argparse.ArgumentParser(description='Builds the Multigraph')
    parser.add_argument('Reads', type=str, nargs='?', default='-',
                        help='Reference file or - for stdin')
    parser.add_argument('-o', dest='output', action='store', default='-',
                        help='Output fasta or - for stdout')
    parser.add_argument('--rev', action='store_true',
                        help='Consider reverse strip')
    args = parser.parse_args()

    if args.Reads == '-':
        f = sys.stdin
    else:
        f = open(args.Reads, 'rb')

    if args.output == '-':
        g = sys.stdout
    else:
        g = open(args.output, 'wb')

    F = [str(s.seq) for s in parse(f, "fasta")]
    if args.rev:
        # TODO: Parse reverse strip fragments
        pass

    for edge in multi_graph(F):
        g.write("%s\n" % " ".join([str(i) for i in edge]))

    if args.Reads != '-':
        f.close()
    if args.output != '-':
        g.close()
