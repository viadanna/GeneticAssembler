#!/usr/bin/python
"""
    Greedy: Implements greedy algorithm to build path in OG for SCS
    Author: Paulo Viadanna
"""
from operator import itemgetter
import argparse
import sys


def greedly(edges):
    vertices = len(set([e[0] for e in edges]))
    results, weigth = edges[0][0:2], edges[0][2]
    while True:
        found = False
        for e in edges:
            # Search valid edge with highest weigth
            if e[0] != results[-1] or e[1] in results:
                continue
            results.append(e[1])
            weigth += e[2]
            found = True
            break
        if len(results) == vertices:
            # Full genome assembled
            return results, weigth, None
        if not found:
            # Should start another contig instead
            return results, weigth, [e for e in edges
                                     if e[0] not in results
                                     and e[1] not in results]


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

    edges, paths = [map(int, e.split(' ')) for e in f.readlines()], []

    while edges is not None:
        path, weigth, edges = greedly(
            sorted(edges, key=itemgetter(2), reverse=True))
        paths.append((path, weigth))

    g.write(str(paths))
    g.write('\n')

    if args.Edges != '-':
        f.close()
    if args.output != '-':
        g.close()
