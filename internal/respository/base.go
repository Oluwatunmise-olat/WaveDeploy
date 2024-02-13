package respository

import (
	"github.com/Oluwatunmise-olat/WaveDeploy/internal/db"
	"gorm.io/gorm"
)

func DBTransaction() *gorm.DB {
	return db.DB
}
