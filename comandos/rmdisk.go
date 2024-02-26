package comandos

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func Mens_errorrm(err error) {
	fmt.Println("Error: ", err)
}

var valorpathrm string = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1"
var pathrm string = ""

func Rmdisk(arre_coman []string) {
	fmt.Println("=========================RMDISK==========================")

	val_path := "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1"
	band_path := true

	val_driveletter := ""
	band_driveletter := false
	band_error := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "driveletter="):
			if band_driveletter {
				fmt.Println("Error: El parametro -driveletter ya fue ingresado.")
				band_error = true
				break
			}
			fmt.Println(val_data)
			val_driveletter = val_data
			band_driveletter = true
		default:
			fmt.Println("Error: Parametro no valido.")
		}
	}

	if !band_error {
		if band_path {
			if band_driveletter {
				val_path = val_path + "/" + val_driveletter + ".dsk"
				fmt.Println(val_path)
				_, e := os.Stat(val_path)
				if e != nil {
					if os.IsNotExist(e) {
						fmt.Println("Error: No existe el disco que desea eliminar.")
						band_path = false
					}
				} else {
					fmt.Print("Desea eliminar el disco [S/N]?: ")
					var opcion string
					fmt.Scanln(&opcion)
					if opcion == "s" || opcion == "S" {
						cmd := exec.Command("/bin/sh", "-c", "rm \""+val_path+"\"")
						cmd.Dir = "/"
						_, err := cmd.Output()

						if err != nil {
							Mens_errorrm(err)
						} else {
							fmt.Println("El Disco fue eliminado exitosamente")
						}
						band_path = false
					} else if opcion == "n" || opcion == "N" {
						fmt.Println("EL disco no fue eliminado")
						band_path = false
					} else {
						fmt.Println("Error: Opcion no valida intentalo de nuevo.")
					}
				}
			} else {
				fmt.Println("Error: No se encontro el valor driveletter")
			}
		} else {
			fmt.Println("Error: No se encontro el path")
		}
	}
}
