package models
import (
	// "github.com/asdine/storm/v3"
	"time"
)
type Picture struct {
	ID int `storm:"id,increment"` // primary key with auto increment
	// Group string `storm:"index"` // this field will be indexed
	Name string // this field will not be indexed
	CreatedAt time.Time `storm:"index"`
	Data []byte
}
func SavePicture(p *Picture) error{
	return Db.Save(p)
}