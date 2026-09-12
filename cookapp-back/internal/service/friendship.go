package service

import (
	"cookapp/internal/model"
	"cookapp/internal/repository"
	"errors"
)

type FriendshipService struct {
	repo *repository.FriendshipRepository
}

func NewFriendshipService(repo *repository.FriendshipRepository) *FriendshipService {
	return &FriendshipService{repo: repo}
}

func (s *FriendshipService) CreateFriendRequest(requesterID int, receiverID int) error {
	err := s.repo.CreateFriendRequest(requesterID, receiverID)
	return err
}

func (s *FriendshipService) UpdateFriendRequestStatus(friendshipID int, status string, userID int) error {
	friendship, err := s.repo.GetFriendshipByID(friendshipID)
	if err != nil {
		return err
	}
	if friendship.ReceiverID != userID {
		return errors.New("forbidden")
	}
	err = s.repo.UpdateFriendRequestStatus(friendshipID, status)
	return err
}

func (s *FriendshipService) DeleteFriend(friendshipID int, userID int) error {
	friendship, err := s.repo.GetFriendshipByID(friendshipID)
	if err != nil {
		return err
	}
	if friendship.RequesterID != userID && friendship.ReceiverID != userID {
		return errors.New("forbidden")
	}
	err = s.repo.DeleteFriend(friendshipID)
	return err
}

func (s *FriendshipService) GetFriends(userID int) ([]*model.Friendship, error) {
	friendship, err := s.repo.GetFriends(userID)
	return friendship, err
}

func (s *FriendshipService) GetPendingFriendRequests(userID int) ([]*model.Friendship, error) {
	friendship, err := s.repo.GetPendingFriendRequests(userID)
	return friendship, err
}