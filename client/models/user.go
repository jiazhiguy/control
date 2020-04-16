package models
import(
	"time"
)
type User struct {
  ID int // primary key
  Group string `storm:"index"` // this field will be indexed
  Email string `storm:"unique"` // this field will be indexed with a unique constraint
  Name string // this field will not be indexed
  Age int `storm:"index"`
  CreatedAt time.Time `storm:"index"`
}