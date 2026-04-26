package dto

// CreateUserRequest é o payload esperado no body do POST /users
// A tag `binding:"required"` é do Gin — valida automaticamente
// A tag `binding:"email"` valida o formato de e-mail
type CreateUserRequest struct {
    Name  string `json:"name" binding:"required,min=2,max=100"`
    Email string `json:"email" binding:"required,email"`
}

// UpdateUserRequest permite atualização parcial (campos opcionais)
// O uso de ponteiros (*string) permite distinguir "campo não enviado" de "campo vazio"
type UpdateUserRequest struct {
    Name  *string `json:"name" binding:"omitempty,min=2,max=100"`
    Email *string `json:"email" binding:"omitempty,email"`
}

// UserResponse é o que a API retorna — sem campos sensíveis
type UserResponse struct {
    ID        uint   `json:"id"`
    Name      string `json:"name"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
}