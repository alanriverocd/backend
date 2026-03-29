package models

import "time"

type User struct {
	UID         string    `json:"uid" firestore:"uid"`
	Email       string    `json:"email" firestore:"email"`
	DisplayName string    `json:"displayName" firestore:"displayName"`
	PhotoURL    string    `json:"photoUrl,omitempty" firestore:"photoUrl,omitempty"`
	CreatedAt   time.Time `json:"createdAt" firestore:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt" firestore:"updatedAt"`
}

type UpdateUserRequest struct {
	DisplayName string `json:"displayName,omitempty"`
	PhotoURL    string `json:"photoUrl,omitempty"`
}
