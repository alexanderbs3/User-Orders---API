package errors

import "net/http"

// AppError é o tipo de erro padronizado da aplicação.
// Em vez de retornar erros genéricos, retornamos sempre um AppError
// que carrega o código HTTP correto e uma mensagem clara.
type AppError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
}

func (e *AppError) Error() string {
    return e.Message
}

// Construtores semânticos — equivalente às exceções nomeadas do Spring
func NotFound(resource string) *AppError {
    return &AppError{
        Code:    http.StatusNotFound,
        Message: resource + " not found",
    }
}

func BadRequest(message string) *AppError {
    return &AppError{
        Code:    http.StatusBadRequest,
        Message: message,
    }
}

func Conflict(message string) *AppError {
    return &AppError{
        Code:    http.StatusConflict,
        Message: message,
    }
}

func Internal(message string) *AppError {
    return &AppError{
        Code:    http.StatusInternalServerError,
        Message: message,
    }
}