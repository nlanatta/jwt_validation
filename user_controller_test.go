package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strconv"
	"testing"
)

type testData struct {
	user          *User
	loginPassword string
}

func ExeRegistration(email string, password string) (*httptest.ResponseRecorder, error) {
	data := getTestData(email, password)
	ur := UserRegisterRequest{User: UserRequest{
		Email:    data.user.Email,
		Password: data.user.HashedPassword,
	}}

	d, err := json.Marshal(&ur)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer(d))
	if err != nil {
		return nil, err
	}

	// Create a response recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the handler function
	Register(recorder, req)
	return recorder, nil
}

func ExeToken(email string, password string) (*httptest.ResponseRecorder, error) {
	login := UserTokenRequest{
		Email:    email,
		Password: password,
	}

	data, err := json.Marshal(&login)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", "/token", bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Create a response recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the handler function
	Token(recorder, req)
	return recorder, nil
}

func ExeGetUser(token string) (*httptest.ResponseRecorder, error) {

	req, err := http.NewRequest("GET", "/api/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// Create a response recorder to capture the response
	recorder := httptest.NewRecorder()

	// Call the handler function
	GetUser(recorder, req)
	return recorder, nil
}

func TestRegister(t *testing.T) {
	email := "test@test.com"
	password := "123456789"
	value, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatal(err)
	}
	recorder, err := ExeRegistration(email, string(value))
	if err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	r := &UserRegisterResponse{}
	json.NewDecoder(recorder.Body).Decode(&r)
	if r.User.Email != email {
		t.Errorf("Expected body '%s', but got '%s'", email, r.User.Email)
	}
	recorder.Flush()
}

func TestLogin(t *testing.T) {
	password := "123456789"
	email := "test@test.com"
	data := getTestData(email, password)
	UserDB().AddUser(data.user)

	recorder, err := ExeToken(data.user.Email, data.loginPassword)
	if err != nil {
		t.Fatal(err)
	}
	if recorder.Code != http.StatusOK {
		t.Errorf("expected status code %d, but got %d", http.StatusOK, recorder.Code)
	}

	r := &TokenResponse{}
	json.NewDecoder(recorder.Body).Decode(&r)
	if r.Token == "" {
		t.Errorf("expected body to be not empty")
	}

	_, err = ValidateToken(r.Token)
	if err != nil {
		t.Fatal("token invalid")
	}
}

func TestGetUser(t *testing.T) {
	password := "123456789"
	email := "test@test.com"
	data := getTestData(email, password)
	UserDB().AddUser(data.user)

	recorder, err := ExeToken(data.user.Email, data.loginPassword)
	if err != nil {
		t.Fatal(err)
	}
	r := &TokenResponse{}
	json.NewDecoder(recorder.Body).Decode(&r)
	if r.Token == "" {
		t.Errorf("expected body to be not empty")
	}

	type testCase struct {
		name           string
		token          string
		data           testData
		expectedStatus int
		validToken     bool
	}
	testCases := []testCase{
		{
			name:           "Valid Token",
			token:          r.Token,
			data:           data,
			expectedStatus: 200,
			validToken:     true,
		},
		{
			name:           "Invalid Token",
			token:          "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMTEzNDQzNDUiLCJyb2xlcyI6WyJVU0VSIl0sImV4cCI6MTY4NjkyMjIwMiwiaXNzIjoiVG9rZW5SZXNwb25zZSJ9.0OaElAsUIbRKS7RFCRsG62EndW5azjML7RBVJRM252MoknqRDrJobO_3-a_LcJEnBwL-9KdXYj93MszzIWQqyA",
			data:           data,
			expectedStatus: 400,
			validToken:     false,
		},
		{
			name:           "Invalid Token",
			token:          createToken("NOROLE", data.user.ID),
			data:           data,
			expectedStatus: 401,
			validToken:     false,
		},
	}

	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			recorder, err = ExeGetUser(v.token)

			if recorder.Code != v.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", http.StatusOK, recorder.Code)
			}

			_, err := ValidateToken(v.token)
			if v.validToken {
				if err != nil {
					t.Fatal("invalid token")
				}
			} else {
				if err == nil {
					t.Fatal("expected an invalid token")
				}
			}

			if recorder.Code == http.StatusOK {
				ur := &UserResponse{}
				json.NewDecoder(recorder.Body).Decode(&ur)

				if ur.Email != v.data.user.Email {
					t.Fatalf("expected %s, got %s", v.data.user.Email, ur.Email)
				}
			}

		})
	}
}

func getTestData(email, password string) testData {
	os.Setenv("TOKEN_EXP", "5")
	os.Setenv("SECRET_KEY", "9191919191919")
	value := GenerateBCryptPasswordUtil([]byte(password))
	hashedPassword := value
	storedPassword := GenerateBCryptPasswordUtil([]byte(hashedPassword))
	u := &User{
		Email:          email,
		HashedPassword: storedPassword,
		ID:             int64(uuid.New().ID()),
	}

	toReturn := testData{
		user:          u,
		loginPassword: hashedPassword,
	}
	return toReturn
}

func createToken(role string, id int64) string {
	token, _ := GenerateToken(strconv.FormatUint(uint64(id), 10), []string{role}, reflect.TypeOf(TokenResponse{}))
	return token
}
