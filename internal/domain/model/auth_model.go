package model

type LoginResponse struct {
	Token     string `json:"token"`
	ExpiresIn string `json:"expires_in"`
}

type LoginRequest struct {
	EmployeeName string `json:"employee_name"`
	Password     string `json:"password"`
}

type LogoutRequest struct {
	Token string `json:"token"`
}
