package main

import pb "github.com/persona-mp3/pwa/broker"

var c = pb.Connection{
	User:     "guest",
	Host:     "localhost",
	Port:     5672,
	Password: "guest",
}

func Consumer() {
	client, err := pb.NewConnection(c)
	if err != nil {
		return
	}

	defer client.Close()

	dq := pb.Queue{
		Name:      "break_prod",
		Durable:   false,
		AutoDel:   false,
		Exclusive: false,
		NoWait:    false,
	}
	q, err := client.DeclareDirectQueue(dq)
	if err != nil {
		return
	}

	// reads messages
	client.Consume(q)

}

func main() {
	Consumer()
}
