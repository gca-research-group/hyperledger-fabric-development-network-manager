package dtos

type UserDto struct {
	ID             int    `form:"id"`
	Name           string `form:"name"`
	Email          string `form:"email"`
	Password       string `form:"password"`
	IsSuper        bool   `form:"isSuper"`
	OrderBy        string `form:"orderBy"`
	OrderDirection string `form:"orderDirection"`
}
