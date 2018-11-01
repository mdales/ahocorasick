ahocorasick
===========

A Golang implementation of the Aho-Corasick string matching algorithm, derived from Cloudflare's implementation.

This version works differently in three important ways:

1. It returns the position of the matches.
2. It returns all matches.
3. It is thread safe, so can be used concurrently.
