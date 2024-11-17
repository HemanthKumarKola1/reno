package repo

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Repo struct {
	cs *gocql.Session
}

var repoSession *Repo

func NewRepo(s *gocql.Session) {
	repoSession = &Repo{s}
}

// Notification represents a notification entity
type Notification struct {
	UserID           uuid.UUID
	NotificationTime time.Time
	NotificationType string
	Message          string
}

// Create a new notification
func (t *Repo) createNotification(w http.ResponseWriter, r *http.Request) {
	var notification Notification
	err := json.NewDecoder(r.Body).Decode(&notification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert into Cassandra
	err = repoSession.cs.Query("INSERT INTO notifications (user_id, notification_time, notification_type, message) VALUES (?, ?, ?, ?)").
		WithContext(context.Background()).
		Bind(notification.UserID, notification.NotificationTime, notification.NotificationType, notification.Message).
		Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update an existing notification (based on user_id and notification_time)
func (t *Repo) updateNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userIDBytes, er := base64.StdEncoding.DecodeString(vars["user_id"])
	if er != nil {
		// Handle error
	}

	userID, e := uuid.ParseBytes(userIDBytes)
	if e != nil {
		// Handle error
	}
	notificationTime, _ := time.Parse(time.RFC3339, vars["notification_time"])

	var updatedNotification Notification
	err := json.NewDecoder(r.Body).Decode(&updatedNotification)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the notification in Cassandra
	err = repoSession.cs.Query("UPDATE notifications SET notification_type = ?, message = ? WHERE user_id = ? AND notification_time = ?").
		WithContext(context.Background()).
		Bind(updatedNotification.NotificationType, updatedNotification.Message, userID, notificationTime).
		Exec()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (r *Repo) repo() {

	// Define the time range for the query
	now := time.Now()
	oneMinuteAgo := now.Add(-1 * time.Minute)

	// Query Cassandra for notifications within the time range
	var notifications []Notification
	iter := repoSession.cs.Query("SELECT * FROM notifications WHERE notification_time >= ? AND notification_time <= ?").
		WithContext(context.Background()).
		Bind(oneMinuteAgo, now).
		Iter()
	var notification Notification
	for iter.Scan(&notification.UserID, &notification.NotificationTime, &notification.NotificationType, &notification.Message) {
		notifications = append(notifications, notification)
	}
	if err := iter.Close(); err != nil {
		panic(err)
	}

	// Process the retrieved notifications
	for _, notification := range notifications {
		// Implement your notification delivery logic here
		fmt.Printf("User ID: %s, Notification Time: %s, Type: %s, Message: %s\n", notification.UserID, notification.NotificationTime, notification.NotificationType, notification.Message)
		// Update notification status in Cassandra or other systems
		// ...
	}
}
