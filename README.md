# Replicas_ReqCount
The purpose of this project is to implement replicas of the same service and each these replicas respond with the number of the served requests by the instance and the cluster. This project was developed in Go, Redis, Docker containers and Docker Swarm.

## Contents
* [How to install](#how-to-install)
* [How to run](#how-to-run)
* [Architecture](#arch)
* [Future Work](#future)

<a name="how-to-install"></a>
* ## How to install
The service was developed to run in docker containers. So, firstly you have to install [**docker**](https://docs.docker.com/engine/install/ubuntu/) and [**docker-compose**](https://docs.docker.com/compose/install/) in your machine.

<a name="how-to-run"></a>
* ## How to run
When docker installation finished, you must add into the project your personal .env file, in which you determine your personal confgirutations for your service, e.g. the service PORT, the service HOSTNAME etc. 

The .env file must have the same structure as follows:

```.env
SERVICE_HOST=localhost
SERVICE_PORT=8083
CONTAINER_DB_HOST=db
CONTAINER_DB_PORT=6379
```

For this project version, the .env file must be placed in two different places. Firstly, you have to set **.env** in the root git folder **/Replicas_ReqCount** and then copy the same **.env** file into the **/routes** package folder( .env in /routes folder is for testing purposes).

The following picture shows us the .env files in our project:

<img width="91" alt="env_replicas" src="https://user-images.githubusercontent.com/25862065/141652423-4505cb70-0a5c-4407-80df-c567987c16a9.png">

After setting the **.env** file you can run the test for our infrustucture to be sure for our service functionality. Run the **testing_infrastructure.sh** script which is building and image of our service and deploying
our service and a [**Redis**](https://redis.io/) in docker containers with no replicas. Then, the script run our test for the http Status code, return the result and then stop these containers. In this version, I had
only one test for checking the http status code in our routes package but in the future, my goal is to implement more tests for our service in order to evaluate the behaviour in as many cases. Moreover, in a future 
version I wll make an automating procedure that will run the tests and if the results do not fail, start automatically the delpoyment.

```
bash testing_infrastructure.sh
```

Finally, now we can run the **start_service.sh** script, which will deploy our sevrice with its replicas. First, it will built our image, the it will initialize the [**Docker Swarm**](https://docs.docker.com/engine/swarm/) if it was not been before, and then it will deploy the the sevice in our docker swarm.

```
bash start_service.sh
```

Now your service is up and running in ***http://localhost:PORT*** where **PORT** is the port you have defined in your .env file.

<a name="arch"></a>
* ## Architecture
There are a variety of different architecture solutions to apply developing this service. In this chapter I will analyze some of the options reffering to
their props and cons and also analyze the choosen architecture and tools.

### Instances
Every replica of our service is a different instance of the same service. One architecture sollution was to built our sollution with a microservice 
architecture meaning that every intance would have its own database in order to save its data. Although, for the project purposes, we do not need to save 
any data expect from the atomic counter for the served requests. 

This counter inform us the number of the requests that have been served by a specific 
instance while it was alive. For any reasons, this instance may be killed or stoped, so our load balancer, in this case Docker Swarm, is responsible to 
generate a totally new instance. This means, that we have an absolutely new instance-replica of that service, which now it will count its own number of 
served requests. So, I decide that we do not need a database for each instance, but we only want to keep this counter into the instance memory, as it is 
depicting the server request number of that and only that instance.

The two different architectures are being presented in the following figures:

![micro](https://user-images.githubusercontent.com/25862065/141655127-9c07a9da-662d-485e-974e-c275ca61d5e0.png)        

<img width="85" alt="ram" src="https://user-images.githubusercontent.com/25862065/141655224-15dea4ee-4748-4c4f-8d9f-299edb3e171f.png">


### Synchronize State
Every replica-instance service must know the total number of served requests of the cluster. One way to do sync the total counter into all clusters is to 
have an internal communication every database instance with each other. This a bad solution as it is increasing the communicatoin overhead as the replicas 
grow in number. 
Another approach in this architecture is to have a master database and the replicas. Each of the replica database must update the total count to the master, 
and with its turn, the master update them all. This approach reduce the internal communication overhead of the first case.

In the choose architecture, in which we do not use databases, so we store the atomic number of served requests and the number of requests served by the cluster into
the memory of each instance. To sync their state, one approach is to use an MQTT broker, that holds a topic, e.g. **'total_count'** and a publisher service. Each of the 
replicas instances subscribe to that topic and each instance send the total count into the **Publisher** service. The Publisher push that data each time it receives 
from one instance to the **'total_count' topic**, and as a result every instance that is a subscribed to that topic, receives and updates the total_count number.

The advantages of that architecture is that , it is fast but has some drwabacks. Let's analyze the case where an instance, for any reason stopped or killed by the load
balancer, then the load balancer will replace it with a new one. It doesn't matter if the new instance lost some of the history transactions as it only needs the latest
clusters total number of requests. The problem occurs, when the load balancer, choose to send the request to new one instance. The new instance, has total cluster number
equal to zero, and forwards this into the Publisher, means that our cluster will receive wrong data.

Another drawback of this solution, is that we do not keep the total number of served requests, and in case the system restarts, we will lose it.
So, I decide to use a database as a master service that will keep the thruth of our system. The data that is going to be saved into our database do not have relation, but 
they are <key,value>. So I choose to use, a NoSQL key,value database, the **Redis**.

Each instance, before sending back the response, asks the database for the total number of requests, then it updates the local total number and sends back this number to the database.
The database ensures, that in case the entire system restarts, we will not lost our data. To do that, we use the Redis persistence mode. A problem that might occur in that 
approahc, is when an instance receives the total number from the database, before updating the local total number and sends it back to the database, another instance 
may catch the first one by receiving the old total number from the database, update its local total number and send it back into the database before the first one. 

To understand better the problem we can see the followin figure:

![case_3_4](https://user-images.githubusercontent.com/25862065/141658202-7e66ed72-02cb-4247-a2a5-991854ea9cbf.png)

To solve this problem, I use an interesting feature of Redis which is called pipelining. Pipelining is  network optimization that allows a Redis client to send multiple requests to the server without waiting for replies and reading all of them at once. So, each instance now, increments the total number and receives it at once, avoiding
race conditions between the instances. A code example is the following:

```golang
func IncrGet_total_count() string {
	var ctx = context.TODO()

	db, err_client := newClient()
	if err_client != nil {
		fmt.Println(fmt.Errorf("Failed to connect to client with error: %s", err_client.Error()))
	}
	pipe := db.Client.TxPipeline()

	total_count := pipe.Incr(ctx, "total_count")

	_, err_pipe := pipe.Exec(ctx)
	if err_pipe != nil {
		panic(err_pipe)
	}

	return strconv.FormatInt(total_count.Val(), 10)
}
```
<a name="future"></a>
* ## Future Work

1. Read more about the Goland language and refactor my code.
2. Develop more tests for our service and automate the procedure. Build the infrastructure only if all tests pass.
3. Use a production orchestrator like Kubernetes.
4. Change my file structure into the golang's standard. [**Golang File Structure**](https://github.com/golang-standards/project-layout)