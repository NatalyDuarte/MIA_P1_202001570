package main

import (
	"bufio"
	"fmt"
	"mimodulo/comandos"
	"os"
	"strings"
)

func main() {
	analizar()
}

func mens_error(err error) {
	fmt.Println("Error: ", err)
}

func analizar() {
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
				split_comando(comando)
			}
		}
	}
}

func split_comando(comando string) {
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
		ejecutar_comando(arre_coman)
	}
}

func ejecutar_comando(arre_coman []string) {
	data := strings.ToLower(arre_coman[0])

	if data == "mkdisk" {
		/*=======================MKDISK================== */
		comandos.Mkdisk(arre_coman)

	} else if data == "rep" {
		/*=======================REP===================== */
		//rep()
	} else if data == "execute" {
		/*=======================EXECUTE================= */
		//execute(arre_coman)
	} else {
		/*=======================ERROR=================== */
		fmt.Println("Error: El comando no fue reconocido.")
	}
}
