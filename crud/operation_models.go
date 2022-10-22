package crud

type UpdateAdminInfo struct {
	Id         uint64
	Surname    *string
	Name       *string
	Patronymic *string
	Email      *string
	Phone      *string
	BarId      *uint64
}

type UpdateOwnerInfo struct {
	Id         uint64
	Surname    *string
	Name       *string
	Patronymic *string
	Email      *string
	Phone      *string
}

type UpdateBar struct {
	Id          uint64
	Email       *string
	Address     *string
	Name        *string
	Description *string
	LogoUrl     *string
	Phone       *string
	IsVisible   *bool
	WorkHours   []WorkHours
}
