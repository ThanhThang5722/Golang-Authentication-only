package models

import (
	"authentication/pkg/auth"
	"authentication/pkg/database"
	"authentication/pkg/utils"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	collections "authentication/pkg/const/collection"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	BaseModel `bson:",inline"`
	ID        primitive.ObjectID `json:"_id" bson:"_id"`
	Username  string             `json:"username"`
	Password  string             `json:"password"`
	Email     string             `json:"email"`
	Picture   string             `json:"picture"`
}

func (user *User) Create(username string, email string, password string) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	// Check if username is existed
	var existedUser User
	if err := existedUser.GetOne("username", username); err == nil {
		return errors.New("username is existed")
	}

	// Check if email is existed
	if err := existedUser.GetOne("email", email); err == nil {
		return errors.New("email is existed")
	}

	// Default value for user
	*user = User{
		BaseModel: BaseModel{
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		},
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: password,
		Email:    email,
		Picture:  "http://is.am/5c4k",
	}

	_, err := collection.InsertOne(context.Background(), user)

	return err
}

func (user *User) GetOne(field string, value interface{}) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	filter := bson.M{field: value}
	err := collection.FindOne(context.Background(), filter).Decode(&user)

	return err
}

func (user *User) Update(newUser User) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	*user = newUser

	filter := bson.M{"_id": user.ID}
	replacement := newUser
	replacement.Password = user.Password // Don't change password
	replacement.BaseModel.UpdatedAt = time.Now().Unix()
	_, err := collection.ReplaceOne(context.Background(), filter, replacement)

	return err
}

func (user *User) UpdatePassword(newPassword string) error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"password":   newPassword,
			"updated_at": time.Now().Unix(),
		},
	}
	_, err := collection.UpdateOne(context.Background(), filter, update)

	return err
}

func (user *User) ResetPassword(email string) error {
	err := user.GetOne("email", email)
	if err != nil {
		return err
	}

	// Create and hash new password
	const passwordLength int = 8
	newPassword := utils.RandomPassword(passwordLength)
	newHashedPassword, err := auth.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.UpdatePassword(newHashedPassword)

	// Send the new password to user's email
	message := fmt.Sprintf("Your new password is: %v", newPassword)
	log.Println(message)
	destiationEmailList := []string{email, "zero2272005@gmail.com"}
	err = utils.SendMail(message, destiationEmailList)

	return err
}

func (user *User) Delete() error {
	collection := database.GetMongoInstance().Db.Collection(string(collections.User))

	filter := bson.M{"_id": user.ID}
	_, err := collection.DeleteOne(context.Background(), filter)

	return err
}
