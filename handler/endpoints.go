package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

type CustomClaims struct {
	UserID int `json:"id"`
	jwt.StandardClaims
}

// SecretKey is a secret key used to sign and verify JWTs. In a real application, this should be kept secure.
var SecretKey = []byte("hello")

// This is just a test endpoint to get you started. Please delete this endpoint.
// (GET /hello)
func (s *Server) Hello(ctx echo.Context, params generated.HelloParams) error {

	paramRepo := repository.GetTestByIdInput{
		Id: string(params.Id),
	}
	fmt.Println(paramRepo)

	out, err := s.Repository.GetTestById(ctx.Request().Context(), paramRepo)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	var resp generated.HelloResponse
	resp.Message = fmt.Sprintf("Hello User %s", out.Name)

	return ctx.JSON(http.StatusOK, resp)
}

func (s *Server) AuthRegister(ctx echo.Context) error {
	// This handler will only be called when the JWT is valid and the JWT contains
	// the scopes required.
	var register generated.Register
	err := ctx.Bind(&register)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, "could not bind request body")
	}

	s.Lock()
	defer s.Unlock()

	err = validateRegister(register)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	hashPass, err := hashPassword(register.Password)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	registerInput := repository.SaveRegisterInput{
		Fullname:     register.Fullname,
		HashPassword: hashPass,
		PhoneNumber:  register.PhoneNumber,
	}

	out, err := s.Repository.SaveRegister(ctx.Request().Context(), registerInput)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	msg := "success"
	resp := generated.RegisterResponse{
		Data: &struct {
			Id *int "json:\"id,omitempty\""
		}{Id: &out.Id},
		Message: &msg,
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) AuthLogin(ctx echo.Context) error {
	// This handler will only be called when the JWT is valid and the JWT contains
	// the scopes required.
	var login generated.Login
	err := ctx.Bind(&login)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, "could not bind request body")
	}

	s.Lock()
	defer s.Unlock()

	repoInput := repository.GetUsersByPhoneInput{
		PhoneNumber: login.PhoneNumber,
	}

	repoOut, err := s.Repository.Login(ctx.Request().Context(), repoInput)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	err = checkPassword(login.Password, repoOut.HashPassword)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	numLoginInput := repository.UpdateNumberLoginInput{
		Id: repoOut.Id,
	}

	err = s.Repository.UpdateNumberLogin(ctx.Request().Context(), numLoginInput)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	token, err := generateToken(repoOut.Id)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	msg := "login success"
	resp := generated.LoginResponse{
		Data: &struct {
			Id  *int    "json:\"id,omitempty\""
			Jwt *string "json:\"jwt,omitempty\""
		}{
			Id:  &repoOut.Id,
			Jwt: &token,
		},
		Message: &msg,
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) GetProfile(ctx echo.Context, id int) error {

	claim, err := claimToken(ctx)
	if err != nil {
		return returnError(ctx, http.StatusForbidden, err.Error())
	}

	if claim.UserID != id {
		return returnError(ctx, http.StatusForbidden, "Forbidden")
	}

	repoInput := repository.GetProfiletByIdInput{
		Id: id,
	}
	repoOut, err := s.Repository.GetProfileById(ctx.Request().Context(), repoInput)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	msg := "get profile success"
	resp := generated.ProfileResponse{
		Data: &struct {
			Fullname    *string "json:\"fullname,omitempty\""
			PhoneNumber *string "json:\"phone_number,omitempty\""
		}{
			Fullname:    &repoOut.Fullname,
			PhoneNumber: &repoOut.PhoneNumber,
		},
		Message: &msg,
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func (s *Server) UpdateProfile(ctx echo.Context) error {

	claim, err := claimToken(ctx)
	if err != nil {
		return returnError(ctx, http.StatusForbidden, err.Error())
	}

	var update generated.UpdateProfile
	err = ctx.Bind(&update)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, "could not bind request body")
	}

	s.Lock()
	defer s.Unlock()

	repoInput := repository.GetProfileByPhoneInput{
		PhoneNumber: update.PhoneNumber,
	}
	repoOut, err := s.Repository.GetProfileByPhone(ctx.Request().Context(), repoInput)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	if claim.UserID != repoOut.Id {
		return returnError(ctx, http.StatusForbidden, "Forbidden")
	}

	updateInput := repository.UpdateProfileInput{
		Id:       repoOut.Id,
		Fullname: update.Fullname,
	}

	err = s.Repository.UpdateProfile(ctx.Request().Context(), updateInput)
	if err != nil {
		return returnError(ctx, http.StatusBadRequest, err.Error())
	}

	msg := "update profile success"
	resp := generated.ProfileResponse{
		Data: &struct {
			Fullname    *string "json:\"fullname,omitempty\""
			PhoneNumber *string "json:\"phone_number,omitempty\""
		}{
			Fullname:    &update.Fullname,
			PhoneNumber: &repoOut.PhoneNumber,
		},
		Message: &msg,
	}

	return ctx.JSON(http.StatusCreated, resp)
}

func returnError(ctx echo.Context, code int, message string) error {
	errResponse := generated.ErrorResponse{
		Message: message,
	}
	return ctx.JSON(code, errResponse)
}

func validateRegister(r generated.Register) error {
	if len(r.Fullname) <= 3 {
		return errors.New("fullname minimum 3 charrater")
	}

	if len(r.Fullname) > 60 {
		return errors.New("fullname max 60 character")
	}

	err := validatePhoneNum(r.PhoneNumber)
	if err != nil {
		return err
	}

	err = validatePassword(r.Password)
	if err != nil {
		return err
	}

	return nil
}

func validatePassword(pass string) error {
	if len(pass) < 6 || len(pass) > 64 {
		return errors.New("password minimum 6 and maximal 64 character")
	}

	var isLower, isUpper, isSym bool

	for _, r := range pass {
		if !isLower && unicode.IsLower(r) {
			isLower = true
		}

		if !isUpper && unicode.IsUpper(r) {
			isUpper = true
		}

		if !isSym && (unicode.IsSymbol(r) || unicode.IsPunct(r)) {
			isSym = true
		}
	}

	isValid := isLower && isUpper && isSym

	if !isValid {
		return errors.New("password must have containing at least 1 capital characters AND 1 number AND 1 special (non alpha-numeric) characters")
	}

	return nil
}

func validatePhoneNum(num string) error {
	normalizedPhoneNumber := strings.ReplaceAll(num, " ", "")
	hasPrefix := strings.HasPrefix(normalizedPhoneNumber, "+62")
	if !hasPrefix {
		return errors.New("phone number must start with the Indonesia country code '+62'")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	salt, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}

	hashedPassword := string(salt)
	return hashedPassword, nil
}

func checkPassword(plaintextPassword, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plaintextPassword))
	return err
}

func generateToken(userID int) (string, error) {
	claims := CustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func claimToken(e echo.Context) (*CustomClaims, error) {
	tokenString := e.Request().Header.Get("Authorization")

	if len(tokenString) < 7 || tokenString[:7] != "Bearer " {
		return nil, errors.New("unauthorized: Missing or invalid Bearer token")
	}
	tokenString = tokenString[7:]

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		return SecretKey, nil
	})

	if err != nil {
		log.Println("Error parsing token:", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("Unauthorized")
}
