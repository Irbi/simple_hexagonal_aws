package usecases

import "test_aws/domain"

type Stat struct {
	Total int
	Valid int
	Outdated int
}

type StatInteractor struct {
	Repository domain.StatRepository
}

func (itr *StatInteractor) Stat() (stat Stat, err error) {
	stat = Stat{}

	stats, err := itr.Repository.Stat()
	if err != nil {
		return stat, err
	}

	stat.Outdated = stats.Outdated
	stat.Valid = stats.Valid
	stat.Total = stat.Valid + stat.Outdated

	return stat, err
}

func (itr *StatInteractor) AddValid() error {
	err := itr.Repository.AddValid()
	return err
}

func (itr *StatInteractor) AddOut() error {
	err := itr.Repository.AddOut()
	return err
}