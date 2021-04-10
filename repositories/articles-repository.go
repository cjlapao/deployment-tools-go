package repositories

import (
	"context"
	"log"
	"time"

	"github.com/cjlapao/deployment-tools-go/entities"

	"github.com/rs/xid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetArticleByID Gets a record by ID
func (a *Repository) GetArticleByID(id string) entities.Article {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var article entities.Article

	filter := bson.D{{Key: "_id", Value: id}}
	collection := a.Factory.Database.Collection(ArticlesCollectionName)

	collection.FindOne(ctx, filter).Decode(&article)

	return article
}

// GetAllArticles Gets all the records from the Article collection
func (a *Repository) GetAllArticles() []entities.Article {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var articles []entities.Article

	filter := bson.D{{}}
	collection := a.Factory.Database.Collection(ArticlesCollectionName)

	cur, err := collection.Find(ctx, filter)
	if err != nil {
		log.Fatal(err)
	}
	for cur.Next(ctx) {
		var article entities.Article
		err := cur.Decode(&article)
		if err != nil {
			log.Fatal(err)
		}

		articles = append(articles, article)
	}

	return articles
}

// UpsertArticle Update/Insert and article in the database
func (a *Repository) UpsertArticle(article entities.Article) *mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(true)
	if len(article.ID) == 0 {
		article.ID = xid.New().String()
	}

	filter := bson.D{{Key: "_id", Value: article.ID}}
	collection := a.Factory.Database.Collection(ArticlesCollectionName)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: article.Title}, {Key: "description", Value: article.Description}, {Key: "content", Value: article.Content}}}}

	upsertResult, err := collection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		log.Fatal(err)
	}

	return upsertResult
}

// UpdateArticle Updates an article record
func (a *Repository) UpdateArticle(article entities.Article) *mongo.UpdateResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Update().SetUpsert(false)
	if len(article.ID) == 0 {
		article.ID = xid.New().String()
	}
	filter := bson.D{{Key: "_id", Value: article.ID}}
	collection := a.Factory.Database.Collection(ArticlesCollectionName)
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: article.Title}, {Key: "description", Value: article.Description}, {Key: "content", Value: article.Content}}}}

	upsertResult, err := collection.UpdateOne(ctx, filter, update, opts)

	if err != nil {
		log.Fatal(err)
	}

	return upsertResult
}

// UpsertManyArticles Bulk Inserts/Updates into the Articles collection
func (a *Repository) UpsertManyArticles(articles []entities.Article) *mongo.BulkWriteResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var operations []mongo.WriteModel
	for i := 0; i < len(articles); i++ {
		if len(articles[i].ID) == 0 {
			articles[i].ID = xid.New().String()
		}
		filter := bson.D{{Key: "_id", Value: articles[i].ID}}
		update := bson.D{{Key: "$set", Value: bson.D{{Key: "title", Value: articles[i].Title}, {Key: "description", Value: articles[i].Description}, {Key: "content", Value: articles[i].Content}}}}
		op := mongo.NewUpdateOneModel()
		op.SetUpsert(true)
		op.SetFilter(filter)
		op.SetUpdate(update)
		operations = append(operations, op)
	}

	collection := a.Factory.Database.Collection(ArticlesCollectionName)
	bulkOptions := options.BulkWriteOptions{}

	upsertResult, err := collection.BulkWrite(ctx, operations, &bulkOptions)

	if err != nil {
		log.Fatal(err)
	}

	return upsertResult
}

// DeleteArticle Deletes a record from the Article collection
func (a *Repository) DeleteArticle(id string) *mongo.DeleteResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: id}}

	collection := a.Factory.Database.Collection(ArticlesCollectionName)

	deleteResult, err := collection.DeleteOne(ctx, filter)

	if err != nil {
		log.Fatal(err)
	}

	return deleteResult
}
