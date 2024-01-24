#!/bin/bash -e

tag="v"$(date '+%y%m%d%H%M')

docker build . -t wowfish:$tag

docker tag wowfish:$tag wowfish/compose:v1.0.0

docker save wowfish/compose:v1.0.0 -o ./deploy/wowfish.tar 
