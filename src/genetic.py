#!/usr/bin/python
"""
    Genetic: Implements genetic algorithm to build path in OG for SCS
    Author: Paulo Viadanna
"""
from datetime import datetime
from operator import itemgetter
from random import choice
import argparse
import sys


def genetic(edges, father, mother, random_index=2):
    vertices = len(set([e[0] for e in edges]))
    starting = choice(edges[0:random_index])
    results, weight = starting[0:2], starting[2]
    while True:
        vertice = results[-1]
        possible = [e for e in [father.get(vertice)] + [mother.get(vertice)] if e]
        if len(possible) < random_index:
            for e in edges:
                if e[0] == vertice and e[1] not in results:
                    possible.append(e)
                if len(possible) == random_index:
                    break
        if len(possible) == 0:
            # Should start another contig
            return results, weight, [e for e in edges
                                     if e[0] not in results
                                     and e[1] not in results]
        index = min([random_index, len(possible)])
        edge = choice(possible[:index])
        results.append(edge[1])
        weight += edge[2]

        if len(results) == vertices:
            # Full genome assembled
            return results, weight, None


def assemble(edges, father, mother, random_index=2):
    paths = []
    while edges is not None:  # Assemble all contigs
        path, weight, edges = genetic(
            sorted(edges, key=itemgetter(2), reverse=True),
            father, mother, random_index)
        paths.append((path, weight))
    return paths, sum([p[1] for p in paths])


def parents_dict(parents):
    if len(parents) == 0:
        return {}, {}
    father = {v: parents[i+1] for i, v in enumerate(parents[0][0][:-1])}
    mother = {v: parents[i+1] for i, v in enumerate(parents[1][0][:-1])}
    return father, mother


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Genetic algorithm for SCS')
    parser.add_argument('Edges', type=unicode, nargs='?', default='-',
                        help='Edges file or - for stdin')
    parser.add_argument('-o', dest='output', action='store', default='-',
                        help='Output fasta or - for stdout')
    parser.add_argument('-c', dest='children', action='store', type=int,
                        default=10, help='Number of children per generation')
    parser.add_argument('-g', dest='generations', action='store', type=int,
                        default=10, help='Number of generations')
    parser.add_argument('-v', dest='verbose', action='store_true',
                        help='Be verbose')
    args = parser.parse_args()

    if args.Edges == '-':
        f = sys.stdin
    else:
        f = open(args.Edges, 'rb')

    if args.output == '-':
        g = sys.stdout
    else:
        g = open(args.output, 'wb')

    if args.verbose:
        def verb(*args):
            sys.stderr.write("%s: %s" % (
                datetime.now().strftime("%H:%M:%S"),
                " ".join(args)))
    else:
        def verb(*args):
            pass

    given_edges = [map(int, e.split(' ')) for e in f.readlines()]

    father, mother, results = {}, {}, []
    for gen in xrange(args.generations):
        children = [assemble(given_edges[:], father, mother)
                    for _ in xrange(args.children)]
        [verb("Generation %d child: %d\n" % (gen, c[1])) for c in children]
        children.sort(key=itemgetter(1), reverse=True)
        father, mother = parents_dict(children)
        results += children
        verb("Generation %d parents: %d, %d\n" % (
            gen+1, children[0][1], children[1][1]))

    result = sorted(results, key=itemgetter(1), reverse=True)[0][0]
    g.write(str(result))
    g.write('\n')

    if args.Edges != '-':
        f.close()
    if args.output != '-':
        g.close()
