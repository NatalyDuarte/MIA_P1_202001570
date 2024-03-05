package comandos

import (
	"fmt"
	"os"
	"strings"
)

func Rep(arre_coman []string) {
	fmt.Println("==========================REP==========================")

	val_path := ""
	band_path := false

	val_id := ""
	band_id := false

	val_ruta := ""
	band_ruta := false

	val_name := ""
	band_name := false

	band_error := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "name="):
			if band_name {
				fmt.Println("Error: El parametro -name ya fue ingresado.")
				band_error = true
				break
			}
			val_name = strings.Replace(val_data, "\"", "", 2)
			val_name = strings.ToLower(val_name)
			if val_name == "mbr" || val_name == "disk" || val_name == "inode" || val_name == "journaling" {
				band_name = true
			} else if val_name == "block" || val_name == "bm_inode" || val_name == "bm_block" || val_name == "tree" {
				band_name = true
			} else if val_name == "sb" || val_name == "file" || val_name == "ls" {
				band_name = true
			} else {
				fmt.Println("Error: El Valor del parametro -name no es valido")
				band_error = true
				break
			}
		case strings.Contains(data, "path="):
			if band_path {
				fmt.Println("Error: El parametro -path ya fue ingresado.")
				band_error = true
				break
			}
			val_path = val_data
			band_path = true
		case strings.Contains(data, "id="):
			if band_id {
				fmt.Println("Error: El parametro -id ya fue ingresado.")
				band_error = true
				break
			}
			val_id = val_data
			band_id = true
		case strings.Contains(data, "ruta="):
			if band_ruta {
				fmt.Println("Error: El parametro -ruta ya fue ingresado.")
				band_error = true
				break
			}
			val_ruta = val_data
			band_ruta = true
			fmt.Print(val_ruta)
		default:
			fmt.Println("Error: Parametro no valido.")
		}
	}

	if !band_error {
		if band_path {
			if band_name {
				if band_id {
					if val_name == "mbr" {
						RepoMbr(val_path, val_id)
					}
				} else {
					fmt.Println("Error: No se encontro el id")
				}
			} else {
				fmt.Println("Error: No se encontro el name")
			}
		} else {
			fmt.Println("Error: No se encontro el path")
		}
	}
}

func RepoMbr(path string, id string) {
	//regex := regexp.MustCompile(`^[a-zA-Z]+`)
	//letters := regex.FindString(id)
	carpeta := ""
	archivo := ""

	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			carpeta = path[:i]
			archivo = path[i+1:]
			break
		}
	}

	fmt.Println("Carpeta:", carpeta)
	fmt.Println("Archivo:", archivo)
	fmt.Println("Disco:", id)

	_, err := os.Stat(carpeta)

	// Si la carpeta no existe, crearla
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				fmt.Println("Error al crear la carpeta:", err)
				return
			}
		} else {
			fmt.Println("Error al verificar la carpeta:", err)
			return
		}
	}

	fmt.Println("La carpeta", carpeta, "ya existe o se ha creado correctamente")

}
