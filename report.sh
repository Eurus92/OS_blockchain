#!/bin/bash
# shellcheck disable=SC2164
cd ./eval
make
cd ../
pdflatex report