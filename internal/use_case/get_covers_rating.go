package use_case

type GetCoversRatingInput struct {
	AudiosIDs []int
}

type GetCoversRatingOutput struct {
	CoversIDsRating []int
}

func (u UseCase) GetCoversRating(input GetCoversRatingInput) (output GetCoversRatingOutput, err error) {
	//TODO implement me
	panic("implement me")
}
