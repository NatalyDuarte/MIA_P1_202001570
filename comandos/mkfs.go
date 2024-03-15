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
		stmine := time.Now().Format("02/01/2006 15:04:05")
		//Superbloque
		copy(superbloque.S_filesystem_type[:], "2")
		copy(superbloque.S_inodes_count[:], []byte(strconv.Itoa(int(numero_estructuras))))
		copy(superbloque.S_blocks_count[:], []byte(strconv.Itoa(int(numero_estructuras*3))))
		copy(superbloque.S_free_blocks_count[:], []byte(strconv.Itoa(int((numero_estructuras*3)-5))))
		copy(superbloque.S_free_inodes_count[:], []byte(strconv.Itoa(int(numero_estructuras-5))))
		copy(superbloque.S_mtime[:], []byte(stmine))
		copy(superbloque.S_umtime[:], []byte(stmine))
		copy(superbloque.S_mnt_count[:], []byte("1"))
		copy(superbloque.S_magic[:], []byte("0xEF53"))
		copy(superbloque.S_inode_size[:], []byte(strconv.Itoa(int(unsafe.Sizeof(inodos)))))
		copy(superbloque.S_block_size[:], []byte(strconv.Itoa(block)))
		copy(superbloque.S_firts_ino[:], []byte("0"))
		copy(superbloque.S_first_blo[:], []byte("0"))
		startBitmapInodes := iniciopartis + int(unsafe.Sizeof(estructuras.Super_bloque{})) + 1
		copy(superbloque.S_bm_inode_start[:], []byte(strconv.Itoa(int(startBitmapInodes))))
		startBitmapBloks := startBitmapInodes + int(numero_estructuras) + 1
		copy(superbloque.S_bm_block_start[:], []byte(strconv.Itoa(int(startBitmapBloks))))
		firstInodeFree := startBitmapBloks + int(3*numero_estructuras) + 1
		firstBlockFree := firstInodeFree + int(n)*int(unsafe.Sizeof(estructuras.Inodo{})) + 1
		copy(superbloque.S_inode_start[:], []byte(strconv.Itoa(int(firstInodeFree))))
		copy(superbloque.S_block_start[:], []byte(strconv.Itoa((int(firstBlockFree)))))
		//Escribir Superbloque
		disco.Seek(int64(iniciopartis), 0)
		err = binary.Write(disco, binary.BigEndian, &superbloque)

		//Bitmap de inodos
		var buffer, buffer2 byte
		buffer = '0'
		buffer2 = '1'
		for i := 0; i < int(n); i++ {
			pos := int(startBitmapInodes) + i
			disco.Seek(int64(pos), 0)
			err = binary.Write(disco, binary.LittleEndian, &buffer)
			if err != nil {
				fmt.Println("Error: Error al escribir el bitmapInodes:", err)
			}
		}
		//insertamos los inodos usados para user.txt
		disco.Seek(int64(int(startBitmapInodes)), 0)
		err = binary.Write(disco, binary.LittleEndian, &buffer2)
		if err != nil {
			fmt.Println("Error: Error al escribir un inodo en bitmapInodes:", err)
		}
		disco.Seek(int64(int(startBitmapInodes)+1), 0)
		err = binary.Write(disco, binary.LittleEndian, &buffer2)
		if err != nil {
			fmt.Println("Error: Error al escribir un inodo en bitmapInodes:", err)
		}
		//Inodo para la carpeta root
		inodoRoot := estructuras.Inodo{}
		stmin := time.Now().Format("02/01/2006 15:04:05")
		copy(inodoRoot.I_uid[:], []byte("1"))
		copy(inodoRoot.I_gid[:], []byte("1"))
		copy(inodoRoot.I_size[:], []byte("27"))
		copy(inodoRoot.I_atime[:], []byte(stmin))
		copy(inodoRoot.I_ctime[:], []byte(stmin))
		copy(inodoRoot.I_mtime[:], []byte(stmin))
		inodoRoot.I_block[0] = 0
		for i := 1; i < 16; i++ {
			inodoRoot.I_block[i] = -1
		}
		copy(inodoRoot.I_type[:], []byte("0"))
		copy(inodoRoot.I_perm[:], []byte("664"))
		disco.Seek(int64(int(firstInodeFree)), 0)
		err = binary.Write(disco, binary.LittleEndian, &inodoRoot)
		if err != nil {
			fmt.Println("Error: Error al escribir un inodo:", err)
		}
		// bloque para la carpeta root
		blocUser := estructuras.Bloque_carpeta{}
		blocUser.B_content[0].B_inodo = 0
		blocUser.B_content[0].B_name[0] = '.'

		blocUser.B_content[1].B_inodo = 0
		blocUser.B_content[1].B_name[0] = '.'
		blocUser.B_content[1].B_name[1] = '.'

		blocUser.B_content[2].B_inodo = 1
		blocUser.B_content[2].B_name = ReturnValueArr("users.txt")

		blocUser.B_content[3].B_inodo = -1
		blocUser.B_content[3].B_name[0] = '.'

		disco.Seek(int64(int(firstBlockFree)), 0)
		err = binary.Write(disco, binary.LittleEndian, &blocUser)
		if err != nil {
			fmt.Println("Error: Error al escribir un bloque:", err)
			return
		}
		//inodo de tipo archivo para users.txt
		copy(inodoRoot.I_uid[:], []byte("1"))
		copy(inodoRoot.I_gid[:], []byte("1"))
		copy(inodoRoot.I_size[:], []byte("27"))
		copy(inodoRoot.I_atime[:], []byte(stmin))
		copy(inodoRoot.I_ctime[:], []byte(stmin))
		copy(inodoRoot.I_mtime[:], []byte(stmin))
		inodoRoot.I_block[0] = 1
		for i := 1; i < 16; i++ {
			inodoRoot.I_block[i] = -1
		}
		copy(inodoRoot.I_type[:], []byte("1"))
		copy(inodoRoot.I_perm[:], []byte("755"))
		disco.Seek(int64(int32(firstInodeFree)+int32(unsafe.Sizeof(estructuras.Inodo{}))), 0)
		err = binary.Write(disco, binary.LittleEndian, &inodoRoot)
		if err != nil {
			fmt.Println("Error: Error al escribir un inodo:", err)
		}
		//bloque archivo para users.txt
		blocUserT := estructuras.Bloque_archivo{}
		blocUserT.B_content = ReturnValueArr64("1,G,root\n1,U,root,root,123\n")
		disco.Seek(int64(int32(firstBlockFree)+int32(unsafe.Sizeof(estructuras.Bloque_archivo{}))), 0)
		err = binary.Write(disco, binary.LittleEndian, &blocUserT)
		if err != nil {
			fmt.Println("Error: Error al escribir un inodo:", err)
		}
		fmt.Println("Ext2 creado exitosamente")
	} else {
		fmt.Println("Error: No existe el id")
	}
	defer func() {
		disco.Close()
	}()
	disco.Close()
}

func ReturnValueArr(value string) [12]byte {
	var tmp [12]byte
	for i := 0; i < 12; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}
func ReturnValueArr64(value string) [64]byte {
	var tmp [64]byte
	for i := 0; i < 64; i++ {
		if i >= len(value) {
			break
		}
		tmp[i] = value[i]
	}
	return tmp
}
