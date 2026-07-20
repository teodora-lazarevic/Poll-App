package services

import (
	"context"

	"github.com/teodora-lazarevic/Poll-App/ent"
	"github.com/teodora-lazarevic/Poll-App/ent/poll"
	"github.com/teodora-lazarevic/Poll-App/ent/polloption"
	"github.com/teodora-lazarevic/Poll-App/ent/user"
	"github.com/teodora-lazarevic/Poll-App/ent/vote"
)

// var (
// 	ErrAlreadyVoted    = errors.New("User has already voted for this poll")
// 	ErrOptionNotInPoll = errors.New("Option does not belong to this poll")
// )

type PollResult struct {
	OptionID   int    `json:"option_id"`
	OptionName string `json:"option_name"`
	VoteCount  int    `json:"vote_count"`
}

type VoterResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type VoteService struct {
	DB *ent.Client
}

func NewVoteService(db *ent.Client) *VoteService {
	return &VoteService{DB: db}
}

func (s *VoteService) CastVote(ctx context.Context, userId, pollId, optionId int) error {
	voter, err := s.DB.User.Get(ctx, userId)
	if ent.IsNotFound(err) {
		return ErrUserNotFound
	}

	fetchedPoll, err := s.DB.Poll.Get(ctx, pollId)
	if ent.IsNotFound(err) {
		return ErrPollNotFound
	}

	option, err := s.DB.PollOption.Query().
		Where(polloption.ID(optionId), polloption.HasPollWith(poll.ID(pollId))).
		First(ctx)

	if ent.IsNotFound(err) {
		return ErrOptionNotInPoll
	} else if err != nil {
		return err
	}

	tx, err := s.DB.Tx(ctx)
	if err != nil {
		return err
	}

	hasVoted, err := s.DB.Vote.Query().
		Where(vote.HasUserWith(user.IDEQ(userId)), vote.HasPollWith(poll.IDEQ(pollId))).
		Exist(ctx)

	if err != nil {
		return err
	}
	if hasVoted {
		tx.Rollback()
		return ErrAlreadyVoted
	}

	_, err = tx.Vote.Create().
		SetUser(voter).
		SetOption(option).
		SetPoll(fetchedPoll).
		Save(ctx)

	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *VoteService) GetPollResults(ctx context.Context, pollId int) ([]PollResult, error) {
	p, err := s.DB.Poll.Query().
		Where(poll.IDEQ(pollId)).
		WithOptions(func(q *ent.PollOptionQuery) {
			q.WithVotes()
		}).
		WithCreator().
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, ErrPollNotFound
	} else if err != nil {
		return nil, err
	}

	var results []PollResult
	for _, option := range p.Edges.Options {
		results = append(results, PollResult{
			OptionID:   option.ID,
			OptionName: option.Text,
			VoteCount:  len(option.Edges.Votes),
		})
	}

	return results, nil
}

func (s *VoteService) GetVotersForOption(ctx context.Context, userId, pollId, optionId int) ([]VoterResponse, error) {
	option, err := s.DB.PollOption.Query().
		Where(polloption.IDEQ(optionId), polloption.HasPollWith(poll.IDEQ(pollId))).
		WithPoll(func(q *ent.PollQuery) {
			q.WithCreator()
		}).
		Only(ctx)

	if ent.IsNotFound(err) {
		return nil, ErrOptionNotInPoll
	} else if err != nil {
		return nil, err
	}

	if option.Edges.Poll.Edges.Creator.ID != userId {
		return nil, ErrUnauthorized
	}

	votes, err := option.QueryVotes().WithUser().All(ctx)
	if err != nil {
		return nil, err
	}

	var voters []VoterResponse
	for _, v := range votes {
		voters = append(voters, VoterResponse{
			ID:       v.Edges.User.ID,
			Username: v.Edges.User.Username,
			Email:    v.Edges.User.Email,
		})
	}

	return voters, nil
}
