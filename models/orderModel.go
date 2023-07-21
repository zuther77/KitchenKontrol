package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID         primitive.ObjectID `bson:"_id"`
	Order_id   string             `json:"order_id"`
	Order_date time.Time          `json:"order_date" validate="required" `
	Created_At time.Time          `json:"created_at"`
	Updated_At time.Time          `json:"updated_at"`
	Table_id   *string            `json:"table_id" validate="required" `
}
