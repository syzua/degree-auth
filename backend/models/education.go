package models

type Education struct {
	CertNo           string   `json:"certNo"`
	Name             string   `json:"name"`
	StudentID        string   `json:"studentId"`
	School           string   `json:"school"`
	Major            string   `json:"major"`
	Degree           string   `json:"degree"`
	GraduationDate   string   `json:"graduationDate"`
	IssueDate        string   `json:"issueDate"`
	Status           string   `json:"status"`
	AuthorizedViewers []string `json:"authorizedViewers"`
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
