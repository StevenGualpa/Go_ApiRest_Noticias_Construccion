package repositories

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"noticias_uteq/models"
	"os"
)

var (
	DB *gorm.DB
)

func InitDB() {
	//dsn := "mysql://s8lcdocvs33ocomi:mzgkw1tckobvzwuk@en1ehf30yom7txe7.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306/ivvbv5zu5z6orru9"
	//dsn := "s8lcdocvs33ocomi:mzgkw1tckobvzwuk@tcp(en1ehf30yom7txe7.cbetxkdyhwsb.us-east-1.rds.amazonaws.com:3306)/ivvbv5zu5z6orru9?charset=utf8mb4&parseTime=True&loc=Local"
	dsnLocal := "postgres://hblzlyfbvyzvqy:ec9d05e02a5bc4a693c3ec36e6561630ef92399f0571dd76bcebd33721f6c152@ec2-34-236-100-103.compute-1.amazonaws.com:5432/d9huao3al0fbf1"
	DBURL := os.Getenv("DATABASE_URL")
	if len(DBURL) < 1 {
		DBURL = dsnLocal
	}
	db, err := gorm.Open(postgres.Open(DBURL), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	DB = db

	fmt.Println("Conectado a la DB")
}

func Migrate() {
	if err := DB.AutoMigrate(
		&models.Usuario{},
		&models.Noticia{},
		&models.Facultad{},
		&models.Multimedia{},
		&models.Revista{},
		&models.UserCode{},
		&models.Facultad{},
		&models.Carrera{},
		&models.UsuarioFacultad{},
		&models.DeviceToken{},
		&models.Analitics{},
		&models.Notificacion{},
		&models.Convocatoria{},
		&models.Sesion{},
		&models.Evento{},
		&models.UsuarioFacultadNombre{},
	); err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Modelos migrados")
}
