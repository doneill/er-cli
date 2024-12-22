package data

// ----------------------------------------------
// tables
// ----------------------------------------------

type Accounts_User struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	RemoteID string `json:"remote_id" gorm:"column:remote_id"`
	Username string `json:"username" gorm:"column:username"`
}

type Event struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	RemoteID  string `json:"remote_id" gorm:"column:remote_id"`
	ProfileID int    `json:"profile_id" gorm:"column:profile_id"`
	Title     string `json:"title" gorm:"column:title"`
	CreatedAt string `json:"created_at" gorm:"created_at"`
	UpdatedAt int64  `json:"updated_at" gorm:"updated_at"`
}

type User_Profile struct {
	ID       int    `gorm:"primaryKey;autoIncrement"`
	RemoteID string `json:"remote_id" gorm:"column:remote_id"`
	Username string `json:"username" gorm:"column:username"`
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
