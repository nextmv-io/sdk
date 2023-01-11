#!/bin/bash
set -e

echo "🐰 cd into input-gen"
cd input-gen
echo "🐰 generating file"
nextmv run local .
echo "🐰 cd into routing app"
cd ../routing
echo "🐰 running routing app"
nextmv run local . -- -hop.runner.input.path ../input-gen/routing-input.json \
  -hop.runner.output.path output.json -hop.solver.limits.duration 10s
