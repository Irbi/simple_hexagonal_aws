package interfaces

import (
	"fmt"
	"test_aws/domain"
)

type LocalDbHandler interface {
	Execute(statement string, args ...interface{}) error
	Query(statement string) Row
}

type Row interface {
	Scan(dest ...interface{}) error
	Next() bool
	Close()
}

type LocalFileRepo struct {
	dbHandler  LocalDbHandler
}

type LocalStatRepo struct {
	dbHandler  LocalDbHandler
}

func NewLocalFileRepo(dbHandler LocalDbHandler) *LocalFileRepo {
	repo := new(LocalFileRepo)
	repo.dbHandler = dbHandler
	return repo
}

func (repo *LocalFileRepo) FindById(id string) (file domain.File, err error) {
	q := fmt.Sprintf("SELECT id, name, checksum, version FROM files WHERE id = '%s' LIMIT 1", id)
	file, err = repo.find(q)

	file.ID = id

	return file, err
}

func (repo *LocalFileRepo) FindByName(name string) (file domain.File, err error) {
	q := fmt.Sprintf("SELECT id, name, checksum, version FROM files WHERE name like '%s' LIMIT 1", name)
	file, err = repo.find(q)

	file.Name = name

	return file, err
}

func (repo *LocalFileRepo) Store(file domain.File) (result bool, err error) {
	result = true

	if file.ID != "" {
		q := fmt.Sprintf("UPDATE files SET checksum = ?, version = ? WHERE id = ?")
		err = repo.dbHandler.Execute(q, file.Checksum, file.Version, file.ID)
	} else {
		file.GenID()
		q := fmt.Sprintf("INSERT INTO files (id, name, checksum, version) VALUES (?, ?, ?, ?)")
		err = repo.dbHandler.Execute(q, file.ID, file.Name, file.Checksum, file.Version)
	}
	if err != nil {
		result = false
	}

	return result, err
}

func (repo *LocalFileRepo) find(query string) (file domain.File, err error) {
	row := repo.dbHandler.Query(query)
	defer row.Close()

	var (
		id string
		name string
		checksum string
		version int
	)
	file = domain.File{}

	isNext := row.Next()
	if isNext {
		err = row.Scan(&id, &name, &checksum, &version);
		if err != nil {
			fmt.Printf("Error on local db: %v", err)
		}
	}

	file.ID = id
	file.Name = name
	file.Checksum = checksum
	file.Version = version

	return file, err
}

func NewLocalStatRepo(dbHandler LocalDbHandler) *LocalStatRepo {
	repo := new(LocalStatRepo)
	repo.dbHandler = dbHandler
	return repo
}

func (repo *LocalStatRepo) AddValid() error {
	q := "UPDATE stat SET valid = valid+1"
	err := repo.dbHandler.Execute(q)
	return err
}

func (repo *LocalStatRepo) AddOut() error {
	q := "UPDATE stat SET outdated = outdated+1"
	err := repo.dbHandler.Execute(q)
	return err
}

func (repo *LocalStatRepo) Stat() (stat domain.Stat, err error) {
	q := "SELECT valid, outdated FROM stat"
	row := repo.dbHandler.Query(q)
	defer row.Close()

	var (
		valid int
		outdated int
	)
	stat = domain.Stat{}

	isNext := row.Next()
	if isNext {
		err = row.Scan(&valid, &outdated);
		if err != nil {
			fmt.Printf("Error on local db: %v", err)
		}
	}

	stat.Valid = valid
	stat.Outdated = outdated

	return stat, err
}