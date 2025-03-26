package dtos

type PeerDto struct {
	ID             int    `form:"id"`
	Name           string `form:"name"`
	Domain         string `form:"domain"`
	OrderBy        string `form:"orderBy"`
	OrderDirection string `form:"orderDirection"`
}
