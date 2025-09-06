package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hewo233/hdu-dx2/db"
	"github.com/hewo233/hdu-dx2/models"
	"github.com/hewo233/hdu-dx2/shared/consts"
	"net/http"
	"strconv"
	"time"
)

func checkUserInFamily(c *gin.Context, familyID uint) {
	phone := c.GetString("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40000,
			"message": "phone is required",
		})
		c.Abort()
		return
	}

	var user models.User
	if err := db.DB.Table("user").Where("phone = ?", phone).First(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40001,
			"message": "user not found: " + err.Error(),
		})
		c.Abort()
		return
	}

	var familyUser models.FamilyUser
	if err := db.DB.Table("family_user").Where("user_id = ? AND family_id = ?", user.ID, familyID).First(&familyUser).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"errno":   40300,
			"message": "user not in family: " + err.Error(),
		})
		c.Abort()
		return
	}
}

type createBillRequest struct {
	Date        string `json:"date" binding:"required"`
	Type        string `json:"type" binding:"required,oneof=income expense"`
	Amount      int    `json:"amount" binding:"required,gt=0"`
	Category    string `json:"category" binding:"required"`
	Description string `json:"description"`
	Object      string `json:"object" binding:"required"`
	Username    string `json:"username" binding:"required"`
}

func CreateBill(c *gin.Context) {
	var req createBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40000,
			"message": "failed to bind CreateBill Request: " + err.Error(),
		})
		c.Abort()
		return
	}

	familyIDStr := c.Param("family_id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40002,
			"message": "invalid family_id: " + err.Error(),
		})
		c.Abort()
		return
	}

	checkUserInFamily(c, uint(familyID))
	if c.IsAborted() {
		return
	}

	timeDate, err := time.Parse(consts.TimeFormat, req.Date)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40001,
			"message": "failed to parse date: " + err.Error(),
		})
		c.Abort()
		return
	}

	bill := &models.Bill{
		Date:        timeDate,
		Type:        req.Type,
		Amount:      req.Amount,
		Category:    req.Category,
		Description: req.Description,
		Object:      req.Object,
		Username:    req.Username,
		FamilyID:    uint(familyID),
	}

	if err := db.DB.Table(consts.BillTable).Create(bill).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to create bill: " + err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":   20000,
		"message": "Create Bill Successfully",
		"data":    bill,
	})
}

func ListBills(c *gin.Context) {

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

	checkUserInFamily(c, uint(familyID))
	if c.IsAborted() {
		return
	}

	var bills []models.Bill
	if err := db.DB.Table(consts.BillTable).Where("family_id=?", uint(familyID)).Find(&bills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to list bills: " + err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":   20000,
		"message": "List Bills Successfully",
		"data":    bills,
	})
}

type SelectBillsRequest struct {
	Type      string `json:"type" binding:"omitempty,oneof=income expense"`
	Category  string `json:"category" binding:"omitempty"`
	Object    string `json:"object" binding:"omitempty"`
	Username  string `json:"username" binding:"omitempty"`
	StartDate string `json:"start_date" binding:"omitempty"`
	EndDate   string `json:"end_date" binding:"omitempty"`
}

func SelectBills(c *gin.Context) {
	var req SelectBillsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40000,
			"message": "failed to bind SelectBills Request: " + err.Error(),
		})
		c.Abort()
		return
	}

	fmt.Printf("Query conditions: %+v\n", req)
	fmt.Println(req)

	familyIDStr := c.Param("family_id")
	familyID, err := strconv.ParseUint(familyIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40001,
			"message": "invalid family_id: " + err.Error(),
		})
		c.Abort()
		return
	}

	checkUserInFamily(c, uint(familyID))
	if c.IsAborted() {
		return
	}

	query := db.DB.Table(consts.BillTable)

	query = query.Where("family_id = ?", uint(familyID))

	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.Category != "" {
		query = query.Where("category = ?", req.Category)
	}
	if req.Object != "" {
		query = query.Where("object = ?", req.Object)
	}
	if req.Username != "" {
		query = query.Where("username = ?", req.Username)
	}
	if req.StartDate != "" {
		startDate, err := time.Parse(consts.TimeFormat, req.StartDate)
		fmt.Println("Parsed start date:", startDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errno":   40001,
				"message": "failed to parse start_date: " + err.Error(),
			})
			c.Abort()
			return
		}
		query = query.Where("date >= ?", startDate)
	}
	if req.EndDate != "" {
		endDate, err := time.Parse(consts.TimeFormat, req.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"errno":   40002,
				"message": "failed to parse end_date: " + err.Error(),
			})
			c.Abort()
			return
		}
		fmt.Println("Parsed end date:", endDate)
		query = query.Where("date <= ?", endDate)
	}

	var bills []models.Bill
	if err := query.Find(&bills).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to select bills: " + err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":   20000,
		"message": "Select Bills Successfully",
		"data":    bills,
	})
}

func DeleteBill(c *gin.Context) {
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

	billIDStr := c.Param("bill_id")
	billID, err := strconv.ParseUint(billIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40001,
			"message": "invalid bill_id: " + err.Error(),
		})
		c.Abort()
		return
	}

	checkUserInFamily(c, uint(familyID))
	if c.IsAborted() {
		return
	}

	var bill models.Bill
	if err := db.DB.Table("bill").Where("id = ? AND family_id = ?", uint(billID), uint(familyID)).First(&bill).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"errno":   40002,
			"message": "bill not found: " + err.Error(),
		})
		c.Abort()
		return
	}

	if err := db.DB.Table("bill").Where("id = ?", uint(billID)).Delete(&models.Bill{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"errno":   50000,
			"message": "failed to delete bill: " + err.Error(),
		})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"errno":   20000,
		"message": "Delete Bill Successfully",
	})
}
