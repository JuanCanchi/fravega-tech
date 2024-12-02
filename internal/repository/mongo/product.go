package mongo

import (
	"context"
	"fmt"
	"fravega-tech/internal/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"sync"
)

type ProductRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Database) *ProductRepository {
	return &ProductRepository{
		collection: db.Collection("products"),
	}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	result, err := r.collection.InsertOne(context.Background(), product)
	if err != nil {
		log.Printf("Error al insertar el producto: %v", err)
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		product.ID = oid.Hex()
	} else {
		log.Printf("Error al convertir InsertedID a ObjectID")
		return fmt.Errorf("no se pudo convertir InsertedID a ObjectID")
	}

	return err
}

func (r *ProductRepository) GetAll(name, category string) ([]domain.Product, error) {
	filter := bson.M{}

	if name != "" {
		filter["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if category != "" {
		filter["categories"] = bson.M{"$regex": category, "$options": "i"}
	}

	cursor, err := r.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var products []domain.Product
	if err := cursor.All(context.Background(), &products); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Update(id string, product *domain.Product) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	filter := bson.M{"_id": objectID}
	update := bson.M{
		"$set": bson.M{
			"name":        product.Name,
			"description": product.Description,
			"price":       product.Price,
			"categories":  product.Categories,
			"updated_at":  product.UpdatedAt,
			"image_url":   product.ImageURL,
		},
	}

	_, err = r.collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)
	if err != nil {
		log.Printf("Error al actualizar el producto: %v", err)
		return err
	}

	return err
}

func (r *ProductRepository) Delete(id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %v", err)
	}

	_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
	return err
}

func (r *ProductRepository) DeleteMany(ids []string) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var errOccurred error

	for _, id := range ids {
		wg.Add(1)

		go func(id string) {
			defer wg.Done()

			objectID, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				mu.Lock()
				errOccurred = fmt.Errorf("invalid ID format: %v", err)
				mu.Unlock()
				return
			}

			_, err = r.collection.DeleteOne(context.Background(), bson.M{"_id": objectID})
			if err != nil {
				mu.Lock()
				errOccurred = fmt.Errorf("failed to delete product %s: %v", id, err)
				mu.Unlock()
			}
		}(id)
	}

	wg.Wait()

	if errOccurred != nil {
		return errOccurred
	}

	return nil
}
