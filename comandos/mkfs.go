package comandos

import (
	"encoding/binary"
	"fmt"
	"math"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

func Mkfs(arre_coman []string) {
	fmt.Println("=========================MKFS==========================")

	val_id := ""
	band_id := false

	val_type := "full"
	band_type := false

	val_fs := "2fs"
	band_fs := false

	band_error := false

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
			fmt.Println("Disco: " + letters + ".dsk")
			val_path = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1/" + letters + ".dsk"
			splitted := regexp.MustCompile(`^[a-zA-Z](\d+)`)
			match := splitted.FindStringSubmatch(val_id)

			if len(match) > 1 {
				numbe := match[1][0:1]
				number = numbe
			}
			fmt.Println("Particion: ", number)
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
		case strings.Contains(data, "fs="):
			if band_fs {
				fmt.Println("Error: El parametro -fs ya fue ingresado.")
				band_error = true
				break
			}
			if val_data == "2fs" {
				val_fs = "2fs"
				band_fs = true
			} else if val_data == "3fs" {
				val_fs = "3fs"
				band_fs = true
			} else {
				fmt.Println("Error: El valor de fs no es el indicado")
				band_error = true
				break
			}

		default:
			fmt.Println("Error: Parametro no valido.")
		}
	}
	if !band_error {
		if band_id {
			if val_fs == "2fs" {
				fmt.Println("Formateando con 2fs")
				Formatearext2(val_path, val_id)

			} else if val_fs == "3fs" {
				fmt.Println("Formateando con 3fs")
			}

		} else {
			fmt.Println("Error: No se ingreso el -id y es obligatorio")
		}

	}

}

func Formatearext2(val_path string, val_id string) {
	band_crea := false
	var empty [100]byte
	mbr := estructuras.Mbr{}
	disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}
	defer disco.Close()
	disco.Seek(0, 0)
	err = binary.Read(disco, binary.BigEndian, &mbr)
	if mbr.Mbr_tamano != empty {
		for i := 0; i < 4; i++ {
			id := string(mbr.Mbr_partition[i].Part_id[:])
			id = strings.Trim(id, "\x00")
			fmt.Println(id)
			fmt.Println(val_id)
			if id == val_id {
				posicion = i
				band_crea = true
			}
		}

	}
	if band_crea {
		TamanioParti := string(mbr.Mbr_partition[posicion].Part_size[:])
		TamanioParti = strings.Trim(TamanioParti, "\x00")
		tamanioparti, err := strconv.Atoi(TamanioParti)
		if err != nil {
			fmt.Println("Error: No se pudo convertir el tamaño de la particion")
		}
		InicioParti := string(mbr.Mbr_partition[posicion].Part_start[:])
		InicioParti = strings.Trim(InicioParti, "\x00")
		iniciopartis, err := strconv.Atoi(InicioParti)
		if err != nil {
			fmt.Println("Error: No se pudo convertir el inicio de la particion")
		}
		superbloque := estructuras.Super_bloque{}
		inodos := estructuras.Inodo{}
		block := 64
		n := (float64(tamanioparti) - float64(unsafe.Sizeof(superbloque))) / (4 + float64(unsafe.Sizeof(inodos)) + 3*float64(block))
		numero_estructuras := math.Floor(n)
		fmt.Println(numero_estructuras)
		stmine := time.Now().Format("02/01/2006 15:04:05")
		//Superbloque
		copy(superbloque.S_filesystem_type[:], "2fs")
		copy(superbloque.S_inodes_count[:], []byte(strconv.Itoa(int(numero_estructuras))))
		copy(superbloque.S_blocks_count[:], []byte(strconv.Itoa(int(numero_estructuras*3))))
		copy(superbloque.S_free_blocks_count[:], []byte(strconv.Itoa(int(numero_estructuras*3))))
		copy(superbloque.S_free_inodes_count[:], []byte(strconv.Itoa(int(numero_estructuras))))
		copy(superbloque.S_mtime[:], []byte(stmine))
		copy(superbloque.S_umtime[:], []byte("0"))
		copy(superbloque.S_mnt_count[:], []byte("1"))
		copy(superbloque.S_magic[:], []byte("0xEF53"))
		copy(superbloque.S_inode_size[:], []byte(strconv.Itoa(int(unsafe.Sizeof(inodos)))))
		copy(superbloque.S_block_size[:], []byte(strconv.Itoa(block)))
		copy(superbloque.S_firts_ino[:], []byte("0"))
		copy(superbloque.S_first_blo[:], []byte("0"))
		copy(superbloque.S_bm_inode_start[:], []byte(strconv.Itoa(int(unsafe.Sizeof(superbloque)))))
		copy(superbloque.S_bm_block_start[:], []byte(strconv.Itoa(int(unsafe.Sizeof(superbloque))+(int(numero_estructuras)))))
		copy(superbloque.S_inode_start[:], []byte(strconv.Itoa(int(unsafe.Sizeof(superbloque))+(int(numero_estructuras))+int(numero_estructuras*3))))
		copy(superbloque.S_block_start[:], []byte(strconv.Itoa(int(unsafe.Sizeof(superbloque))+(int(numero_estructuras))+int(numero_estructuras*3)+int(numero_estructuras*float64(unsafe.Sizeof(inodos))))))
		//Escribir Super Boot
		disco.Seek(int64(iniciopartis), 0)
		err = binary.Write(disco, binary.BigEndian, &superbloque)
		fmt.Println("Superbloque creado exitosamente")
	} else {
		fmt.Println("Error: No existe el id")
	}
	defer func() {
		disco.Close()
	}()
	disco.Close()
}
