package travel

type Repository interface {
	Create(travel *Travel) error
	FindByID(id uint) (*Travel, error)
	FindByUserID(userID uint) ([]*Travel, error)
	Update(travel *Travel) error
	FindPublicTravels() ([]*PublicTravel, error)
	FindPublicByID(id uint) (*PublicTravel, error)
	FindCopySource(sourceTravelID uint) (*CopySource, error)
	Delete(id uint) error
	CopyTravel(originalTravel *Travel, newUserID uint) (*Travel, error)
}
