package repository

import (
	"backend-el-ternak/internal/config"
	"backend-el-ternak/internal/models"
	"fmt"

	"gorm.io/gorm"
)

func CreateUser(user *models.User) error {
	return config.DB.Transaction(func(tx *gorm.DB) error {

        if user.IsPJ {
            if user.KandangID == nil || *user.KandangID == 0 {
                return fmt.Errorf("kandang_id wajib diisi jika user menjadi penanggung jawab")
            }

            kandangID := *user.KandangID

            var oldPJ models.User
            err := tx.
                Where("kandang_id = ? AND is_pj = ?", kandangID, true).
                First(&oldPJ).Error

            if err == nil {
                if oldPJ.ID != user.ID {
                    if err := tx.Model(&oldPJ).Update("is_pj", false).Error; err != nil {
                        return err
                    }
                }
            } else if err != gorm.ErrRecordNotFound {
                return err
            }
        }

        if err := tx.Create(user).Error; err != nil {
            return err
        }

        return nil
    })
}

func GetUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserById(id uint) (*models.UserSummary, error) {
	var user models.UserSummary
	err := config.DB.Model(&models.User{}).
	Select("id", "username", "role", "is_active", "is_pj", "kandang_id").
	Where("id = ?", id).
	First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetAllUser() ([]models.UserSummary, error) {
	var users []models.UserSummary
	err := config.DB.Model(&models.User{}).
		Select(`users.id, users.username, users.role, users.is_active, users.is_pj,
		        users.kandang_id, kandangs.nama AS nama_kandang`).
		Joins(`LEFT JOIN kandangs ON users.kandang_id = kandangs.id`).
		Scan(&users).Error
	
	if err != nil {
		return nil, err
	}

	return users, nil
}

func GetUserByRole(role string) ([]models.UserSummary, error) {
	var users []models.UserSummary

	err := config.DB.Model(&models.User{}).
		Select(`users.id, users.username, users.role, users.is_active, users.is_pj,
		        users.kandang_id, kandangs.nama AS nama_kandang`).
		Joins(`LEFT JOIN kandangs ON users.kandang_id = kandangs.id`).
		Where("users.role = ?", role).
		Scan(&users).Error

	if err != nil {
		return nil, err
	}

	return users, nil
}


func UpdateUserByUsername(username string, newData map[string]interface{}) error {
    return config.DB.Transaction(func(tx *gorm.DB) error {
        var user models.User
        if err := tx.Where("username = ?", username).First(&user).Error; err != nil {
            return fmt.Errorf("user %s not found", username)
        }

        isBecomingPJ, ok := newData["is_pj"]
        if ok && isBecomingPJ == true {
            kandangID, hasKandang := newData["kandang_id"]
            if !hasKandang {
                kandangID = user.KandangID
            }

            var oldPJ models.User
            err := tx.
                Where("kandang_id = ? AND is_pj = ?", kandangID, true).
                First(&oldPJ).Error

            if err == nil {
                if oldPJ.ID != user.ID {
                    if err := tx.Model(&oldPJ).Update("is_pj", false).Error; err != nil {
                        return err
                    }
                }
            } else if err != gorm.ErrRecordNotFound {
                return err
            }
        }

        result := tx.Model(&user).Updates(newData)
        if result.Error != nil {
            return result.Error
        }
        if result.RowsAffected == 0 {
            return fmt.Errorf("not found")
        }

        return nil
    })
}


func DeleteByUsername(username string) error  {
	result := config.DB.Where("username = ?", username).Delete(&models.User{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("not found")
	}

	return nil
}