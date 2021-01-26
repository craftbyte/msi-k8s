package main

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Task is a todo app item that pertains to a task a user has to complete
type Task struct {
	ID   primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name string             `json:"name" bson:"name"`
	Done bool               `json:"done" bson:"done"`
}

// NewTask is the JSON body we receive when creating a new task
type NewTask struct {
	Name string `json:"name"`
}

// UpdatedTask is the body we receive when updating a task's status
type UpdatedTask struct {
	Done bool `json:"done"`
}

func taskDBWithCtx(ctx context.Context) (*mongo.Database, error) {
	clientOptions := options.Client().ApplyURI(getEnv("MONGO_URL", "mongodb://localhost:27017/"))
	mongoAuth := false
	mongoUsername, mongoAuth := os.LookupEnv("MONGO_USERNAME")
	mongoPassword, mongoAuth := os.LookupEnv("MONGO_PASSWORD")
	if mongoAuth {
		credential := options.Credential{
			Username:      mongoUsername,
			Password:      mongoPassword,
			AuthMechanism: "SCRAM-SHA-256",
		}
		clientOptions.SetAuth(credential)
	}
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}
	return client.Database("msi-docker"), nil
}

func taskDB() (*mongo.Database, error) {
	ctx := context.Background()
	return taskDBWithCtx(ctx)
}

func createTaskWithCtx(ctx context.Context, db *mongo.Database, task *Task) error {
	res, err := db.Collection("tasks").InsertOne(ctx, task)
	task.ID = res.InsertedID.(primitive.ObjectID)
	return err
}

func createTask(db *mongo.Database, task *Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return createTaskWithCtx(ctx, db, task)
}

func getAllTasksWithCtx(ctx context.Context, db *mongo.Database) ([]*Task, error) {
	filter := bson.D{{}}
	return filterTasksWithCtx(ctx, db, filter)
}

func getAllTasks(db *mongo.Database) ([]*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return getAllTasksWithCtx(ctx, db)
}

func filterTasksWithCtx(ctx context.Context, db *mongo.Database, filter interface{}) ([]*Task, error) {
	var tasks []*Task

	cur, err := db.Collection("tasks").Find(ctx, filter)
	if err != nil {
		return tasks, err
	}

	for cur.Next(ctx) {
		var t Task
		err := cur.Decode(&t)
		if err != nil {
			return tasks, err
		}

		tasks = append(tasks, &t)
	}

	if err := cur.Err(); err != nil {
		return tasks, err
	}
	cur.Close(ctx)

	if len(tasks) == 0 {
		return make([]*Task, 0), nil
	}

	return tasks, nil
}

func filterTasks(db *mongo.Database, filter interface{}) ([]*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return filterTasksWithCtx(ctx, db, filter)
}

func updateTaskWithCtx(ctx context.Context, db *mongo.Database, id primitive.ObjectID, done bool) (*Task, error) {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	update := bson.D{primitive.E{Key: "$set", Value: bson.D{
		primitive.E{Key: "done", Value: done},
	}}}

	t := &Task{}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	return t, db.Collection("tasks").FindOneAndUpdate(ctx, filter, update, &opt).Decode(t)
}

func updateTask(db *mongo.Database, id primitive.ObjectID, done bool) (*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return updateTaskWithCtx(ctx, db, id, done)
}

func deleteTaskWithCtx(ctx context.Context, db *mongo.Database, id primitive.ObjectID) error {
	filter := bson.D{primitive.E{Key: "_id", Value: id}}

	res, err := db.Collection("tasks").DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
func deleteTask(db *mongo.Database, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return deleteTaskWithCtx(ctx, db, id)
}
