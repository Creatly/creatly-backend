package domain

type PaginationQuery struct {
	Skip  int64 `form:"skip"`
	Limit int64 `form:"limit"`
}

type SearchQuery struct {
	Search string `form:"search"`
}

type StudentFiltersQuery struct {
	RegisterDateFrom  string `form:"registerDateFrom"`
	RegisterDateTo    string `form:"registerDateTo"`
	LastVisitDateFrom string `form:"lastVisitDateFrom"`
	LastVisitDateTo   string `form:"lastVisitDateTo"`
	Verified          *bool  `form:"verified"`
}

type GetStudentsQuery struct {
	PaginationQuery
	SearchQuery
	StudentFiltersQuery
}

func (p PaginationQuery) GetSkip() *int64 {
	if p.Skip == 0 {
		return nil
	}

	return &p.Skip
}

func (p PaginationQuery) GetLimit() *int64 {
	if p.Limit == 0 {
		return nil
	}

	return &p.Limit
}
