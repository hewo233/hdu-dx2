package route

import (
	"github.com/gin-gonic/gin"
	"github.com/hewo233/hdu-dx2/handler"
	"github.com/hewo233/hdu-dx2/middleware"
	"github.com/hewo233/hdu-dx2/shared/consts"
)

var R *gin.Engine

func InitRoute() {
	R = gin.New()

	R.Use(gin.Logger(), gin.Recovery())
	R.Use(middleware.CorsMiddleware())

	R.GET("/ping", handler.Ping)

	auth := R.Group("/auth")
	{
		auth.POST("/register", handler.UserRegister)
		auth.POST("/login", handler.UserLogin)
	}

	user := R.Group("/user")
	user.Use(middleware.JWTAuth(consts.User))
	{
		user.GET("/info/:phone", handler.GetUserInfoByPhone)
		user.POST("/update", handler.ModifyUserSelf)
		// Test
		user.GET("/list", handler.ListUser)

		user.GET("/family", handler.ListUserFamily)
	}

	family := R.Group("/family")
	family.Use(middleware.JWTAuth(consts.User))
	{
		family.POST("/create", handler.CreateFamily)
		family.POST("/join", handler.AddUserToFamily)
		family.GET("/members/:family_id", handler.ListFamilyMember)
		family.GET("/list", handler.ListUserFamily)
	}

	financial := R.Group("/financial")
	financial.Use(middleware.JWTAuth(consts.User))
	{
		financial.POST("/bill/create/:family_id", handler.CreateBill)
		financial.GET("/bill/list/:family_id", handler.ListBills)
		financial.GET("/bill/select/:family_id", handler.SelectBills)
		financial.DELETE("/bill/delete/:family_id/:bill_id", handler.DeleteBill)
	}
}
