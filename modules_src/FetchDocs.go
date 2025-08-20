package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"


	"firebase.google.com/go/v4/auth"
)

type RequestPayload struct {
	Query  map[string]interface{} `json:"query"`
	Paging struct {
		Page  int64 `json:"page"`
		Limit int64 `json:"limit"`
	} `json:"paging"`
	Projection map[string]interface{} `json:"projection"`
	Collection string `json:"collection"`
}

type ResponsePayload struct {
	TotalDocs    int64                      `json:"total_docs"`
	CurrentPage  int64                      `json:"current_page"`
	TotalEntries int                        `json:"total_entries"`
	Data         []map[string]interface{}   `json:"data"`
}

func MapToBsonM(input map[string]interface{}) (bson.M, error) {
	data, err := bson.Marshal(input)
	if err != nil {
		return nil, err
	}
	var result bson.M
	err = bson.Unmarshal(data, &result)
	return result, err
}

func convertDatesDeep(data interface{}) interface{} {
    switch val := data.(type) {
    case map[string]interface{}:
        for k, v := range val {
            if k == "$date" {
                if s, ok := v.(string); ok {
                    if t, err := time.Parse(time.RFC3339, s); err == nil {
                        return t
                    }
                }
            } else {
                val[k] = convertDatesDeep(v)
            }
        }
        return val
    case []interface{}:
        for i, v := range val {
            val[i] = convertDatesDeep(v)
        }
        return val
    default:
        return val
    }
}

func FetchDocs(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
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
	if Role["Role"] != "admin" && Role["Role"] != "staff" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Ensure valid defaults
	if payload.Paging.Limit <= 0 {
		payload.Paging.Limit = 10
	}
	if payload.Paging.Page <= 0 {
		payload.Paging.Page = 1
	}

	skip := (payload.Paging.Page - 1) * payload.Paging.Limit

	proj,err := MapToBsonM(payload.Projection)
	if err != nil {
		http.Error(w, "Invalid projection", http.StatusBadRequest)
		return
	}

	findOptions := options.Find().
	SetSkip(skip).
	SetLimit(payload.Paging.Limit).
	SetProjection(proj)
	collection = db.Database("KVM").Collection(payload.Collection)
	filter := convertDatesDeep(payload.Query)
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		http.Error(w, "Error reading results", http.StatusInternalServerError)
		return
	}

	count, err := collection.CountDocuments(ctx, payload.Query)
	if err != nil {
		http.Error(w, "Error counting documents", http.StatusInternalServerError)
		return
	}

	response := ResponsePayload{
		TotalDocs:    count,
		CurrentPage:  payload.Paging.Page,
		TotalEntries: len(results),
		Data:         results,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
