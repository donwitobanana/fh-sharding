#!/bin/bash

curl -v -X 'POST' \
  'http://localhost:8080' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d @./test_data.json