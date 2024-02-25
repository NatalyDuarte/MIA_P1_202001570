package comandos

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mimodulo/estructuras"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Mens_error(err error) {
	fmt.Println("Error: ", err)
}

var valorpath string = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1"
var path string = ""
var alfabeto []string = []string{
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j",
	"k", "l", "m", "n", "ñ", "o", "p", "q", "r", "s",
	"t", "u", "v", "w", "x", "y", "z",
}
var posicion = 0

func Mkdisk(arre_coman []string) {
	fmt.Println("=========================MKDISK==========================")

	val_size := 0
	val_fit := "ff"
	val_unit := "m"
	val_path := "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1"

	band_size := false
	band_fit := false
	band_unit := false
	band_path := true
	band_error := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "size="):
			if band_size {
				fmt.Println("Error: El parametro -size ya fue ingresado.")
				band_error = true
				break
			}

			band_size = true

			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size

			if err != nil {
				Mens_error(err)
			}

			if val_size < 0 {
				band_error = true
				fmt.Println("Error: El parametro -size es negativo.")
				break
			}
		case strings.Contains(data, "fit="):
			if band_fit {
				fmt.Println("Error: El parametro -fit ya fue ingresado.")
				band_error = true
				break
			}

			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" {
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				band_fit = true
				val_fit = "w"
			} else {
				fmt.Println("Error: El Valor del parametro -fit no es valido.")
				band_error = true
				break
			}
		case strings.Contains(data, "unit="):
			if band_unit {
				fmt.Println("Error: El parametro -unit ya fue ingresado.")
				band_error = true
				break
			}

			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			if val_unit == "k" || val_unit == "m" {
				band_unit = true
			} else {
				fmt.Println("Error: El Valor del parametro -unit no es valido.")
				band_error = true
				break
			}
		default:
			fmt.Println("Error: Parametro no valido.")
		}
	}

	if !band_error {
		if band_path {
			if band_size {
				total_size := 1024
				master_boot_record := estructuras.Mbr{}

				Crear_disco(val_path)

				fecha := time.Now()
				str_fecha := fecha.Format("02/01/2006 15:04:05")

				copy(master_boot_record.Mbr_fecha_creacion[:], str_fecha)

				rand.Seed(time.Now().UnixNano())
				min := 0
				max := 100
				num_random := rand.Intn(max-min+1) + min

				copy(master_boot_record.Mbr_dsk_signature[:], strconv.Itoa(int(num_random)))

				if band_unit {
					if val_unit == "m" {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
						total_size = val_size * 1024
					} else {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024)))
						total_size = val_size
					}
				} else {
					copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
					total_size = val_size * 1024
				}

				str_total_size := strconv.Itoa(total_size)

				cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=\""+path+"\" bs=1024 count="+str_total_size)
				cmd.Dir = "/"
				_, err := cmd.Output()

				if err != nil {
					fmt.Println("Aqu", err)
					Mens_error(err)
				}

				disco, err := os.OpenFile(path, os.O_RDWR, 0660)

				if err != nil {
					fmt.Println("Aqu1", err)
					Mens_error(err)
				}

				mbr_byte := Struct_a_bytes(master_boot_record)

				newpos, err := disco.Seek(0, os.SEEK_SET)

				if err != nil {
					fmt.Println("Aqu2", err)
					Mens_error(err)
				}

				_, err = disco.WriteAt(mbr_byte, newpos)

				if err != nil {
					Mens_error(err)
				}

				disco.Close()
			} else {
				fmt.Println("Error: No se encontro el valor size")
			}
		} else {
			fmt.Println("Error: No se encontro el path")
		}
	}
}

func Crear_disco(ruta string) {
	nome := NombreDisk()
	if nome == "error" {
		fmt.Println("Error: No se encontro el nombre del disco")
	} else {
		ruta = ruta + "/" + nome
	}
	path = ruta
	aux, err := filepath.Abs(ruta)

	if err != nil {
		fmt.Println("Aqui", err)
		Mens_error(err)
	}

	cmd1 := exec.Command("/bin/sh", "-c", "sudo mkdir -p '"+filepath.Dir(aux)+"'")
	cmd1.Dir = "/"
	_, err1 := cmd1.Output()

	if err1 != nil {
		fmt.Println("Aqui1", err)
		Mens_error(err)
	}

	cmd2 := exec.Command("/bin/sh", "-c", "sudo chmod -R 777 '"+filepath.Dir(aux)+"'")
	cmd2.Dir = "/"
	_, err2 := cmd2.Output()

	if err2 != nil {
		fmt.Println("Aqui2", err)
		Mens_error(err)
	}

	if _, err := os.Stat(filepath.Dir(aux)); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			fmt.Println("[FAILURE] No se pudo crear el disco.")
		}
	}
}
func Struct_a_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)

	if err != nil && err != io.EOF {
		Mens_error(err)
	}

	return buf.Bytes()
}

func NombreDisk() string {
	nombre := ""
	if posicion < 0 || posicion >= len(alfabeto) {
		nombre = "error"
	} else {
		nombre = alfabeto[posicion]
		nombre = nombre + ".dsk"
		posicion++
	}
	return nombre
}
