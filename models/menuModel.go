package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID         primitive.ObjectID `bson:"_id"`
	Name       string             `json:"name" validate:"required,min=2,max=100"`
	Category   string             `json:"category" validate:"required,min=2,max=100"`
	Start_date *time.Time         `json:"start_date"`
	End_date   *time.Time         `json:"end_date"`
	Created_At time.Time          `json:"created_at"`
	Updated_At time.Time          `json:"updated_at"`
	Menu_id    string             `json:"food_id"`
}
