package comandos

import (
	"encoding/binary"
	"fmt"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unsafe"
)

func Login(arre_coman []string) {
	fmt.Println("=========================LOGIN==========================")

	val_rutadis := "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1/"
	band_path := true

	val_user := ""
	val_pass := ""
	val_id := ""

	band_user := false
	band_pass := false
	band_id := false

	band_error := false
	band_enc := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "user="):
			if band_user {
				fmt.Println("Error: El parametro -user ya fue ingresado.")
				band_error = true
				break
			}
			val_user = val_data
			band_user = true
		case strings.Contains(data, "pass="):
			if band_pass {
				fmt.Println("Error: El parametro -pass ya fue ingresado.")
				band_error = true
				break
			}
			val_pass = val_data
			band_pass = true
		case strings.Contains(data, "id="):
			if band_id {
				fmt.Println("Error: El parametro -id ya fue ingresado.")
				band_error = true
				break
			}
			val_id = val_data
			band_id = true
		default:
			fmt.Println("Error: Parametro no valido.")
		}
	}

	if !band_error {
		if band_path {
			if band_user {
				if band_pass {
					if band_id {
						regex := regexp.MustCompile(`^[a-zA-Z]+`)
						letters := regex.FindString(val_id)
						val_rutadis = val_rutadis + letters + ".dsk"
						fmt.Println("Disco:", letters)
						mbr := estructuras.Mbr{}
						sb := estructuras.Super_bloque{}
						disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)
						var empty [100]byte
						if err != nil {
							Mens_error(err)
						}
						defer func() {
							disco.Close()
						}()
						disco.Seek(0, 0)
						err = binary.Read(disco, binary.BigEndian, &mbr)
						if mbr.Mbr_tamano != empty {

						}
						numero_parti := 0
						for i := 0; i < 4; i++ {
							s_part_id := string(mbr.Mbr_partition[i].Part_id[:])
							s_part_id = strings.Trim(s_part_id, "\x00")

							if s_part_id == val_id {
								numero_parti = i
								band_enc = true
								break
							}

						}

						if band_enc {
							s_part_startas := string(mbr.Mbr_partition[numero_parti].Part_start[:])
							s_part_startas = strings.Trim(s_part_startas, "\x00")
							part_starta, err := strconv.Atoi(s_part_startas)
							if err != nil {
								Mens_error(err)
							}
							disco.Seek(int64(part_starta), 0)
							err = binary.Read(disco, binary.BigEndian, &sb)
							if sb.S_filesystem_type != empty {
								inodeA := estructuras.Inodo{}
								s_inodo_start := string(sb.S_inode_start[:])
								s_inodo_start = strings.Trim(s_inodo_start, "\x00")
								inodo_start, err := strconv.Atoi(s_inodo_start)
								if err != nil {
									Mens_error(err)
								}
								disco.Seek(int64(int32(inodo_start)+int32(unsafe.Sizeof(estructuras.Inodo{}))), 0)
								err = binary.Read(disco, binary.LittleEndian, &inodeA)
								if err != nil {
									fmt.Println("Error: Al leer un Inode en el archivo ")
								}
								fmt.Println(val_user)
								fmt.Println(val_pass)
							}
						}

					} else {
						fmt.Println("Error: No se encontro el valor id")
					}
				} else {
					fmt.Println("Error: No se encontro el valor pass")
				}

			} else {
				fmt.Println("Error: No se encontro el valor user")
			}
		} else {
			fmt.Println("Error: No se encontro el path")
		}
	}
}
