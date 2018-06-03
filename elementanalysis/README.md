# Element Analysis

Using the current default profile ([`cectl` style profile](https://github.com/ghchinoy/cectl)), provide information about the Element catalog


As of 11/22/2017

* 1482 objects
* 130 elements


## Usage

`elementanalysis` will use the current default CE Environment profile and output the aggregate total number of HTTP methods (per HTTP method), aggregate count of resource paths, and total schema count across all Elements.

Optionally, provide a list of Element keys, and stats for only those will be provided, example: `elementanalysis zuorav2`