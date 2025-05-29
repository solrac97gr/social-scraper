package database

type InfluencerRepository interface {
	SaveInfluencerAnalysis(influencer *InfluencerAnalysis) error
	GetInfluencerAnalysisByLink(link string) (*InfluencerAnalysis, error)
}
