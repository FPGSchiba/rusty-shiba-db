# Rusty Storage solution

## General

In general following files will be present:
1. `/collections.rsc`
   1. A list of collection id and name mappings

## Collections

Collections are represented as Folders in the FS structure of the OS.
For each collection two things are given:

1. `/{colleciton-id}.rsc`
   1. A Schema file to determine the schema of the collection.
   2. If no schema was given this file will be empty
2. `/{collection-id}.rsi`
   1. An Index file to store the current B-Tree index (Node structure only keys)
   2. To map Nodes to files within the Collection folder
3. `/{collection-id}/*.rsd`
   1. A number of data files within the collection folder