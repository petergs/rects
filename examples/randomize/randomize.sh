#!/bin/bash
if [ "$1" = "" ]; then
    echo "Pass an integer parameter (n) to this script to generate (n) random images."
    exit 0
fi

dn=$(dirname -- "$0")
i=1
while [ $i -le $1 ]
do
    rects -e $dn/colors.yml -o $dn/image$i.png
    ((i++))
done