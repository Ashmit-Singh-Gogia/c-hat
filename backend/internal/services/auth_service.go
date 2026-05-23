package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"os"
	"time"

	"github.com/ashmit-singh-gogia/c-hat/internal/models"
	"github.com/ashmit-singh-gogia/c-hat/internal/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService struct {
	userRepo *repositories.UserRepository
}

func NewAuthService(userRepo *repositories.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) FindOrCreateUser(dto *UserDTO) (*models.User, error) {
	user, err := s.userRepo.FindByProvider(dto.Provider, dto.ProviderID)
	if err == nil {
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// user doesn't exist → create a new one
	user = models.User{
		Provider:   dto.Provider,
		ProviderID: dto.ProviderID,
		Username:   dto.Username,
		Email:      dto.Email,
		Avatar:     dto.Avatar,
		IsVerified: true, // since we trust the provider, we can mark the user as verified
	}
	if err := s.userRepo.CreateUser(&user); err != nil {
		return nil, errors.New("failed to create new user")
	}
	return &user, nil
}

func (s *AuthService) ProcessUserLogin(dto *UserDTO) (*models.User, error) {
	return s.FindOrCreateUser(dto)
}

func (s *AuthService) RegisterLocalUser(username, email, password string) error {
	// 1. Input Validation
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	cleanedEmail, err := ValidateAndCleanEmail(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	if username == "" || len(username) < 3 || len(username) > 20 {
		return errors.New("username must be between 3 and 20 characters")
	}

	// 2. Duplicate Email Check
	_, err = s.userRepo.GetUserByEmail(cleanedEmail)
	if err == nil {
		return errors.New("user with this email already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 3. Password Hashing
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return errors.New("failed to hash password")
	}
	newUser := &models.User{
		Username:   username,
		Email:      cleanedEmail,
		Password:   hashedPassword,
		Provider:   "local",
		ProviderID: cleanedEmail,
		IsVerified: false, //	User is not verified until they click the link in the email
	}

	if err := s.userRepo.DB.Create(newUser).Error; err != nil {
		return err
	}

	// Email Verification
	err = s.ResendVerificationEmail(newUser.Email)
	if err != nil {
		log.Printf("CRITICAL: User created, but failed to send verification email to %s: %v", newUser.Email, err)
	} else {
		log.Printf("Verification email sent to %s", newUser.Email)
	}
	return nil
}

func (s *AuthService) LoginLocalUser(email, password string) (*models.User, error) {
	cleanedEmail, err := ValidateAndCleanEmail(email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	user, err := s.userRepo.GetUserByEmail(cleanedEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user with this email does not exist")
		}
		return nil, err
	}

	if !user.IsVerified {
		return nil, errors.New("unverified_account")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("incorrect password")
	}
	// Pass a jwt token or session cookie here in here
	return &user, nil
}

func (s *AuthService) VerifyEmail(token string) error {
	user, err := s.userRepo.GetUserByVerificationToken(token)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired verification token")
		}
		return err
	}

	if user.TokenExpiresAt == nil || time.Now().After(*user.TokenExpiresAt) {
		return errors.New("verification token has expired")
	}

	user.IsVerified = true
	user.VerificationToken = ""
	user.TokenExpiresAt = nil

	if err := s.userRepo.DB.Save(&user).Error; err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ResendVerificationEmail(email string) error {
	cleanedEmail, err := ValidateAndCleanEmail(email)
	if err != nil {
		return errors.New("invalid email format")
	}

	user, err := s.userRepo.GetUserByEmail(cleanedEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user with this email does not exist")
		}
		return err
	}

	if user.IsVerified {
		return errors.New("email is already verified")
	}

	// Generate a new verification token
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return errors.New("failed to generate secure token")
	}
	newToken := hex.EncodeToString(tokenBytes)
	tokenExpiry := time.Now().Add(24 * time.Hour)

	user.VerificationToken = newToken
	user.TokenExpiresAt = &tokenExpiry

	if err := s.userRepo.DB.Save(&user).Error; err != nil {
		return err
	}

	from := os.Getenv("MAIL")
	password := os.Getenv("PASSWD")
	toList := []string{cleanedEmail}
	host := "smtp.gmail.com"
	port := "587"
	msg := fmt.Sprintf("Subject: Email Verification\n\nPlease verify your email by clicking the following link: http://localhost:5173/verify-email?token=%s", newToken)
	body := []byte(msg)
	auth := smtp.PlainAuth("", from, password, host)
	err = smtp.SendMail(host+":"+port, auth, from, toList, body)
	if err != nil {
		fmt.Printf("Failed to send verification email to %s: %v\n", cleanedEmail, err)
		return fmt.Errorf("failed to send verification email: %v", err)
	}
	return nil
}

func HashPassword(password string) (string, error) { // bcrypt is a popular password hashing algorithm that is designed to be slow and computationally expensive, making it resistant to brute-force attacks.
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func ValidateAndCleanEmail(email string) (string, error) {
	addr, err := mail.ParseAddress(email)
	if err != nil {
		return "", err
	}
	return addr.Address, nil // Return the cleaned email address
}
