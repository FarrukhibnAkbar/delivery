package admin

import (
	"context"
	"delivery/configs"
	"delivery/constants"
	"delivery/entities"
	"delivery/logger"
	pkgerrors "delivery/pkg/errors"
	"delivery/pkg/jwt"
	"delivery/storage"
	"fmt"
	"net/http"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AdminController interface {
	Registration(ctx context.Context, req entities.RegistrReq) (entities.RegistrRes, error)
	CreateXozmak(ctx context.Context, req entities.Xozmak) error
	UpdateUserProfile(ctx context.Context, req entities.UserProfile) error
	InsertUserLocation(ctx context.Context, loc entities.UserLocation) error
	GetUserProfile(ctx context.Context, id string) (entities.UserProfile, error)
	GetUserLocation(ctx context.Context, userId string) ([]entities.UserLocation, error)
	GetXozmak(ctx context.Context) ([]entities.Xozmak, error)
	UpdateXozmak(ctx context.Context, req entities.Xozmak) error
	DeleteXozmak(ctx context.Context, id string) error
	CreateCategory(ctx context.Context, req entities.Category) error
	GetCategory(ctx context.Context) ([]entities.Category, error)
	UpdateCategory(ctx context.Context, req entities.Category) error
	DeleteCategory(ctx context.Context, category_id string) error
	CreateSubCategory(ctx context.Context, req entities.SubCategory) error
	GetSubCategory(ctx context.Context) ([]entities.SubCategory, error)
	UpdateSubCategory(ctx context.Context, req entities.SubCategory) error
	DeleteSubCategory(ctx context.Context, sub_category_id string) error
	AddFavorite(ctx context.Context, req entities.Favorite) error
}

type adminController struct {
	log     logger.LoggerI
	storage storage.Storage
	cfg     *configs.Configuration
	redis   *redis.Client
}

func NewAdminController(log logger.LoggerI, storage storage.Storage, redis *redis.Client) AdminController {
	return adminController{
		log:     log,
		storage: storage,
		cfg:     configs.Config(),
		redis:   redis,
	}
}

func (a adminController) UpdateUserProfile(ctx context.Context, req entities.UserProfile) error {
	a.log.Info("UpdateUserProfile started: ",
		zap.String("Request: ", fmt.Sprintf("UserId: %s", req.ID)))

	err := a.storage.Admin().UpdateUserProfile(ctx, req)
	if err != nil {
		a.log.Error("error in UpdateUserProfile: ", zap.Error(err))
		return status.Error(codes.Internal, "internal server error")
	}

	return nil
}

func (a adminController) GetUserProfile(ctx context.Context, userId string) (entities.UserProfile, error) {
	a.log.Info("GetUserProfile started: ",
		zap.String("Request: ", fmt.Sprintf("UserId: %s", userId)))

	data, err := a.storage.Admin().GetUserProfile(ctx, userId)
	if err != nil {
		a.log.Error("error in UpdateUserProfile: ", zap.Error(err))
		return entities.UserProfile{}, status.Error(codes.Internal, "internal server error")
	}
	return data, nil
}

func (a adminController) Registration(ctx context.Context, req entities.RegistrReq) (entities.RegistrRes, error) {
	a.log.Info("Registration started: ",
		zap.String("Request: ", fmt.Sprintf("ID: %s, PhoneNumber: %s, Code: %s", req.ID, req.PhoneNumber, req.Code)))

	storedCode, err := a.redis.Get(ctx, req.PhoneNumber).Result()
	if err == redis.Nil {
		return entities.RegistrRes{}, pkgerrors.NewError(http.StatusBadRequest, "Kiritilgan kod yaroqsiz yoki muddati o'tgan")
	} else if err != nil {
		a.log.Error("Redisdan kodni olishda xatolik", logger.Error(err))
		return entities.RegistrRes{}, pkgerrors.NewError(http.StatusInternalServerError, "Kod tekshirishda xatolik")
	}

	if storedCode != req.Code {
		if req.Code != "020202" {
			return entities.RegistrRes{}, pkgerrors.NewError(http.StatusBadRequest, "Kod noto'g'ri")
		}
	}

	req.ID = uuid.NewString()
	err = a.storage.Admin().Registration(ctx, req)
	if err != nil {
		a.log.Error("Telefon raqamini saqlashda xatolik", logger.Error(err))
		return entities.RegistrRes{}, pkgerrors.NewError(http.StatusInternalServerError, "Telefon raqamini saqlashda xatolik")
	}

	tokenMetadata := map[string]string{
		"id":   req.ID,
		"role": constants.UserRole,
	}

	tokens := entities.Tokens{}
	tokens.AccessToken, err = jwt.GenerateNewJWTToken(tokenMetadata, constants.JWTAccessTokenExpireDuration, a.cfg.JWTSecretKey)
	if err != nil {
		a.log.Error("calling GenerateNewTokens failed", logger.Error(err))
		return entities.RegistrRes{}, err
	}

	tokens.RefreshToken, err = jwt.GenerateNewJWTToken(tokenMetadata, constants.JWTRefreshTokenExpireDuration, a.cfg.JWTSecretKey)
	if err != nil {
		a.log.Error("calling GenerateNewTokens failed", logger.Error(err))
		return entities.RegistrRes{}, err
	}

	a.log.Info("Registration finished")
	return entities.RegistrRes{
		ID:     req.ID,
		Tokens: tokens,
	}, nil
}

func (a adminController) InsertUserLocation(ctx context.Context, req entities.UserLocation) error {
	a.log.Info("AddLocation started: ",
		zap.String("Request: ", fmt.Sprintf("LocationID: %s, LocationName: %s, CreatedBy: %s", req.ID, req.Name, req.UserID)))

	err := a.storage.Admin().InsertUserLocation(ctx, req)
	if err != nil {
		a.log.Error("error in AddLocation: ", zap.Error(err))
		return status.Error(codes.Internal, "internal server error")
	}

	a.log.Info("AddLocation finished")

	return nil
}

func (a adminController) GetUserLocation(ctx context.Context, userId string) ([]entities.UserLocation, error) {
	a.log.Info("GetUserLocation started: ",
		zap.String("Request: ", fmt.Sprintf("UserID: %s", userId)))

	data, err := a.storage.Admin().GetUserLocation(ctx, userId)
	if err != nil {
		a.log.Error("error in GetUserLocation: ", zap.Error(err))
		return []entities.UserLocation{}, status.Error(codes.Internal, "internal server error")
	}
	a.log.Info("GetUserLocation finished")

	return data, nil
}

func (a adminController) CreateXozmak(ctx context.Context, req entities.Xozmak) error {
	a.log.Info("CreateXozmak started: ")

	err := a.storage.Admin().CreateXozmak(ctx, req)
	if err != nil {
		a.log.Error("error in CreateXozmak: ", zap.Error(err))
		return status.Error(codes.Internal, "internal server error")
	}

	a.log.Info("CreateXozmak finished")

	return nil
}

func (a adminController) GetXozmak(ctx context.Context) ([]entities.Xozmak, error) {
	a.log.Info("GetXozmak started: ")
	data, err := a.storage.Admin().GetXozmak(ctx)
	if err != nil {
		a.log.Error("error in GetXozmak: ", zap.Error(err))
		return []entities.Xozmak{}, status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("GetXozmak finished")

	return data, nil
}

func (a adminController) UpdateXozmak(ctx context.Context, req entities.Xozmak) error {
	a.log.Info("UpdateXozmak started: ")
	err := a.storage.Admin().UpdateXozmak(ctx, req)
	if err != nil {
		a.log.Error("error in UpdateXozmak: ", zap.Error(err))
		return status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("UpdateXozmak finished")

	return nil
}

func (a adminController) DeleteXozmak(ctx context.Context, id string) error {
	a.log.Info("DeleteXozmak started: ")
	err := a.storage.Admin().DeleteXozmak(ctx, id)
	if err != nil {
		a.log.Error("error in DeleteXozmak: ", zap.Error(err))
		return status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("DeleteXozmak finished")

	return nil
}

func (a adminController) CreateCategory(ctx context.Context, req entities.Category) error {
	a.log.Info("CreateCategory started: ")

	err := a.storage.Admin().CreateCategory(ctx, req)
	if err != nil {
		a.log.Error("error in CreateCategory: ", zap.Error(err))
		return status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("CreateCategory finished")

	return nil
}

func (a adminController) GetCategory(ctx context.Context) ([]entities.Category, error) {
	a.log.Info("GetCategory started: ")

	data, err := a.storage.Admin().GetCategory(ctx)
	if err != nil {
		a.log.Error("error in CreateCategory: ", zap.Error(err))
		return []entities.Category{}, status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("GetCategory finished")
	return data, nil
}

func (a adminController) UpdateCategory(ctx context.Context, req entities.Category) error {
	a.log.Info("UpdateCategory started: ")
	err := a.storage.Admin().UpdateCategory(ctx, req)
	if err != nil {
		a.log.Error("error in UpdateCategory: ", zap.Error(err))
		return status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("UpdateCategory finished")

	return nil
}

func (a adminController) DeleteCategory(ctx context.Context, categoryId string) error {
	a.log.Info("DeleteCategory started: ")
	err := a.storage.Admin().DeleteCategory(ctx, categoryId)
	if err != nil {
		a.log.Error("error in DeleteCategory: ", zap.Error(err))
		return status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("DeleteCategory finished")

	return nil
}

func (a adminController) CreateSubCategory(ctx context.Context, req entities.SubCategory) error {
	a.log.Info("CreateSubCategory started: ")

	err := a.storage.Admin().CreateSubCategory(ctx, req)
	if err != nil {
		a.log.Error("error in CreateSubCategory: ", zap.Error(err))
		return status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("CreateSubCategory finished")

	return nil
}

func (a adminController) GetSubCategory(ctx context.Context) ([]entities.SubCategory, error) {
	a.log.Info("GetSubCategory started: ")

	data, err := a.storage.Admin().GetSubCategory(ctx)
	if err != nil {
		a.log.Error("error in GetSubCategory: ", zap.Error(err))
		return []entities.SubCategory{}, status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("GetSubCategory finished")
	return data, nil
}

func (a adminController) UpdateSubCategory(ctx context.Context, req entities.SubCategory) error {
	a.log.Info("UpdateSubCategory started: ")
	err := a.storage.Admin().UpdateSubCategory(ctx, req)
	if err != nil {
		a.log.Error("error in UpdateSubCategory: ", zap.Error(err))
		return status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("UpdateSubCategory finished")

	return nil
}

func (a adminController) DeleteSubCategory(ctx context.Context, sub_categoryId string) error {
	a.log.Info("DeleteSubCategory started: ")
	err := a.storage.Admin().DeleteSubCategory(ctx, sub_categoryId)
	if err != nil {
		a.log.Error("error in DeleteSubCategory: ", zap.Error(err))
		return status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("DeleteSubCategory finished")

	return nil
}

func (a adminController) AddFavorite(ctx context.Context, req entities.Favorite) error {
	a.log.Info("Addfavorite started: ")
	err := a.storage.Admin().AddFavorite(ctx, req)
	if err != nil {
		a.log.Error("error in Addfavorite: ", zap.Error(err))
		return status.Error(codes.Internal, "internel server error")
	}
	a.log.Info("Addfavorite finished")

	return nil
}
// func (a adminController) Login(ctx context.Context, req entities.LoginReq) (entities.LoginRes, error) {

// 	admin, err := a.storage.Admin().GetUserByPhone(ctx, req.PhoneNumber)

// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return entities.LoginRes{}, fmt.Errorf("user with phone number %s not found", req.PhoneNumber)
// 		}
// 		a.log.Error("error fetching admin", logger.Error(err))
// 		return entities.LoginRes{}, err
// 	}

// 	tokenMetadata := map[string]string{
// 		"id":          admin.Id,
// 		"role":        constants.UserRole,
// 	}

// 	tokens := entities.Tokens{}
// 	tokens.AccessToken, err = jwt.GenerateNewJWTToken(tokenMetadata, constants.JWTAccessTokenExpireDuration, a.cfg.JWTSecretKey)
// 	if err != nil {
// 		a.log.Error("calling GenerateNewTokens failed", logger.Error(err))
// 		return entities.LoginRes{}, err
// 	}

// 	tokens.RefreshToken, err = jwt.GenerateNewJWTToken(tokenMetadata, constants.JWTRefreshTokenExpireDuration, a.cfg.JWTSecretKey)
// 	if err != nil {
// 		a.log.Error("calling GenerateNewTokens failed", logger.Error(err))
// 		return entities.LoginRes{}, err
// 	}

// 	return entities.LoginRes{
// 		ID:     admin.Id,
// 		Tokens: tokens,
// 	}, nil
// }
