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

var valorpath = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1/A.dsk"

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
		comandos.Rep(arre_coman)
		//Mostrar_mkdisk(valorpath)
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
	} else if data == "unmount" {
		/*========================PAUSE================== */
		comandos.Unmount(arre_coman)
	} else if data == "mkfs" {
		/*========================PAUSE================== */
		comandos.Mkfs(arre_coman)
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

func Mostrar_mkdisk(ruta string) {
	var empty [100]byte
	mbr_empty := estructuras.Mbr{}

	disco, err := os.OpenFile(ruta, os.O_RDWR, 0660)

	if err != nil {
		mens_error(err)
	}
	defer func() {
		disco.Close()
	}()

	mbr2 := comandos.Struct_a_bytes(mbr_empty)
	sstruct := len(mbr2)

	lectura := make([]byte, sstruct)
	_, err = disco.ReadAt(lectura, 0)

	if err != nil && err != io.EOF {
		mens_error(err)
	}
	mbr := comandos.Bytes_a_struct_mbr(lectura)

	if err != nil {
		mens_error(err)
	}

	if mbr.Mbr_tamano != empty {

		fmt.Print("Mbr_tamano: ")
		fmt.Println(string(mbr.Mbr_tamano[:]))
		fmt.Print("Mbr_fecha_creacion: ")
		fmt.Println(string(mbr.Mbr_fecha_creacion[:]))
		fmt.Print("Mbr_dsk_signature: ")
		fmt.Println(string(mbr.Mbr_dsk_signature[:]))
		fmt.Print("Dsk_fit: ")
		fmt.Println(string(mbr.Dsk_fit[:]))
		for i := 0; i < 4; i++ {
			fmt.Println("--------------------Particion: " + strconv.Itoa(i) + "--------------------")
			fmt.Print("Status: " + string(mbr.Mbr_partition[i].Part_status[:]) + " ")
			fmt.Print("Type: " + string(mbr.Mbr_partition[i].Part_type[:]) + " ")
			fmt.Print("Fit: " + string(mbr.Mbr_partition[i].Part_fit[:]) + " ")
			fmt.Print("part_start: " + string(mbr.Mbr_partition[i].Part_start[:]) + " ")
			fmt.Print("Part_s: " + string(mbr.Mbr_partition[i].Part_size[:]) + " ")
			fmt.Print("Name: " + string(mbr.Mbr_partition[i].Part_name[:]) + " ")
			fmt.Print("part_correlativo: " + string(mbr.Mbr_partition[i].Part_correlative[:]) + " ")
			fmt.Print("part_id: " + string(mbr.Mbr_partition[i].Part_id[:]) + " ")
			fmt.Println("--------------------")
			s_part_type := string(mbr.Mbr_partition[i].Part_type[:])
			s_part_type = strings.Trim(s_part_type, "\x00")
			if s_part_type == "e" {
				mostrar_todos_los_ebr(ruta, mbr.Mbr_partition[i])
			}

		}
	}
}

func mostrar_todos_los_ebr(ruta string, particion_extendida estructuras.Partition) {
	var empty [100]byte
	s_part_start := string(particion_extendida.Part_start[:])
	s_part_start = strings.Trim(s_part_start, "\x00")
	part_start, err := strconv.Atoi(s_part_start)
	if err != nil {
		mens_error(err)
	}
	ebr := obtener_ebr(ruta, part_start)
	bandera_bucle := false
	mbr2 := comandos.Struct_a_bytes(ebr)
	sstruct := len(mbr2)
	fmt.Println("Es aqui: ", sstruct)
	for !bandera_bucle {

		if ebr.Part_status != empty {
			fmt.Println("=====================================================================================================")
			fmt.Println("Part_mount: ", string(ebr.Part_status[:]), "Part_fit: ", string(ebr.Part_fit[:]), "Part_start: ", string(ebr.Part_start[:]), "Part_s: ", string(ebr.Part_size[:]), "Part_next: ", string(ebr.Part_next[:]), "Part_name: ", string(ebr.Part_name[:]))
			fmt.Println("=====================================================================================================")
		} else {
			break
		}
		s_part_next := string(ebr.Part_next[:])
		s_part_next = strings.Trim(s_part_next, "\x00")
		part_next, err := strconv.Atoi(s_part_next)
		if err != nil {
			mens_error(err)
		}
		if part_next == -1 {
			bandera_bucle = true
		} else {
			ebr = obtener_ebr(ruta, part_next)
		}

	}

}
func obtener_ebr(ruta string, part_start int) estructuras.Ebr {
	var empty [100]byte
	ebr_empty := estructuras.Ebr{}

	disco, err := os.OpenFile(ruta, os.O_RDWR, 0660)

	if err != nil {
		mens_error(err)
	}
	defer func() {
		disco.Close()
	}()

	ebr2 := comandos.Struct_a_bytes(ebr_empty)
	sstruct := len(ebr2)

	lectura := make([]byte, sstruct)
	_, err = disco.ReadAt(lectura, int64(part_start))

	if err != nil && err != io.EOF {
		mens_error(err)
	}

	ebr := comandos.Bytes_a_struct_ebr(lectura)

	if err != nil {
		mens_error(err)
	}

	if ebr.Part_next != empty {
		return ebr
	} else {
		return ebr_empty
	}
}
