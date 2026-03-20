package repositories

type RepositoryHealth struct {
	access *RepositoryAccess
}

func (r RepositoryHealth) GetHealth() string {
	return "all okay"
}

type RepositoryHealthMethods interface {
	GetHealth() string
}

func NewRepositoryHealth(access *RepositoryAccess) RepositoryHealthMethods {
	return &RepositoryHealth{
		access: access,
	}
}
