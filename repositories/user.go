package repositories

const SqliteDBDsn = "user.db?cache=shared"

type UserRepository struct {
}

func NewUser() *UserRepository {
	return &UserRepository{}
}
