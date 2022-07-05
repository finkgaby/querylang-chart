package main

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
	srv "querylang-chart/server"
	"runtime"
)

const defaultPort = "8080"

const (
	subSubjectName = "QUERY.unserialized"
	pubSubjectName = "QUERY.serialized"
	streamName     = "QUERY"
	streamSubjects = "QUERY.*"
)

type Query struct {
	Query     string
	SubjectId string
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
	nc, _ := nats.Connect("localhost:4222")
	log.Println("Creates JetStreamContext")
	js, err := nc.JetStream()
	checkErr(err)

	createStream(js)

	log.Printf("Create durable consumer monitor on subject:%q", subSubjectName)
	js.Subscribe(subSubjectName, func(msg *nats.Msg) {
		log.Printf("message incoming")
		msg.Ack()
		var query Query
		err := json.Unmarshal(msg.Data, &query)
		checkErr(err)
		log.Printf("Subscriber fetched msg.Data:%s from subSubjectName:%q", string(msg.Data), msg.Subject)
		reviewQuery(js, query)
	}, nats.Durable("monitor"), nats.ManualAck())

	runtime.Goexit()
}

func createStream(js nats.JetStreamContext) {
	log.Printf("DeleteStream old stream: %q", streamName)
	js.DeleteStream(streamName)

	stream, _ := js.StreamInfo(streamName)
	if stream == nil {
		log.Printf("creating new stream %q and subjects %q", streamName, streamSubjects)
		_, err := js.AddStream(&nats.StreamConfig{
			Name:     streamName,
			Subjects: []string{streamSubjects},
		})
		checkErr(err)
	}
}

func reviewQuery(js nats.JetStreamContext, query Query) {
	desQuery := srv.DeserializeQuery(query.Query)
	query.Query = desQuery.Query
	queryJSON, _ := json.Marshal(query)
	_, err := js.Publish(pubSubjectName, queryJSON)
	checkErr(err)
	log.Printf("Published queryJSON:%s to subjectName:%q", string(queryJSON), pubSubjectName)
}

func checkErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
