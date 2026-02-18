package application

import (
	"context"
	"errors"
	"time"

	"github.com/sample-provider/buy-credit-api/internal/domain/repository"
	"github.com/sample-provider/buy-credit-api/internal/infrastructure/auth"
)

type AuthUseCase struct {
	partnerRepo repository.PartnerRepository
	jwtService  *auth.JWTService
}

type AuthRequest struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
	TokenType   string `json:"tokenType"`
	ExpiresIn   int    `json:"expiresIn"`
}

func NewAuthUseCase(partnerRepo repository.PartnerRepository, jwtService *auth.JWTService) *AuthUseCase {
	return &AuthUseCase{
		partnerRepo: partnerRepo,
		jwtService:  jwtService,
	}
}

func (uc *AuthUseCase) Authenticate(ctx context.Context, req AuthRequest) (*AuthResponse, error) {
	partner, err := uc.partnerRepo.FindByClientID(ctx, req.ClientID)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if partner.ClientSecret != req.ClientSecret {
		return nil, errors.New("invalid credentials")
	}

	expiresIn := 1 * time.Hour
	token, err := uc.jwtService.GenerateToken(partner.ID, partner.ClientID, expiresIn)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int(expiresIn.Seconds()),
	}, nil
}
