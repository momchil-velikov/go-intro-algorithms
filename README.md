# go-intro-algorithms
Collection of algorithms and data structures from CLRS, written in Go

This is a dump of the code I wrote during my own refreshement course in data structures,
based on:
 - video lectures of MIT 6.046J / 18.410J Introduction to Algorithms (SMA 5503) course, https://www.youtube.com/playlist?list=PL8B24C31197EC371C
 - Introduction to Algorithms, 3rd ed., http://www.amazon.co.uk/Introduction-Algorithms-T-Cormen/dp/0262533057

Contained within are implementations of:

  - binary heaps, heap sort
  - quicksort with several different partitioning strategies
  - counting sort and radix sort
  - open-addressing hash table
  - treaps
  - red-black trees
  - skip lists
  - B-trees
  - depth-first search in a directed graph
  - topological sort
  - strongly connected components: algorithms of Tarjan and Kosaraju
  - minimum spanning tree: algorithms of Prim and Kruskal
  - single source shortest paths: algorithms of Dijkstra and Bellman-Ford
  - all pairs shortest paths: algorithms of Floyd-Warshal and Johnson
  - maximum flow: algorithm of Edmonds-Karp
  - few other things :)
   
Extras:
  - a small library for reading and writing files in a very useful subset of the Graphviz format
  - an utility for generation of random irreducible control-flow graphs
  - some sample graphs on Graphs/data directory
  
Some of the programs have Go test suites; these are considerd pretty solid.

Use of this source code is governed by a BSD-style license that can be found in the COPYING file.
