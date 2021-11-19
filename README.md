# Replicas_ReqCount
The purpose of this project is to implement replicas of the same service and each these replicas respond with the number of the served requests by the instance and the cluster. This project was developed in Go, Redis, Docker containers and Docker Swarm.

## Contents
* [How to install](#how-to-install)
* [How to run](#how-to-run)
* [Architecture](#arch)
* [Future Work](#future)
* [Update](#updates)

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
ContainerDBhost=db
ContainerDBport=6379
```

For this project version, the .env file must be placed in two different places. Firstly, you have to set **.env** in the root git folder **/Replicas_ReqCount** and then copy the same **.env** file into the **/routes** package folder( .env in /routes folder is for testing purposes).

The following picture shows us the .env files in our project:

<img width="91" alt="env_replicas" src="https://user-images.githubusercontent.com/25862065/141652423-4505cb70-0a5c-4407-80df-c567987c16a9.png">

After setting the **.env** file you can run the test for our infrustucture to be sure for our service functionality. Run the **testing_infrastructure.sh** script which is building an image of our service and deploying
our service and a [**Redis**](https://redis.io/) in docker containers with no replicas. Since, containers are up, test for the http Status code 
runs, returns the result and then stop the containers. In this version, I had only one test for checking the http status code in our routes package but in the future, my goal is to implement more tests for our service in order to evaluate the behaviour in as many cases. Moreover, in a future  version I wll make an automating procedure that will run the tests and if the results do not fail, then it will start automatically the delpoyment.

```
bash testing_infrastructure.sh
```

Finally, we can run the **start_service.sh** script, which will deploy our sevrice with its replicas. Firstly, it will built our image and it will 
initialize the [**Docker Swarm**](https://docs.docker.com/engine/swarm/) if it has not been before, and then it will deploy the sevice in our docker swarm.

```
bash start_service.sh
```

Now your service is up and running in ***http://localhost:PORT*** where **PORT** is the port you have defined in your .env file.

<a name="arch"></a>
* ## Architecture
There is a variety of different architecture solutions to apply developing this service. In this chapter I will analyze some of the options reffering to
their props and cons and also analyze the choosen architecture and tools.

### Instances
Every replica of our service is a different instance of the same service. One architecture solution is to built our service with a microservice 
architecture. This means that every intance would have its own database in order to save its data. Although, for the project purposes, we do not 
need to save any data expect from the atomic counter for the served requests. 

This counter inform us the number of the requests that have been served by a specific 
instance while it was alive. For any reasons, this instance may be killed or stoped, so our load balancer, in this case Docker Swarm, is 
responsible to generate a totally new instance. This means, that we have an absolutely new instance-replica of that service that has not served 
any request yet and it will count its own number of served requests. So, I decide that we do not need a database for each instance, but we only want to keep this counter into the instance memory, as it is depicting the number of served request of that instance.

The two different architectures are being presented in the following figures:

![micro](https://user-images.githubusercontent.com/25862065/141655127-9c07a9da-662d-485e-974e-c275ca61d5e0.png)        

<img width="85" alt="ram" src="https://user-images.githubusercontent.com/25862065/141655224-15dea4ee-4748-4c4f-8d9f-299edb3e171f.png">


### Synchronize State
Every replica-instance service must also know the total number of served requests of the cluster. One way to sync the total counter into all 
clusters is by having an internal communication each database with each other. This is a bad solution as it is increasing the 
communication overhead as the replicas grow in number. 
Another approach is to have a master database and the rest replicas acting as slaves. Each of the slave-database must update the total count to
the master, and with the master with its turn, update them all. This approach reduce the internal communication overhead of the first case.

In my architecture, I do not use databases, so we store the atomic number of served requests and the total number of requests served by the 
cluster into the memory of each instance. To sync their state, one approach is to use an MQTT broker, that holds a topic, e.g. **'total_count'** 
and a publisher service. Each of the replicas instances subscribe to that topic and every time they receive a request, they increase and send the 
total count into the **Publisher** service. The Publisher, each time it receives the total count from one instance service, push it into the 
**'total_count' topic**, and as a result every instance which is subscribed to that topic, receives and updates it local total count number.

The advantages of using an MQTT broker is speed but it has some important drwabacks. In case where an instance, for any reason stopped or killed 
by the load balancer, then the load balancer will replace it with a new one. It doesn't matter if the new instance has lost some of the history 
transactions as it only needs the latest total number of requests. The problem occurs, when the load balancer, choose to send the request to the 
new instance before that instance receives and updates is total number. So, having the new instance a total cluster number equal to zero, it will 
forwards this into the Publisher, which means that the rest of the services will receive 0 total count.

Another drawback of this solution, is that in case of a system restart, total count will be lost. So, I decide to use a database as a master
service. Data do not have relation, as it is only a the value of the total requests. So I choose to use **Redis** which is a NoSQL key,value 
database.

Each instance, for every received request, gets the total number of requests from Redis, increments it by one, updates the local total number and 
sends the new total count number back to the database. The database ensures, that in case the entire system restarts, we will not lost our data. 
To do that, we use the Redis persistence mode. Although, a race condition problem might occur, when an instance receives the total number from the 
database, and before updating the database with is local total count, another instance service also asks the daabase and manage to update first 
the database. 

The following figures describes the above problem:

![case_3_4](https://user-images.githubusercontent.com/25862065/141658202-7e66ed72-02cb-4247-a2a5-991854ea9cbf.png)

I solved the race condition problem by using an interesting feature of Redis which is called pipelining. Pipelining is a network optimization that 
allows a Redis client to send multiple requests to the server without waiting for replies and reading all of them at once. Moreover, each instance 
now, increments the total number and receives it at one transaction, avoiding race conditions between the instances. A code example is the following:

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

<a name="updates"></a>
* ## Updates

1. **Code style**
Use of golang.org/x/lint/golint linter.

2. **context.Background()**
Set it only in one place in the code.

3. **string formatting**
Use of fmt.Sprinf functionality to build and append strings.

4. **Redis error handling**
If INCR key fail, retry a second time, otherwise send back to the client error response StatusNotImplemented.

5. **Do not use global variables(e.g. PORT)**
Setting the code to read once the configurations for our service, like service PORT, and then pass it 
as an argument to every function that needs that information.

6. **Get the hostname once**
Set the hostname once in our code, and then pass it as an argument to our router handle function, instead of getting it in every request.

7. **Local Counter concurrency**
Protect the local counter from concurrent writes(race conditions) by adding sync with mutexes.

8. **Finish Test**
Compare also the response body. 