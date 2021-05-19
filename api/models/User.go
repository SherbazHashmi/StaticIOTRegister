package models

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"html"
	"log"
	"strings"
	"time"
)

type User struct {
	ID            uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Nickname      string    `gorm:"size:255; not null' unique" json:"nickname"`
	Email         string    `gorm:"size:100; not null;unique" json:"email"`
	Password      string    `gorm:"size:100; not null;" json:"password"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	AccountActive bool      `gorm:"default:true" json:"active"`
}

type FieldValidation struct {
	Label string
	Value string
}

type ValidationAction struct {
	Field               FieldValidation
	RequiredValidations []string
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Nickname = html.EscapeString(strings.TrimSpace(u.Nickname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	fieldValidations := map[string]ValidationAction{
		"nickname": {
			Field:               FieldValidation{Label: "nickname", Value: ""},
			RequiredValidations: []string{"presence"},
		},
		"password": {
			Field:               FieldValidation{Label: "password", Value: u.Password},
			RequiredValidations: []string{"presence"},
		},
		"email": {
			Field:               FieldValidation{Label: "email", Value: u.Email},
			RequiredValidations: []string{"presence", "email_validate"},
		},
	}

	// Determines which validations are required for each action
	actionValidationMap := map[string][]string{
		"update": {
			"nickname", "password", "email",
		},
		"login": {
			"email", "password",
		},
		"default": {
			"nickname", "password", "email",
		},
	}

	// Populates a list of validation actions
	fieldsToValidate := actionValidationMap[action]
	var invalidFields []string
	for _, fieldToValidate := range fieldsToValidate {
		// retrieve validation action
		fieldToValidate, isPresent := fieldValidations[fieldToValidate]
		if !isPresent {
			print("Unable to validate field")
		}

		// Iterate through required actions and perform them
		if validateField(fieldToValidate) {
			invalidFields = append(invalidFields, fieldToValidate.Field.Label)
		}
	}

	if len(invalidFields) > 0 {
		return errors.New("Following fields are invalid: " + strings.Join(invalidFields, ", ")[:1])
	}
	return nil
}

func validateField(validationAction ValidationAction) bool {
	requiredValidations := validationAction.RequiredValidations
	for _, requiredValidation := range requiredValidations {
		switch requiredValidation {
		case "presence":
			if validationAction.Field.Value == "" {
				return false
			}
		case "email_validate":
			if checkmail.ValidateFormat(validationAction.Field.Value) != nil {
				return false
			}
		}
	}
	return true
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	// Setup Error Object
	var err error

	// Create User with Debugging and Pull Out Error
	err = db.Debug().Create(&u).Error

	// Check and return error if present with empty user.
	if err != nil {
		return &User{}, err
	}

	// return user object without error
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	// setup error object
	var err error

	// allocate slice for users
	users := []User{}

	// query for users assuming errors, give a reference to the users object
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error

	if err != nil {
		return &[]User{}, err
	}
	return &users, err

}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	// Setup error object
	var err error

	// Query database for the user id, pull out potential error pointer
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error

	// If there is an error, return an empty user with an error message.
	if err != nil {
		return &User{}, err
	}

	// If there is no record return that as an error
	// SMELL, will this even be reached if the error is caught above?
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User not found")
	}

	return u, err
}

func (u *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error) {
	// perform pre update checks (hash password)
	err := u.BeforeSave()

	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"nickname":  u.Nickname,
			"email":     u.Email,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	err = db.Debug().Model(&User{}).Where("id = ? ", uid).Take(&u).Error

	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) DeleteUser(db *gorm.DB, uid uint32) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

func (u *User) ChangeUserActiveStatus(db *gorm.DB, uid uint32, active bool) (*User, error) {
	err := u.BeforeSave()

	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"disabled": active,
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}

	return u, nil
}
