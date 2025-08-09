// mockTokenIssuer for testing.
type mockTokenIssuer struct {
	issueErr    error
	validateErr error
	token       string
	username    string
}

func (m *mockTokenIssuer) IssueToken(_ UserInfo) (string, error) {
	if m.issueErr != nil {
		return "", m.issueErr
	}
	return m.token, nil
}

func (m *mockTokenIssuer) ValidateToken(_ string) (string, error) {
	if m.validateErr != nil {
		return "", m.validateErr
	}
	return m.username, nil
}