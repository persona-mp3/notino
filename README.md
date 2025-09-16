# Notino

Notino is a lightweight micro-service architecture that utilizes an API endpoint for creating new users. 
Following an event-driven system with notifications, new users are sent to RabbitMQ and then to a cosumer client 
to send out verification emails, using an external API [SendGrid](https://sendgrid.com/en-us)

## Installation
Make sure you have Docker and the GoLang compiler installed. You can find the instructions at [Docker](https://docker.com) and [GoLang](https://go.dev) respectively

```bash
git clone https://github.com/persona-mp3/notino.git
cd notino
```


## Usage
RabbitMQ is the message broker and the database, MySQL. You can find more about RabbitMQ [here](https://rabbitmq.com)


First, we need to run the docker file to be able to set these applications, instead of installing them locally
```bash
docker run -d --rm --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:4-management
```
That pulls the latest rabbitmq image from DockerHub, and spins up a running container in your terminal, named `rabbitmq`. 
Running, in interactive mode `-it`, your host's port is mapped to  Docker's port on ports `5672` and `15672`

If all of that is working, you can hit this into your browser: `http://localhost:15672`, and you should see a RabbitMQ Dashboard
You can use "guest" as both login credentials, unless you pre-configured yours.


And then for the database, we just use the latest version of MySQL
```bash
docker run --name mysql-container -e MYSQL_ROOT_PASSWORD=<password> -p 3306:3306 -d mysql
docker exec mysql_container bash 
mysql -u root -p # you will be prompted to enter a password
```
Pulls the latest docker version of mysql, maps host's port 3306 to Docker's and sets `root` password using `-e`
And then copy the schema into the terminal


Now that RabbitMQ and Docker has been setup, you need to include those details into the `.env` file for the server and client to connect

One last thing before using Notino, is to get the API key, if you want to. It's not neccessary for the application to work, but no emails 
will be sent out in that event.


The server is built in `Go` with a pre-configured endpoint as `http://localhost:<port>/users/create`, and the client in `Typescript`
To start the server run at port 8080:
```bash
go run main.go :8080
```

And then in another terminal:
```bash
npm run dev
```
This starts the client script.


And then you can use `curl` to test the API endpoint
```bash
curl -H '{Content-Type: application/json}' \
-d '{"firstName":"Joe", "lastName":"Roegan", "email": "<youremail>", "userName":"persona-mp3@github.com"}' \ 
-X POST \
http://localhost:<port>/users/create
```

The server should instantly send a response back, and then on the client, you sould see a new notification

