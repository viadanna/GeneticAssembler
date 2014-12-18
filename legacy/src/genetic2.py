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


def genetic(edges, father, mother, cutoff):
    """
    Builds hamiltonian path for the given weigthed edges
    Will choose a random edge in given list and parents if given
    Returns the list of vertices, weight and edges not added for lack of a
    possible path
    """
    vertices = len(set([e[0] for e in edges]))
    starting = choice([e for e in edges if e[2] > edges[0][2] * cutoff])
    results, weight = starting[0:2], starting[2]

    while len(results) < vertices:
        # Loop until all vertices are added
        vertice = results[-1]
        if father and vertice == father[0]:
            results += father[1:]
        elif mother and vertice == mother[0]:
            results += mother[1:]
        else:  # Mutation needed
            possible = [e for e in edges
                        if e[0] == vertice and e[1] not in results]
            best = max([e[2] for e in possible])
            possible = [e for e in possible if e[2] > best * cutoff]
            if len(possible) == 0:
                # Should start another contig
                return results, weight, [e for e in edges
                                        if e[0] not in results
                                        and e[1] not in results]
            # Choose a random edge from possible ones
            edge = choice(possible)
            results.append(edge[1])
            weight += edge[2]

    # Full genome assembled
    return results, weight, None


def assemble(edges, father, mother, cutoff=0.9):
    """
    Runs the assembly algorithm until all edges are added
    Returns contigs and the sum of all contigs' weights
    """
    paths = []
    while edges is not None:  # Assemble all contigs
        path, weight, edges = genetic(
            sorted(edges, key=itemgetter(2), reverse=True),
            father, mother, cutoff)
        paths.append((path, weight))
    return paths, sum([p[1] for p in paths])


def find_gene(edges, parent, other=None, fraction=8):
    " Find the best sequence from a parent "
    best, med = [], 0
    for contig in parent[0]:
        size = len(contig[0])/fraction
        for i, v in enumerate(contig[0][:-1]):
            if other and v in other:
                continue
            total = 0
            previous = v
            for k, w in enumerate(contig[0]):
                if k <= i:
                    continue
                if other and w in other:
                    break
                total += edges["%d-%d" % (previous, w)]
                if total / (k - i) > med and k-i >= size:
                    med = total / (k - i)
                    best = contig[0][i:k+1]
                previous = w
    return best


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
    edges_dict = {"%d-%d" % (e[0], e[1]): e[2] for e in given_edges}

    father, mother, results = {}, {}, []
    for gen in xrange(args.generations):
        children = []
        for _ in xrange(args.children):
            child = assemble(given_edges[:], father, mother)
            children.append(child)
            verb("Generation %d child: %d\n" % (gen, child[1]))
        children.sort(key=itemgetter(1), reverse=True)
        father = find_gene(edges_dict, children[0]) if len(children) > 0 else None
        mother = find_gene(edges_dict, children[1], father) if len(children) > 1 else None
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
