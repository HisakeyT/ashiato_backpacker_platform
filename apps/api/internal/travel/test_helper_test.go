package travel

type MockRepository struct {
	CreateFunc            func(travel *Travel) error
	FindByIDFunc          func(id uint) (*Travel, error)
	FindByUserIDFunc      func(userID uint) ([]*Travel, error)
	UpdateFunc            func(travel *Travel) error
	FindPublicTravelsFunc func() ([]*PublicTravel, error)
	FindPublicByIDFunc    func(id uint) (*PublicTravel, error)
	FindCopySourceFunc    func(sourceTravelID uint) (*CopySource, error)
	DeleteFunc            func(id uint) error
	CopyTravelFunc        func(originalTravel *Travel, newUserID uint) (*Travel, error)
}

func (m *MockRepository) Create(travel *Travel) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(travel)
	}
	return nil
}

func (m *MockRepository) FindByID(id uint) (*Travel, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(id)
	}

	return nil, nil
}

func (m *MockRepository) FindByUserID(userID uint) ([]*Travel, error) {
	if m.FindByUserIDFunc != nil {
		return m.FindByUserIDFunc(userID)
	}

	return nil, nil
}

func (m *MockRepository) Update(travel *Travel) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(travel)
	}

	return nil
}

func (m *MockRepository) FindPublicTravels() ([]*PublicTravel, error) {
	if m.FindPublicTravelsFunc != nil {
		return m.FindPublicTravelsFunc()
	}

	return nil, nil
}

func (m *MockRepository) FindPublicByID(id uint) (*PublicTravel, error) {
	if m.FindPublicByIDFunc != nil {
		return m.FindPublicByIDFunc(id)
	}

	return nil, nil
}

func (m *MockRepository) FindCopySource(sourceTravelID uint) (*CopySource, error) {
	if m.FindCopySourceFunc != nil {
		m.FindCopySourceFunc(sourceTravelID)
	}

	return nil, nil
}

func (m *MockRepository) Delete(id uint) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(id)
	}
	return nil
}

func (m *MockRepository) CopyTravel(originalTravel *Travel, newUserID uint) (*Travel, error) {
	if m.CopyTravelFunc != nil {
		return m.CopyTravelFunc(originalTravel, newUserID)
	}

	return nil, nil
}
