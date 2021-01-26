package entities

// Boundaries represents the boundary to which a user is constrained.
type Boundaries struct {
	Model
	Type                     string `json:"type"`
	StatsRetention           int64  `json:"stats_retention"`
	SubscribersLimit         int64  `json:"subscribers_limit"`
	CampaignsLimit           int64  `json:"campaigns_limit"`
	TemplatesLimit           int64  `json:"templates_limit"`
	GroupsLimit              int64  `json:"groups_limit"`
	TeamMembersLimit         int64  `json:"team_members_limit"`
	ScheduleCampaignsEnabled bool   `json:"schedule_campaigns_enabled"`
	OrganizationsEnabled     bool   `json:"organizations_enabled"`
	SAMLEnabled              bool   `json:"saml_enabled"`
}
