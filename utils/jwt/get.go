package jwt

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hewo233/hdu-dx2/db"
	models "github.com/hewo233/hdu-dx2/models"
	"github.com/hewo233/hdu-dx2/shared/consts"
)

func GetPhoneFromJWT(c *gin.Context) (string, *models.User, error) {
	phone := c.GetString("phone")

	user := models.NewUser()

	result := db.DB.Table(consts.UserTable).Where("phone = ?", phone).First(user)
	if result.Error != nil {
		return "", nil, result.Error
	}
	if result.RowsAffected == 0 {
		return "", nil, errors.New("user not found")
	}

	return user.Phone, user, nil
}
