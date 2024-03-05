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
var Listlogi []estructuras.ListLogica

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
	if !pathc {
		if !band_error {
			if band_delete {
				if band_path {
					if band_name {

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

			} else {
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
							copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
						} else {
							s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
							s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
							i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

							s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
							s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
							i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

							copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant+1))
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
							copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
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
							copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
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
	expresionRegular := regexp.MustCompile(`\d+`)
	numerosComoString := ""

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
						coincidencias := expresionRegular.FindAllString(nombre, -1)
						numerosComoString = strings.Join(coincidencias, ",")
						s_dsk_fit := string(master_boot_record.Dsk_fit[:])
						s_dsk_fit = strings.Trim(s_dsk_fit, "\x00")

						if s_dsk_fit == "f" {
							copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "e")
							copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], aux_fit)

							if num_particion == 0 {
								mbr_empty_byte := Struct_a_bytes(mbr_empty)
								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
							} else {
								s_part_start_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_start[:])
								s_part_start_ant = strings.Trim(s_part_start_ant, "\x00")
								i_part_start_ant, _ := strconv.Atoi(s_part_start_ant)

								s_part_size_ant := string(master_boot_record.Mbr_partition[num_particion-1].Part_size[:])
								s_part_size_ant = strings.Trim(s_part_size_ant, "\x00")
								i_part_size_ant, _ := strconv.Atoi(s_part_size_ant)

								copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], strconv.Itoa(i_part_start_ant+i_part_size_ant+1))
							}

							copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], strconv.Itoa(size_bytes))
							copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
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

							extended_boot_record := estructuras.Ebr{}

							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_fit[:], aux_fit)
							copy(extended_boot_record.Part_start[:], s_part_start)
							copy(extended_boot_record.Part_size[:], "0")
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], "")
							ebr_byte := Struct_a_bytes(extended_boot_record)
							f.Write(ebr_byte)

							ebr_empty_byte := Struct_a_bytes(ebr_empty)

							pos_actual, _ := f.Seek(0, os.SEEK_CUR)
							f.Seek(int64(pos_actual+1), os.SEEK_SET)

							for i := 1; i < (size_bytes - int(binary.Size(ebr_empty_byte))); i++ {
								f.Write([]byte{1})
							}

							fmt.Println("La Particion extendida fue creada exitosamente")
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
								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
							} else {
								s_part_start_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_start[:])
								s_part_start_best_ant = strings.Trim(s_part_start_best_ant, "\x00")
								i_part_start_best_ant, _ = strconv.Atoi(s_part_start_best_ant)

								s_part_size_best_ant = string(master_boot_record.Mbr_partition[best_index-1].Part_size[:])
								s_part_size_best_ant = strings.Trim(s_part_size_best_ant, "\x00")
								i_part_size_best_ant, _ = strconv.Atoi(s_part_size_best_ant)

								copy(master_boot_record.Mbr_partition[best_index].Part_start[:], strconv.Itoa(i_part_start_best_ant+i_part_size_best_ant+1))
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

							extended_boot_record := estructuras.Ebr{}

							copy(extended_boot_record.Part_status[:], "0")
							copy(extended_boot_record.Part_fit[:], aux_fit)
							copy(extended_boot_record.Part_start[:], s_part_start_best)
							copy(extended_boot_record.Part_size[:], "0")
							copy(extended_boot_record.Part_next[:], "-1")
							copy(extended_boot_record.Part_name[:], "")
							ebr_byte := Struct_a_bytes(extended_boot_record)
							f.Write(ebr_byte)

							pos_actual, _ := f.Seek(0, os.SEEK_CUR)
							f.Seek(int64(pos_actual+1), os.SEEK_SET)

							ebr_empty_byte := Struct_a_bytes(mbr_empty)

							for i := 1; i < (size_bytes - int(binary.Size(ebr_empty_byte))); i++ {
								f.Write([]byte{1})
							}

							fmt.Println(" La Particion extendida fue creada exitosamente")
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
								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(int(binary.Size(mbr_empty_byte))+1))
							} else {
								s_part_start_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_start[:])
								s_part_start_worst_ant = strings.Trim(s_part_start_worst_ant, "\x00")
								i_part_start_worst_ant, _ = strconv.Atoi(s_part_start_worst_ant)
								s_part_size_worst_ant = string(master_boot_record.Mbr_partition[worst_index-1].Part_size[:])
								s_part_size_worst_ant = strings.Trim(s_part_size_worst_ant, "\x00")
								i_part_size_worst_ant, _ = strconv.Atoi(s_part_size_worst_ant)

								copy(master_boot_record.Mbr_partition[worst_index].Part_start[:], strconv.Itoa(i_part_start_worst_ant+i_part_size_worst_ant+1))
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
							f.Seek(int64(pos_actual+1), os.SEEK_SET)

							ebr_empty_byte := Struct_a_bytes(mbr_empty)

							for i := 1; i < (size_bytes - int(binary.Size(ebr_empty_byte))); i++ {
								f.Write([]byte{1})
							}

							fmt.Println("La Particion extendida fue creada exitosamente")
						}

					} else {
						fmt.Println("Error: Ya no hay espacio disponible")
					}
				} else {
					fmt.Println("Error: Ya no hay particiones disponibles")
				}
			} else {
				fmt.Println("Error: Solo puede haber una particion extendida por disco")
			}
		} else {
			fmt.Println("Error: El disco se encuentra vacio")
		}
		f.Close()
	}
}
func crear_particion_logixa(direccion string, nombre string, size int, fit string, unit string) {
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
		fmt.Println("Error: No existe el disco duro con ese nombre")
	} else {
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)
		lectura := make([]byte, sstruct)
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

		if master_boot_record.Mbr_tamano != empty {
			s_part_type := ""
			num_extendida := -1

			for i := 0; i < 4; i++ {
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				if s_part_type == "e" {
					num_extendida = i
					break
				}
			}
			if num_extendida != -1 {
				s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
				s_part_start = strings.Trim(s_part_start, "\x00")
				i_part_start, _ := strconv.Atoi(s_part_start)

				cont := i_part_start
				f.Seek(int64(cont), os.SEEK_SET)

				ebr2 := Struct_a_bytes(ebr_empty)
				sstruct := len(ebr2)

				lectura := make([]byte, sstruct)
				f.Read(lectura)
				extended_boot_record := Bytes_a_struct_ebr(lectura)

				s_part_size_ext := string(extended_boot_record.Part_size[:])
				s_part_size_ext = strings.Trim(s_part_size_ext, "\x00")

				if s_part_size_ext == "0" {
					s_part_size := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ := strconv.Atoi(s_part_size)

					if i_part_size < size_bytes {
						fmt.Println("Error: La particion logica no puede exceder el espacio de la particion extendida")
					} else {
						copy(extended_boot_record.Part_status[:], "0")
						copy(extended_boot_record.Part_fit[:], aux_fit)

						pos_actual, _ := f.Seek(0, os.SEEK_CUR)

						copy(extended_boot_record.Part_start[:], strconv.Itoa(int(pos_actual)-int(binary.Size(ebr_empty))+1))
						copy(extended_boot_record.Part_size[:], strconv.Itoa(size_bytes))
						copy(extended_boot_record.Part_next[:], "-1")
						copy(extended_boot_record.Part_name[:], nombre)

						s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
						s_part_start = strings.Trim(s_part_start, "\x00")
						i_part_start, _ := strconv.Atoi(s_part_start)

						ebr_byte := Struct_a_bytes(extended_boot_record)
						f.Seek(int64(i_part_start), os.SEEK_SET)
						f.Write(ebr_byte)

						fmt.Println("La particion logica fue creada exitosamente")
					}
				} else {
					s_part_size := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
					s_part_size = strings.Trim(s_part_size, "\x00")
					i_part_size, _ := strconv.Atoi(s_part_size)

					s_part_start := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
					s_part_start = strings.Trim(s_part_start, "\x00")
					i_part_start, _ := strconv.Atoi(s_part_start)

					espacio_particion := i_part_size + i_part_start

					fmt.Println("[ESPACIO PARTICION] ", espacio_particion, " Bytes")
					fmt.Println("[ESPACIO NECESARIO] ", size_bytes, " Bytes")

					s_part_next := string(extended_boot_record.Part_next[:])
					s_part_next = strings.Trim(s_part_next, "\x00")
					i_part_next, _ := strconv.Atoi(s_part_next)
					pos_actual, _ := f.Seek(0, os.SEEK_CUR)

					for (i_part_next != -1) && (int(pos_actual) < (i_part_size + i_part_start)) {
						f.Seek(int64(i_part_next), os.SEEK_SET)

						ebr2 := Struct_a_bytes(ebr_empty)
						sstruct := len(ebr2)

						lectura := make([]byte, sstruct)
						f.Read(lectura)

						pos_actual, _ = f.Seek(0, os.SEEK_CUR)

						extended_boot_record = Bytes_a_struct_ebr(lectura)

						s_part_next = string(extended_boot_record.Part_next[:])
						s_part_next = strings.Trim(s_part_next, "\x00")
						i_part_next, _ = strconv.Atoi(s_part_next)
					}

					s_part_start_ext := string(extended_boot_record.Part_start[:])
					s_part_start_ext = strings.Trim(s_part_start_ext, "\x00")
					i_part_start_ext, _ := strconv.Atoi(s_part_start_ext)

					s_part_size_ext := string(extended_boot_record.Part_size[:])
					s_part_size_ext = strings.Trim(s_part_size_ext, "\x00")
					i_part_size_ext, _ := strconv.Atoi(s_part_size_ext)

					s_part_size_mbr := string(master_boot_record.Mbr_partition[num_extendida].Part_size[:])
					s_part_size_mbr = strings.Trim(s_part_size_mbr, "\x00")
					i_part_size_mbr, _ := strconv.Atoi(s_part_size_mbr)

					s_part_start_mbr := string(master_boot_record.Mbr_partition[num_extendida].Part_start[:])
					s_part_start_mbr = strings.Trim(s_part_start_mbr, "\x00")
					i_part_start_mbr, _ := strconv.Atoi(s_part_start_mbr)

					espacio_necesario := i_part_start_ext + i_part_size_ext + size_bytes
					fmt.Println("Next prueba")
					fmt.Println(i_part_start_ext + i_part_size_ext)

					if espacio_necesario <= (i_part_size_mbr + i_part_start_mbr) {
						copy(extended_boot_record.Part_next[:], strconv.Itoa(i_part_start_ext+i_part_size_ext))

						pos_actual, _ = f.Seek(0, os.SEEK_CUR)

						f.Seek(int64(int(pos_actual)-int(binary.Size(ebr_empty))), os.SEEK_SET)
						ebr_byte := Struct_a_bytes(extended_boot_record)
						f.Write(ebr_byte)

						f.Seek(int64(i_part_start_ext+i_part_size_ext), os.SEEK_SET)
						copy(extended_boot_record.Part_status[:], "0")
						copy(extended_boot_record.Part_fit[:], aux_fit)

						pos_actual, _ = f.Seek(0, os.SEEK_CUR)

						copy(extended_boot_record.Part_start[:], strconv.Itoa(int(pos_actual)))
						copy(extended_boot_record.Part_size[:], strconv.Itoa(size_bytes))
						copy(extended_boot_record.Part_next[:], "-1")
						copy(extended_boot_record.Part_name[:], nombre)

						ebr_byte = Struct_a_bytes(extended_boot_record)
						f.Write(ebr_byte)
						fmt.Println("La Particion logica fue creada exitosamente")
					} else {
						fmt.Println("Error: La particion logica a crear excede el espacio disponible de la particion extendida")
					}
				}
			} else {
				fmt.Println("Error: No se puede crear una particion logica si no hay una extendida.")
			}

		} else {
			fmt.Println("Error: El disco se encuentra vacio")
		}
		f.Close()
	}
}

func existe_particion(direccion string, nombre string) bool {
	extendida := -1
	mbr_empty := estructuras.Mbr{}
	ebr_empty := estructuras.Ebr{}
	var empty [100]byte
	fin := false

	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err == nil {
		mbr2 := Struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		lectura := make([]byte, sstruct)
		f.Seek(0, os.SEEK_SET)
		f.Read(lectura)

		master_boot_record := Bytes_a_struct_mbr(lectura)

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

				ebr2 := Struct_a_bytes(ebr_empty)
				sstruct := len(ebr2)

				for !fin {
					lectura := make([]byte, sstruct)
					n_leidos, _ := f.Read(lectura)

					pos_actual, _ := f.Seek(0, os.SEEK_CUR)

					if n_leidos == 0 && (pos_actual >= int64(i_part_size+i_part_start)) {
						fin = true
						break
					}

					extended_boot_record := Bytes_a_struct_ebr(lectura)
					sstruct = len(lectura)

					if extended_boot_record.Part_size == empty {
						fin = true
						break
					} else {
						s_part_name = string(extended_boot_record.Part_name[:])
						s_part_name = strings.Trim(s_part_name, "\x00")

						if s_part_name == nombre {
							f.Close()
							return true
						}

						s_part_next := string(extended_boot_record.Part_next[:])
						s_part_next = strings.Trim(s_part_next, "\x00")

						if s_part_next != "-1" {
							f.Close()
							return false
						}
					}
				}
			}
		} else {
			fmt.Println("[ERROR] el disco se encuentra vacio...")
		}
	}

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
