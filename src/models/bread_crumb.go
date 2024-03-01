package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BreadCrumb struct {
	FromIp        string             `json:"fromIp,omitempty"`
	ByUser        primitive.ObjectID `json:"byUser,omitempty"`
	AtTime        primitive.DateTime `json:"atTime,omitempty"`
	CorrelationId string             `json:"correlationId,omitempty"`
}

func NewBreadCrumb(ip string, userIdHex string, corrId string) (*BreadCrumb, error) {
	this := &BreadCrumb{}
	// Create oid for User ID
	userId, err := primitive.ObjectIDFromHex(userIdHex)
	if err != nil {
		return nil, err
	}

	this.FromIp = ip
	this.ByUser = userId
	this.CorrelationId = corrId
	this.AtTime = primitive.NewDateTimeFromTime(time.Now())
	return this, nil
}

func (crumb *BreadCrumb) AsBson() bson.M {
	return bson.M{
		"fromIp":        crumb.FromIp,
		"byUser":        crumb.ByUser,
		"atTime":        crumb.AtTime,
		"correlationId": crumb.CorrelationId,
	}
}
