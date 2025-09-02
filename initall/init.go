package initall

import (
	"github.com/hewo233/hdu-dx2/db"
	"github.com/hewo233/hdu-dx2/utils/jwt"
)

func Init() {
	db.Init()
	jwt.InitJWTKey()
}
