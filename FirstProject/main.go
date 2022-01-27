package main

import (
	"database/sql"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func ConexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contrasenia := ""
	Nombre := "connectionparking"

	conexion, err := sql.Open(Driver, Usuario+":"+Contrasenia+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion
}

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crearEmpleado", FormEmpleado)
	http.HandleFunc("/guardarEmpleado", GuardarEmpleado)
	http.HandleFunc("/editar", EditarEmpleado)
	http.HandleFunc("/eliminar", EliminarEmpleado)
	http.ListenAndServe(":8080", nil)
}

var templates = template.Must(template.ParseGlob("templates/*"))

var conexionEstablecida = ConexionBD()

type Empleado struct {
	IdPropietario string
	IdTipoDoc     string
	TipoDocumento string
	Documento     string
	Nombre        string
	Apellido      string
	Telefono      string
}

var empleado = Empleado{}
var rspEmpleado = []Empleado{}

func Inicio(response http.ResponseWriter, request *http.Request) {

	query := "SELECT p.IdPropietario, td.Nombre, p.Ndocumento,  p.Nombre, p.Apellido," +
		"p.Telefono FROM propietario p join tipo_documento td on p.TipoDocumento = td.IdTipo"
	sentenciaSQL, err := conexionEstablecida.Query(query)

	if err != nil {
		panic(err.Error())
	}
	for sentenciaSQL.Next() {
		err = sentenciaSQL.Scan(
			&empleado.IdPropietario, &empleado.TipoDocumento, &empleado.Documento,
			&empleado.Nombre, &empleado.Apellido, &empleado.Telefono)

		if err != nil {
			panic(err.Error())
		}
		rspEmpleado = append(rspEmpleado, empleado)
	}
	templates.ExecuteTemplate(response, "inicio", rspEmpleado)
	rspEmpleado = []Empleado{}
}

func FormEmpleado(response http.ResponseWriter, request *http.Request) {
	empleado = Empleado{}
	empleado.TipoDocumento = "Seleccione un tipo de documento"
	templates.ExecuteTemplate(response, "crearEmpleado", empleado)
}

func GuardarEmpleado(response http.ResponseWriter, request *http.Request) {

	if request.Method == "POST" {
		idPropietario := request.FormValue("idPropietario")
		tipoDocumento := request.FormValue("tipoDoc")
		documento := request.FormValue("identificación")
		nombre := request.FormValue("nombre")
		apellido := request.FormValue("apellido")
		telefono := request.FormValue("telefono")

		if idPropietario == "" {
			sentenciaSQL, err := conexionEstablecida.Prepare("INSERT INTO propietario(Ndocumento, TipoDocumento, Nombre, Apellido, Telefono) VALUES(?, ?, ?, ?, ?)")
			if err != nil {
				panic(err.Error())
			}
			sentenciaSQL.Exec(documento, tipoDocumento, nombre, apellido, telefono)
		} else {
			sentenciaSQL, err := conexionEstablecida.Prepare("UPDATE propietario SET Ndocumento=?, TipoDocumento=?, Nombre=?, Apellido=?, Telefono=? WHERE IdPropietario = ?")
			if err != nil {
				panic(err.Error())
			}
			sentenciaSQL.Exec(documento, tipoDocumento, nombre, apellido, telefono, idPropietario)
		}

		http.Redirect(response, request, "/", 301)
	}
}

func EliminarEmpleado(response http.ResponseWriter, request *http.Request) {

	idEmpleado := request.URL.Query().Get("id")

	sentenciaSQL, err := conexionEstablecida.Prepare("DELETE FROM propietario WHERE IdPropietario = ?")

	if err != nil {
		panic(err.Error())
	}

	sentenciaSQL.Exec(idEmpleado)
	http.Redirect(response, request, "/", 301)
}

func EditarEmpleado(response http.ResponseWriter, request *http.Request) {

	idEmpleado := request.URL.Query().Get("id")

	query := "SELECT p.IdPropietario, td.IdTipo, td.Nombre, p.Ndocumento,  p.Nombre, p.Apellido," +
		"p.Telefono FROM propietario p join tipo_documento td on p.TipoDocumento = td.IdTipo" +
		" WHERE IdPropietario = ?"

	sentenciaSQL, err := conexionEstablecida.Query(query, idEmpleado)

	if err != nil {
		panic(err.Error())
	}

	for sentenciaSQL.Next() {
		err = sentenciaSQL.Scan(
			&empleado.IdPropietario, &empleado.IdTipoDoc, &empleado.TipoDocumento, &empleado.Documento,
			&empleado.Nombre, &empleado.Apellido, &empleado.Telefono)

		if err != nil {
			panic(err.Error())
		}

	}

	templates.ExecuteTemplate(response, "crearEmpleado", empleado)

}
