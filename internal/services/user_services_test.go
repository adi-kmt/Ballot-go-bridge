package services

import (
	"fmt"
	"testing"
)

func getServicesWithMockRepo() *UserService {
	return NewUserService(&FakeRepo{})
}

func TestAddUserIfAlreadyExists(t *testing.T) {
	service := getServicesWithMockRepo()

	_, err := service.AddUser("user1", "password", true)

	if err == nil {
		t.Errorf("Something went wrong")
	} else if err.Message == "User already exists" {
		fmt.Println("User already exists")
	} else {
		t.Errorf("Testcase failed")
	}
}

func TestAddUser(t *testing.T) {
	service := getServicesWithMockRepo()
	token, err := service.AddUser("userName", "password", true)

	if err != nil {
		t.Errorf("Testcase failed")
	}
	fmt.Println(token)
}

func TestValidateUserUserDoesnotExist(t *testing.T) {
	service := getServicesWithMockRepo()
	_, err := service.ValidateUser("userName", "password")
	if err == nil {
		t.Errorf("Something went wrong")
	} else if err.Message == "User not found" {
		fmt.Println("User not found")
	} else {
		t.Errorf("Testcase failed")
	}
}

func TestValidateUserWrongPassword(t *testing.T) {
	service := getServicesWithMockRepo()
	_, err := service.ValidateUser("user1", "password")
	if err == nil {
		t.Errorf("Something went wrong")
	} else if err.Message == "Invalid password" {
		fmt.Println("Invalid password")
	} else {
		t.Errorf("Testcase failed")
	}
}

func TestValidateUser(t *testing.T) {
	service := getServicesWithMockRepo()
	_, err := service.ValidateUser("user1", "password1")
	if err != nil {
		t.Errorf("Testcase failed")
	}
	fmt.Println("We have a valid user")
}
