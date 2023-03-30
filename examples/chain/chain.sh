#!/bin/bash
# chain multiple calls to rects with different configs and 
# pass each stage to the next invocation
dn=$(dirname -- "$0")
rects -c $dn/config1.json -o $dn/stage1.png
rects -c $dn/config2.json -i $dn/stage1.png -o $dn/stage2.png