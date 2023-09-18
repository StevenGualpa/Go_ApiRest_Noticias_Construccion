package routers

import (
	"github.com/gofiber/fiber/v2"
	"noticias_uteq/handlers"
)

func SetupRouter(app *fiber.App) {

	//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
	//aqui vamos de nuevo aqui

	eventosRouter := app.Group("/eventos")
	handlersrouters := handlers.EventoHandler{}
	// Crear un nuevo evento
	eventosRouter.Post("Insert/", handlersrouters.CreateEvento)
	// Obtener todos los eventos
	eventosRouter.Get("GetAll/", handlersrouters.GetAllEventos)
	// Obtener un evento por ID
	eventosRouter.Get("GetById/:id", handlersrouters.GetEventoByID)
	// Actualizar un evento
	eventosRouter.Put("Update/:id", handlersrouters.UpdateEvento)
	// Eliminar un evento
	eventosRouter.Delete("Delete/:id", handlersrouters.DeleteEvento)
	// Obtener todos los eventos por fecha
	eventosRouter.Get("GetEventosByFecha/fecha/:fecha", handlersrouters.GetEventosByFecha)
	// Obtener todos los eventos por ID de usuario
	eventosRouter.Get("GetByIdUsuario/usuario/:idUsuario", handlersrouters.GetEventosByIDUsuario)

	noticiasRouter := app.Group("/noticias")
	handlersNoticias := handlers.NoticiasHandler{}

	noticiasRouter.Get("/byUsuario/:usuarioID", handlersNoticias.GetNoticiasByUsuario) // Nueva ruta
	noticiasRouter.Get("/GetAll", handlersNoticias.GetNoticias)
	//noticiasRouter.Post("/register", handlersNoticias.CreateNoticia)
	noticiasRouter.Post("/registerfiltro", handlersNoticias.CreateNoticiaFiltro)
	noticiasRouter.Post("/registerfiltroconvocatoria", handlersNoticias.CreateConvocatoriaFiltro)
	noticiasRouter.Delete("/delete/:ID", handlersNoticias.Delete)
	noticiasRouter.Get("/notificaciones/:usuarioID", handlersNoticias.GetNotificacionesByUsuario)
	noticiasRouter.Get("/convocatorias", handlersNoticias.GetAllConvocatorias)                           // Obtener todas las convocatorias
	noticiasRouter.Get("/convocatoria/:id", handlersNoticias.GetConvocatoriaByID)                        // Obtener una convocatoria por su ID
	noticiasRouter.Put("/updateconvocatoria/:id", handlersNoticias.UpdateConvocatoria)                   // Actualizar una convocatoria
	noticiasRouter.Get("/convocatoriafill/:id", handlersNoticias.GetConvocatoriasByUsuario)              // Obtener una convocatoria por su ID
	noticiasRouter.Put("/updatenotificacion/:notificacionID", handlersNoticias.UpdateNotificacionStatus) // Obtener una convocatoria por su ID

	multimediaRouter := app.Group("/multimedia")
	handlersMultimedia := handlers.MultimediaHandler{}
	multimediaRouter.Post("/register", handlersMultimedia.CreateMultimedia)
	multimediaRouter.Get("/getAll", handlersMultimedia.GetAllMultimedia)
	multimediaRouter.Get("/getAll/:contentType", handlersMultimedia.GetMultimediaByContentType)
	multimediaRouter.Put("/update/:id", handlersMultimedia.UpdateMultimedia)
	multimediaRouter.Delete("/delete/:id", handlersMultimedia.DeleteMultimedia)

	revistaRouter := app.Group("/revista")
	handlersRevista := handlers.RevistaHandler{}
	revistaRouter.Get("/", handlersRevista.GetRevistas)
	revistaRouter.Post("/", handlersRevista.PostRevista)
	revistaRouter.Delete("/:id", handlersRevista.DeleteRevista)

	/*Rutas para La Administracion de Usuarios*/
	authRouter := app.Group("/usuarios")
	handlersAuth := handlers.UsuarioHandler{}
	authRouter.Get("/getAllUsers", handlersAuth.GetUsers)
	authRouter.Get("/getRoleUsers/:role/", handlersAuth.GetUsersByRole)
	authRouter.Get("/getAllUsersCode", handlersAuth.GetAllUserCodes)
	authRouter.Post("/login", handlersAuth.Login)

	authRouter.Post("/registercorreocode", handlersAuth.RegisterUserCode)
	authRouter.Put("/updateUser/:id", handlersAuth.UpdateUser)
	authRouter.Put("/verifyUser/:id/:code", handlersAuth.VerifyUserCode)
	/*Reenviar Codigo*/
	authRouter.Put("/generateNewCode/:id", handlersAuth.GenerateNewCode)
	authRouter.Put("/ResetPassword/:email", handlersAuth.ResetPassword)
	/*Verificar Codigo*/
	authRouter.Post("/VerifyUserCodeCorreo", handlersAuth.VerifyUserCodeCorreo)
	authRouter.Post("/ChangePassword", handlersAuth.ChangePassword)

	authRouter.Post("/google/login", handlers.GoogleLoginHandler)
	//authRouter.Get("/google/login", handlers.GoogleLoginHandler)
	authRouter.Get("/google/callback", handlers.GoogleCallbackHandler)

	sesionRouter := app.Group("/sesiones")
	handlerssesion := handlers.AuthHandler{}
	sesionRouter.Get("/getAll", handlerssesion.GetUsersAuth)
	sesionRouter.Post("/google/signin", handlerssesion.GoogleSignInHandler)

	/*Actualizar Calve con correo*/

	authRouter.Delete("/deleteUser/:id", handlersAuth.DeleteUser)
	authRouter.Get("/deviceTokens", handlersAuth.GetAllDeviceTokens)

	/*Rutas para La Administracion de Facultades*/
	facudRouter := app.Group("/facultades")
	handlersFac := handlers.FacultadHandler{}
	facudRouter.Post("/register", handlersFac.CreateFacultad)
	facudRouter.Put("/update/:id", handlersFac.UpdateFacultad)
	facudRouter.Delete("/delete/:id", handlersFac.DeleteFacultad)
	facudRouter.Get("/getAll", handlersFac.GetAllFacultades)

	/*Rutas para La Administracion de Carreras*/
	carreraRouter := app.Group("/carreras")
	handlerscarrera := handlers.CarreraHandler{}
	carreraRouter.Post("/register", handlerscarrera.CreateCarrera)
	carreraRouter.Put("/update/:id", handlerscarrera.UpdateCarrera)
	carreraRouter.Delete("/delete/:id", handlerscarrera.DeleteCarrera)
	carreraRouter.Get("/getAll", handlerscarrera.GetAllCarreras)
	carreraRouter.Get("/getAll/:facultadID", handlerscarrera.GetCarrerasByFacultadID)

	/*Rutas para La Administracion de Preferencias*/
	preferenciaRouter := app.Group("/preferencias")
	handlerspreferencia := handlers.UsuarioFacultadHandler{}
	preferenciaRouter.Post("/register", handlerspreferencia.CreateUsuarioFacultad)
	preferenciaRouter.Get("/getAll/:usuarioID", handlerspreferencia.GetFacultadPreferencesByUsuario)
	preferenciaRouter.Post("/upsert", handlerspreferencia.UpsertUsuarioFacultad)

	preferenciaRouter.Post("/registerrp", handlerspreferencia.CreateUsuarioFacultadRp)
	preferenciaRouter.Get("/getAllrp/:usuarioID", handlerspreferencia.GetFacultadPreferencesByUsuarioRp)
	preferenciaRouter.Post("/upsertrp", handlerspreferencia.UpsertUsuarioFacultadRp)

	preferenciaRouter.Delete("/delete/:id", handlerspreferencia.DeleteUsuarioFacultad)

	/*Rutas para la parte estadistica*/
	estadisticRouter := app.Group("/estadisticas")
	handlersestadistics := handlers.AnaliticsHandler{}
	estadisticRouter.Post("/insert", handlersestadistics.InsertAnalitics)
	estadisticRouter.Get("/GetAll", handlersestadistics.GetAllAnalitics)
	estadisticRouter.Get("/GetAllSeccion/:seccion", handlersestadistics.GetAnaliticsBySeccion)
	estadisticRouter.Get("/GetAllFecha/:fecha", handlersestadistics.GetAnaliticsByFecha)
	estadisticRouter.Get("/GetAllSeccionFecha/:seccion/:fecha", handlersestadistics.GetAnaliticsBySeccionAndFecha)
	estadisticRouter.Get("/GetAllSeccionNombreCount/:seccion", handlersestadistics.GetNombreCountBySeccion)

	app.Get("/hola", func(c *fiber.Ctx) error {
		return c.SendString("Hola mundoo")
	})
}
