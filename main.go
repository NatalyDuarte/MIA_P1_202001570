package main

import (
	"bufio"
	"fmt"
	"mimodulo/comandos"
	"os"
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
		/*========================UNMOUNT================== */
		comandos.Unmount(arre_coman)
	} else if data == "mkfs" {
		/*========================MKFS================== */
		comandos.Mkfs(arre_coman)
	} else if data == "login" {
		/*========================MKFS================== */
		comandos.Login(arre_coman)
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
