package repository

import (
	"cookapp/internal/model"
	"database/sql"
)

type FriendshipRepository struct {
	db *sql.DB
}

func NewFriendshipRepository(db *sql.DB) *FriendshipRepository {
	return &FriendshipRepository{db: db}
}

func (r *FriendshipRepository) CreateFriendRequest(requesterID int, receiverID int) error {
	_, err := r.db.Exec("INSERT INTO friendship (requester_id, receiver_id, status) VALUES ($1, $2, 'pending')", requesterID, receiverID)
	return err
}

func (r *FriendshipRepository) UpdateFriendRequestStatus(friendshipID int, status string) error {

	_, err := r.db.Exec("UPDATE friendship SET status = $1 WHERE id = $2", status, friendshipID)
	return err
}

func (r *FriendshipRepository) DeleteFriend(friendshipID int) error {
	_, err := r.db.Exec("DELETE FROM friendship WHERE id = $1", friendshipID)
	return err
}

func (r *FriendshipRepository) GetFriends(userID int) ([]*model.Friendship, error) {
	rows, err := r.db.Query(`SELECT 
		friendship.id,
		friendship.requester_id,
		friendship.receiver_id,
		friendship.status,
		friendship.created_at,
		CASE WHEN friendship.requester_id = $1 THEN receiver.username ELSE requester.username END as friend_username
		FROM friendship
		JOIN users AS requester ON friendship.requester_id = requester.id
		JOIN users AS receiver ON friendship.receiver_id = receiver.id
		WHERE (friendship.requester_id = $2 OR friendship.receiver_id = $3) 
		AND friendship.status = 'accepted'`, userID, userID, userID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []*model.Friendship
	for rows.Next() {
		f := &model.Friendship{}
		err := rows.Scan(&f.ID, &f.RequesterID, &f.ReceiverID, &f.Status, &f.CreatedAt, &f.FriendUsername)
		if err != nil {
			return nil, err
		}
		friendships = append(friendships, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if friendships == nil {
		friendships = []*model.Friendship{}
	}

	return friendships, nil
}

func (r *FriendshipRepository) GetFriendshipByID(friendshipID int) (*model.Friendship, error) {
	friendship := &model.Friendship{}
	err := r.db.QueryRow("SELECT id, requester_id, receiver_id, status, created_at FROM friendship WHERE id = $1", friendshipID).Scan(
		&friendship.ID,
		&friendship.RequesterID,
		&friendship.ReceiverID,
		&friendship.Status,
		&friendship.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return friendship, nil
}

func (r *FriendshipRepository) GetPendingFriendRequests(userID int) ([]*model.Friendship, error) {
	rows, err := r.db.Query(`SELECT 
		friendship.id,
		friendship.requester_id,
		friendship.receiver_id,
		friendship.status,
		friendship.created_at,
		requester.username as friend_username
		FROM friendship
		JOIN users AS requester ON friendship.requester_id = requester.id
		WHERE friendship.receiver_id = $1
		AND friendship.status = 'pending'`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friendships []*model.Friendship
	for rows.Next() {
		f := &model.Friendship{}
		err := rows.Scan(&f.ID, &f.RequesterID, &f.ReceiverID, &f.Status, &f.CreatedAt, &f.FriendUsername)
		if err != nil {
			return nil, err
		}
		friendships = append(friendships, f)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if friendships == nil {
		friendships = []*model.Friendship{}
	}

	return friendships, nil
}
