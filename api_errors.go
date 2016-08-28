package main

// APIErrorCause - причина ошибки API
type APIErrorCause struct {
    Target          string                  `json:"target"`
    Сause           string                  `json:"cause"`
}

// APIError - ошибка API
type APIError struct {
    Code            int                     `json:"code"`
    Сauses          []APIErrorCause         `json:"causes"`
    Description     string                  `json:"description"`
}

// ErrorCodeToString - описание ошибки 
func ErrorCodeToString(code int) string {
    switch code {
        case 28:  return "Not find all request params"
        case 29:  return "Request params not validate"
        case 41:  return "Error create data"
        case 44:  return "There are no results"
        case 400: return "Bad request"
        case 401: return "Not authorized"
        case 409: return "Conflict"        
        case 404: return "Not Found"
        case 500: return "Internal Server Error"
        case 503: return "Service Unavailable"
        default:  return "Unknown"
    }
}

// NewSimpleAPIError - создание простых ошибок
func NewSimpleAPIError(code int) APIError {    
    return APIError{ Code: code,  Сauses: make([]APIErrorCause,0), Description: ErrorCodeToString(code)}
}

// NewAPIError - создание ошибок с причинами
func NewAPIError(code int, causes []APIErrorCause) APIError {    
    return APIError{ Code: code,  Сauses: causes, Description: ErrorCodeToString(code)}
}