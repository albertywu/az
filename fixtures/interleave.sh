#!/usr/bin/env bash

for i in {1..100}; do
  sleep 0.01
  if [ $(($i % 2)) -eq 0 ]; then
    echo $i       # write odd numbers to stderr
  else
    echo $i 1>&2  # write even numbers to stdout
  fi
done;
