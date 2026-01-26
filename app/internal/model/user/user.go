package user

import (
	"time"
)

type User struct {
	ID          string    `bson:"_id,omitempty" json:"id"`
	Username    string    `bson:"username" json:"username"`
	Email       string    `bson:"email" json:"email"`
	Password    string    `bson:"password" json:"-"`
	CreatedAt   time.Time `bson:"createdAt" json:"createdAt"`
	LastLoginAt time.Time `bson:"lastLoginAt" json:"lastLoginAt"`
	IsActive      bool      `bson:"isActive" json:"isActive"`
	Roles         []string  `bson:"roles" json:"roles"`
	MFASecret     string    `bson:"mfaSecret,omitempty" json:"-"`
	MFAEnabled    bool      `bson:"mfaEnabled" json:"mfaEnabled"`
	RecoveryCodes []string  `bson:"recoveryCodes,omitempty" json:"-"`
}
