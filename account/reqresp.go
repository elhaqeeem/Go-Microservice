package account

type (
	CreateUserRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	CreateUserResponce struct {
		Ok string `json:"ok"`
	}

	GetUserRequest struct {
		Email string `json:"email"`
	}
	GetUserResponce struct {
		Id string `json:"id"`
	}
)
