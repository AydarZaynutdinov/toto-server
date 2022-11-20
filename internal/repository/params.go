package repository

type QueryParameters struct {
	Packages    []string
	CountryCode []string
	Percentile  int
	Limit       int
	Offset      int
}
