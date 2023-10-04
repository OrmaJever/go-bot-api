package SelectUser

type User struct {
	tableName struct{} `pg:"alias:u"`
	Id        int64    `json:"id"`
	TgId      int64    `json:"tg_id"`
	ChatId    int64    `json:"chat_id"`
	FirstName string   `json:"first_name"`
	Username  string   `json:"username"`
	RealName  string   `json:"real_name"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

type Customize struct {
	tableName struct{} `pg:"alias:c"`
	Id        int32    `json:"id"`
	UserId    int32    `json:"user_id"`
	Image     string   `json:"image"`
	Text      string   `json:"text"`
	User      User     `json:"user" pg:"rel:has-one"`
}

type SelectedUser struct {
	tableName   struct{}  `pg:"alias:su"`
	Id          int32     `json:"id"`
	TgId        int64     `json:"tg_id"`
	Type        int8      `json:"type"`
	ChatId      int64     `json:"chat_id"`
	CustomizeId int32     `json:"customize_id"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	User        User      `json:"user" pg:"rel:has-one,fk:tg_id,join_fk:tg_id"`
	Customize   Customize `json:"customize" pg:"rel:has-one"`
}
