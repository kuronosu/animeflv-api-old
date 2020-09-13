package db

// Serial represents a sequence document
type Serial struct {
	ID  string `bson:"_id"`
	Seq int    `bson:"seq"`
}
