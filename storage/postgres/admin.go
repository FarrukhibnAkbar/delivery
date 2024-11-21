package postgres

import (
	"context"
	"delivery/constants"
	"delivery/entities"
	e "delivery/errors"
	"errors"
	"fmt"

	"github.com/jackc/pgconn"

	_ "github.com/lib/pq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type adminRepo struct {
	db *gorm.DB
}

func NewAdmin(db *gorm.DB) *adminRepo {
	return &adminRepo{db: db}
}

func (a adminRepo) CreateXozmak(ctx context.Context, req entities.Xozmak) error {
	res := a.db.WithContext(ctx).Table("xozmaks").Create(&req)
	if res.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(res.Error, &pgErr) && pgErr.Code == constants.PGUniqueKeyViolationCode {
			return fmt.Errorf("error in CreateXozmak: %w", constants.ErrXozmakAlreadyExists)
		}
		return fmt.Errorf("error in CreateXozmak: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("error in CreateXozmak: %w", constants.ErrRowsAffectedIsZero)
	}
	return nil
}

func (a adminRepo) GetXozmak(ctx context.Context) ([]entities.Xozmak, error) {
	var xozmak []entities.Xozmak
	err := a.db.Table("xozmaks").Where("state=?", constants.Active).Find(&xozmak).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.Xozmak{}, errors.New("xozmak not found")
		}
		return []entities.Xozmak{}, err
	}
	return xozmak, nil
}

func (a adminRepo) UpdateXozmak(ctx context.Context, req entities.Xozmak) error {
	result := a.db.WithContext(ctx).Table("xozmaks").Where("id = ?", req.ID).Updates(req)

	if result.Error != nil {
		return fmt.Errorf("failed to update xozmak data: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, xozmak with ID %v not found or no changes made", req.ID)
	}

	return nil
}

func (a adminRepo) DeleteXozmak(ctx context.Context, id string) error {

	res := a.db.WithContext(ctx).Table("xozmaks").Where("id = ?", id).Update("state", constants.InActive)
	if res.Error != nil {
		return fmt.Errorf("failed to delete xozmak data: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, xozmak with ID %v not found or no changes made", id)
	}

	return nil
}

func (a adminRepo) Registration(ctx context.Context, req entities.RegistrReq) error {
	res := a.db.WithContext(ctx).Table("users").Create(&req)
	if res.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(res.Error, &pgErr) && pgErr.Code == constants.PGUniqueKeyViolationCode {
			return e.ErrAccountAlreadyExists
		}
		return fmt.Errorf("error in Signup: %w", res.Error)
	}
	return nil
}

func (a *adminRepo) UpdateUserProfile(ctx context.Context, updateData entities.UserProfile) error {
	res := a.db.WithContext(ctx).Table("users").Where("id = ?", updateData.ID).Updates(updateData)
	if res.Error != nil {
		return fmt.Errorf("failed to update user data: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, userProfile with id %v not found or no changes made", updateData.ID)
	}
	return nil
}

func (a *adminRepo) InsertUserLocation(ctx context.Context, req entities.UserLocation) error {
	res := a.db.WithContext(ctx).Table("users_locations").Create(&req)
	if res.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(res.Error, &pgErr) && pgErr.Code == constants.PGUniqueKeyViolationCode {
			return fmt.Errorf("error in InsertUserLocation: %w", constants.ErrXozmakAlreadyExists)
		}
		return fmt.Errorf("error in InsertUserLocation: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("error in InsertUserLocation: %w", constants.ErrRowsAffectedIsZero)
	}
	return nil
}

func (a *adminRepo) GetUserProfile(ctx context.Context, userId string) (entities.UserProfile, error) {
	var usersData entities.UserProfile
	err := a.db.Where("id = ?", userId).Table("users").First(&usersData).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.UserProfile{}, errors.New("user not found")
		}
		return entities.UserProfile{}, err
	}
	return usersData, nil
}

func (a *adminRepo) GetUserLocation(ctx context.Context, userId string) ([]entities.UserLocation, error) {
	var userLocation []entities.UserLocation

	err := a.db.Where("user_id = ?", userId).Table("users_locations").Find(&userLocation).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.UserLocation{}, errors.New("user locations not found")
		}
		return []entities.UserLocation{}, err
	}
	return userLocation, nil
}

func (a *adminRepo) CreateCategory(ctx context.Context, req entities.Category) error {
	res := a.db.WithContext(ctx).Table("category").Create(&req)
	if res.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(res.Error, &pgErr) && pgErr.Code == constants.PGUniqueKeyViolationCode {
			return fmt.Errorf("error in CreateCategory: %w", constants.ErrXozmakAlreadyExists)
		}
		return fmt.Errorf("error in CreateCategory: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("error in CreateCategory: %w", constants.ErrRowsAffectedIsZero)
	}
	return nil
}

func (a *adminRepo) GetCategory(ctx context.Context) ([]entities.Category, error) {
	var categoryList []entities.Category
	err := a.db.Table("category").Where("state=?", constants.Active).Find(&categoryList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.Category{}, errors.New("user not found")
		}
		return []entities.Category{}, err
	}
	return categoryList, nil
}

func (a *adminRepo) UpdateCategory(ctx context.Context, req entities.Category) error {
	res := a.db.WithContext(ctx).Table("category").Where("id = ?", req.ID).Updates(req)
	if res.Error != nil {
		return fmt.Errorf("failed to update category: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, category with id %v not found or no changes made", req.ID)
	}
	return nil
}

func (a adminRepo) DeleteCategory(ctx context.Context, categoryId string) error {

	res := a.db.WithContext(ctx).Table("category").Where("id = ?", categoryId).Update("state", constants.InActive)
	if res.Error != nil {
		return fmt.Errorf("failed to delete category data: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, category with ID %v not found or no changes made", categoryId)
	}
	return nil
}

func (a *adminRepo) CreateSubCategory(ctx context.Context, req entities.SubCategory) error {
	res := a.db.WithContext(ctx).Table("sub_category").Create(&req)
	if res.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(res.Error, &pgErr) && pgErr.Code == constants.PGUniqueKeyViolationCode {
			return fmt.Errorf("error in CreateSubCategory: %w", constants.ErrXozmakAlreadyExists)
		}
		return fmt.Errorf("error in CreateSubCategory: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("error in CreateSubCategory: %w", constants.ErrRowsAffectedIsZero)
	}
	return nil
}

func (a *adminRepo) GetSubCategory(ctx context.Context) ([]entities.SubCategory, error) {
	var categoryList []entities.SubCategory
	err := a.db.Table("sub_category").Where("state=?", constants.Active).Find(&categoryList).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entities.SubCategory{}, errors.New("user not found")
		}
		return []entities.SubCategory{}, err
	}
	return categoryList, nil
}

func (a *adminRepo) UpdateSubCategory(ctx context.Context, req entities.SubCategory) error {
	res := a.db.WithContext(ctx).Table("sub_category").Where("id = ?", req.ID).Updates(req)
	if res.Error != nil {
		return fmt.Errorf("failed to update subCategory: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, subCategory with id %v not found or no changes made", req.ID)
	}
	return nil
}

func (a adminRepo) DeleteSubCategory(ctx context.Context, sub_categoryId string) error {

	res := a.db.WithContext(ctx).Table("sub_category").Where("id = ?", sub_categoryId).Update("state", constants.InActive)
	if res.Error != nil {
		return fmt.Errorf("failed to delete category data: %w", res.Error)
	}

	if res.RowsAffected == 0 {
		return fmt.Errorf("no rows affected, sub_category with ID %v not found or no changes made", sub_categoryId)
	}
	return nil
}

func (a adminRepo) AddFavorite(ctx context.Context, req entities.Favorite) error {
    if err := a.db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "user_id"}, {Name: "product_id"}},
        DoUpdates: clause.Assignments(map[string]interface{}{
            "is_favorited": gorm.Expr("NOT favorites.is_favorited"), 
        }),
    }).Create(&req).Error; err != nil {
        return fmt.Errorf("error in Addfavorite upsert: %v", err)
    }
    return nil
}


//func(a *adminRepo) GetProfile(ctx context.Context)()

// func (a adminRepo) GetUserByPhone(ctx context.Context, phoneNumber string) (entities.LoginPostgres, error) {
// 	user := entities.LoginPostgres{}

// 	res := a.db.WithContext(ctx).Table("users").Where("phone_number = ?", phoneNumber).First(&user)

// 	if res.Error != nil {
// 		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
// 			return entities.LoginPostgres{}, fmt.Errorf("no any user found with phone number %s", phoneNumber)
// 		}
// 		return entities.LoginPostgres{}, fmt.Errorf("error in GetUserByPhone: %w", res.Error)
// 	}

// 	if res.RowsAffected == 0 {
// 		return entities.LoginPostgres{}, fmt.Errorf("no user found with phone number %s", phoneNumber)
// 	}

// 	return user, nil
// }
