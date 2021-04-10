package entities

// User entity
type User struct {
	ID        string `json:"id" bson:"_id"`
	Email     string `json:"email" bson:"email"`
	FirstName string `json:"firstName" bson:"firstName"`
	LastName  string `json:"lastName" bson:"lastName"`
	Password  string `json:"password" bson:"password"`
}
