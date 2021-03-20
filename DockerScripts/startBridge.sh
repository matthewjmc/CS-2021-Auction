#!/bin/bash
echo "Starting Bridge"
sudo docker network create -d macvlan \
  --subnet=10.0.48.0/20 \
  --ip-range=10.0.50.1/24 \
  --gateway=10.0.56.1 \
  --aux-address="my-router=10.0.50.222" \
  -o parent=eth1 macvlan
