package main

import (
	"bufio"
	"fmt"
	"io"
	"mimodulo/comandos"
	"mimodulo/estructuras"
	"os"
	"strconv"
	"strings"
)

var valorpath = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1/B.dsk"

func main() {
	Analizar()
}

func mens_error(err error) {
	fmt.Println("Error: ", err)
}

func Analizar() {
	finalizar := false
	fmt.Println(" ===============================================================")
	fmt.Println("                       Proyecto 1                               ")
	fmt.Println("               Nataly Guzman - 202001570                        ")
	fmt.Println("                    Bienvenido                                  ")
	fmt.Println(" ===============================================================")
	reader := bufio.NewReader(os.Stdin)
	for !finalizar {
		fmt.Print("Ingrese Comando: ")
		comando, _ := reader.ReadString('\n')
		if strings.Contains(comando, "exit") {
			finalizar = true
		} else if strings.Contains(comando, "EXIT") {
			finalizar = true
		} else {
			if comando != "" && comando != "exit\n" && comando != "EXIT\n" {
				Split_comando(comando)
			}
		}
	}
}

func Split_comando(comando string) {
	var arre_coman []string
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)
	band_comentario := false
	if strings.Contains(comando, "pause") {
		arre_coman = append(arre_coman, comando)
	} else if strings.Contains(comando, "PAUSE") {
		arre_coman = append(arre_coman, comando)
	} else if strings.Contains(comando, "#") {
		band_comentario = true
		fmt.Println(comando)
	} else {
		arre_coman = strings.Split(comando, " -")
	}

	if !band_comentario {
		Ejecutar_comando(arre_coman)
	}
}

func Ejecutar_comando(arre_coman []string) {
	data := strings.ToLower(arre_coman[0])

	if data == "mkdisk" {
		/*=======================MKDISK================== */
		comandos.Mkdisk(arre_coman)
	} else if data == "rmdisk" {
		/*=======================RMDISK================== */
		comandos.Rmdisk(arre_coman)
	} else if data == "fdisk" {
		/*=======================FDISK=================== */
		// faltan cosas
		comandos.Fdisk(arre_coman)
	} else if data == "rep" {
		/*=======================REP===================== */
		rep()
	} else if data == "execute" {
		/*=======================EXECUTE================= */
		Execute(arre_coman)
	} else if data == "mkfile" {
		/*========================PAUSE================== */
		//pause()
	} else if data == "mount" {
		/*========================MOUNT================== */
		comandos.Mount(arre_coman)
	} else if data == "pause" {
		/*========================PAUSE================== */
		pause()
	} else {
		/*=======================ERROR=================== */
		fmt.Println("Error: El comando no fue reconocido.")
	}
}

func pause() {
	fmt.Print("[MENSAJE] Presiona enter para continuar...")
	fmt.Scanln()
}

func Execute(arre_coman []string) {
	fmt.Print("==============EXECUTE=======================")
	val_path := ""

	band_path := false
	band_error := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "path="):
			if band_path {
				fmt.Println("Error: El parametro -path ya fue ingresado.")
				band_error = true
				break
			}

			band_path = true

			val_path = strings.Replace(val_data, "\"", "", 2)
		default:
			fmt.Println("Error: Parametro no valido.")
		}
	}

	if !band_error {
		if band_path {
			file, err := os.Open(*&val_path)
			if err != nil {
				fmt.Println("Error al abrir el archivo:", err)
				return
			}
			defer file.Close()
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				if line == "" {
					continue
				}
				if strings.HasPrefix(line, "#") {
					continue
				}
				fmt.Println("\n Comando:", line)
				Split_comando(line)
			}
			if err := scanner.Err(); err != nil {
				fmt.Println("Error al leer el archivo:", err)
			}
		}
	}
}

/* PARA MIENTRAS */
func rep() {
	var empty [100]byte
	mbr_empty := estructuras.Mbr{}
	//ebr_empty := estructuras.Ebr{}
	fmt.Println("=================REP===================")
	// Apertura de archivo
	disco, err := os.OpenFile(valorpath, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		mens_error(err)
	}

	// Calculo del tamano de struct en bytes
	mbr2 := comandos.Struct_a_bytes(mbr_empty)
	//ebr2 := comandos.Struct_a_bytes(ebr_empty)
	sstruct := len(mbr2)
	//srect := len(ebr2)

	// Lectrura del archivo binario desde el inicio
	lectura := make([]byte, sstruct)
	_, err = disco.ReadAt(lectura, 0)

	//lect := make([]byte, srect)
	//_, err = disco.ReadAt(lect, 0)
	// ERROR
	if err != nil && err != io.EOF {
		mens_error(err)
	}

	// Conversion de bytes a struct
	mbr := comandos.Bytes_a_struct_mbr(lectura)
	//ebr := comandos.Bytes_a_struct_ebr(lect)

	// ERROR
	if err != nil {
		mens_error(err)
	}

	if mbr.Mbr_tamano != empty {
		fmt.Print("Tamaño: ")
		fmt.Println(string(mbr.Mbr_tamano[:]))
		fmt.Print("Fecha: ")
		fmt.Println(string(mbr.Mbr_fecha_creacion[:]))
		fmt.Print("Signature: ")
		fmt.Println(string(mbr.Mbr_dsk_signature[:]))
		fmt.Print("Fit: ")
		fmt.Println(string(mbr.Dsk_fit[:]))
		for i := 0; i < 4; i++ {
			fmt.Println("Particion " + strconv.Itoa(i))
			fmt.Print("Status: ")
			fmt.Println(string(mbr.Mbr_partition[i].Part_status[:]))
			fmt.Print("Type: ")
			fmt.Println(string(mbr.Mbr_partition[i].Part_type[:]))
			fmt.Print("Name: ")
			fmt.Println(string(mbr.Mbr_partition[i].Part_name[:]))
			fmt.Print("Start: ")
			fmt.Println(string(mbr.Mbr_partition[i].Part_start[:]))
			fmt.Print("size: ")
			fmt.Println(string(mbr.Mbr_partition[i].Part_size[:]))
		}

	}
	disco.Close()
}
