package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"mongodbv2/model"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

func init() {
	err := godotenv.Load(".env")
	errorHandle(err)
	connectionString := os.Getenv("MONGODB_URI")
	const dbName = "netflix"
	const collectionName = "watchList"

	//set Client Option
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	errorHandle(err)
	fmt.Println("mongodb connection success....")
	collection = client.Database(dbName).Collection(collectionName)
	fmt.Println("collection instance is ready") // Confirmation message when collection is ready
}

func insertOneMovie(movie model.Netflix) {
	inserted, err := collection.InsertOne(context.Background(), movie)
	errorHandle(err)
	fmt.Println("inserted 1 movie with id: ", inserted.InsertedID)

}

func updateOneMovie(movieID string) {
	id, err := primitive.ObjectIDFromHex(movieID)
	errorHandle(err)
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"watched": true}}
	updateResult, err := collection.UpdateOne(context.Background(), filter, update)
	errorHandle(err)
	fmt.Println("Updated Movie Count: ", updateResult.ModifiedCount)

}

func deleteOneMovie(movieId string) {
	id, err := primitive.ObjectIDFromHex(movieId)
	errorHandle(err)
	filter := bson.M{"_id": id}
	deletedResult, err := collection.DeleteOne(context.Background(), filter)
	errorHandle(err)
	fmt.Println("Deleted movie Count : ", deletedResult.DeletedCount)
}

func deleteAllMovie() int64 {
	filter := bson.D{{}}
	deletedResult, err := collection.DeleteMany(context.Background(), filter)
	errorHandle(err)
	fmt.Println("All Deleted Movie Count: ", deletedResult.DeletedCount)
	return deletedResult.DeletedCount
}

func getAllmovies() []primitive.M {
	filter := bson.D{{}}
	cur, err := collection.Find(context.Background(), filter)
	errorHandle(err)
	var movies []primitive.M
	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		errorHandle(err)
		movies = append(movies, movie)
	}
	defer cur.Close(context.Background())
	return movies
}

// HTTP handler functions

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allMovies := getAllmovies()
	err := json.NewEncoder(w).Encode(allMovies)
	errorHandle(err)

}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "POST")
	var movie model.Netflix
	err := json.NewDecoder(r.Body).Decode(&movie)
	errorHandle(err)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)

}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "PUT")
	params := mux.Vars(r)
	updateOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "DELETE")
	params := mux.Vars(r)
	deleteOneMovie(params["id"])
	json.NewEncoder(w).Encode(params["id"])

}

func DeleteAllMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Headers", "DELETE")
	count := deleteAllMovie()
	json.NewEncoder(w).Encode(count)
}
func errorHandle(err error) {
	if err != nil {
		fmt.Println(err)

	}
}
