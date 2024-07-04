package handlers

// createUserRequestBody is the type of the "CreateUser"
// endpoint request body.
type createUserRequestBody struct {
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
}
