#!/usr/bin/env bash

echo "foo"      # write something to stdout
echo "bar" 1>&2 # write something to stderr
exit 1