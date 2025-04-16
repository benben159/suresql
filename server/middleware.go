package server

import (
	"net/http"

	"github.com/medatechnology/suresql"

	"github.com/medatechnology/simplehttp"
)

const (
	API_KEY_STRING     = "API_KEY"
	CLIENT_ID_STRING   = "CLIENT_ID"
	TOKEN_TABLE_STRING = "token"
)

// AuthMiddleware verifies API key and client ID from request headers
func MiddlewareAPIKeyHeader() simplehttp.MedaMiddleware {
	return simplehttp.WithName("APIKeyClientID", APIKeyClientIDHeader())
}

func APIKeyClientIDHeader() simplehttp.MedaMiddlewareFunc {
	return func(next simplehttp.MedaHandlerFunc) simplehttp.MedaHandlerFunc {
		return func(ctx simplehttp.MedaContext) error {

			// Get headers, make sure you already have MiddlewareHeaderParser before invoking this middleware
			state := NewMiddlewareState(ctx, "APIKey/ClientID")

			// Validate API key
			apiKey := ctx.GetHeader(API_KEY_STRING)
			if apiKey == "" {
				return state.SetError("API key required", nil, http.StatusUnauthorized).LogAndResponse("API key not provided", nil, true)
			}

			if suresql.CurrentNode.InternalConfig.APIKey != apiKey {
				return state.SetError("Invalid API key", nil, http.StatusUnauthorized).LogAndResponse("Invalid API key", nil, true)
			}

			// Validate Client ID
			clientID := ctx.GetHeader(CLIENT_ID_STRING)
			if clientID == "" {
				return state.SetError("Client ID required", nil, http.StatusUnauthorized).LogAndResponse("Client ID not provided", nil, true)
			}

			if suresql.CurrentNode.InternalConfig.ClientID != clientID {
				return state.SetError("Invalid Client ID", nil, http.StatusUnauthorized).LogAndResponse("Invalid Client ID", nil, true)
			}

			// Continue to next handler
			return next(ctx)
		}
	}
}

// TokenValidationMiddleware verifies that a valid token is present
func MiddlwareTokenCheck() simplehttp.MedaMiddleware {
	return simplehttp.WithName("token checker", TokenValidationFromTTL())
}

func TokenValidationFromTTL() simplehttp.MedaMiddlewareFunc {
	return func(next simplehttp.MedaHandlerFunc) simplehttp.MedaHandlerFunc {
		return func(ctx simplehttp.MedaContext) error {
			// Get headers, make sure you already have MiddlewareHeaderParser before invoking this middleware
			state := NewMiddlewareState(ctx, "token")

			// Get token from Authorization header
			header := state.Header
			token := header.Authorization.Token
			if token == "" {
				// NOTE: do we need this? Try from query parameters if it's not passed from the header??
				token = ctx.GetQueryParam(TOKEN_TABLE_STRING)
			}

			// Double check if token is empty
			if token == "" {
				return state.SetError("Authentication token required", nil, http.StatusUnauthorized).LogAndResponse("no token", nil, true)
			}

			// Validate token
			tok, valid := TokenStore.TokenExist(token)
			if !valid {
				return state.SetError("Invalid or expired token", nil, http.StatusUnauthorized).LogAndResponse("no token", nil, true)
			}

			// Set username in context for use in handlers
			ctx.Set(TOKEN_TABLE_STRING, tok)
			// Continue to next handler
			return next(ctx)
		}
	}
}
