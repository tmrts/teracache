<h1 align=center>
Teracache
<br>
<a href="http://travis-ci.org/tmrts/teracache"><img alt="build-status" src="https://img.shields.io/badge/build-passing-brightgreen.svg?style=flat-square" /></a>
<a href="https://github.com/tmrts/teracache/blob/master/LICENSE" ><img alt="License" src="https://img.shields.io/badge/license-Apache%20License%202.0-E91E63.svg?style=flat-square"/></a>
<a href="https://github.com/tmrts/teracache/releases" ><img alt="Release Version" src="https://img.shields.io/badge/release-v0.0.1-blue.svg?style=flat-square"/></a>
<a href="https://godoc.org/github.com/tmrts/teracache" ><img alt="Documentation" src="https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square"/></a>
</h1>

Teracache is a scalable, decentralized, highly-available in-memory cache for
read-asymmetrical workflows.

## Workflow

The workflow diagram depicts a system with multiple topics and nodes where each
color represents the topic they're in.

![workflow-diagram](https://cdn.rawgit.com/tmrts/teracache/master/workflow-diagram.svg)

## What is Read-Asymmetrical?

In a read-asymmetrical workflow, reads are frequent while the writes are rare.
For example a newsfeed will resize and optimize a thumbnail when a new article
arrives and then serve the optimized thumbnail to thousands or millions of
users.

In such cases, it is very profitable to use a cache that doesn't sacrifice
performance to allow write operations. Teracache does exactly this trade-off.

## How to use invalidations?

Cache entries are immutable and no entry is considered stale. Entries
are evicted only to admit more popular entries.

## Inspired By

- [groupcache](https://github.com/golang/groupcache)
- [SWIM Membership Protocol](https://www.cs.cornell.edu/~asdas/research/dsn02-swim.pdf)
- [Microsoft Orleans](https://www.microsoft.com/en-us/research/wp-content/uploads/2016/02/Orleans-MSR-TR-2014-41.pdf)
- [Amazon Dynamo](http://s3.amazonaws.com/AllThingsDistributed/sosp/amazon-dynamo-sosp2007.pdf)
- [Consistent Hashing and Random Trees](https://www.akamai.com/es/es/multimedia/documents/technical-publication/consistent-hashing-and-random-trees-distributed-caching-protocols-for-relieving-hot-spots-on-the-world-wide-web-technical-publication.pdf)
