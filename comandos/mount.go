package comandos

import (
	"encoding/binary"
	"fmt"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strconv"
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
						mbr := estructuras.Mbr{}
						disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

						if err != nil {
							Mens_error(err)
						}
						disco.Seek(0, 0)
						err = binary.Read(disco, binary.BigEndian, &mbr)
						posicion := 0

						if err != nil {
							Mens_error(err)
						}
						if mbr.Mbr_tamano != empty {
							for i := 0; i < 4; i++ {
								name := string(mbr.Mbr_partition[i].Part_name[:])
								name = strings.Trim(name, "\x00")
								types := string(mbr.Mbr_partition[i].Part_type[:])
								types = strings.Trim(types, "\x00")
								status := string(mbr.Mbr_partition[i].Part_status[:])
								status = strings.Trim(status, "\x00")
								if name == val_name {
									if types == "p" {
										if status != "1" {
											band_enc = true
											posicion = i
											break
										} else {
											fmt.Println("Error: No se puede montar la particion ya que se encuentra montada")
										}

									}

								}
							}

							if band_enc {
								fmt.Println("Particion encontrada")
								copy(mbr.Mbr_partition[posicion].Part_status[:], "1")

								copy(mbr.Mbr_partition[posicion].Part_id[:], []byte(val_driveletter+numerosComoString+termina))

								disco.Seek(0, os.SEEK_SET)
								err = binary.Write(disco, binary.BigEndian, &mbr)

								if err != nil {
									Mens_error(err)
								}

								fmt.Println("Particion montada exitosamente")
								imprimir(val_path)

							} else {
								fmt.Println("Error: No existe la particion en el disco o la particion no es primaria o ya esta montada la particion")
							}
							disco.Close()
						}
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
func imprimir(ruta string) {
	var empty [100]byte
	mbr := estructuras.Mbr{}

	disco, err := os.OpenFile(ruta, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}
	defer func() {
		disco.Close()
	}()

	disco.Seek(0, 0)
	err = binary.Read(disco, binary.BigEndian, &mbr)

	if err != nil {
		Mens_error(err)
	}

	if mbr.Mbr_tamano != empty {
		for i := 0; i < 4; i++ {
			s_part_status := string(mbr.Mbr_partition[i].Part_status[:])
			s_part_status = strings.Trim(s_part_status, "\x00")

			if s_part_status == "1" {
				fmt.Println("--------------------Particion: " + strconv.Itoa(i) + "--------------------")
				fmt.Println("Part_Status: " + string(mbr.Mbr_partition[i].Part_status[:]) + " ")
				fmt.Println("Part_Type: " + string(mbr.Mbr_partition[i].Part_type[:]) + " ")
				fmt.Println("Part_Fit: " + string(mbr.Mbr_partition[i].Part_fit[:]) + " ")
				fmt.Println("Part_start: " + string(mbr.Mbr_partition[i].Part_start[:]) + " ")
				fmt.Println("Part_size: " + string(mbr.Mbr_partition[i].Part_size[:]) + " ")
				fmt.Println("Part_Name: " + string(mbr.Mbr_partition[i].Part_name[:]) + " ")
				fmt.Println("Part_correlativo: " + string(mbr.Mbr_partition[i].Part_correlative[:]) + " ")
				fmt.Println("Part_id: " + string(mbr.Mbr_partition[i].Part_id[:]) + " ")
			}
		}
	}
}
