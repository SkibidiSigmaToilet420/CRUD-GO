package main

import (
	"database/sql"
	/*"database/sql/driver"*/
	"fmt"

	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contraseña := ""
	Nombre := "crud"

	conexion, err := sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+Nombre)
	if err != nil {
		panic(err.Error())
	}
	return conexion

}

var platillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", Inicio)
	http.HandleFunc("/crear", Crear)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/borrar", Borrar)
	http.HandleFunc("/editar", Editar)
	http.HandleFunc("/actualizar", Actualizar)

	fmt.Println("Servidor corriendo...")
	fmt.Println("http://localhost:3000")
	http.ListenAndServe(":3000", nil)

}

func Borrar(w http.ResponseWriter, r *http.Request) {
	idProducto := r.URL.Query().Get("id")
	/* fmt.Println(idProducto) */

	conexionEstablecida := conexionBD()

	borrarRegistros, err := conexionEstablecida.Prepare("DELETE FROM productos WHERE id=?")

	if err != nil {
		panic(err.Error)
	}

	borrarRegistros.Exec(idProducto)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

type Productos struct {
	Id     int
	Nombre string
	Precio float64
	Estado string
}

func Inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()

	registros, err := conexionEstablecida.Query("SELECT * FROM productos")

	if err != nil {
		panic(err.Error)
	}
	producto := Productos{}
	arregloProducto := []Productos{}

	for registros.Next() {
		var id int
		var nombre, estado string
		var precio float64

		err = registros.Scan(&id, &nombre, &precio, &estado)
		if err != nil {
			panic(err.Error)
		}
		producto.Id = id
		producto.Nombre = nombre
		producto.Precio = precio
		producto.Estado = estado

		arregloProducto = append(arregloProducto, producto)

	}
	//fmt.Println(arregloProducto)

	/* fmt.Fprintf(w, "HOLA CONCHETUMARE") */
	platillas.ExecuteTemplate(w, "inicio", arregloProducto)
}

func Editar(w http.ResponseWriter, r *http.Request) {
	idProducto := r.URL.Query().Get("id")
	fmt.Println(idProducto)

	conexionEstablecida := conexionBD()
	registro,_ := conexionEstablecida.Query("SELECT * FROM productos WHERE id=?", idProducto)

	producto := Productos{}
	for registro.Next() {
		var id int
		var nombre, estado string
		var precio float64

		err := registro.Scan(&id, &nombre, &precio, &estado)
		if err != nil {
			panic(err.Error)
		}

		producto.Id = id
		producto.Nombre = nombre
		producto.Precio = precio
		producto.Estado = estado

		log.Printf("Actualizando empleado con ID: %s, Nombre: %s, estado: %s", fmt.Sprintf("%d", id), nombre, estado)

	}
	fmt.Println(producto)
	erro:=platillas.ExecuteTemplate(w, "editar",producto)
	if erro != nil {
        http.Error(w, erro.Error(), http.StatusInternalServerError)
        log.Println("Error al edtiar la página de crear los detalle:", erro)
    }


}

func Crear(w http.ResponseWriter, r *http.Request) {
	/* fmt.Fprintf(w, "HOLA CONCHETUMARE") */
	platillas.ExecuteTemplate(w, "crear", nil)
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		nombre := r.FormValue("nombre")
		precio := r.FormValue("precio")
		estado := r.FormValue("estado")

		conexionEstablecida := conexionBD()

		insertarRegistros, err := conexionEstablecida.Prepare("INSERT INTO productos(nombre,precio,estatus) VALUES(?,?,?) ")

		if err != nil {
			panic(err.Error)
		}

		insertarRegistros.Exec(nombre, precio, estado)

		http.Redirect(w, r, "/", http.StatusSeeOther)

	}

}

func Actualizar(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		// Aquí puedes manejar la solicitud POST
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		precio := r.FormValue("precio")
		estado := r.FormValue("estado")

		establecerConexion := conexionBD()

		modificar, erro := establecerConexion.Prepare("UPDATE productos SET nombre=?,precio=?,estatus=? WHERE id=?")
		if erro != nil {
			panic(erro.Error())
		}
		modificar.Exec(nombre,precio,estado, id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
	/* if r.Method == "POST" {

		id := r.FormValue("Id")
		nombre := r.FormValue("Nombre")
		precio := r.FormValue("Precio")
		estado := r.FormValue("Estado")

		conexionEstablecida := conexionBD()


		modificarRegistros, err := conexionEstablecida.Prepare("UPDATE productos SET Nombre=?,Precio=?,Estado=? WHERE id=? ")

		if err != nil {
			panic(err.Error)
		}

		modificarRegistros.Exec(nombre, precio, estado, id)

		http.Redirect(w, r, "/", http.StatusSeeOther)

	} */

}

