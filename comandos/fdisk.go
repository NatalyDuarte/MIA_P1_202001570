package comandos

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
	"mimodulo/estructuras"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var valorpathfd string = "/home/nataly/Documentos/Mia lab/Proyecto1/MIA_P1_202001570/Discos/MIA/P1"

func Fdisk(commandArray []string) {
	fmt.Println("=========================FDISK==========================")

	val_size := 0
	val_unit := "k"
	val_path := ""
	val_type := "p"
	val_fit := "w"
	val_name := ""
	val_delete := ""
	val_add := 0

	band_size := false
	band_unit := false
	band_path := false
	band_type := false
	band_fit := false
	band_name := false
	band_error := false
	band_delete := false
	band_add := false

	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]
		switch {
		case strings.Contains(data, "size="):
			if band_size {
				fmt.Println("Error: El parametro -size ya fue ingresado")
				band_error = true
				break
			}
			band_size = true
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size
			if err != nil {
				Mens_error(err)
				band_error = true
			}
			if val_size < 0 {
				band_error = true
				fmt.Println("Error: El parametro -size es negativo")
				break
			}
		case strings.Contains(data, "driveletter="):
			if band_path {
				fmt.Println("Error: El parametro -driveletter ya fue ingresado")
				band_error = true
				break
			}
			band_path = true
			val_path = valorpathfd + "/" + val_data + ".dsk"
		case strings.Contains(data, "name="):
			if band_name {
				fmt.Println("Error: El parametro -name ya fue ingresado")
				band_error = true
				break
			}
			band_name = true
			val_name = strings.Replace(val_data, "\"", "", 2)
		case strings.Contains(data, "unit="):
			if band_unit {
				fmt.Println("Error: El parametro -unit ya fue ingresado")
				band_error = true
				break
			}
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)
			if val_unit == "b" || val_unit == "k" || val_unit == "m" {
				band_unit = true
			} else {
				fmt.Println("Error: El Valor del parametro -unit no es valido")
				band_error = true
				break
			}
		case strings.Contains(data, "type="):
			if band_type {
				fmt.Println("Error: El parametro -type ya fue ingresado")
				band_error = true
				break
			}
			val_type = strings.Replace(val_data, "\"", "", 2)
			val_type = strings.ToLower(val_type)
			if val_type == "p" || val_type == "e" || val_type == "l" {
				band_type = true
			} else {
				fmt.Println("Error: El Valor del parametro -type no es valido")
				band_error = true
				break
			}
		case strings.Contains(data, "fit="):
			if band_fit {
				fmt.Println("Error: El parametro -fit ya fue ingresado")
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
				fmt.Println("Error: El Valor del parametro -fit no es valido")
				band_error = true
				break
			}
		case strings.Contains(data, "delete="):
			if band_delete {
				fmt.Println("Error: El parametro -delete ya fue ingresado")
				band_error = true
				break
			}
			val_delete = strings.Replace(val_data, "\"", "", 2)
			val_delete = strings.ToLower(val_delete)

			if val_delete == "full" {
				band_delete = true
			} else {
				fmt.Println("Error: El Valor del parametro -delete no es valido")
				band_error = true
				break
			}
		case strings.Contains(data, "add="):
			if band_add {
				fmt.Println("Error: El parametro -add ya fue ingresado")
				band_error = true
				break
			}
			band_add = true
			aux_add, err := strconv.Atoi(val_data)
			val_add = aux_add
			fmt.Println("Add: ", val_add)
			if err != nil {
				Mens_error(err)
				band_error = true
			}

		default:
			fmt.Println("Error: Parametro no valido")
		}
	}
	pathc := existe_particion(val_path, val_name)
	if pathc {
		if band_delete {
			if band_path {
				if band_name {
					fmt.Print("Desea eliminar la particion [S/N]?: ")
					var opcion string
					fmt.Scanln(&opcion)
					if opcion == "s" || opcion == "S" {
						eliminar(val_path, val_name)
					} else if opcion == "n" || opcion == "N" {
						fmt.Println("La particion no fue eliminado")
					} else {
						fmt.Println("Error: Opcion no valida intentalo de nuevo.")
					}
				} else {
					fmt.Println("Error: El parametro -name no fue ingresado")
				}
			} else {
				fmt.Println("Error: El parametro -driveletter no fue ingresado")
			}
		} else if band_add {
			if band_path {
				if band_name {
					if band_unit {

					} else {
						fmt.Println("Error: El parametro -unit no fue ingresado")
					}
				} else {
					fmt.Println("Error: El parametro -name no fue ingresado")
				}
			} else {
				fmt.Println("Error: El parametro -driveletter no fue ingresado")
			}

		}

	} else if !pathc {
		if !band_error {
			if band_size {
				if band_path {
					if band_name {
						if band_type {
							if val_type == "p" {
								// Primaria
								crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
							} else if val_type == "e" {
								// Extendida
								crear_particion_extendida(val_path, val_name, val_size, val_fit, val_unit)
							} else {
								// Logica
								crear_particion_logixa(val_path, val_name, val_size, val_fit, val_unit)
							}
						} else {
							// Si no lo indica se tomara como Primaria
							crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
						}
					} else {
						fmt.Println("Error: El parametro -name no fue ingresado")
					}
				} else {
					fmt.Println("Error: El parametro -driveletter no fue ingresado")
				}
			} else {
				fmt.Println("Error: El parametro -size no fue ingresado")
			}

		}
	} else {
		fmt.Println("Error: El nombre de la particion ya existe")
	}
}
func crear_particion_primaria(direccion string, nombre string, size int, fit string, unit string) {
	aux_unit := ""
	aux_path := direccion
	size_bytes := 1024
	expresionRegular := regexp.MustCompile(`\d+`)
	numerosComoString := ""
	aux_fit := ""
	mbr_empty := estructuras.Mbr{}
	var empty [100]byte

	if fit != "" {
		aux_fit = fit
	} else {
		aux_fit = "w"
	}

	if unit != "" {
		aux_unit = unit
		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			size_bytes = size * 1024
		} else {
			size_bytes = size * 1024 * 1024
		}
	} else {
		size_bytes = size * 1024
	}
	f, err := os.OpenFile(aux_path, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	} else {
		band_particion := false
		num_particion := 0
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)
		lectura := make([]byte, sstruct)
		_, err = f.ReadAt(lectura, 0)

		if err != nil && err != io.EOF {
			Mens_error(err)
		}

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if err != nil {
			Mens_error(err)
		}

		if master_boot_record.Mbr_tamano != empty {
			s_part_start := ""

			for i := 0; i < 4; i++ {
				s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
				s_part_start = strings.Trim(s_part_start, "\x00")

				if s_part_start == "-1" {
					band_particion = true
					num_particion = i
					break
				}
			}

			if band_particion {
				espacio_usado := 0
				for i := 0; i < 4; i++ {
					s_size := string(master_boot_record.Mbr_partition[i].Part_size[:])
					s_size = strings.Trim(s_size, "\x00")
					i_size, err := strconv.Atoi(s_size)

					s_part_status := string(master_boot_record.Mbr_partition[i].Part_status[:])
					s_part_status = strings.Trim(s_part_status, "\x00")

					if err != nil {
						Mens_error(err)
					}

					if s_part_status != "1" {
						espacio_usado += i_size
					}
				}

				s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
				s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
				i_tamaño_disco, err2 := strconv.Atoi(s_tamaño_disco)

				if err2 != nil {
					Mens_error(err)
				}

				espacio_disponible := i_tamaño_disco - espacio_usado

				if espacio_disponible >= size_bytes {
					s_dsk_fit := string(master_boot_record.Dsk_fit[:])
					s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")
					coincidencias := expresionRegular.FindAllString(nombre, -1)
					numerosComoString = strings.Join(coincidencias, ",")
					if s_dsk_fit == "f" {
						copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
						copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "p")
						copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

						if num_particion == 0 {
							mbr_empty_byte := Struct_a_bytes(mbr_empty)
							copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))))
						} else {
							s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
							s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
							i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

							s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
							s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
							i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

							copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant))
						}

						copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
						copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], nombre)
						copy(master_boot_record.Mbr_partition[num_particion].Part_correlative[:], []byte(numerosComoString))
						copy(master_boot_record.Mbr_partition[num_particion].Part_id[:], "")
						mbr_byte := Struct_a_bytes(master_boot_record)
						f.Seek(0, os.SEEK_SET)
						f.Write(mbr_byte)

						s_part_start = string(master_boot_record.Mbr_partition[num_particion].Part_start[:])
						s_part_start = strings.Trim(s_part_start, "\x00")
						i_part_start, _ := strconv.Atoi(s_part_start)

						f.Seek(int64(i_part_start), os.SEEK_SET)

						for i := 1; i < size_bytes; i++ {
							f.Write([]byte{1})
						}

						fmt.Println("La Particion primaria fue creada exitosamente")
					} else if s_dsk_fit == "b" {
						best_index := num_particion
						s_part_start_act := ""
						s_part_status_act := ""
						s_part_size_act := ""
						i_part_size_act := 0
						s_part_start_best := ""
						i_part_start_best := 0
						s_part_start_best_ant := ""
						i_part_start_best_ant := 0
						s_part_size_best := ""
						i_part_size_best := 0
						s_part_size_best_ant := ""
						i_part_size_best_ant := 0

						for i := 0; i < 4; i++ {
							s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
							s_part_start_act = strings.Trim(s_part_start_act, "\x00")

							s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
							s_part_status_act = strings.Trim(s_part_status_act, "\x00")

							s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
							s_part_size_act = strings.Trim(s_part_size_act, "\x00")
							i_part_size_act, _ = strconv.Atoi(s_part_size_act)

							if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
								if i != num_particion {
									s_part_size_best = string(master_boot_record.Mbr_partition[best_index].Part_size[:])
									s_part_size_best = strings.Trim(s_part_size_best, "\x00")
									i_part_size_best, _ = strconv.Atoi(s_part_size_best)

									s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])

									s_part_size_act = strings.Trim(s_part_size_act, "\x00")
									i_part_size_act, _ = strconv.Atoi(s_part_size_act)

									if i_part_size_best > i_part_size_act {
										best_index = i
										break
									}
								}
							}
						}

						copy(master_boot_record.Mbr_partition[best_index].Part_type[:], "p")
						copy(master_boot_record.Mbr_partition[best_index].Part_fit[:], aux_fit)

						if best_index == 0 {
							mbr_empty_byte := Struct_a_bytes(mbr_empty)
							copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))))
						} else {
							s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])
							s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
							i_part_start_best_ant, _ = strconv.Atoi(s_part_start_best_ant)

							s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])
							s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
							i_part_size_best_ant, _ = strconv.Atoi(s_part_size_best_ant)

							copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant))
						}

						copy(master_boot_record.Mbr_partition[best_index].Part_size[:], strconv.Itoa(size_bytes))
						copy(master_boot_record.Mbr_partition[best_index].Part_status[:], "0")
						copy(master_boot_record.Mbr_partition[best_index].Part_name[:], nombre)
						copy(master_boot_record.Mbr_partition[num_particion].Part_correlative[:], []byte(numerosComoString))
						copy(master_boot_record.Mbr_partition[num_particion].Part_id[:], "")

						mbr_byte := Struct_a_bytes(master_boot_record)

						f.Seek(0, os.SEEK_SET)
						f.Write(mbr_byte)

						s_part_start_best = string(master_boot_record.Mbr_partition[best_index].Part_start[:])

						s_part_start_best = strings.Trim(s_part_start_best, "\x00")
						i_part_start_best, _ = strconv.Atoi(s_part_start_best)

						f.Seek(int64(i_part_start_best), os.SEEK_SET)

						for i := 1; i < size_bytes; i++ {
							f.Write([]byte{1})
						}

						fmt.Println("La Particion primaria fue creada exitosamente")
					} else {
						worst_index := num_particion

						s_part_start_act := ""
						s_part_status_act := ""
						s_part_size_act := ""
						i_part_size_act := 0
						s_part_start_worst := ""
						i_part_start_worst := 0
						s_part_start_worst_ant := ""
						i_part_start_worst_ant := 0
						s_part_size_worst := ""
						i_part_size_worst := 0
						s_part_size_worst_ant := ""
						i_part_size_worst_ant := 0

						for i := 0; i < 4; i++ {
							s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])
							s_part_start_act = strings.Trim(s_part_start_act, "\x00")

							s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])
							s_part_status_act = strings.Trim(s_part_status_act, "\x00")

							s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
							s_part_size_act = strings.Trim(s_part_size_act, "\x00")
							i_part_size_act, _ = strconv.Atoi(s_part_size_act)

							if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
								if i != num_particion {
									s_part_size_worst = string(master_boot_record.Mbr_partition[worst_index].Part_size[:])
									s_part_size_worst = strings.Trim(s_part_size_worst, "\x00")
									i_part_size_worst, _ = strconv.Atoi(s_part_size_worst)

									s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])
									s_part_size_act = strings.Trim(s_part_size_act, "\x00")
									i_part_size_act, _ = strconv.Atoi(s_part_size_act)

									if i_part_size_worst < i_part_size_act {
										worst_index = i
										break
									}
								}
							}
						}

						copy(master_boot_record.Mbr_partition[worst_index].Part_type[:], "p")
						copy(master_boot_record.Mbr_partition[worst_index].Part_fit[:], aux_fit)

						if worst_index == 0 {
							mbr_empty_byte := Struct_a_bytes(mbr_empty)
							copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))))
						} else {
							s_part_start_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_start[:])

							s_part_start_worst_ant = strings.Trim(s_part_start_worst_ant, "\x00")
							i_part_start_worst_ant, _ = strconv.Atoi(s_part_start_worst_ant)

							s_part_size_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_size[:])
							s_part_size_worst_ant = strings.Trim(s_part_size_worst_ant, "\x00")
							i_part_size_worst_ant, _ = strconv.Atoi(s_part_size_worst_ant)

							copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(i_part_start_worst_ant+i_part_size_worst_ant))
						}

						copy(master_boot_record.Mbr_partition[worst_index].Part_size[:], strconv.Itoa(size_bytes))
						copy(master_boot_record.Mbr_partition[worst_index].Part_status[:], "0")
						copy(master_boot_record.Mbr_partition[worst_index].Part_name[:], nombre)
						copy(master_boot_record.Mbr_partition[num_particion].Part_correlative[:], []byte(numerosComoString))
						copy(master_boot_record.Mbr_partition[num_particion].Part_id[:], "")

						mbr_byte := Struct_a_bytes(master_boot_record)

						f.Seek(0, os.SEEK_SET)
						f.Write(mbr_byte)

						s_part_start_worst = string(master_boot_record.Mbr_partition[worst_index].Part_start[:])
						s_part_start_worst = strings.Trim(s_part_start_worst, "\x00")
						i_part_start_worst, _ = strconv.Atoi(s_part_start_worst)

						f.Seek(int64(i_part_start_worst), os.SEEK_SET)

						for i := 1; i < size_bytes; i++ {
							f.Write([]byte{1})
						}

						fmt.Println("La Particion primaria fue creada exitosamente")
					}
				}

				if err != nil {
					Mens_error(err)
				}

			} else {
				fmt.Println("Se excede de las cuatro particiones")
			}
		} else {
			fmt.Println("Error: El disco se encuentra vacio")
		}
		f.Close()
	}
}

func Bytes_a_struct_mbr(s []byte) estructuras.Mbr {
	p := estructuras.Mbr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		Mens_error(err)
	}

	return p
}

func crear_particion_extendida(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024

	mbr_empty := estructuras.Mbr{}
	ebr_empty := estructuras.Ebr{}
	var empty [100]byte

	if fit != "" {
		aux_fit = fit
	} else {
		aux_fit = "w"
	}

	if unit != "" {
		aux_unit = unit

		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			size_bytes = size * 1024
		} else {
			size_bytes = size * 1024 * 1024
		}
	} else {
		size_bytes = size * 1024
	}

	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	} else {
		band_particion := false
		band_extendida := false
		num_particion := 0

		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			s_part_type := ""

			for i := 0; i < 4; i++ {
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					band_extendida = true
					break
				}
			}

			if !band_extendida {
				s_part_start := ""
				s_part_status := ""
				s_part_size := ""
				i_part_size := 0

				for i := 0; i < 4; i++ {
					s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
					s_part_start = strings.Trim(s_part_start, "\x00")

					s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
					s_part_status = strings.Trim(s_part_status, "\x00")

					s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ = strconv.Atoi(s_part_size)

					if s_part_start == "-1" || (s_part_status == "1" && i_part_size >= size_bytes) {
						band_particion = true
						num_particion = i
						break
					}
				}

				if band_particion {
					espacio_usado := 0

					for i := 0; i < 4; i++ {
						s_part_status = string(master_boot_record.Mbr_partition[i].Part_status[:])
						s_part_status = strings.Trim(s_part_status, "\x00")

						if s_part_status != "1" {
							s_part_size = string(master_boot_record.Mbr_partition[i].Part_size[:])

							s_part_size = strings.Trim(s_part_size, "\x00")
							i_part_size, _ = strconv.Atoi(s_part_size)

							espacio_usado += i_part_size
						}
					}

					s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
					s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
					i_tamaño_disco, _ := strconv.Atoi(s_tamaño_disco)

					espacio_disponible := i_tamaño_disco - espacio_usado

					if espacio_disponible >= size_bytes {
						s_dsk_fit := string(master_boot_record.Dsk_fit[:])
						s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

						if s_dsk_fit == "f" {
							copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "e")
							copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

							if num_particion == 0 {
								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))))
							} else {
								s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])

								s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
								i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

								s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])

								s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
								i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant))
							}

							copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], nombre)

							mbr_byte := Struct_a_bytes(master_boot_record)

							f.Seek(0, os.SEEK_SET)
							f.Write(mbr_byte)

							s_part_start = string(master_boot_record.Mbr_partition[num_particion].Part_start[:])

							s_part_start = strings.Trim(s_part_start, "\x00")
							i_part_start, _ := strconv.Atoi(s_part_start)

							f.Seek(int64(i_part_start), os.SEEK_SET)

							extended_boot_record := estructuras.Ebr{}
							copy(extended_boot_record.Part_fit[:], aux_fit)
							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_start[:], s_part_start)
							copy(extended_boot_record.Part_size[:], "0")
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], "")
							ebr_byte := Struct_a_bytes(extended_boot_record)
							f.Write(ebr_byte)

							ebr_empty_byte := Struct_a_bytes(ebr_empty)

							pos_actual, _ := f.Seek(0, os.SEEK_CUR)
							f.Seek(int64(pos_actual), os.SEEK_SET)

							for i := 1; i < (size_bytes - int(binary.Size(ebr_empty_byte))); i++ {
								f.Write([]byte{1})
							}

							fmt.Println("La Particion extendida fue creada con exitosamente")
						} else if s_dsk_fit == "b" {
							best_index := num_particion

							s_part_start_act := ""
							s_part_status_act := ""
							s_part_size_act := ""
							i_part_size_act := 0
							s_part_start_best := ""
							i_part_start_best := 0
							s_part_start_best_ant := ""
							i_part_start_best_ant := 0
							s_part_size_best := ""
							i_part_size_best := 0
							s_part_size_best_ant := ""
							i_part_size_best_ant := 0

							for i := 0; i < 4; i++ {
								s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])

								s_part_start_act = strings.Trim(s_part_start_act, "\x00")

								s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])

								s_part_status_act = strings.Trim(s_part_status_act, "\x00")

								s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])

								s_part_size_act = strings.Trim(s_part_size_act, "\x00")
								i_part_size_act, _ = strconv.Atoi(s_part_size_act)

								if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
									if i != num_particion {
										s_part_size_best = string(master_boot_record.Mbr_partition[best_index].Part_size[:])

										s_part_size_best = strings.Trim(s_part_size_best, "\x00")
										i_part_size_best, _ = strconv.Atoi(s_part_size_best)

										s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])

										s_part_size_act = strings.Trim(s_part_size_act, "\x00")
										i_part_size_act, _ = strconv.Atoi(s_part_size_act)

										if i_part_size_best > i_part_size_act {
											best_index = i
											break
										}
									}
								}
							}

							copy(master_boot_record.Mbr_partition[best_index].Part_type[:], "e")
							copy(master_boot_record.Mbr_partition[best_index].Part_fit[:], aux_fit)

							if best_index == 0 {
								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))))
							} else {
								s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])

								s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
								i_part_start_best_ant, _ = strconv.Atoi(s_part_start_best_ant)

								s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])

								s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
								i_part_size_best_ant, _ = strconv.Atoi(s_part_size_best_ant)

								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant))
							}

							copy(master_boot_record.Mbr_partition[best_index].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[best_index].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[best_index].Part_name[:], nombre)

							mbr_byte := Struct_a_bytes(master_boot_record)

							f.Seek(0, os.SEEK_SET)
							f.Write(mbr_byte)

							s_part_start_best = string(master_boot_record.Mbr_partition[best_index].Part_start[:])

							s_part_start_best = strings.Trim(s_part_start_best, "\x00")
							i_part_start_best, _ = strconv.Atoi(s_part_start_best)

							f.Seek(int64(i_part_start_best), os.SEEK_SET)

							extended_boot_record := estructuras.Ebr{}
							copy(extended_boot_record.Part_fit[:], aux_fit)
							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_start[:], s_part_start_best)
							copy(extended_boot_record.Part_size[:], "0")
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], "")
							ebr_byte := Struct_a_bytes(extended_boot_record)
							f.Write(ebr_byte)

							pos_actual, _ := f.Seek(0, os.SEEK_CUR)
							f.Seek(int64(pos_actual), os.SEEK_SET)

							ebr_empty_byte := Struct_a_bytes(mbr_empty)

							for i := 1; i < (size_bytes - int(binary.Size(ebr_empty_byte))); i++ {
								f.Write([]byte{1})
							}

							fmt.Println("La Particion extendida fue creada con exitosamente")
						} else {
							worst_index := num_particion

							s_part_start_act := ""
							s_part_status_act := ""
							s_part_size_act := ""
							i_part_size_act := 0
							s_part_start_worst := ""
							i_part_start_worst := 0
							s_part_start_worst_ant := ""
							i_part_start_worst_ant := 0
							s_part_size_worst := ""
							i_part_size_worst := 0
							s_part_size_worst_ant := ""
							i_part_size_worst_ant := 0

							for i := 0; i < 4; i++ {
								s_part_start_act = string(master_boot_record.Mbr_partition[i].Part_start[:])

								s_part_start_act = strings.Trim(s_part_start_act, "\x00")

								s_part_status_act = string(master_boot_record.Mbr_partition[i].Part_status[:])

								s_part_status_act = strings.Trim(s_part_status_act, "\x00")

								s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])

								s_part_size_act = strings.Trim(s_part_size_act, "\x00")
								i_part_size_act, _ = strconv.Atoi(s_part_size_act)

								if s_part_start_act == "-1" || (s_part_status_act == "1" && i_part_size_act >= size_bytes) {
									if i != num_particion {
										s_part_size_worst = string(master_boot_record.Mbr_partition[worst_index].Part_size[:])

										s_part_size_worst = strings.Trim(s_part_size_worst, "\x00")
										i_part_size_worst, _ = strconv.Atoi(s_part_size_worst)

										s_part_size_act = string(master_boot_record.Mbr_partition[i].Part_size[:])

										s_part_size_act = strings.Trim(s_part_size_act, "\x00")
										i_part_size_act, _ = strconv.Atoi(s_part_size_act)

										if i_part_size_worst < i_part_size_act {
											worst_index = i
											break
										}
									}
								}
							}

							copy(master_boot_record.Mbr_partition[worst_index].Part_type[:], "e")
							copy(master_boot_record.Mbr_partition[worst_index].Part_fit[:], aux_fit)

							if worst_index == 0 {
								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))))
							} else {
								s_part_start_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_start[:])

								s_part_start_worst_ant = strings.Trim(s_part_start_worst_ant, "\x00")
								i_part_start_worst_ant, _ = strconv.Atoi(s_part_start_worst_ant)

								s_part_size_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_size[:])

								s_part_size_worst_ant = strings.Trim(s_part_size_worst_ant, "\x00")
								i_part_size_worst_ant, _ = strconv.Atoi(s_part_size_worst_ant)

								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(i_part_start_worst_ant+i_part_size_worst_ant))
							}

							copy(master_boot_record.Mbr_partition[worst_index].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[worst_index].Part_status[:], "0")
							copy(master_boot_record.Mbr_partition[worst_index].Part_name[:], nombre)

							mbr_byte := Struct_a_bytes(master_boot_record)

							f.Seek(0, os.SEEK_SET)
							f.Write(mbr_byte)

							s_part_start_worst = string(master_boot_record.Mbr_partition[worst_index].Part_start[:])

							s_part_start_worst = strings.Trim(s_part_start_worst, "\x00")
							i_part_start_worst, _ = strconv.Atoi(s_part_start_worst)

							f.Seek(int64(i_part_start_worst), os.SEEK_SET)

							extended_boot_record := estructuras.Ebr{}
							copy(extended_boot_record.Part_fit[:], aux_fit)
							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_start[:], s_part_start_worst)
							copy(extended_boot_record.Part_size[:], "0")
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], "")
							ebr_byte := Struct_a_bytes(extended_boot_record)
							f.Write(ebr_byte)

							pos_actual, _ := f.Seek(0, os.SEEK_CUR)
							f.Seek(int64(pos_actual), os.SEEK_SET)

							ebr_empty_byte := Struct_a_bytes(mbr_empty)

							for i := 1; i < (size_bytes - int(binary.Size(ebr_empty_byte))); i++ {
								f.Write([]byte{1})
							}

							fmt.Println(" La Particion extendida fue creada con exitosamente")
						}
					} else {
						fmt.Println("Error: La particion que desea crear excede el espacio disponible")
					}
				} else {
					fmt.Println("Error: La suma de particiones primarias y extendidas no debe exceder de 4 particiones")
				}
			} else {
				fmt.Println("Error: Solo puede haber una particion extendida por disco")
			}
		} else {
			fmt.Println("Error: el disco se encuentra vacio")
		}
		f.Close()
	}
}

func existe_particion(direccion string, nombre string) bool {

	mbr_empty := estructuras.Mbr{}
	var empty [100]byte

	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		_, err = f.ReadAt(lectura, 0)

		if err != nil && err != io.EOF {
			Mens_error(err)
		}

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if err != nil {
			Mens_error(err)
		}

		if master_boot_record.Mbr_tamano != empty {
			s_part_name := ""
			s_part_type := ""

			for i := 0; i < 4; i++ {
				s_part_name = string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")

				if s_part_name == nombre {
					f.Close()
					return true
				}

				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					s_part_start := string(master_boot_record.Mbr_partition[i].Part_start[:])
					s_part_start = strings.Trim(s_part_start, "\x00")
					part_start, err := strconv.Atoi(s_part_start)
					if err != nil {
						Mens_error(err)
					}
					fmt.Println(direccion)
					ebr := Mostrarebr(direccion, part_start)
					band_salir := false
					for !band_salir {
						if ebr.Part_status != empty {
							s_part_name := string(ebr.Part_name[:])
							s_part_name = strings.Trim(s_part_name, "\x00")
							if s_part_name == nombre {
								return true
							}
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
							band_salir = true
						} else {
							ebr = Mostrarebr(direccion, part_next)
						}
					}
				}
			}
		} else {
			fmt.Println("Error: el disco se encuentra vacio")
		}
	} else {
		Mens_error(err)
	}
	defer func() {
		f.Close()
	}()

	f.Close()
	return false
}

func Bytes_a_struct_ebr(s []byte) estructuras.Ebr {
	p := estructuras.Ebr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	if err != nil && err != io.EOF {
		Mens_error(err)
	}

	return p
}

func obtener_ebr(ruta string, part_start int) estructuras.Ebr {
	var empty [100]byte
	ebr_empty := estructuras.Ebr{}

	// Apertura de archivo
	disco, err := os.OpenFile(ruta, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		Mens_error(err)
	}
	defer func() {
		disco.Close()
	}()

	// Calculo del tamano de struct en bytes
	ebr2 := Struct_a_bytes(ebr_empty)
	sstruct := len(ebr2)

	// Lectrura del archivo binario desde el inicio
	lectura := make([]byte, sstruct)
	_, err = disco.ReadAt(lectura, int64(part_start))

	// ERROR
	if err != nil && err != io.EOF {
		Mens_error(err)
	}

	// Conversion de bytes a struct
	ebr := Bytes_a_struct_ebr(lectura)

	// ERROR
	if err != nil {
		Mens_error(err)
	}

	if ebr.Part_next != empty {
		return ebr
	} else {
		return ebr_empty
	}
}
func Mostrarebr(ruta string, part_start int) estructuras.Ebr {
	var empty [100]byte
	ebr_empty := estructuras.Ebr{}

	disco, err := os.OpenFile(ruta, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}
	defer func() {
		disco.Close()
	}()

	ebr2 := Struct_a_bytes(ebr_empty)
	sstruct := len(ebr2)

	lectura := make([]byte, sstruct)
	_, err = disco.ReadAt(lectura, int64(part_start))
	if err != nil && err != io.EOF {
		Mens_error(err)
	}
	ebr := Bytes_a_struct_ebr(lectura)

	if err != nil {
		Mens_error(err)
	}

	if ebr.Part_next != empty {
		return ebr
	} else {
		return ebr_empty
	}
}

func Mostrarmbr(ruta string) estructuras.Mbr {

	mbr_empty := estructuras.Mbr{}
	var empty [100]byte
	f, err := os.OpenFile(valorpath, os.O_RDWR, 0660)

	defer func() {
		f.Close()
	}()

	if err == nil {
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			return master_boot_record
		} else {
			return mbr_empty
		}
	}
	return mbr_empty
}

func crear_particion_logixa(direccion string, nombre string, size int, fit string, unit string) {
	aux_fit := ""
	aux_unit := ""
	size_bytes := 1024
	band_crear := false

	if fit != "" {
		aux_fit = fit
	} else {
		aux_fit = "w"
	}

	if unit != "" {
		aux_unit = unit

		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			size_bytes = size * 1024
		} else {
			size_bytes = size * 1024 * 1024
		}
	} else {
		size_bytes = size * 1024
	}
	ebr_empty := estructuras.Ebr{}
	partition_extendida := 0
	band_eliminada := false

	mbr_empty := estructuras.Mbr{}
	var empty [100]byte
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			s_part_type := ""

			for i := 0; i < 4; i++ {

				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					partition_extendida = i
					s_part_s := string(master_boot_record.Mbr_partition[i].Part_start[:])
					s_part_s = strings.Trim(s_part_s, "\x00")

					break
				}
			}
		}
		if partition_extendida != 0 {
			s_part_start := string(master_boot_record.Mbr_partition[partition_extendida].Part_start[:])
			s_part_start = strings.Trim(s_part_start, "\x00")
			i_part_start, err := strconv.Atoi(s_part_start)
			if err != nil {
				Mens_error(err)
			}
			band_salir := false
			for !band_salir {
				ebr_emp := estructuras.Ebr{}

				if err != nil {
					Mens_error(err)
				}
				defer func() {
					f.Close()
				}()

				ebr2 := Struct_a_bytes(ebr_emp)
				sstruct := len(ebr2)

				lectura := make([]byte, sstruct)
				_, err = f.ReadAt(lectura, int64(i_part_start))

				if err != nil && err != io.EOF {
					Mens_error(err)
				}

				ebr := Bytes_a_struct_ebr(lectura)

				if err != nil {
					Mens_error(err)
				}
				if ebr == ebr_empty {
					band_salir = true
					break
				}
				s_part_next := string(ebr.Part_next[:])
				s_part_next = strings.Trim(s_part_next, "\x00")
				i_part_next, err := strconv.Atoi(s_part_next)
				if err != nil {
					Mens_error(err)
				}
				s_part_mount := string(ebr.Part_status[:])
				s_part_mount = strings.Trim(s_part_mount, "\x00")

				if i_part_next == -1 {
					copy(ebr_empty.Part_status[:], "0")
					copy(ebr_empty.Part_fit[:], []byte(aux_fit))
					copy(ebr_empty.Part_start[:], strconv.Itoa(i_part_start))
					copy(ebr_empty.Part_size[:], strconv.Itoa(size_bytes))
					copy(ebr_empty.Part_name[:], nombre)
					nuevo_part_next := i_part_start + size_bytes
					copy(ebr_empty.Part_next[:], strconv.Itoa(nuevo_part_next))
					band_crear = true
					band_salir = true
				} else if s_part_mount == "E" {
					s_part_s := string(ebr.Part_size[:])
					s_part_s = strings.Trim(s_part_s, "\x00")
					i_part_s, err := strconv.Atoi(s_part_s)
					if err != nil {
						Mens_error(err)
					}
					if i_part_s >= size_bytes {
						copy(ebr_empty.Part_status[:], "0")
						copy(ebr_empty.Part_fit[:], []byte(aux_fit))
						copy(ebr_empty.Part_start[:], strconv.Itoa(i_part_start))
						copy(ebr_empty.Part_size[:], strconv.Itoa(size_bytes))
						copy(ebr_empty.Part_name[:], []byte(nombre))
						copy(ebr_empty.Part_next[:], ebr.Part_next[:])
						band_eliminada = true
						band_crear = true
						band_salir = true
					}

				} else {
					i_part_start = i_part_next

				}

			}
			if band_crear {
				s_part_next := string(ebr_empty.Part_next[:])
				s_part_next = strings.Trim(s_part_next, "\x00")
				nuevo_part_next, err := strconv.Atoi(s_part_next)
				if err != nil {
					Mens_error(err)
				}

				ebr_byte := Struct_a_bytes(ebr_empty)
				newpos, err := f.Seek(int64(i_part_start), os.SEEK_SET)
				if err != nil {
					Mens_error(err)
				}
				_, err = f.WriteAt(ebr_byte, newpos)

				if err != nil {
					Mens_error(err)
				}

				//s_part_type := ""

				fmt.Println("Se creo correctamente la particion logica")

				if !band_eliminada {
					ebr_emptys := estructuras.Ebr{}
					copy(ebr_emptys.Part_status[:], "0")
					copy(ebr_emptys.Part_fit[:], "")
					copy(ebr_emptys.Part_start[:], "0")
					copy(ebr_emptys.Part_size[:], "0")
					copy(ebr_emptys.Part_next[:], "-1")
					copy(ebr_emptys.Part_name[:], "")

					ebr_bytes := Struct_a_bytes(ebr_emptys)
					posi, err := f.Seek(int64(nuevo_part_next), os.SEEK_SET)
					if err != nil {
						Mens_error(err)
					}
					_, err = f.WriteAt(ebr_bytes, posi)

					if err != nil {
						Mens_error(err)
					}
				}
			} else {
				fmt.Println("No se pudo crear la particion logica")
			}

		} else {
			fmt.Println("No se puede crear la extension logica porque no existe particion extendida")
		}

	}
}

func eliminar(direccion string, nombre string) {
	//particion_extendida := 0
	ebr_empty := estructuras.Ebr{}
	partition_extendida := 0

	mbr_empty := estructuras.Mbr{}
	var empty [100]byte
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			for i := 0; i < 4; i++ {

				s_part_name := string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")
				s_part_types := string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_types = strings.Trim(s_part_types, "\x00")

				if s_part_name == nombre {
					if s_part_types == "p" {
						if err != nil {
							Mens_error(err)
						}
						emptyPartstatus := make([]byte, len(master_boot_record.Mbr_partition[i].Part_status))
						copy(master_boot_record.Mbr_partition[i].Part_status[:], emptyPartstatus)
						copy(master_boot_record.Mbr_partition[i].Part_status[:], "E")
						emptyParttype := make([]byte, len(master_boot_record.Mbr_partition[i].Part_type))
						copy(master_boot_record.Mbr_partition[i].Part_type[:], emptyParttype)
						copy(master_boot_record.Mbr_partition[i].Part_type[:], "0")
						emptyPartfit := make([]byte, len(master_boot_record.Mbr_partition[i].Part_fit))
						copy(master_boot_record.Mbr_partition[i].Part_fit[:], emptyPartfit)
						copy(master_boot_record.Mbr_partition[i].Part_fit[:], "0")
						emptyPartstart := make([]byte, len(master_boot_record.Mbr_partition[i].Part_start))
						copy(master_boot_record.Mbr_partition[i].Part_start[:], emptyPartstart)
						copy(master_boot_record.Mbr_partition[i].Part_start[:], "0")
						emptyPartsize := make([]byte, len(master_boot_record.Mbr_partition[i].Part_size))
						copy(master_boot_record.Mbr_partition[i].Part_size[:], emptyPartsize)
						copy(master_boot_record.Mbr_partition[i].Part_size[:], "0")
						emptyPartname := make([]byte, len(master_boot_record.Mbr_partition[i].Part_name))
						copy(master_boot_record.Mbr_partition[i].Part_name[:], emptyPartname)
						copy(master_boot_record.Mbr_partition[i].Part_name[:], "0")
						emptyPartcorre := make([]byte, len(master_boot_record.Mbr_partition[i].Part_correlative))
						copy(master_boot_record.Mbr_partition[i].Part_correlative[:], emptyPartcorre)
						copy(master_boot_record.Mbr_partition[i].Part_correlative[:], "0")
						emptyPartid := make([]byte, len(master_boot_record.Mbr_partition[i].Part_id))
						copy(master_boot_record.Mbr_partition[i].Part_id[:], emptyPartid)
						copy(master_boot_record.Mbr_partition[i].Part_id[:], "0")
						mbr_byte := Struct_a_bytes(master_boot_record)

						newpos, err := f.Seek(0, os.SEEK_SET)

						if err != nil {
							Mens_error(err)
						}

						_, err = f.WriteAt(mbr_byte, newpos)

						if err != nil {
							Mens_error(err)
						}
						fmt.Println("Se elimino correctamente la particion primaria")

					} else if s_part_types == "e" {
						//particion_extendida = i
						copy(master_boot_record.Mbr_partition[i].Part_status[:], "E")
						copy(master_boot_record.Mbr_partition[i].Part_type[:], "0")
						copy(master_boot_record.Mbr_partition[i].Part_fit[:], "0")
						copy(master_boot_record.Mbr_partition[i].Part_start[:], "0")
						copy(master_boot_record.Mbr_partition[i].Part_size[:], "0")
						copy(master_boot_record.Mbr_partition[i].Part_name[:], "0")
						copy(master_boot_record.Mbr_partition[i].Part_correlative[:], "0")
						copy(master_boot_record.Mbr_partition[i].Part_id[:], "0")
						s_part_start := string(master_boot_record.Mbr_partition[partition_extendida].Part_start[:])
						s_part_start = strings.Trim(s_part_start, "\x00")
						i_part_start, err := strconv.Atoi(s_part_start)
						if err != nil {
							Mens_error(err)
						}
						band_salir := false
						for !band_salir {
							ebr_emp := estructuras.Ebr{}

							if err != nil {
								Mens_error(err)
							}
							defer func() {
								f.Close()
							}()

							ebr2 := Struct_a_bytes(ebr_emp)
							sstruct := len(ebr2)

							lectura := make([]byte, sstruct)
							_, err = f.ReadAt(lectura, int64(i_part_start))

							if err != nil && err != io.EOF {
								Mens_error(err)
							}

							ebr := Bytes_a_struct_ebr(lectura)

							if err != nil {
								Mens_error(err)
							}

							if ebr == ebr_empty {
								band_salir = true
								break
							}
							s_part_next := string(ebr.Part_next[:])
							s_part_next = strings.Trim(s_part_next, "\x00")
							if s_part_next == "-1" {
								copy(ebr_empty.Part_status[:], "E")
								copy(ebr_empty.Part_fit[:], "0")
								copy(ebr_empty.Part_start[:], "0")
								copy(ebr_empty.Part_size[:], "0")
								copy(ebr_empty.Part_name[:], "0")
								copy(ebr_empty.Part_next[:], "0")
							} else {
								copy(ebr_empty.Part_status[:], "E")
								copy(ebr_empty.Part_fit[:], "0")
								copy(ebr_empty.Part_start[:], "0")
								copy(ebr_empty.Part_size[:], "0")
								copy(ebr_empty.Part_name[:], "0")
								copy(ebr_empty.Part_next[:], "0")
							}

						}
					}

				}
			}

		}

	}
	defer func() {
		f.Close()
	}()
	f.Close()
}
