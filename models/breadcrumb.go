package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

type BreadCrumb struct {
	FromIp        string `json:"fromIp,omitempty"`
	ByUser        string `json:"byUser,omitempty"`
	AtTime        string `json:"atTime,omitempty"`
	CorrelationId string `json:"correlationId,omitempty"`
}

func NewBreadCrumb(ip string, user string, corrId string) *BreadCrumb {
	this := &BreadCrumb{}
	this.FromIp = ip
	this.ByUser = user
	this.CorrelationId = corrId
	this.AtTime = time.Now().Format("2006-01-02 15:04:05")
	return this
}

func (this *BreadCrumb) AsBson() bson.M {
	return bson.M{
		"fromIp":        this.FromIp,
		"byUser":        this.ByUser,
		"atTime":        this.AtTime,
		"correlationId": this.CorrelationId,
	}
}
