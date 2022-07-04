package main

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"log"
	srv "querylang-chart/server"
	"runtime"
)

const defaultPort = "8080"

const (
	subSubjectName = "QUERY.unserialized"
	pubSubjectName = "QUERY.serialized"
)

type Query struct {
	query     string
	subjectId string
}

func main() {
	//port := os.Getenv("PORT")
	//if port == "" {
	//	port = defaultPort
	//}
	//
	//srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	//
	//http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	//http.Handle("/query", srv)
	//
	//log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	//log.Fatal(http.ListenAndServe(":"+port, nil))

	log.Println("Connect to NATS")
	nc, _ := nats.Connect("demo.nats.io")
	log.Println("Creates JetStreamContext")
	js, err := nc.JetStream()
	checkErr(err)
	log.Printf("Create durable consumer monitor on subject:%q", subSubjectName)
	js.Subscribe(subSubjectName, func(msg *nats.Msg) {
		msg.Ack()
		var query Query
		err := json.Unmarshal(msg.Data, &query)
		checkErr(err)
		log.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		reviewQuery(js, query)
	}, nats.Durable("monitor"), nats.ManualAck())

	runtime.Goexit()
}

func reviewQuery(js nats.JetStreamContext, query Query) {
	desQuery := srv.DeserializeQuery(query.query)
	query.query = desQuery.Query
	queryJSON, _ := json.Marshal(query.query)
	_, err := js.Publish(fmt.Sprintf("(%s_%s)", pubSubjectName, query.subjectId), queryJSON)
	checkErr(err)
	log.Printf("Published queryJSON:%s to subjectName:%q", string(queryJSON), pubSubjectName)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
