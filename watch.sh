#!/bin/bash

# Surveillance des changements de fichiers avec entr
find . -name "*.go" | entr ./deploy.sh