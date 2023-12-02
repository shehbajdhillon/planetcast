// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package database

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	"github.com/tabbed/pqtype"
)

type MembershipType string

const (
	MembershipTypeOWNER  MembershipType = "OWNER"
	MembershipTypeADMIN  MembershipType = "ADMIN"
	MembershipTypeMEMBER MembershipType = "MEMBER"
)

func (e *MembershipType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = MembershipType(s)
	case string:
		*e = MembershipType(s)
	default:
		return fmt.Errorf("unsupported scan type for MembershipType: %T", src)
	}
	return nil
}

type NullMembershipType struct {
	MembershipType MembershipType
	Valid          bool // Valid is true if MembershipType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullMembershipType) Scan(value interface{}) error {
	if value == nil {
		ns.MembershipType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.MembershipType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullMembershipType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.MembershipType), nil
}

type TeamType string

const (
	TeamTypePERSONAL TeamType = "PERSONAL"
	TeamTypeTEAM     TeamType = "TEAM"
)

func (e *TeamType) Scan(src interface{}) error {
	switch s := src.(type) {
	case []byte:
		*e = TeamType(s)
	case string:
		*e = TeamType(s)
	default:
		return fmt.Errorf("unsupported scan type for TeamType: %T", src)
	}
	return nil
}

type NullTeamType struct {
	TeamType TeamType
	Valid    bool // Valid is true if TeamType is not NULL
}

// Scan implements the Scanner interface.
func (ns *NullTeamType) Scan(value interface{}) error {
	if value == nil {
		ns.TeamType, ns.Valid = "", false
		return nil
	}
	ns.Valid = true
	return ns.TeamType.Scan(value)
}

// Value implements the driver Valuer interface.
func (ns NullTeamType) Value() (driver.Value, error) {
	if !ns.Valid {
		return nil, nil
	}
	return string(ns.TeamType), nil
}

type Project struct {
	ID          int64
	TeamID      int64
	Title       string
	SourceMedia string
	Created     time.Time
}

type SubscriptionPlan struct {
	ID                   int64
	TeamID               int64
	StripeSubscriptionID sql.NullString
	RemainingCredits     int64
	Created              time.Time
}

type Team struct {
	ID               int64
	Slug             string
	Name             string
	StripeCustomerID sql.NullString
	TeamType         TeamType
	Created          time.Time
}

type TeamMembership struct {
	ID             int64
	TeamID         int64
	UserID         int64
	MembershipType MembershipType
	Created        time.Time
}

type Transformation struct {
	ID             int64
	ProjectID      int64
	TargetLanguage string
	TargetMedia    string
	Transcript     pqtype.NullRawMessage
	IsSource       bool
	Status         string
	Progress       float64
	Created        time.Time
}

type Userinfo struct {
	ID       int64
	Email    string
	FullName string
	Created  time.Time
}
