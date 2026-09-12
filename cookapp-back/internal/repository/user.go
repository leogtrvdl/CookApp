package repository

import (
	"cookapp/internal/model"
	"database/sql"
	"math/rand"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(username string, email string, hashedPassword string) error {
	avatars := []string{
		"https://pub-c99b03c3b3594edf8d525f1ea103ef88.r2.dev/RobotJaune.svg",
		"https://pub-c99b03c3b3594edf8d525f1ea103ef88.r2.dev/RobotRouge.svg",
		"https://pub-c99b03c3b3594edf8d525f1ea103ef88.r2.dev/RobotViolet.svg",
	}
	defaultAvatar := avatars[rand.Intn(len(avatars))]

	_, err := r.db.Exec("INSERT INTO users (username, email, password, avatar_path) VALUES ($1, $2, $3, $4)", username, email, hashedPassword, defaultAvatar)
	return err
}

func (r *UserRepository) GetUserByUsername(username string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow("SELECT id, username, email, password FROM users WHERE username = $1", username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(userID int) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow("SELECT id, username, email, password, COALESCE(bio, ''), COALESCE(avatar_path, ''), created_at FROM users WHERE id = $1", userID).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Bio,
		&user.AvatarPath,
		&user.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(userID int, username string, email string, bio string) error {

	_, err := r.db.Exec("UPDATE users SET username = $1, email = $2, bio = $3 WHERE id = $4", username, email, bio, userID)
	return err
}

func (r *UserRepository) DeleteUser(userID int) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = $1", userID)
	return err
}

func (r *UserRepository) UpdateAvatar(userID int, avatarPath string) error {
	_, err := r.db.Exec("UPDATE users SET avatar_path = $1 WHERE id = $2", avatarPath, userID)
	return err
}

func (r *UserRepository) SearchUsers(username string) ([]*model.User, error) {
	rows, err := r.db.Query("SELECT id, username, COALESCE(avatar_path, '') FROM users WHERE LOWER(username) LIKE $1", "%"+username+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var user []*model.User
	for rows.Next() {
		u := &model.User{}
		err := rows.Scan(&u.ID, &u.Username, &u.AvatarPath)
		if err != nil {
			return nil, err
		}
		user = append(user, u)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return user, nil
}
