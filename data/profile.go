package data

// UserProfile provides a basic user definition
type UserProfile struct {
	User                     string   `json:"user"`
	Password                 string   `json:"password"`
	UserType                 string   `json:"userType"`
	Age                      string   `json:"age"`
	PostalCode               string   `json:"postalCode"`
	SelectedAgencies         []string `json:"selectedAgencies"`
	SelectedProjects         []string `json:"selectedProjects"`
	NotificationEmail        string   `json:"notificationEmail"`
	EnableEmailNotifications bool     `json:"enableEmailNotifications"`
	NotificationSMS          string   `json:"notificationSMS"`
	EnableSMSNotifications   bool     `json:"enableSMSNotifications"`
}
