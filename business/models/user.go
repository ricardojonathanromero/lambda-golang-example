package models

import "time"

type UserDB struct {
	ID        string    `dynamodbav:"Id" json:"id"`
	Name      string    `dynamodbav:"Name" json:"name"`
	Lastname  string    `dynamodbav:"Lastname" json:"lastname"`
	Age       int32     `dynamodbav:"Age" json:"age"`
	Email     string    `dynamodbav:"Email" json:"email"`
	CreatedAt time.Time `dynamodbav:"CreatedAt" json:"created_at"`
	UpdatedAt time.Time `dynamodbav:"UpdatedAt" json:"updated_at"`
}
