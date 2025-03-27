package models

import "github.com/go-playground/validator/v10"

// Validator instance
var validate = validator.New()

// ValidateUser checks if the user fields are valid
func (u *User) ValidateUser() error {
	return validate.Struct(u)
}

// ValidateBlog checks if the blog fields are valid
func (b *Blog) ValidateBlog() error {
	return validate.Struct(b)
}

// ValidateComment checks if the comment fields are valid
func (c *Comment) ValidateComment() error {
	return validate.Struct(c)
}

// ValidateLike checks if the like fields are valid
func (l *Like) ValidateLike() error {
	return validate.Struct(l)
}

// ValidateFollow checks if the follow fields are valid
func (f *Follow) ValidateFollow() error {
	return validate.Struct(f)
}
