# !/bin/bash
sudo docker-compose -f no-replicas.yml up --build -d

sudo docker exec -it req_counter_service /bin/bash -c "cd routes; go test;" 

sudo docker container stop req_counter_service redis_db