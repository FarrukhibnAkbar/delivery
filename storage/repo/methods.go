package repo

import (
	"context"
	"delivery/entities"
)

// IAdminStorage account storage interface
type IAdminStorage interface {
	CreateXozmak(ctx context.Context, req entities.Xozmak) error
	Registration(ctx context.Context, req entities.RegistrReq) error
	UpdateUserProfile(ctx context.Context, updateData entities.UserProfile) error
	InsertUserLocation(ctx context.Context, req entities.UserLocation) error
	GetUserProfile(ctx context.Context, id string)(entities.UserProfile, error)
	GetUserLocation(ctx context.Context, userId string)([]entities.UserLocation, error)
	GetXozmak(ctx context.Context)([]entities.Xozmak, error)
	UpdateXozmak(ctx context.Context, req entities.Xozmak) error
	DeleteXozmak(ctx context.Context, id string) error
	CreateCategory(ctx context.Context, req entities.Category) error
	GetCategory(ctx context.Context)([]entities.Category, error)
	UpdateCategory(ctx context.Context, req entities.Category) error
	DeleteCategory(ctx context.Context, categoryId string) error
	CreateSubCategory(ctx context.Context, req entities.SubCategory) error
	GetSubCategory(ctx context.Context)([]entities.SubCategory, error)
	UpdateSubCategory(ctx context.Context, req entities.SubCategory) error
	DeleteSubCategory(ctx context.Context, sub_categoryId string) error
	AddFavorite(ctx context.Context, req entities.Favorite) error
}
