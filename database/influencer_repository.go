package database

type AllInfluencerAnalysis struct {
	TotalCount int64                 `json:"total_count" bson:"total_count"`
	Analyses   []*InfluencerAnalysis `json:"analyses" bson:"analyses"`
	Pagination struct {
		Page  int64 `json:"page" bson:"page"`
		Limit int64 `json:"limit" bson:"limit"`
	} `json:"pagination" bson:"pagination"`
}

type InfluencerRepository interface {
	SaveInfluencerAnalysis(influencer *InfluencerAnalysis) error
	GetInfluencerAnalysisByLink(link string) (*InfluencerAnalysis, error)
	DeleteExpiredAnalyses() error
	GetAllInfluencerAnalyses(page int, limit int) (AllInfluencerAnalysis, error)
}

type UserRepository interface {
	SaveUser(user *User) error
	SaveUserToken(token *UserToken) error
	SaveUserProfile(profile *UserProfile) error
	UpdateUserProfile(userID string, profile *UserProfile) error
	UpdateUser(user *User) error
	GetUserByID(userID string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserTokenByUserID(userID string) (*UserToken, error)
	DeleteExpiredTokens() error
}
