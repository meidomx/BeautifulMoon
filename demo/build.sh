#!/bin/sh
go build -x
echo ""
echo ""
echo "========> starting..."
GODEBUG=cgocheck=0 ./demo
echo "========> done."
