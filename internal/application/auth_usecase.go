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
	APIKey    string `json:"apiKey"`
	APISecret string `json:"apiSecret"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}

func NewAuthUseCase(partnerRepo repository.PartnerRepository, jwtService *auth.JWTService) *AuthUseCase {
	return &AuthUseCase{
		partnerRepo: partnerRepo,
		jwtService:  jwtService,
	}
}

func (uc *AuthUseCase) Authenticate(ctx context.Context, req AuthRequest) (*AuthResponse, error) {
	partner, err := uc.partnerRepo.FindByClientID(ctx, req.APIKey)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if partner.ClientSecret != req.APISecret {
		return nil, errors.New("invalid credentials")
	}

	expiresIn := 1 * time.Hour
	token, err := uc.jwtService.GenerateToken(partner.ID, partner.ClientID, expiresIn)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: token,
	}, nil
}
