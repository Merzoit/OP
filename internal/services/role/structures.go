package role

type Role struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Permissions string `json:"permissions"`
}
