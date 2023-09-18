package repositories

import (
	"errors"
	"fmt"
	"noticias_uteq/models"
	"time"
)

// Registramos un Usuario
func RegisterUsuario(usuario *models.Usuario) error {
	exists, err := EmailExists(usuario.Email)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("El correo electrónico ya está registrado")
	}

	return DB.Save(&usuario).Error
}

// Verificamos que le correo no exista
func EmailExists(email string) (bool, error) {
	var count int64
	err := DB.Model(&models.Usuario{}).Where("email = ?", email).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Regitrsamos nuevo codigo
func CreateUserCode(userCode *models.UserCode) error {
	return DB.Create(&userCode).Error
}

/*Metodo para Iniciar Sesion*/
func Login(usuario *models.Usuario, deviceToken string) error {
	result := DB.First(&usuario, "email = ? AND password = ?", usuario.Email, usuario.Password)
	if result.RowsAffected < 1 {
		return errors.New("Usuario o contraseña no válido")
	}

	if !usuario.Verificado {
		return errors.New("Cuenta no verificada")
	}

	// Almacenar el token del dispositivo después de un inicio de sesión exitoso
	dToken := &models.DeviceToken{
		Token:     deviceToken,
		UsuarioID: usuario.ID,
	}
	DB.Save(&dToken)

	return nil
}

/*Metodo para Obtener Todos los Usuarios*/
func GetAllUsers() ([]models.Usuario, int, error) {
	var usuarios []models.Usuario
	result := DB.Find(&usuarios)

	// Comprueba si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return usuarios, len(usuarios), nil
}

/*Metodo para Obtener Todos los Usuarios*/
func GetUsuariosByRole(rol string) ([]models.Usuario, int, error) {
	var usuarios []models.Usuario
	result := DB.Where("rol = ?", rol).Find(&usuarios)

	// Comprueba si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return usuarios, len(usuarios), nil
}

/*Metodo para Actualizar un Usuario*/
func UpdateUsuario(id string, updates *models.Usuario) error {
	// Crea una estructura de usuario para almacenar el usuario existente
	var usuario models.Usuario

	// Busca al usuario en la base de datos usando solo el ID
	result := DB.First(&usuario, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}

	// Asigna los campos actualizados al usuario recuperado de la base de datos
	usuario.Password = updates.Password
	// No se cambian
	// usuario.Email = usuario.Email
	// usuario.Rol = usuario.Rol
	// usuario.Verificado = usuario.Verificado

	// Actualiza al usuario
	result = DB.Save(&usuario)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

/*Metodo para Eliminar un Usuario*/
func DeleteUsuario(id string) error {
	// Crea un objeto de usuario vacío
	usuario := models.Usuario{}

	// Busca al usuario en la base de datos utilizando su ID
	result := DB.First(&usuario, id)

	// Comprueba si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return result.Error
	}

	// Elimina al usuario
	result = DB.Delete(&usuario)

	// Comprueba si se produjo algún error durante la eliminación
	if result.Error != nil {
		return result.Error
	}

	return nil
}

/*Gestion de Codigos*/
func GetAllUserCodes() ([]models.UserCode, int, error) {
	var userCodes []models.UserCode
	result := DB.Find(&userCodes)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return userCodes, len(userCodes), nil
}

/*Metodo para Obtener un UserCode por UserID*/
func GetUserCodeByUserID(userID uint) (*models.UserCode, error) {
	var userCode models.UserCode
	result := DB.First(&userCode, "user_id = ?", userID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &userCode, nil
}

/*Metodo para Eliminar un UserCode*/
func DeleteUserCode(userCode *models.UserCode) error {
	result := DB.Delete(&userCode)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Obtener usuario por id
func GetUsuarioByID(id uint) (*models.Usuario, error) {
	var usuario models.Usuario
	result := DB.First(&usuario, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &usuario, nil
}

// Obtener usuario por correo electrónico
func GetUsuarioByEmail(email string) (*models.Usuario, error) {
	var usuario models.Usuario
	result := DB.Where("email = ?", email).First(&usuario)
	if result.Error != nil {
		return nil, result.Error
	}
	return &usuario, nil
}

/*Metodo para Obtener Todos los Tokens de los Dispositivos*/
func GetAllDeviceTokens() ([]models.DeviceToken, error) {
	var deviceTokens []models.DeviceToken
	result := DB.Find(&deviceTokens)

	// Comprueba si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return nil, result.Error
	}

	return deviceTokens, nil
}

func VerifyAndExpireUserCode1(userID uint, code string) (*models.UserCode, *models.Usuario, error) {
	var userCode models.UserCode
	result := DB.First(&userCode, "user_id = ? AND code = ?", userID, code)
	if result.Error != nil {
		return nil, nil, result.Error
	}
	if userCode.Status == "expired" || time.Now().After(userCode.ExpiryTime) {
		return nil, nil, errors.New("El código ha expirado")
	}

	var usuario models.Usuario
	if err := DB.First(&usuario, userID).Error; err != nil {
		return nil, nil, err
	}

	usuario.Verificado = true
	if err := DB.Save(&usuario).Error; err != nil {
		return nil, nil, err
	}

	userCode.Status = "expired"
	result = DB.Save(&userCode)
	if result.Error != nil {
		return nil, nil, result.Error
	}

	return &userCode, &usuario, nil
}

// Solicitar reseteo de contraseña: genera un código y lo almacena.
func RequestPasswordReset(email string) (*models.Usuario, error) {
	usuario := &models.Usuario{}
	if err := DB.First(usuario, "email = ?", email).Error; err != nil {
		return nil, errors.New("Usuario no encontrado.")
	}

	return usuario, nil
}

// Validar el código de reseteo
func ValidatePasswordResetCode(userID uint, code string) (bool, error) {
	var userCode models.UserCode
	if err := DB.First(&userCode, "user_id = ? AND code = ?", userID, code).Error; err != nil {
		return false, err
	}

	if userCode.Status != "active" || userCode.ExpiryTime.Before(time.Now()) {
		return false, errors.New("El código es inválido o ha expirado.")
	}

	return true, nil
}

// Actualizar la contraseña del usuario
func UpdateUserPassword(userID uint, newPassword string) error {
	usuario := &models.Usuario{}
	if err := DB.First(usuario, "id = ?", userID).Error; err != nil {
		return err
	}

	usuario.Password = newPassword
	if err := DB.Save(usuario).Error; err != nil {
		return err
	}

	return nil
}

// Función para guardar el UserCode en la base de datos
func SaveUserCode(userCode *models.UserCode) error {
	return DB.Save(userCode).Error
}

// Metodos para la sesion
func CreateSesion(sesion *models.Sesion) error {
	// Crear una nueva sesión
	return DB.Save(&sesion).Error
}

func GetSesionById(id uint) (*models.Sesion, error) {
	// Crear un objeto Sesion vacío
	var sesion models.Sesion

	// Buscar la sesión en la base de datos utilizando su ID
	result := DB.First(&sesion, id)

	// Comprobar si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return nil, result.Error
	}

	return &sesion, nil
}

func GetSesionByCodigo(codigo string) (*models.Sesion, error) {
	// Crear un objeto Sesion vacío
	var sesion models.Sesion

	// Buscar la sesión en la base de datos utilizando el código
	result := DB.First(&sesion, "codigo = ?", codigo)

	// Comprobar si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return nil, result.Error
	}

	return &sesion, nil
}

func UpdateSesion(id uint, sesion *models.Sesion) error {
	// Crear un objeto Sesion vacío
	var sesionToUpdate models.Sesion

	// Buscar la sesión en la base de datos
	result := DB.First(&sesionToUpdate, id)

	// Comprobar si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return result.Error
	}

	// Actualizar los campos necesarios en el objeto sesionToUpdate con los datos del objeto sesion
	sesionToUpdate.Email = sesion.Email
	sesionToUpdate.EmailVerified = sesion.EmailVerified
	sesionToUpdate.FamilyName = sesion.FamilyName
	sesionToUpdate.GivenName = sesion.GivenName
	sesionToUpdate.Locale = sesion.Locale
	sesionToUpdate.Name = sesion.Name
	sesionToUpdate.Picture = sesion.Picture
	sesionToUpdate.Sub = sesion.Sub
	sesionToUpdate.Codigo = sesion.Codigo

	// Guardar el objeto sesion actualizado
	result = DB.Save(&sesionToUpdate)

	// Comprobar si se produjo algún error durante la actualización
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteSesion(id uint) error {
	// Crear un objeto Sesion vacío
	sesion := models.Sesion{}

	// Buscar la sesión en la base de datos utilizando su ID
	result := DB.First(&sesion, id)

	// Comprobar si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return result.Error
	}

	// Eliminar la sesión
	result = DB.Delete(&sesion)

	// Comprobar si se produjo algún error durante la eliminación
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func GetAllSesiones() ([]models.Sesion, int, error) {
	var sesiones []models.Sesion
	result := DB.Find(&sesiones)

	// Comprobar si se produjo algún error durante la búsqueda
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return sesiones, len(sesiones), nil
}

func Logingoogle(codigo string, deviceToken string) (*models.Usuario, error) {
	// Obtén el correo electrónico de la sesión usando el código
	var sesion models.Sesion
	result := DB.First(&sesion, "codigo = ?", codigo)
	if result.Error != nil {
		fmt.Println("Error buscando sesión:", result.Error)
		return nil, errors.New("Código de sesión no válido")
	}
	if result.RowsAffected == 0 {
		fmt.Println("No se encontró ninguna sesión con ese código")
		return nil, errors.New("Código de sesión no encontrado")
	}

	// Busca el usuario usando el correo electrónico obtenido de la sesión
	var usuario models.Usuario
	result = DB.First(&usuario, "email = ?", sesion.Email)
	if result.Error != nil {
		fmt.Println("Error buscando usuario:", result.Error)
		return nil, errors.New("Usuario no encontrado salio zero resultado")
	}
	if result.RowsAffected == 0 {
		fmt.Println("No se encontró usuario con el correo electrónico:", sesion.Email)
		return nil, errors.New("Usuario no encontrado y el correo tampoco")
	}

	if !usuario.Verificado {
		fmt.Println("La cuenta no está verificada")
		return nil, errors.New("Cuenta no verificada")
	}

	// Almacenar el token del dispositivo después de un inicio de sesión exitoso
	dToken := &models.DeviceToken{
		Token:     deviceToken,
		UsuarioID: usuario.ID,
	}
	DB.Save(&dToken)

	return &usuario, nil
}
