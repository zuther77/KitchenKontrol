package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Invoice struct {
	ID               primitive.ObjectID `bson:"_id" `
	Invoice_id       string             `json:"invoice_id" `
	Order_id         string             `json:"order_id" `
	Created_At       time.Time          `json:"created_at"`
	Updated_At       time.Time          `json:"updated_at"`
	Payment_due_date time.Time          `json:"payment_due_date"`
	Payment_method   *string            `json:"payment_method" validate:"eq=CARD|eq=CASH|eq="`
	Payment_status   *string            `json:"payment_status" validate:"required,eq=PAID|eq=PENDING" `
}
