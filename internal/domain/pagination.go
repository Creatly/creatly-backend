package domain

type PaginationQuery struct {
	Skip  int64 `form:"skip"`
	Limit int64 `form:"limit"`
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
