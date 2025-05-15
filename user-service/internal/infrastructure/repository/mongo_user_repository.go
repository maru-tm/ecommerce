package repository

import (
	"context"
	"fmt"

	"user-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db *mongo.Database
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *domain.User) (*domain.User, error) {
	collection := r.db.Collection("users")

	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user with ID '%s': %w", user.ID, err)
	}

	return user, nil
}

func (r *userRepository) GetUserByID(id string) (*domain.User, error) {
	collection := r.db.Collection("users")

	var user domain.User
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user with ID '%s': %w", id, err)
	}

	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*domain.User, error) {
	collection := r.db.Collection("users")

	var user domain.User
	err := collection.FindOne(context.Background(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user with username '%s': %w", username, err)
	}

	return &user, nil
}

func (r *userRepository) ListUsers() ([]domain.User, error) {
	collection := r.db.Collection("users")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer cursor.Close(context.Background())

	var users []domain.User
	for cursor.Next(context.Background()) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepository) UpdateUser(user *domain.User) (*domain.User, error) {
	collection := r.db.Collection("users")

	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": user.ID},
		bson.M{"$set": bson.M{
			"username":      user.Username,
			"password_hash": user.PasswordHash,
			"email":         user.Email,
			"full_name":     user.FullName,
			"status":        user.Status,
			"updated_at":    user.UpdatedAt,
		}},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update user with ID '%s': %w", user.ID, err)
	}

	return user, nil
}

func (r *userRepository) DeleteUser(id string) error {
	collection := r.db.Collection("users")

	_, err := collection.DeleteOne(context.Background(), bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete user with ID '%s': %w", id, err)
	}

	return nil
}
