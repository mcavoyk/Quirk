#!/usr/bin/env bash

kill $(lsof -iTCP:5005 -sTCP:LISTEN -t)
cat test.log
rm test.log