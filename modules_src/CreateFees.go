package main

import (
	// "fmt"
	"encoding/json"
	"net/http"
	// "io"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"firebase.google.com/go/v4/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"fmt"
)

type Fees struct {
	UID string `json:"UID" bson:"UID"`
	Amount string `json:"Amount" bson:"Amount"`
	ToBePaid string `json:"ToBePaid" bson:"ToBePaid"`
	TimeStamp string `json:"TimeStamp" bson:"TimeStamp"`
}

func CreateFees(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
	fmt.Println("HTTP Method:", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := authClient.VerifyIDToken(ctx, r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Token verified:", token.UID)

	collection := db.Database("KVM").Collection("Users")
	var Role bson.M
	err = collection.FindOne(ctx, bson.M{"UID": token.UID}, options.FindOne().SetProjection(bson.M{"Role": 1, "_id": 0})).Decode(&Role)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if Role["Role"] == nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}
	if Role["Role"] != "admin" && Role["Role"] != "DEVELOPER" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var data Fees
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	data.TimeStamp = time.Now().Format(time.RFC3339)

	collection = db.Database("KVM").Collection("Fees")
	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		http.Error(w, "Insert error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Fees created successfully",
		"insertedID": result.InsertedID,
	})

}
