package models

import (
	"time"
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
