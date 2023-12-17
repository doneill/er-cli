package data

// ----------------------------------------------
// table stucts
// ----------------------------------------------

type Accounts_User struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	RemoteID  string `json:"remote_id" gorm:"column:remote_id"`
	Username  string `json:"username" gorm:"column:username"`
	FirstName string `json:"first_name" gorm:"column:first_name"`
	LastName  string `json:"last_name" gorm:"column:last_name"`
	Email     string `json:"email" gorm:"column:email"`
}

// ----------------------------------------------
// table configurations
// ----------------------------------------------

type Tabler interface {
	TableName() string
}

func (Accounts_User) TableName() string {
	return "accounts_user"
}
