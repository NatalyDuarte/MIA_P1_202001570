package comandos

import (
	"fmt"
	"io"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strings"
)

func Mkfs(arre_coman []string) {
	fmt.Println("=========================MKFS==========================")

	val_id := ""
	band_id := false

	val_type := ""
	band_type := false

	band_error := false
	band_crea := false
	regex := regexp.MustCompile(`^[a-zA-Z]+`)
	val_path := ""
	var number string

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "id="):
			if band_id {
				fmt.Println("Error: El parametro -id ya fue ingresado.")
				band_error = true
				break
			}
			val_id = val_data
			band_id = true
			letters := regex.FindString(val_id)
			val_path = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1/" + letters + ".dsk"
			splitted := regexp.MustCompile(`^[a-zA-Z](\d+)`)
			match := splitted.FindStringSubmatch(val_id)

			if len(match) > 1 {
				numbe := match[1][0:1]
				number = numbe
			}
			fmt.Println(number)
		case strings.Contains(data, "type="):
			if band_type {
				fmt.Println("Error: El parametro -type ya fue ingresado.")
				band_error = true
				break
			}
			val_type = val_data
			band_type = true
			if val_type != "full" {
				fmt.Println("Error: El type debe ser full")
				band_error = true
				break
			}

		default:
			fmt.Println("Error: Parametro no valido.")
		}
	}
	var empty [100]byte
	mbr_empty := estructuras.Mbr{}
	disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}
	mbr2 := Struct_a_bytes(mbr_empty)
	sstruct := len(mbr2)
	lectura := make([]byte, sstruct)
	_, err = disco.ReadAt(lectura, 0)
	if err != nil && err != io.EOF {
		Mens_error(err)
	}
	mbr := Bytes_a_struct_mbr(lectura)
	if err != nil {
		Mens_error(err)
	}

	if mbr.Mbr_tamano != empty {
		for i := 0; i < 4; i++ {
			id := string(mbr.Mbr_partition[i].Part_id[:])
			id = strings.Trim(id, "\x00")
			if id == val_id {
				posicion = i
				band_crea = true
			}
		}

	}

	if !band_error {
		if band_id {
			if band_crea {

			} else {
				fmt.Println("Error: No existe el id")
			}

		} else {
			fmt.Println("Error: No se ingreso el -id y es obligatorio")
		}

	}
	disco.Close()

}
