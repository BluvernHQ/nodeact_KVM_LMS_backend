package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "io"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"firebase.google.com/go/v4/auth"
)

type Attendance struct {
	Session      string `json:"Session"`
	UIDs 		[]string `json:"UIDs"`
	MarkedAt    string `json:"MarkedAt"`
}

func MarkAttendance(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
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
	if Role["Role"] != "Admin" && Role["Role"] != "Staff" && Role["Role"] != "DEVELOPER" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var data Attendance
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}


	collection = db.Database("KVM").Collection("Attendance")
	var docs []interface{}

	docs = append(docs, bson.M{
		"Session": data.Session,
		"UID": token.UID,
		"Role": "staff",
		"MarkedAt": data.MarkedAt,
	})

	for _, uid := range data.UIDs {
		docs = append(docs, bson.M{
			"Session":      data.Session,
			"UID":          uid,
			"Role": "student",
			"MarkedAt":     data.MarkedAt,
		})
	}

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		http.Error(w, "Insert error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Attendances inserted successfully",
	})
}
