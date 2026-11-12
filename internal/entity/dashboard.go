package entity

// Settings live in a key/value table so the next project can add switches
// without a migration.
const (
	SettingRegistrationApproval = "registration_approval"
)

type Setting struct {
	Key   string
	Value string
}

func RegistrationApprovalEnabled(value string) bool { return value == "true" }

// DashboardData is the admin overview payload.
type DashboardData struct {
	TotalUsers   int
	ActiveUsers  int
	PendingUsers int
	RolesCount   int
	RoleData     []RoleDatum
	SignupData   []SignupDatum
	RecentUsers  []User
}

type RoleDatum struct {
	Name  string
	Label string
	Users int
}

type SignupDatum struct {
	Date  string
	Users int
}
