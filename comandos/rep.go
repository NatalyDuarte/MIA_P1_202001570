package comandos

import (
	"fmt"
	"io"
	"log"
	"mimodulo/estructuras"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var val_rutadis string = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1/"

func Rep(arre_coman []string) {
	fmt.Println("==========================REP==========================")

	val_path := ""
	band_path := false

	val_id := ""
	band_id := false

	val_ruta := ""
	band_ruta := false

	val_name := ""
	band_name := false

	band_error := false

	for i := 1; i < len(arre_coman); i++ {
		aux_data := strings.SplitAfter(arre_coman[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		switch {
		case strings.Contains(data, "name="):
			if band_name {
				fmt.Println("Error: El parametro -name ya fue ingresado.")
				band_error = true
				break
			}
			val_name = strings.Replace(val_data, "\"", "", 2)
			val_name = strings.ToLower(val_name)
			if val_name == "mbr" || val_name == "disk" || val_name == "inode" || val_name == "journaling" {
				band_name = true
			} else if val_name == "block" || val_name == "bm_inode" || val_name == "bm_block" || val_name == "tree" {
				band_name = true
			} else if val_name == "sb" || val_name == "file" || val_name == "ls" {
				band_name = true
			} else {
				fmt.Println("Error: El Valor del parametro -name no es valido")
				band_error = true
				break
			}
		case strings.Contains(data, "path="):
			if band_path {
				fmt.Println("Error: El parametro -path ya fue ingresado.")
				band_error = true
				break
			}
			val_path = val_data
			band_path = true
		case strings.Contains(data, "id="):
			if band_id {
				fmt.Println("Error: El parametro -id ya fue ingresado.")
				band_error = true
				break
			}
			val_id = val_data
			band_id = true
		case strings.Contains(data, "ruta="):
			if band_ruta {
				fmt.Println("Error: El parametro -ruta ya fue ingresado.")
				band_error = true
				break
			}
			val_ruta = val_data
			band_ruta = true
			fmt.Print("Ruta: ", val_ruta)
		default:
			fmt.Println("Error: Parametro no valido.")
		}
	}

	if !band_error {
		if band_path {
			if band_name {
				if band_id {
					if val_name == "mbr" {
						RepoMbr(val_path, val_id)
					}
				} else {
					fmt.Println("Error: No se encontro el id")
				}
			} else {
				fmt.Println("Error: No se encontro el name")
			}
		} else {
			fmt.Println("Error: No se encontro el path")
		}
	}
}

func RepoMbr(path string, id string) {
	//regex := regexp.MustCompile(`^[a-zA-Z]+`)
	//letters := regex.FindString(id)
	carpeta := ""
	archivo := ""

	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			carpeta = path[:i]
			archivo = path[i+1:]
			break
		}
	}

	fmt.Println("Carpeta:", carpeta)
	fmt.Println("Archivo:", archivo)
	fmt.Println("Disco:", id)
	val_rutadis = val_rutadis + id + ".dsk"

	_, err := os.Stat(carpeta)

	// Si la carpeta no existe, crearla
	if err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(carpeta, 0755)
			if err != nil {
				fmt.Println("Error al crear la carpeta:", err)
				return
			}
		} else {
			fmt.Println("Error al verificar la carpeta:", err)
			return
		}
	}

	fmt.Println("La carpeta", carpeta, "ya existe o se ha creado correctamente")

	var buffer string
	buffer += "digraph G{\nsubgraph cluster{\nlabel=\"REPORTE DE MBR\"\ntbl[shape=box,label=<\n<table border='0' cellborder='1' cellspacing='0' width='300'  height='200' >\n"
	var empty [100]byte
	mbr_empty := estructuras.Mbr{}

	bandera_extendida := false
	numero_exten := 0

	disco, err := os.OpenFile(val_rutadis, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}
	defer func() {
		disco.Close()
	}()

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
		s_part_tamaño := string(mbr.Mbr_tamano[:])
		s_part_tamaño = strings.Trim(s_part_tamaño, "\x00")
		buffer += "<tr> <td width='150' bgcolor=\"pink\"><b>Mbr_tamano</b></td><td width='150' bgcolor=\"pink\">" + s_part_tamaño + "</td>  </tr>\n"
		s_part_crea := string(mbr.Mbr_fecha_creacion[:])
		s_part_crea = strings.Trim(s_part_crea, "\x00")
		buffer += "<tr>  <td bgcolor=\"pink\"><b>Mbr_fecha_creacion</b></td><td bgcolor=\"pink\">" + s_part_crea + "</td>  </tr>\n"
		s_part_sig := string(mbr.Mbr_dsk_signature[:])
		s_part_sig = strings.Trim(s_part_sig, "\x00")
		buffer += "<tr>  <td bgcolor=\"pink\"><b>Mbr_dsk_signature</b></td><td bgcolor=\"pink\">" + s_part_sig + "</td>  </tr>\n"
		s_part_fit := string(mbr.Dsk_fit[:])
		s_part_fit = strings.Trim(s_part_fit, "\x00")
		buffer += "<tr>  <td bgcolor=\"pink\"><b>Dsk_fit</b></td><td bgcolor=\"pink\">" + s_part_fit + "</td>  </tr>"
		for i := 0; i < 4; i++ {
			s_part_statas := string(mbr.Mbr_partition[i].Part_status[:])
			s_part_statas = strings.Trim(s_part_statas, "\x00")
			s_part_startas := string(mbr.Mbr_partition[i].Part_start[:])
			s_part_startas = strings.Trim(s_part_startas, "\x00")
			if s_part_statas != "E" && s_part_startas != "-1" {
				buffer += "<tr>"
				buffer += "<td colspan=\"2\" bgcolor=\"purple\">Particion</td>"
				buffer += "</tr>"
				s_part_stat := string(mbr.Mbr_partition[i].Part_status[:])
				s_part_stat = strings.Trim(s_part_stat, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_status" + "</b></td><td bgcolor=\"pink\">" + s_part_stat + "</td>  </tr>\n"
				s_part_type := string(mbr.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_type" + "</b></td><td bgcolor=\"pink\">" + s_part_type + "</td>  </tr>\n"
				s_part_fi := string(mbr.Mbr_partition[i].Part_fit[:])
				s_part_fi = strings.Trim(s_part_fi, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_fit" + "</b></td><td bgcolor=\"pink\">" + s_part_fi + "</td>  </tr>\n"
				s_part_st := string(mbr.Mbr_partition[i].Part_start[:])
				s_part_st = strings.Trim(s_part_st, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_start " + "</b></td><td bgcolor=\"pink\">" + s_part_st + "</td>  </tr>\n"
				s_part_siz := string(mbr.Mbr_partition[i].Part_size[:])
				s_part_siz = strings.Trim(s_part_siz, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_size" + "</b></td><td bgcolor=\"pink\">" + s_part_siz + "</td>  </tr>\n"
				s_part_name := string(mbr.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")
				buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_name" + "</b></td><td bgcolor=\"pink\">" + s_part_name + "</td>  </tr>\n"
				if s_part_type == "e" {
					numero_exten = i
					bandera_extendida = true
				}
			}

		}
		buffer += ("</table>\n>];\n}")
		if bandera_extendida {
			s_part_start := string(mbr.Mbr_partition[numero_exten].Part_start[:])
			s_part_start = strings.Trim(s_part_start, "\x00")
			part_start, err := strconv.Atoi(s_part_start)
			if err != nil {
				Mens_error(err)
			}
			ebr := obtener_ebr(val_rutadis, part_start)
			bandera_bucle := false

			u := 0
			buffer += ("subgraph cluster_1{\n label=\"Particiones Logicas dentro de la extendida\"\ntbl_1[shape=box, label=<\n<table border='0' cellborder='1' cellspacing='0'  width='300' height='160' >\n")
			for !bandera_bucle {
				if ebr.Part_status != empty {

					s_part_ebnexts := string(ebr.Part_next[:])
					s_part_ebnexts = strings.Trim(s_part_ebnexts, "\x00")
					if s_part_ebnexts != "-1" {
						buffer += "<tr>"
						buffer += "<td width='150' colspan=\"2\" bgcolor=\"magenta\">Particion</td>"
						buffer += "</tr>"
						s_part_ebfit := string(ebr.Part_fit[:])
						s_part_ebfit = strings.Trim(s_part_ebfit, "\x00")
						buffer += "<tr>  <td width='150' bgcolor=\"pink\"><b>Part_fit" + "</b></td><td width='150' bgcolor=\"pink\">" + s_part_ebfit + "</td>  </tr>\n"
						s_part_ebst := string(ebr.Part_start[:])
						s_part_ebst = strings.Trim(s_part_ebst, "\x00")
						buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_start " + "</b></td><td bgcolor=\"pink\">" + s_part_ebst + "</td>  </tr>\n"
						s_part_ebnext := string(ebr.Part_next[:])
						s_part_ebnext = strings.Trim(s_part_ebnext, "\x00")
						buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_next" + "</b></td><td bgcolor=\"pink\">" + s_part_ebnext + "</td>  </tr>\n"
						s_part_ebsiz := string(ebr.Part_size[:])
						s_part_ebsiz = strings.Trim(s_part_ebsiz, "\x00")
						buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_size " + "</b></td><td bgcolor=\"pink\">" + s_part_ebsiz + "</td>  </tr>\n"
						s_part_ebname := string(ebr.Part_name[:])
						s_part_ebname = strings.Trim(s_part_ebname, "\x00")
						buffer += "<tr>  <td bgcolor=\"pink\"><b>Part_name " + "</b></td><td bgcolor=\"pink\">" + s_part_ebname + "</td>  </tr>\n"
					}

					u += 1
				} else {
					break
				}
				s_part_next := string(ebr.Part_next[:])
				s_part_next = strings.Trim(s_part_next, "\x00")
				part_next, err := strconv.Atoi(s_part_next)
				if err != nil {
					Mens_error(err)
				}
				if part_next == -1 {
					bandera_bucle = true
				} else {
					ebr = obtener_ebr(val_rutadis, part_next)
				}

			}
			if u != 0 {
				buffer += "</table>\n>];\n}"
			}
		}
		buffer += "}\n"
		file, err2 := os.Create("Mbr.dot")
		if err2 != nil && !os.IsExist(err) {
			log.Fatal(err2)
		}
		defer file.Close()

		err = os.Chmod("Mbr.dot", 0777)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}

		_, err = file.WriteString(buffer)
		if err != nil && !os.IsExist(err) {
			log.Fatal(err)
		}
		cmd := exec.Command("dot", "-Tpng", "Mbr.dot", "-o", path)
		err = cmd.Run()
		if err != nil {
			fmt.Errorf("no se pudo generar la imagen: %v", err)
		}
		val_rutadis = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1/"

	}

}
