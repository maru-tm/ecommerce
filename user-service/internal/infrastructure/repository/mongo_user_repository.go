package repository

import (
	"context"
	"fmt"
	"log"

	"user-service/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type userRepository struct {
	db         *mongo.Database
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) domain.UserRepository {
	return &userRepository{
		db:         db,
		collection: db.Collection("users"),
	}
}

func (r *userRepository) CreateUser(user *domain.User) (*domain.User, error) {
	ctx := context.Background()
	log.Printf("[CreateUser] Creating user ID=%s", user.ID)

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		err = fmt.Errorf("failed to create user with ID '%s': %w", user.ID, err)
		log.Printf("[CreateUser] Error: %v", err)
		return nil, err
	}

	log.Printf("[CreateUser] Successfully created user ID=%s", user.ID)
	return user, nil
}

func (r *userRepository) GetUserByID(id string) (*domain.User, error) {
	ctx := context.Background()
	log.Printf("[GetUserByID] Fetching user ID=%s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = fmt.Errorf("invalid ObjectID format for user ID '%s': %w", id, err)
		log.Printf("[GetUserByID] Error: %v", err)
		return nil, err
	}

	var user domain.User
	err = r.collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("[GetUserByID] User ID=%s not found", id)
			return nil, nil
		}
		err = fmt.Errorf("failed to fetch user with ID '%s': %w", id, err)
		log.Printf("[GetUserByID] Error: %v", err)
		return nil, err
	}

	log.Printf("[GetUserByID] Found user ID=%s", id)
	return &user, nil
}

func (r *userRepository) GetUserByUsername(username string) (*domain.User, error) {
	ctx := context.Background()
	log.Printf("[GetUserByUsername] Fetching user username=%s", username)

	var user domain.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			log.Printf("[GetUserByUsername] User username=%s not found", username)
			return nil, nil
		}
		err = fmt.Errorf("failed to fetch user with username '%s': %w", username, err)
		log.Printf("[GetUserByUsername] Error: %v", err)
		return nil, err
	}

	log.Printf("[GetUserByUsername] Found user username=%s", username)
	return &user, nil
}

func (r *userRepository) ListUsers() ([]domain.User, error) {
	ctx := context.Background()
	log.Println("[ListUsers] Fetching all users")

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		err = fmt.Errorf("failed to fetch users: %w", err)
		log.Printf("[ListUsers] Error: %v", err)
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Printf("[ListUsers] Cursor close error: %v", err)
		}
	}()

	var users []domain.User
	for cursor.Next(ctx) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			err = fmt.Errorf("failed to decode user: %w", err)
			log.Printf("[ListUsers] Error: %v", err)
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		err = fmt.Errorf("cursor iteration error: %w", err)
		log.Printf("[ListUsers] Error: %v", err)
		return nil, err
	}

	log.Printf("[ListUsers] Fetched %d users", len(users))
	return users, nil
}

func (r *userRepository) UpdateUser(user *domain.User) (*domain.User, error) {
	ctx := context.Background()
	log.Printf("[UpdateUser] Updating user ID=%s", user.ID)

	objectID, err := primitive.ObjectIDFromHex(user.ID)
	if err != nil {
		err = fmt.Errorf("invalid ObjectID format for user ID '%s': %w", user.ID, err)
		log.Printf("[UpdateUser] Error: %v", err)
		return nil, err
	}

	update := bson.M{
		"$set": bson.M{
			"username":      user.Username,
			"password_hash": user.PasswordHash,
			"email":         user.Email,
			"full_name":     user.FullName,
			"status":        user.Status,
			"updated_at":    user.UpdatedAt,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		err = fmt.Errorf("failed to update user with ID '%s': %w", user.ID, err)
		log.Printf("[UpdateUser] Error: %v", err)
		return nil, err
	}

	if result.MatchedCount == 0 {
		log.Printf("[UpdateUser] No user found with ID=%s", user.ID)
		return nil, fmt.Errorf("no user found with ID %s", user.ID)
	}

	log.Printf("[UpdateUser] Successfully updated user ID=%s", user.ID)
	return user, nil
}

func (r *userRepository) DeleteUser(id string) error {
	ctx := context.Background()
	log.Printf("[DeleteUser] Deleting user ID=%s", id)

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		err = fmt.Errorf("invalid ObjectID format for user ID '%s': %w", id, err)
		log.Printf("[DeleteUser] Error: %v", err)
		return err
	}

	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		err = fmt.Errorf("failed to delete user with ID '%s': %w", id, err)
		log.Printf("[DeleteUser] Error: %v", err)
		return err
	}

	if result.DeletedCount == 0 {
		log.Printf("[DeleteUser] No user found to delete with ID=%s", id)
		return fmt.Errorf("no user found to delete with ID %s", id)
	}

	log.Printf("[DeleteUser] Successfully deleted user ID=%s", id)
	return nil
}
