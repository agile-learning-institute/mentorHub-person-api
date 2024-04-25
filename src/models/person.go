package models

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserName    string             `json:"userName,omitempty"`
	FirstName   string             `json:"firstName,omitempty"`
	LastName    string             `json:"lastName,omitempty"`
	Roles       []string           `json:"roles,omitempty"`
	Description string             `json:"description,omitempty"`
	Status      string             `json:"status,omitempty"`
	Title       string             `json:"title,omitempty"`
	Cadence     string             `json:"cadence,omitempty"`
	Email       string             `json:"eMail,omitempty"`
	GitHub      string             `json:"gitHub,omitempty"`
	Phone       string             `json:"phone,omitempty"`
	Device      string             `json:"device,omitempty"`
	Location    string             `json:"location,omitempty"`
	MentorId    primitive.ObjectID `json:"mentorId,omitempty"`
	PartnerId   primitive.ObjectID `json:"partnerId,omitempty"`
	LastSaved   *BreadCrumb        `json:"lastSaved,omitempty"`
}

func (p Person) String() string {
	// Use a StringBuilder to efficiently build the string representation
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Person[ID=%s, UserName=%s, FirstName=%s, LastName=%s, Email=%s, Phone=%s, Roles=[%s]]",
		p.ID.Hex(), p.UserName, p.FirstName, p.LastName, p.Email, p.Phone, strings.Join(p.Roles, ", ")))

	return sb.String()
}
