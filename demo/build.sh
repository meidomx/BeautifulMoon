#!/bin/sh
rm -f demo.exe demo
go build -x
echo ""
echo ""
echo "========> starting..."
GODEBUG=cgocheck=0,gctrace=1 ./demo
echo "========> done."
