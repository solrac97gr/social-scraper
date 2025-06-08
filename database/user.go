package database

import "time"

type Role string

const (
	SuperAdminRole Role = "super_admin"
	AdminRole      Role = "admin"
	UserRole       Role = "user"
)

type Subscription string

const (
	FreeSubscription       Subscription = "free"
	PremiumSubscription    Subscription = "premium"
	EnterpriseSubscription Subscription = "enterprise"
)

type User struct {
	ID              string       `json:"id" bson:"_id"`
	Username        string       `json:"username" bson:"username"`
	Email           string       `json:"email" bson:"email"`
	Password        string       `json:"password" bson:"password"`
	Role            Role         `json:"role" bson:"role"`                           // Role of the user, e.g., super_admin, admin, user
	Subscription    Subscription `json:"subscription_type" bson:"subscription_type"` // Subscription type of the user, e.g., free, premium, enterprise
	ProfileComplete bool         `json:"profile_complete" bson:"profile_complete"`   // Indicates if the user profile is complete
	CreatedAt       time.Time    `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at" bson:"updated_at"`
}

type UserProfile struct {
	ID          string    `json:"id" bson:"_id"`
	UserID      string    `json:"user_id" bson:"user_id"`           // ID of the user associated with the profile
	FirstName   string    `json:"first_name" bson:"first_name"`     // First name of the user
	LastName    string    `json:"last_name" bson:"last_name"`       // Last name of the user
	PhoneNumber string    `json:"phone_number" bson:"phone_number"` // Phone number of the user
	Address     string    `json:"address" bson:"address"`           // Address of the user
	ProfilePic  string    `json:"profile_pic" bson:"profile_pic"`   // URL of the user's profile picture
	CompanyName string    `json:"company_name" bson:"company_name"` // Company name of the user
	CompanyID   string    `json:"company_id" bson:"company_id"`     // Company ID associated with the user
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`     // Creation time of the profile
	UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`     // Last update time of the profile
}

type UserToken struct {
	ID        string    `json:"id" bson:"_id"`
	UserID    string    `json:"user_id" bson:"user_id"`       // ID of the user associated with the token
	Token     string    `json:"token" bson:"token"`           // JWT token for the user
	IsValid   bool      `json:"is_valid" bson:"is_valid"`     // Indicates if the token is valid
	CreatedAt time.Time `json:"created_at" bson:"created_at"` // Creation time of the token
	ExpiresAt time.Time `json:"expires_at" bson:"expires_at"` // Expiration time of the token
}
