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
	SessionId string   `json:"SessionId"`
	UIDs      []string `json:"UIDs"`
	MarkedAt  string   `json:"MarkedAt"`
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
	if Role["Role"] != "admin" && Role["Role"] != "staff" && Role["Role"] != "DEVELOPER" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var data Attendance
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Parse MarkedAt string to time.Time
	var markedAtTime time.Time
	if data.MarkedAt != "" {
		markedAtTime, err = time.Parse(time.RFC3339, data.MarkedAt)
		if err != nil {
			// Try other common formats if RFC3339 fails
			markedAtTime, err = time.Parse("2006-01-02T15:04:05.000Z", data.MarkedAt)
			if err != nil {
				http.Error(w, "Invalid date format for MarkedAt. Use RFC3339 format (e.g., 2024-01-15T10:30:00.000Z)", http.StatusBadRequest)
				return
			}
		}
	} else {
		// If no MarkedAt provided, use current time
		markedAtTime = time.Now()
	}

	// Check for duplicate attendance for the same session on the same day
	// Normalize the date to compare only the date part (ignore time)
	dateOnly := time.Date(markedAtTime.Year(), markedAtTime.Month(), markedAtTime.Day(), 0, 0, 0, 0, markedAtTime.Location())
	nextDay := dateOnly.AddDate(0, 0, 1)

	collection = db.Database("KVM").Collection("Attendance")

	// Delete existing attendance records for this session on this date
	_, err = collection.DeleteMany(ctx, bson.M{
		"Session": data.SessionId,
		"MarkedAt": bson.M{
			"$gte": dateOnly,
			"$lt":  nextDay,
		},
	})
	if err != nil {
		http.Error(w, "Database error while removing existing attendance: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var docs []interface{}

	docs = append(docs, bson.M{
		"Session":  data.SessionId,
		"UID":      token.UID,
		"Role":     "staff",
		"MarkedAt": markedAtTime,
	})

	for _, uid := range data.UIDs {
		docs = append(docs, bson.M{
			"Session":  data.SessionId,
			"UID":      uid,
			"Role":     "student",
			"MarkedAt": markedAtTime,
		})
	}

	_, err = collection.InsertMany(ctx, docs)
	if err != nil {
		http.Error(w, "Insert error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Attendance marked successfully for session on " + dateOnly.Format("2006-01-02"),
	})
}
