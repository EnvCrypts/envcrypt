package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/vijayvenkatj/envcrypt/internal/config"
	"github.com/vijayvenkatj/envcrypt/internal/helpers"
)

type OIDCClaims struct {
	Subject    string
	Repository string
	Ref        string
	Issuer     string
	RawClaims  map[string]any
}

type OIDCVerifier interface {
	VerifyToken(ctx context.Context, rawIDToken string) (*OIDCClaims, error)
}
type GithubOIDCVerifier struct {
	verifier *oidc.IDTokenVerifier
}

func NewGithubOIDCVerifier(ctx context.Context, audience string) (OIDCVerifier, error) {
	provider, err := oidc.NewProvider(ctx, "https://token.actions.githubusercontent.com")
	if err != nil {
		return nil, fmt.Errorf("failed to init github provider: %w", err)
	}

	v := provider.Verifier(&oidc.Config{
		ClientID: audience,
	})

	return &GithubOIDCVerifier{verifier: v}, nil
}

func (g *GithubOIDCVerifier) VerifyToken(
	ctx context.Context,
	rawToken string,
) (*OIDCClaims, error) {

	idToken, err := g.verifier.Verify(ctx, rawToken)
	if err != nil {
		return nil, fmt.Errorf("invalid oidc token: %w", err)
	}

	var ghClaims struct {
		Sub        string `json:"sub"`
		Repository string `json:"repository"`
		Ref        string `json:"ref"`
		Iss        string `json:"iss"`
	}
	if err := idToken.Claims(&ghClaims); err != nil {
		return nil, fmt.Errorf("bad claims: %w", err)
	}

	var raw map[string]any
	if err := idToken.Claims(&raw); err != nil {
		return nil, fmt.Errorf("failed to decode raw claims: %w", err)
	}

	return &OIDCClaims{
		Subject:    ghClaims.Sub,
		Repository: ghClaims.Repository,
		Ref:        ghClaims.Ref,
		Issuer:     ghClaims.Iss,
		RawClaims:  raw,
	}, nil
}

func (handler *Handler) GitHubOIDCLogin(w http.ResponseWriter, r *http.Request) {

	var req config.GithubOIDCLoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helpers.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	defer r.Body.Close()

	claims, err := handler.OIDC.VerifyToken(r.Context(), req.IDToken)
	if err != nil {
		helpers.WriteError(w, http.StatusUnauthorized, "invalid OIDC token")
		return
	}

	repoPrincipal := claims.Subject
	if repoPrincipal == "" {
		helpers.WriteError(w, http.StatusUnauthorized, "missing repo identity")
		return
	}

	sessionID, projectID, err := handler.Services.SessionService.Create(
		r.Context(),
		repoPrincipal,
	)
	if err != nil {
		helpers.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := config.GithubOIDCLoginResponse{
		SessionID: *sessionID,
		ProjectID: *projectID,
	}

	helpers.WriteResponse(w, http.StatusOK, resp)
}
