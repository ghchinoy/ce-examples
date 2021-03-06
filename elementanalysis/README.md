# Element Analysis

Using the current default profile ([`cectl` style profile](https://github.com/ghchinoy/cectl)), provide information about the Element catalog


As of 11/22/2017

* 1482 objects
* 130 elements


## Usage

`elementanalysis` will use the current default CE Environment profile and output the aggregate total number of HTTP methods (per HTTP method), aggregate count of resource paths, and total schema count across all Elements.

Optionally, provide a list of Element keys, and stats for only those will be provided, example: `elementanalysis zuorav2`

## Example output

Running `elementanalysis` with the default profile yields, for example:

```
$ elementanalysis
   PUT   168
  POST   937
 PATCH   669
   GET  2443
DELETE   750
 Paths  2722
Schema 12334
```