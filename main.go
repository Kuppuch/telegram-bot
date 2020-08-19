package main

import (
	//"cloud.google.com/go/firestore"
	"context"
	fb "firebase.google.com/go"
	"fmt"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"time"

	//"cloud.google.com/go/firestore"
/*	fb "firebase.google.com/go"
	"google.golang.org/api/option"*/
	"log"
)

type FirestoreValue struct {
	CreateTime time.Time `json:"createTime"`
	// Fields is the data for this value. The type depends on the format of your
	// database. Log the interface{} value and inspect the result to see a JSON
	// representation of your database fields.
	Fields     interface{} `json:"fields"`
	Name       string      `json:"name"`
	UpdateTime time.Time   `json:"updateTime"`
}

// FirestoreEvent is the payload of a Firestore event.
type FirestoreEvent struct {
	OldValue   FirestoreValue `json:"oldValue"`
	Value      FirestoreValue `json:"value"`
	UpdateMask struct {
		FieldPaths []string `json:"fieldPaths"`
	} `json:"updateMask"`
}

func DocChange(ctx context.Context, e FirestoreEvent) error {
	log.Println(e.Value.Name)
	return nil
}

func main() {
	ctx := context.Background()
	sa := option.WithCredentialsFile("/home/kuppuch/go/test-e9d05-firebase-adminsdk-r03nj-34e9a8db2e.json")
	conf := &fb.Config{ProjectID: "test-e9d05"}
	app, err := fb.NewApp(ctx, conf, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// write your query here to match your trigger path
	q := client.Collection("users")
	qsnapIter := q.Snapshots(ctx)
	// Listen forever for changes to the query's results.
	lastValues := map[string]FirestoreValue{}
	for {
		qsnap, err := qsnapIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		for _, c := range qsnap.Changes {
			// You might check the type of change and only call the function with the right subset of events
			newValue := FirestoreValue{
				CreateTime: c.Doc.CreateTime,
				UpdateTime: c.Doc.UpdateTime,
				Name:       fmt.Sprintf("projects/%s/databases/(default)/documents/%s/%s", "test-e9d05", "users", c.Doc.Ref.ID),
				Fields:     c.Doc.Data(),
			}
			e := FirestoreEvent{
				Value: newValue,
				// have not sorted out how to mock UpdateMask
			}

			if val, ok := lastValues[c.Doc.Ref.ID]; ok {
				e.OldValue = val
				fmt.Println ("!!!!!!!!!Какое-то новое значение ", val)
			}
			lastValues[c.Doc.Ref.ID] = newValue


			// call the cloud function with the change
			DocChange(context.Background(), e)
		}
	}
}
