package modeltests

import (
	"github.com/SherbazHashmi/goblog/api/models"
	"gopkg.in/go-playground/assert.v1"
	"log"
	"testing"
)

func TestFindAllUsers(t *testing.T) {
	err := refreshUserTable()

	if err != nil {
		log.Fatal(err)
	}

	errs := seedUsers()

	if len(errs) > 0 {
		for _, err := range errs {
			log.Fatalf("[ERR] %v", err)
		}

	}

	users, err := userInstance.FindAllUsers(server.DB)

	if err != nil {
		t.Errorf("this is the error getting the users; %v\n", err)
		return
	}

	assert.Equal(t, len(*users), 2)
}

func TestSaveUser(t *testing.T) {
	err := refreshUserTable()

	if err != nil {
		log.Fatal(err)
	}

	newUser := models.User{
		ID:       1,
		Email:    "test@gmail.com",
		Nickname: "test",
		Password: "password",
	}

	savedUser, err := newUser.SaveUser(server.DB)

	if err != nil {
		t.Errorf("unable to get users: %v\n,", err)
	}

	assert.Equal(t, newUser.ID, savedUser.ID)
	assert.Equal(t, newUser.Email, savedUser.Email)
	assert.Equal(t, newUser.Nickname, savedUser.Nickname)
}

func TestUserById(t *testing.T) {
	err := refreshUserTable()

	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()

	if err != nil {
		log.Fatalf("can't seed users table: %v", err)
	}

	foundUser, err := userInstance.FindUserByID(server.DB, user.ID)

	if err != nil {
		log.Fatalf("unable to find the user with ID: %d,\n%v", user.ID, err)
	}

	assert.Equal(t, foundUser.ID, user.ID)
	assert.Equal(t, foundUser.Email, user.Email)
	assert.Equal(t, foundUser.Nickname, user.Nickname)
}

func TestUpdateAUser(t *testing.T) {
	err := refreshUserTable()

	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()

	if err != nil {
		log.Fatalf("unable to seed users %v", err)
	}

	userUpdate := models.User{
		ID:       1,
		Nickname: "modiUpdate",
		Email:    "modifiedEmail@gmail.com",
		Password: "password",
	}

	updatedUser, err := userUpdate.UpdateUser(server.DB, user.ID)

	if err != nil {
		t.Errorf("unable to update user %d, %v", user.ID, err)
		return
	}

	assert.Equal(t, updatedUser.ID, userUpdate.ID)
	assert.Equal(t, updatedUser.Email, userUpdate.Email)
	assert.Equal(t, updatedUser.Nickname, userUpdate.Nickname)
}

func TestDeleteAUser(t *testing.T) {
	err := refreshUserTable()

	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()

	if err != nil {
		log.Fatalf("unable to seed user ID: %d, %s ", user.ID, err)
		return
	}

	isDeleted, err := userInstance.DeleteUser(server.DB, user.ID)

	if err != nil {
		t.Errorf(" unable to delete user ID:%d,\n %v", user.ID, err)
		return
	}
	assert.Equal(t, isDeleted, int64(1))
}
