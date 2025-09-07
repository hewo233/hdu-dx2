package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hewo233/hdu-dx2/db"
	"github.com/hewo233/hdu-dx2/models"
	"github.com/hewo233/hdu-dx2/shared/consts"
	"github.com/hewo233/hdu-dx2/utils/jwt"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type createFamilyRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func CreateFamily(c *gin.Context) {
	var req createFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40000,
			"message": "failed to bind CreateFamily Request: " + err.Error(),
		})
		c.Abort()
		return
	}

	_, _, err := jwt.GetPhoneFromJWT(c)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errno":   40101,
				"message": "Unauthorized, user in jwt not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errno":   50007,
				"message": "failed to get user info: " + err.Error(),
			})
		}
		c.Abort()
		return
	}

	family := models.NewFamily()
	family.Name = req.Name
	family.Password = req.Password

	if err := db.DB.Table(consts.FamilyTable).Create(family).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to create family: " + err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":    20000,
		"message":  "family created successfully",
		"familyId": family.ID,
	})
}

type addUserToFamilyRequest struct {
	FamilyID uint   `json:"family_id" binding:"required"`
	UserID   uint   `json:"user_id" binding:"required"`
	Role     string `json:"role" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AddUserToFamily(c *gin.Context) {
	var req addUserToFamilyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40000,
			"message": "failed to bind AddUserToFamily Request: " + err.Error(),
		})
		c.Abort()
		return
	}

	_, _, err := jwt.GetPhoneFromJWT(c)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errno":   40101,
				"message": "Unauthorized, user in jwt not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errno":   50007,
				"message": "failed to get user info: " + err.Error(),
			})
		}
		c.Abort()
		return
	}

	findFamily := models.NewFamily()

	// check family exists
	if err := db.DB.Table(consts.FamilyTable).Where("id = ?", req.FamilyID).First(findFamily).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"errno":   40004,
				"message": "this family does not exist",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to query database: " + err.Error(),
		})
		c.Abort()
		return
	}

	if findFamily.Password != req.Password {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40007,
			"message": "family password is incorrect",
		})
		c.Abort()
		return
	}

	findUser := models.NewUser()

	// check user exists
	if err := db.DB.Table(consts.UserTable).Where("id = ?", req.UserID).First(findUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"errno":   40005,
				"message": "this user does not exist",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to query database: " + err.Error(),
		})
		c.Abort()
		return
	}

	if findUser.ID != req.UserID {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40010,
			"message": "user cannot modify other user's family",
		})
		c.Abort()
		return
	}

	// check if user already in family
	result := db.DB.Table(consts.FamilyUserTable).Where("user_id = ? AND family_id = ?", req.UserID, req.FamilyID).Limit(1).Find(&models.FamilyUser{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to query database: " + result.Error.Error(),
		})
		c.Abort()
		return
	}
	if result.RowsAffected > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40006,
			"message": "this user is already in the family",
		})
		c.Abort()
		return
	}

	// add user to family

	familyUser := map[string]interface{}{
		"user_id":   req.UserID,
		"family_id": req.FamilyID,
		"role":      req.Role,
	}

	if err := db.DB.Table(consts.FamilyUserTable).Create(familyUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to add user to family: " + err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":   20000,
		"message": "user added to family successfully",
	})
}

func ListFamilyMember(c *gin.Context) {
	familyIDStr := c.Param("family_id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40000,
			"message": "invalid family_id: " + err.Error(),
		})
		c.Abort()
		return
	}

	// only family member can list family members
	_, user, err := jwt.GetPhoneFromJWT(c)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errno":   40101,
				"message": "Unauthorized, user in jwt not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errno":   50007,
				"message": "failed to get user info: " + err.Error(),
			})
		}
		c.Abort()
		return
	}

	// check if user is in family
	result := db.DB.Table(consts.FamilyUserTable).Where("user_id = ? AND family_id = ?", user.ID, familyID).Limit(1).Find(&models.FamilyUser{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to query database: " + result.Error.Error(),
		})
		c.Abort()
		return
	}
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40008,
			"message": "user is not in the family",
		})
		c.Abort()
		return
	}

	findFamily := models.NewFamily()

	// check family exists
	if err := db.DB.Table(consts.FamilyTable).Where("id = ?", uint(familyID)).First(findFamily).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{
				"errno":   40004,
				"message": "this family does not exist",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to query database: " + err.Error(),
		})
		c.Abort()
		return
	}

	type FamilyMember struct {
		UserID   uint   `json:"user_id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	var members []FamilyMember

	if err := db.DB.Table(consts.FamilyUserTable).Where("family_id = ?", uint(familyID)).
		Select("family_user.user_id, \"user\".username, family_user.role").
		Joins("LEFT JOIN \"user\" ON family_user.user_id = \"user\".id").
		Scan(&members).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to query database: " + err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":   20000,
		"message": "success",
		"members": members,
	})
}

func ListAllFamilies(c *gin.Context) {
	_, _, err := jwt.GetPhoneFromJWT(c)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"errno":   40101,
				"message": "Unauthorized, user in jwt not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"errno":   50007,
				"message": "failed to get user info: " + err.Error(),
			})
		}
		c.Abort()
		return
	}

	var families []models.Family
	if err := db.DB.Table(consts.FamilyTable).Find(&families).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to query database: " + err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":    20000,
		"message":  "success",
		"families": families,
	})

}
