#!/bin/sh

sudo docker build -t req_counter_service .

sudo docker swarm init

sudo docker stack deploy -c docker-compose.yml replicas