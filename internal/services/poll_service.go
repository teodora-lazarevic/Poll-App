package services

import (
	"context"
	"errors"

	"github.com/teodora-lazarevic/Poll-App/ent"
	"github.com/teodora-lazarevic/Poll-App/ent/poll"
	"github.com/teodora-lazarevic/Poll-App/ent/polloption"
	"github.com/teodora-lazarevic/Poll-App/ent/user"
)

type PollService struct {
	DB *ent.Client
}

var (
	ErrUnauthorized       = errors.New("Unauthorized action")
	ErrDuplicateOption    = errors.New("Duplicate option text")
	ErrPollNotFound       = errors.New("Poll not found")
	ErrPollOptionNotFound = errors.New("Poll option not found")
)

func NewPollService(db *ent.Client) *PollService {
	return &PollService{DB: db}
}

func (s *PollService) ListPolls(ctx context.Context) ([]*ent.Poll, error) {
	return s.DB.Poll.Query().
		WithOptions().
		WithCreator().
		All(ctx)
}

func (s *PollService) GetPollById(ctx context.Context, pollId int) (*ent.Poll, error) {
	return s.DB.Poll.Query().
		Where(poll.IDEQ(pollId)).
		WithOptions().
		WithCreator().
		Only(ctx)
}

func (s *PollService) CreatePoll(ctx context.Context, userId int, title, description string, options []string) (*ent.Poll, error) {
	creator, err := s.DB.User.Query().Where(user.ID(userId)).First(ctx)
	if ent.IsNotFound(err) {
		return nil, ErrUnauthorized
	}

	tx, err := s.DB.Tx(ctx)
	if err != nil {
		return nil, err
	}

	newPoll, err := tx.Poll.Create().
		SetTitle(title).
		SetDescription(description).
		SetCreator(creator).
		Save(ctx)

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	seen := make(map[string]bool)
	for _, optionText := range options {
		if seen[optionText] {
			tx.Rollback()
			return nil, ErrDuplicateOption
		}
		seen[optionText] = true

		_, err := tx.PollOption.Create().
			SetText(optionText).
			SetPoll(newPoll).
			Save(ctx)

		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	return newPoll, tx.Commit()
}

func (s *PollService) AddPollOption(ctx context.Context, userId, pollId int, text string) (*ent.PollOption, error) {
	fetchedPoll, err := s.GetPollById(ctx, pollId)
	if err != nil {
		return nil, err
	}

	if fetchedPoll.Edges.Creator.ID != userId {
		return nil, ErrUnauthorized
	}

	exists, err := s.DB.PollOption.Query().
		Where(polloption.HasPollWith(poll.ID(pollId)), polloption.TextEQ(text)).Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrDuplicateOption
	}

	return s.DB.PollOption.Create().
		SetText(text).
		SetPoll(fetchedPoll).
		Save(ctx)
}

func (s *PollService) DeletePoll(ctx context.Context, userId, pollId int) error {
	fetchedPoll, err := s.GetPollById(ctx, pollId)
	if err != nil {
		return err
	}

	if fetchedPoll.Edges.Creator.ID != userId {
		return ErrUnauthorized
	}

	return s.DB.Poll.DeleteOneID(pollId).Exec(ctx)
}

func (s *PollService) DeletePollOption(ctx context.Context, userId, pollId, optionId int) error {
	fetchedPoll, err := s.GetPollById(ctx, pollId)
	if err != nil {
		return err
	}

	if fetchedPoll.Edges.Creator.ID != userId {
		return ErrUnauthorized
	}

	deletedCnt, err := s.DB.PollOption.
		Delete().
		Where(polloption.ID(optionId), polloption.HasPollWith(poll.ID(pollId))).
		Exec(ctx)
	if err != nil {
		return err
	}

	if deletedCnt == 0 {
		return ErrPollOptionNotFound
	}

	return nil
}
