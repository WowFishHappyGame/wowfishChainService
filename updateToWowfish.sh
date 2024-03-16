#!/bin/bash -e

git pull

#cp -r ./* ../../wowfish/wowfishChainService

rsync -av --exclude='.git' ./ ../../wowfish/wowfishChainService/

cd ../../wowfish/wowfish


