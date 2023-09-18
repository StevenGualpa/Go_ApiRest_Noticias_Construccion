package repositories

import (
	"context"
	"errors"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	"google.golang.org/api/option"
	"gorm.io/gorm"
	"noticias_uteq/models"
)

func Noticias() ([]models.Noticia, int, error) {
	var noticias []models.Noticia

	err := DB.Preload("Usuario").Preload("Tags").Find(&noticias).Error
	if err != nil {
		return noticias, 0, err
	}

	return noticias, len(noticias), nil
}

// este es metodo qye usamos
func GetNoticiasByUsuarioID(usuarioID uint) ([]models.Noticia, error) {
	var noticias []models.Noticia
	var usuarioFacultades []models.UsuarioFacultad

	// Buscar las preferencias del usuario
	DB.Where("usuario_id = ?", usuarioID).Find(&usuarioFacultades)

	// Si el usuario no tiene preferencias, retornar todas las noticias
	if len(usuarioFacultades) == 0 {
		DB.Find(&noticias)
		return noticias, nil
	}

	// Si el usuario tiene preferencias, filtrar las noticias basadas en esas preferencias
	for _, uf := range usuarioFacultades {
		var facultad models.Facultad
		DB.First(&facultad, uf.FacultadID)

		var filteredNoticias []models.Noticia
		DB.Where("categoria = ? OR categoria = 'Institucional'", facultad.Nombre).Find(&filteredNoticias)
		noticias = append(noticias, filteredNoticias...)
	}

	return noticias, nil
}

func RegisterNoticia(noticia *models.Noticia) error {
	return DB.Save(&noticia).Error
}

func DeleteNoticia(id int) error {
	noticia := models.Noticia{}
	noticia.ID = uint(id)

	// Eliminamos el registro de noticias
	return DB.Delete(&noticia).Error
}

type NoticiaRepository struct {
	DB *gorm.DB
}

func (r *NoticiaRepository) CreateNoticia(noticia *models.Noticia) error {
	result := r.DB.Create(noticia)
	return result.Error
}

var app *firebase.App

func init() {
	opt := option.WithCredentialsFile("credenciales/notificauteq-19631-firebase-adminsdk-2qwrx-f87ee5e850.json")
	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Errorf("error initializing app: %v", err)
	}
}

func CreateNoticia(noticia *models.Noticia) error {
	result := DB.Create(&noticia)
	if result.Error != nil {
		return errors.New("no se pudo crear la noticia")
	}
	return nil
}

func GetAllDeviceTokenss() ([]models.DeviceToken, error) {
	var deviceTokens []models.DeviceToken
	result := DB.Find(&deviceTokens)
	if result.Error != nil {
		return nil, errors.New("no se pudo obtener los tokens de los dispositivos")
	}
	return deviceTokens, nil
}

func SendNotification(noticia *models.Noticia) error {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return errors.New("error getting Messaging client")
	}

	deviceTokens, err := GetAllDeviceTokenss()
	if err != nil {
		return err
	}

	// Usamos un mapa para almacenar los tokens únicos
	uniqueTokens := make(map[string]bool)
	for _, deviceToken := range deviceTokens {
		uniqueTokens[deviceToken.Token] = true
	}

	// Enviamos la notificación solo a los tokens únicos
	for token := range uniqueTokens {
		message := &messaging.Message{
			Notification: &messaging.Notification{
				Title: "Nueva Noticia",
				Body:  noticia.Titulo,
			},
			Token: token,
		}

		response, err := client.Send(ctx, message)
		if err != nil {
			fmt.Printf("error sending message: %v\n", err)
			continue
		}

		fmt.Printf("Successfully sent message: %v\n", response)
	}

	return nil
}

func DeleteDeviceToken(token string) error {
	deviceToken := models.DeviceToken{}
	res := DB.Where("ID = ?", token).Delete(&deviceToken)
	if res.Error != nil {
		return res.Error
	}
	return nil
}

/*Filtro de Notficaciones*/

func SaveNotification(titulo string, contenido string, usuarioID uint) error {
	notif := models.Notificacion{
		Titulo:    titulo,
		Contenido: contenido,
		UsuarioID: usuarioID,
	}
	return DB.Create(&notif).Error
}

func SendNotificationFiltro(noticia *models.Noticia) error {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return errors.New("error getting Messaging client")
	}

	deviceTokens, err := GetAllDeviceTokenss()
	if err != nil {
		return err
	}

	processedTokens := make(map[string]bool)

	for _, deviceToken := range deviceTokens {
		if processedTokens[deviceToken.Token] {
			continue // Skip if the token has already been processed
		}

		shouldSend := false

		// Si la categoría es Institucional, enviar la notificación sin restricción
		if noticia.Categoria == "Institucional" {
			shouldSend = true
		} else {
			var usuarioFacultades []models.UsuarioFacultad
			DB.Where("usuario_id = ?", deviceToken.UsuarioID).Find(&usuarioFacultades)

			if len(usuarioFacultades) == 0 {
				shouldSend = true
			} else {
				for _, uf := range usuarioFacultades {
					var facultad models.Facultad
					DB.First(&facultad, uf.FacultadID)
					if noticia.Categoria == facultad.Nombre {
						shouldSend = true
						break
					}
				}
			}
		}

		if shouldSend {
			sendMsg(client, ctx, noticia, deviceToken.Token)
			processedTokens[deviceToken.Token] = true

			// Si la notificación se envía, también la guardamos en la base de datos
			err := SaveNotification("noticia", noticia.Titulo, deviceToken.UsuarioID)
			if err != nil {
				fmt.Println("Error al guardar la notificación:", err)
			}
		}
	}

	return nil
}

func sendMsg(client *messaging.Client, ctx context.Context, noticia *models.Noticia, token string) bool {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Nueva Noticia",
			Body:  noticia.Titulo,
		},
		Token: token,
	}

	_, err := client.Send(ctx, message)
	if err != nil {
		fmt.Printf("error sending message: %v\n", err)
		return false
	}
	return true
}

func CreateNoticiaFiltro(noticia *models.Noticia) error {
	result := DB.Create(&noticia)
	if result.Error != nil {
		return errors.New("no se pudo crear la noticia")
	}
	return nil
}

// aqui
func GetNotificacionesByUsuario(usuarioID uint) ([]models.Notificacion, error) {
	var notificaciones []models.Notificacion
	err := DB.Where("usuario_id = ? AND status = 'active'", usuarioID).Order("created_at desc").Find(&notificaciones).Error
	return notificaciones, err
}

func UpdateNotificacionStatusToDisable(notificacionID uint) error {
	notificacion := models.Notificacion{}
	err := DB.First(&notificacion, notificacionID).Error
	if err != nil {
		return err
	}

	notificacion.Status = "disable"
	err = DB.Save(&notificacion).Error
	return err
}

// Para convocatoria
func sendMsgconvocatoria(client *messaging.Client, ctx context.Context, convocatoria *models.Convocatoria, token string) bool {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: "Nueva Convocatoria",
			Body:  convocatoria.Titulo,
		},
		Token: token,
	}

	_, err := client.Send(ctx, message)
	if err != nil {
		fmt.Printf("error sending message: %v\n", err)
		return false
	}
	return true
}

func CreateConvocatoria(convocatoria *models.Convocatoria) error {
	result := DB.Create(&convocatoria)
	if result.Error != nil {
		return errors.New("no se pudo crear la convocatoria")
	}
	return nil
}

func SendNotificationFiltroconvocatgoria(convocatoria *models.Convocatoria) error {
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return errors.New("error getting Messaging client")
	}

	deviceTokens, err := GetAllDeviceTokenss()
	if err != nil {
		return err
	}

	processedTokens := make(map[string]bool)

	for _, deviceToken := range deviceTokens {
		if processedTokens[deviceToken.Token] {
			continue // Skip if the token has already been processed
		}

		shouldSend := false

		// Si el destinatario es "Institucional", enviar la notificación sin restricción
		if convocatoria.Destinatario == "Institucional" {
			shouldSend = true
		} else {
			var usuarioFacultades []models.UsuarioFacultad
			DB.Where("usuario_id = ?", deviceToken.UsuarioID).Find(&usuarioFacultades)

			if len(usuarioFacultades) == 0 {
				shouldSend = true
			} else {
				for _, uf := range usuarioFacultades {
					var facultad models.Facultad
					DB.First(&facultad, uf.FacultadID)
					if convocatoria.Destinatario == facultad.Nombre {
						shouldSend = true
						break
					}
				}
			}
		}

		if shouldSend {
			sendMsgconvocatoria(client, ctx, convocatoria, deviceToken.Token)
			processedTokens[deviceToken.Token] = true

			// Si la notificación se envía, también la guardamos en la base de datos
			err := SaveNotification("convocatoria", convocatoria.Titulo, deviceToken.UsuarioID)
			if err != nil {
				fmt.Println("Error al guardar la notificación:", err)
			}
		}
	}

	return nil
}

// Obtener todas las convocatorias
func GetAllConvocatorias() ([]models.Convocatoria, error) {
	var convocatorias []models.Convocatoria
	result := DB.Find(&convocatorias)
	if result.Error != nil {
		return nil, result.Error
	}
	return convocatorias, nil
}

// Obtener una convocatoria por su ID
func GetConvocatoriaByID(id uint) (*models.Convocatoria, error) {
	var convocatoria models.Convocatoria
	result := DB.First(&convocatoria, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &convocatoria, nil
}

// Actualizar una convocatoria por su ID
func UpdateConvocatoria(id uint, convocatoriaUpdates models.Convocatoria) error {
	var convocatoria models.Convocatoria
	result := DB.First(&convocatoria, id)
	if result.Error != nil {
		return result.Error
	}

	result = DB.Model(&convocatoria).Updates(convocatoriaUpdates)
	if result.Error != nil {
		return result.Error
	}

	// Reenviamos la convocatoria como una notificación
	// Aquí debes asegurarte de que el método SendNotificationFiltro pueda manejar objetos de tipo Convocatoria
	if err := SendNotificationFiltroconvocatgoria(&convocatoria); err != nil {
		fmt.Println("Error al reenviar la notificación:", err)
		// Aquí podrías querer manejar este error de una manera más específica.
	}

	return nil
}

func GetConvocatoriasByUsuarioID(usuarioID uint) ([]models.Convocatoria, error) {
	var convocatorias []models.Convocatoria
	var usuarioFacultades []models.UsuarioFacultad

	// Buscar las preferencias del usuario
	DB.Where("usuario_id = ?", usuarioID).Find(&usuarioFacultades)

	// Si el usuario no tiene preferencias, retornar todas las convocatorias
	if len(usuarioFacultades) == 0 {
		DB.Find(&convocatorias)
		return convocatorias, nil
	}

	// Si el usuario tiene preferencias, filtrar las convocatorias basadas en esas preferencias
	for _, uf := range usuarioFacultades {
		var facultad models.Facultad
		DB.First(&facultad, uf.FacultadID)

		var filteredConvocatorias []models.Convocatoria
		DB.Where("destinatario = ? OR destinatario = 'Institucional'", facultad.Nombre).Find(&filteredConvocatorias)
		convocatorias = append(convocatorias, filteredConvocatorias...)
	}

	return convocatorias, nil
}
