package app

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/solrac97gr/telegram-followers-checker/database"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserID      = errors.New("invalid user ID")
	ErrInvalidUser        = errors.New("invalid user")
	ErrInvalidUserToken   = errors.New("invalid user token")
	ErrInvalidUserProfile = errors.New("invalid user profile")
)

type UserApp struct {
	Repository database.UserRepository
	JWTSecret  string
}

type Claims struct {
	UserID       string                `json:"user_id"`
	Email        string                `json:"email"`
	Role         database.Role         `json:"role"`
	Subscription database.Subscription `json:"subscription"`
	jwt.RegisteredClaims
}

func NewUserApp(repository database.UserRepository, jwtSecret string) *UserApp {
	if repository == nil {
		return nil
	}
	return &UserApp{
		Repository: repository,
		JWTSecret:  jwtSecret,
	}
}

func (u *UserApp) GetUserByID(userID string) (*database.User, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	return u.Repository.GetUserByID(userID)
}
func (u *UserApp) GetUserTokenByUserID(userID string) (*database.UserToken, error) {
	if userID == "" {
		return nil, ErrInvalidUserID
	}
	return u.Repository.GetUserTokenByUserID(userID)
}

func (u *UserApp) SaveUser(user string, email string, password string, confirmationPassword string) error {
	if user == "" || email == "" || password == "" || confirmationPassword == "" {
		return ErrInvalidUser
	}
	if !isValidPassword(password) {
		return errors.New("password must be at least 8 characters long and include upper/lowercase letters and numbers")
	}
	if password != confirmationPassword {
		return errors.New("passwords do not match")
	}
	if len(email) < 5 || !IsValidEmail(email) {
		return errors.New("invalid email format")
	}
	// Normalize the email to ensure consistency
	email, err := NormalizeEmail(email)
	if err != nil {
		return err
	}

	userObj := &database.User{
		ID:              user,
		Email:           email,
		Password:        HashPassword(password),
		Role:            database.UserRole,         // Default role is User
		Subscription:    database.FreeSubscription, // Default subscription is Free
		ProfileComplete: false,                     // Default profile completion status
		CreatedAt:       time.Now(),                // Set the current time as created at
		UpdatedAt:       time.Now(),                // Set the current time as updated at
	}

	if err := u.Repository.SaveUser(userObj); err != nil {
		return err
	}
	return nil
}
func (u *UserApp) SaveUserToken(userID string, token string) error {
	if userID == "" || token == "" {
		return ErrInvalidUserToken
	}
	if len(token) < 20 {
		return errors.New("token must be at least 20 characters long")
	}
	userToken := &database.UserToken{
		ID:        userID,
		Token:     token,
		IsValid:   true,                           // Assuming the token is valid when saved
		CreatedAt: time.Now(),                     // You can set the current time here
		ExpiresAt: time.Now().Add(72 * time.Hour), // Set the expiration time to 24 hours from now
	}
	return u.Repository.SaveUserToken(userToken)
}
func (u *UserApp) SaveUserProfile(userID string, firstName string, lastName string, phoneNumber string, address string, profilePic string) error {
	profile := &database.UserProfile{
		ID:          userID,
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phoneNumber,
		Address:     address,
		ProfilePic:  profilePic,
		CompanyName: "",         // CompanyName will be set later now is not possible to have company name
		CompanyID:   "",         // CompanyID will be set later now is not possible to have company ID
		CreatedAt:   time.Now(), // Set the current time as created at
		UpdatedAt:   time.Now(), // Set the current time as updated at
	}
	return u.Repository.SaveUserProfile(profile)
}

func (u *UserApp) UpdateUserProfile(userID string, profile *database.UserProfile) error {
	if userID == "" || profile == nil {
		return ErrInvalidUserProfile
	}
	return u.Repository.UpdateUserProfile(userID, profile)
}

func (u *UserApp) UpdateUser(user *database.User) error {
	if user == nil {
		return ErrInvalidUser
	}
	return u.Repository.UpdateUser(user)
}

func (u *UserApp) DeleteExpiredTokens() error {
	return u.Repository.DeleteExpiredTokens()
}

func (u *UserApp) AuthenticateUser(email string, password string) (*database.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password are required")
	}

	// Normalize the email to ensure consistency
	normalizedEmail, err := NormalizeEmail(email)
	if err != nil {
		return nil, errors.New("invalid email format")
	}

	// Get user by email
	user, err := u.Repository.GetUserByEmail(normalizedEmail)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Remove sensitive information before returning
	user.Password = ""

	return user, nil
}

func (u *UserApp) GenerateToken(user *database.User) (string, error) {
	if user == nil {
		return "", errors.New("user cannot be nil")
	}

	if u.JWTSecret == "" {
		return "", errors.New("JWT secret is not configured")
	}

	// Create the claims
	claims := &Claims{
		UserID:       user.ID,
		Email:        user.Email,
		Role:         user.Role,
		Subscription: user.Subscription,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   user.ID,
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token
	tokenString, err := token.SignedString([]byte(u.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func IsValidEmail(email string) bool {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash)
}

func NormalizeEmail(email string) (string, error) {
	emailInLowerCase := strings.ToLower(email)
	parts := strings.Split(emailInLowerCase, "@")

	if len(parts) != 2 {
		return "", errors.New("invalid email format")
	}

	localPart := parts[0]
	domainPart := parts[1]

	if plusIndex := strings.Index(localPart, "+"); plusIndex != -1 {
		localPart = localPart[:plusIndex]
	}

	switch domainPart {
	case "gmail.com", "googlemail.com":
		localPart = strings.ReplaceAll(localPart, ".", "")
	}

	return localPart + "@" + domainPart, nil
}

func isValidPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	hasUpper := false
	hasLower := false
	hasNumber := false
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= 'a' && char <= 'z' {
			hasLower = true
		} else if char >= '0' && char <= '9' {
			hasNumber = true
		}
	}
	return hasUpper && hasLower && hasNumber
}

// GenerateJWT generates a new JWT token for a given user ID
func (u *UserApp) GenerateJWT(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("invalid user ID")
	}

	// Define token expiration time
	expirationTime := time.Now().Add(72 * time.Hour) // Token valid for 72 hours

	// Create the JWT claims, which includes the user ID and expiration time
	claims := &jwt.RegisteredClaims{
		Issuer:    "your-app-name", // Set the issuer of the token
		Subject:   userID,          // Set the subject of the token (user ID)
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Create the JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with your secret key
	tokenString, err := token.SignedString([]byte("your-secret-key"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseJWT parses the JWT token and returns the user ID
func (u *UserApp) ParseJWT(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("invalid token")
	}

	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the token's signing method is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// Return the secret key for verification
		return []byte("your-secret-key"), nil
	})
	if err != nil {
		return "", err
	}

	// Extract the user ID from the token claims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || claims.Subject == "" {
		return "", errors.New("invalid token claims")
	}

	return claims.Subject, nil
}
