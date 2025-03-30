package entities

import (
	"github.com/HunDun0Ben/bs_server/app/entities/insect"
)

type Result struct {
	ID         string        `bson:"_id,omitempty"`
	ImageID    string        `bson:"imageId"`
	Insect     insect.Insect `bson:"insect"`
	Confidence float64       `bson:"confidence"`
}
