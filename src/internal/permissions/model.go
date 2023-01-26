package permissions

type Permissions struct {
	Id     int64  `json:"id"`
	Status string `json:"status"`
}

const (
	BannedPermission = "banned"
	AdminPermission  = "admin"
	ReadPermission   = "read-only"
)

var ArrayOfPermissions = []string{BannedPermission, AdminPermission, ReadPermission}
