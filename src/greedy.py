#!/usr/bin/python
"""
    Greedy: Implements greedy algorithm to build path in OG for SCS
    Author: Paulo Viadanna
"""
from operator import itemgetter
import argparse
import sys


def greedly(edges):
    results, weigth = edges[0][0:2], edges[0][2]
    while True:
        found = False
        for e in edges:
            if e[0] in results or e[1] in results:
                continue
            results.append(e[1])
            weigth += e[2]
        if not found:
            break
    return results, weigth


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Greedy algorithm for SCS')
    parser.add_argument('Edges', type=unicode, nargs='?', default='-',
                        help='Edges file or - for stdin')
    parser.add_argument('-o', dest='output', action='store', default='-',
                        help='Output fasta or - for stdout')
    args = parser.parse_args()

    if args.Edges == '-':
        f = sys.stdin
    else:
        f = open(args.Edges, 'rb')

    if args.output == '-':
        g = sys.stdout
    else:
        g = open(args.output, 'wb')

    edges = [map(int, e.split(' ')) for e in f.readlines()]
    path, weigth = greedly(sorted(edges, key=itemgetter(2), reverse=True))

    g.write(str(path))

    if args.Edges != '-':
        f.close()
    if args.output != '-':
        g.close()
