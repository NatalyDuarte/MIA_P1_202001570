package comandos

import (
	"fmt"
	"io"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strings"
)

var valorpathmo string = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1/"

func Mount(arre_coman []string) {
	fmt.Println("=========================MOUNT==========================")

	val_path := ""
	band_path := true
	expresionRegular := regexp.MustCompile(`\d+`)
	numerosComoString := ""

	val_driveletter := ""
	val_name := ""
	termina := "70"
	band_driveletter := false
	band_name := false
	band_error := false
	band_enc := false

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
			val_driveletter = val_data
			band_driveletter = true
		case strings.Contains(data, "name="):
			if band_name {
				fmt.Println("Error: El parametro -name ya fue ingresado.")
				band_error = true
				break
			}
			val_name = val_data
			band_name = true
			coincidencias := expresionRegular.FindAllString(val_name, -1)
			numerosComoString = strings.Join(coincidencias, ",")
		default:
			fmt.Println("Error: Parametro no valido.")
		}
	}

	if !band_error {
		if band_path {
			if band_driveletter {
				if band_name {
					fmt.Println(val_driveletter + numerosComoString + termina)
					val_path = valorpathmo + val_driveletter + ".dsk"
					if archivoExiste(val_path) {
						var empty [100]byte
						mbr_empty := estructuras.Mbr{}
						disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

						// ERROR
						if err != nil {
							Mens_error(err)
						}

						// Calculo del tamano de struct en bytes
						mbr2 := Struct_a_bytes(mbr_empty)
						sstruct := len(mbr2)

						// Lectrura del archivo binario desde el inicio
						lectura := make([]byte, sstruct)
						_, err = disco.ReadAt(lectura, 0)

						// ERROR
						if err != nil && err != io.EOF {
							Mens_error(err)
						}

						// Conversion de bytes a struct
						mbr := Bytes_a_struct_mbr(lectura)

						// ERROR
						if err != nil {
							Mens_error(err)
						}
						if mbr.Mbr_tamano != empty {
							for i := 0; i < 4; i++ {
								name := string(mbr.Mbr_partition[i].Part_name[:])
								name = strings.Trim(name, "\x00")
								if name == val_name {
									band_enc = true
								}
							}

							if band_enc {
								fmt.Println("Encontrado")
							} else {
								fmt.Println("Error: No existe la particion en el disco")
							}
							disco.Close()
						}
						disco.Close()
					} else {
						fmt.Println("Error: El archivo no existe.")
					}

				} else {
					fmt.Println("Error: No se encontro el valor de name")
				}
			} else {
				fmt.Println("Error: No se encontro el valor driveletter")
			}
		} else {
			fmt.Println("Error: No se encontro el path")
		}
	}
}
func archivoExiste(ruta string) bool {
	_, err := os.Stat(ruta)
	return !os.IsNotExist(err)
}
