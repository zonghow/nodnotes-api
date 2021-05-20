package models

import (
	"database/sql/driver"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NodeID uuid.UUID

// StringToNodeID -> parse string to NodeID
func StringToNodeID(s string) (NodeID, error) {
	id, err := uuid.Parse(s)
	return NodeID(id), err
}

//String -> String Representation of Binary16
func (my NodeID) String() string {
	return uuid.UUID(my).String()
}

//GormDataType -> sets type to binary(16)
func (my NodeID) GormDataType() string {
	return "binary(16)"
}

func (my NodeID) MarshalJSON() ([]byte, error) {
	s := uuid.UUID(my)
	str := "\"" + s.String() + "\""
	return []byte(str), nil
}

func (my *NodeID) UnmarshalJSON(by []byte) error {
	s, err := uuid.ParseBytes(by)
	*my = NodeID(s)
	return err
}

// Scan --> tells GORM how to receive from the database
func (my *NodeID) Scan(value interface{}) error {

	bytes, _ := value.([]byte)
	parseByte, err := uuid.FromBytes(bytes)
	*my = NodeID(parseByte)
	return err
}

// Value -> tells GORM how to save into the database
func (my NodeID) Value() (driver.Value, error) {
	return uuid.UUID(my).MarshalBinary()
}

type NodeModel struct {
	gorm.Model
	ID     NodeID `gorm:"primarykey"`
	UserID uint   `gorm:"not null"`
	Value  string `gorm:"type:text;"`
	IsRoot bool
}

func (n NodeModel) TableName() string {
	return "node"
}

func (n NodeModel) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	n.ID = NodeID(id)
	return err
}
