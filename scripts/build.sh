#!/bin/sh

# Usage: ./build $(pwd)/build /usr 1.0

BUILDDIR=$1
PREFIX=$2
VERSION=$3

prepare() {
	PREV_DIR=$(pwd)
	mkdir $BUILDDIR
	cd $BUILDDIR
	mkdir -p etc/systemd/system $PREFIX/bin
	cd $PREV_DIR
}

build() {
	prepare $BUILDDIR $PREFIX
	cd ./cmd
	go build -o $BUILDDIR/$PREFIX/bin/scheduled
	cd ..
	cp straggle/* $BUILDDIR/
}

build
