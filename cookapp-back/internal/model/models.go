package model

import "time"

type User struct {
	ID         int       `json:"id"`
	Username   string    `json:"username"`
	Email      string    `json:"email"`
	Password   string    `json:"-"`
	Bio        string    `json:"bio"`
	AvatarPath string    `json:"avatar_path"`
	CreatedAt  time.Time `json:"created_at"`
}

type Recipe struct {
	ID            int           `json:"id"`
	UserID        int           `json:"user_id"`
	Title         string        `json:"title"`
	ImagePath     string        `json:"image_path"`
	Ingredients   []*Ingredient `json:"ingredients"`
	Steps         []*Step       `json:"steps"`
	LikeCount     int           `json:"like_count"`
	FavoriteCount int           `json:"favorite_count"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type Ingredient struct {
	ID       int    `json:"id"`
	RecipeID int    `json:"recipe_id"`
	Name     string `json:"name"`
	Quantity string `json:"quantity"`
	Unit     string `json:"unit"`
}

type Step struct {
	ID          int    `json:"id"`
	RecipeID    int    `json:"recipe_id"`
	Order       int    `json:"order"`
	Description string `json:"description"`
}

type Friendship struct {
	ID             int       `json:"id"`
	RequesterID    int       `json:"requester_id"`
	ReceiverID     int       `json:"receiver_id"`
	Status         string    `json:"status"`
	FriendUsername string    `json:"friend_username"`
	CreatedAt      time.Time `json:"created_at"`
}

type Like struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	RecipeID  int       `json:"recipe_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Favorite struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	RecipeID  int       `json:"recipe_id"`
	CreatedAt time.Time `json:"created_at"`
}
