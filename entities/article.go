package entities

// Article entity
type Article struct {
	ID          string `json:"id" bson:"_id"`
	Title       string `json:"title"  bson:"title"`
	Description string `json:"description"  bson:"description"`
	Content     string `json:"content"  bson:"content"`
}
