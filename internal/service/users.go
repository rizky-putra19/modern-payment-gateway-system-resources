package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/hypay-id/backend-dashboard-hypay/config"
	"github.com/hypay-id/backend-dashboard-hypay/internal"
	"github.com/hypay-id/backend-dashboard-hypay/internal/dto"
	"github.com/hypay-id/backend-dashboard-hypay/internal/entity"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/helper"
	"github.com/hypay-id/backend-dashboard-hypay/internal/pkg/slog"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	userRepoReads  internal.UserReadsRepositoryItf
	userRepoWrites internal.UserWritesRepositoryItf
	cfg            config.App
}

func NewUser(
	userRepoReads internal.UserReadsRepositoryItf,
	userRepoWrites internal.UserWritesRepositoryItf,
	cfg config.App,
) *User {
	return &User{
		userRepoReads:  userRepoReads,
		userRepoWrites: userRepoWrites,
		cfg:            cfg,
	}
}

func (u *User) Login(payload dto.LoginPayload) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	user, err := u.userRepoReads.GetUserByUsername(payload.Username)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if user.UserID < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "user not found",
		}
		return resp, errors.New("user not found")
	}

	if !comparePasswords(user.Password, []byte(payload.Password)) {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusBadRequest,
			ResponseMessage: "invalid password",
		}
		return resp, errors.New("invalid password")
	}

	userPermissions, err := u.userRepoReads.GetPermissionByRoleId(user.RoleId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, errors.New("user not have permissions")
	}
	user.Permissions = userPermissions

	token, err := u.generateJWTToken(user)
	if err != nil {
		slog.Infof("%v failed to generate token with error message %v", payload.Username, err.Error())
		return dto.ResponseDto{}, errors.New(err.Error())
	}

	return dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "succcessfully login",
		Data: map[string]interface{}{
			"user":  user,
			"token": token,
		},
	}, nil
}

func (u *User) GetUserData(username string) (entity.User, error) {
	user, err := u.userRepoReads.GetUserByUsername(username)
	if err != nil {
		return entity.User{}, errors.New("user not found")
	}

	return user, nil
}

func (u *User) GetListRolesSvc() (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listRoles, err := u.userRepoReads.GetRolesRepo()
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            listRoles,
	}

	return resp, nil
}

func (u *User) GetListUserMerchantSvc(merchantId string) (dto.ResponseDto, error) {
	var resp dto.ResponseDto

	listUsers, err := u.userRepoReads.GetListUserByMerchantIdRepo(merchantId)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	if len(listUsers) < 1 {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusOK,
			ResponseMessage: "data not found",
			Data:            listUsers,
		}
		return resp, nil
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: "success retrieve data",
		Data:            listUsers,
	}

	return resp, nil
}

func (u *User) InviteUserMerchantSvc(payload dto.InviteMerchantUserDto) (dto.ResponseDto, error) {
	var resp dto.ResponseDto
	cfgAppMail := u.cfg.AppPassMail
	username := payload.Email
	password := helper.GenerateRandomString(8)
	pin := helper.GenerateRandomPinNumericString(6)

	passHash, _ := helper.HashString([]byte(password))
	pinHash, _ := helper.HashString([]byte(pin))

	credentialForRepo := dto.EmailDataHtmlDto{
		Password: passHash,
		Pin:      pinHash,
		Username: payload.Email,
	}

	// create user
	idUser, err := u.userRepoWrites.CreateUsersMerchantRepo(payload, credentialForRepo)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	// credential for send email
	credentialSendMail := dto.EmailDataHtmlDto{
		Password: password,
		Pin:      pin,
		Username: username,
	}

	// send email
	err = helper.SendEmailInviteUser(credentialSendMail, payload.Email, cfgAppMail)
	if err != nil {
		resp = dto.ResponseDto{
			ResponseCode:    http.StatusUnprocessableEntity,
			ResponseMessage: err.Error(),
		}
		return resp, err
	}

	resp = dto.ResponseDto{
		ResponseCode:    http.StatusOK,
		ResponseMessage: fmt.Sprintf("Success create user %v with id: %v", payload.Email, idUser),
	}

	return resp, nil
}

func (u *User) generateJWTToken(user entity.User) (string, error) {
	claims := &dto.Claims{
		Username: user.Username,
		UserType: user.UserType,
		RoleName: user.RoleName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 1).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(u.cfg.JWTSecret))
}

func comparePasswords(hashedPwd string, plainPwd []byte) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		slog.Errorw("compare password failed", "stack_trace", err.Error())
		return false
	}

	return true
}
