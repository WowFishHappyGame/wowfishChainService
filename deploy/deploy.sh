#!/bin/bash

docker load -i wowfish.tar

docker-compose -f docker-compose.yaml up -d