package dtos

type PeerDto struct {
	ID             int    `form:"id"`
	Name           string `form:"name"`
	Domain         string `form:"domain"`
	Port           int    `form:"port"`
	OrderBy        string `form:"orderBy"`
	OrderDirection string `form:"orderDirection"`
}
