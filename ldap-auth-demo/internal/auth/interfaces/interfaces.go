package interfaces

// UserInfo interface for token issuance
type UserInfo interface {
	GetUsername() string
	GetEmail() string
	GetFullName() string
	GetGroups() []string
	GetDepartment() string
}

// TokenIssuer interface for token issuance - allows different implementations
type TokenIssuer interface {
	IssueToken(user UserInfo) (string, error)
	ValidateToken(token string) (string, error)
}