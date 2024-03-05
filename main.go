package main

import (
	"bufio"
	"fmt"
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

func rep() {
	extendida := -1
	mbr_empty := estructuras.Mbr{}
	ebr_empty := estructuras.Ebr{}
	var empty [100]byte
	fin := false

	f, err := os.OpenFile(valorpath, os.O_RDWR, 0660)

	if err == nil {
		mbr2 := comandos.Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		master_boot_record := comandos.Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			s_part_name := ""
			s_part_type := ""
			fmt.Print("Tamaño: ")
			fmt.Println(string(master_boot_record.Mbr_tamano[:]))
			fmt.Print("Fecha: ")
			fmt.Println(string(master_boot_record.Mbr_fecha_creacion[:]))
			fmt.Print("Signature: ")
			fmt.Println(string(master_boot_record.Mbr_dsk_signature[:]))
			fmt.Print("Fit: ")
			fmt.Println(string(master_boot_record.Dsk_fit[:]))
			for i := 0; i < 4; i++ {
				fmt.Println("--------------------------------Particion " + strconv.Itoa(i) + "-----------------------------------")
				fmt.Print("Status: ")
				fmt.Println(string(master_boot_record.Mbr_partition[i].Part_status[:]))
				fmt.Print("Type: ")
				fmt.Println(string(master_boot_record.Mbr_partition[i].Part_type[:]))
				fmt.Print("Name: ")
				fmt.Println(string(master_boot_record.Mbr_partition[i].Part_name[:]))
				fmt.Print("Start: ")
				fmt.Println(string(master_boot_record.Mbr_partition[i].Part_start[:]))
				fmt.Print("size: ")
				fmt.Println(string(master_boot_record.Mbr_partition[i].Part_size[:]))
				fmt.Print("correlativo: ")
				fmt.Println(string(master_boot_record.Mbr_partition[i].Part_correlative[:]))
				fmt.Print("id: ")
				fmt.Println(string(master_boot_record.Mbr_partition[i].Part_id[:]))

				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					extendida = i
				}
			}

			if extendida != -1 {
				s_part_start := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
				s_part_start = strings.Trim(s_part_start, "\x00")
				i_part_start, _ := strconv.Atoi(s_part_start)

				s_part_size := string(master_boot_record.Mbr_partition[extendida].Part_size[:])
				s_part_size = strings.Trim(s_part_size, "\x00")
				i_part_size, _ := strconv.Atoi(s_part_size)

				f.Seek(int64(i_part_start), os.SEEK_SET)

				ebr2 := comandos.Struct_a_bytes(ebr_empty)
				sstruct := len(ebr2)

				for !fin {
					lectura := make([]byte, sstruct)
					n_leidos, _ := f.Read(lectura)

					pos_actual, _ := f.Seek(0, os.SEEK_CUR)

					if n_leidos == 0 && (pos_actual >= int64(i_part_size+i_part_start)) {
						fin = true
						break
					}

					extended_boot_record := comandos.Bytes_a_struct_ebr(lectura)
					sstruct = len(lectura)

					if extended_boot_record.Part_size == empty {
						fin = true
						break
					} else {
						s_part_name = string(extended_boot_record.Part_name[:])
						s_part_name = strings.Trim(s_part_name, "\x00")
						fmt.Println("========LOGICAS========")
						fmt.Print("Status: ")
						fmt.Println(string(extended_boot_record.Part_status[:]))
						fmt.Print("Name: ")
						fmt.Println(string(extended_boot_record.Part_name[:]))
						fmt.Print("Start: ")
						fmt.Println(string(extended_boot_record.Part_start[:]))
						fmt.Print("size: ")
						fmt.Println(string(extended_boot_record.Part_size[:]))

						s_part_next := string(extended_boot_record.Part_next[:])
						s_part_next = strings.Trim(s_part_next, "\x00")

						if s_part_next != "-1" {
							f.Close()

						}
					}
				}
			}
		} else {
			fmt.Println("[ERROR] el disco se encuentra vacio...")
		}
	}

	f.Close()

}
