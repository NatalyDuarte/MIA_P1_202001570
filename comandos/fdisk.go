package comandos

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"mimodulo/estructuras"
	"os"
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

	mbr_empty := estructuras.Mbr{}
	var empty [100]byte

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

				if s_part_start == "-1" && band_particion == false {
					band_particion = true
					num_particion = i
				}
			}

			if band_particion {
				espacio_usado := 0
				for i := 0; i < 4; i++ {
					s_size := string(master_boot_record.Mbr_partition[i].Part_size[:])
					s_size = strings.Trim(s_size, "\x00")
					i_size, err := strconv.Atoi(s_size)

					if err != nil {
						Mens_error(err)
					}

					espacio_usado += i_size
				}

				s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
				s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
				i_tamaño_disco, err2 := strconv.Atoi(s_tamaño_disco)

				if err2 != nil {
					Mens_error(err)
				}

				espacio_disponible := i_tamaño_disco - espacio_usado

				fmt.Println("Espacion disponible: ", espacio_disponible, " Bytes")
				fmt.Println("Espacio necesario: ", size_bytes, " Bytes")
				fmt.Println(num_particion)

				if espacio_disponible >= size_bytes {
					fmt.Println("Particion primaria creada con exito")
					copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
					copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "p")
					copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], []byte(fit))
					copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], []byte(strconv.Itoa(espacio_disponible-size_bytes)))
					copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], []byte(strconv.Itoa(size_bytes)))
					copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], []byte(nombre))
					mbr_byte := Struct_a_bytes(master_boot_record)

					newpos, err := f.Seek(0, os.SEEK_SET)

					if err != nil {
						Mens_error(err)
					}

					_, err = f.WriteAt(mbr_byte, newpos)

					if err != nil {
						Mens_error(err)
					}
				}

				if err != nil {
					Mens_error(err)
				}

			}
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
	aux_unit := ""
	aux_path := direccion
	size_bytes := 1024
	band_er := false
	mbr_empty := estructuras.Mbr{}
	var empty [100]byte

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

				if s_part_start == "-1" && band_particion == false {
					band_particion = true
					num_particion = i
				}
			}

			if band_particion {
				espacio_usado := 0
				for i := 0; i < 4; i++ {
					s_size := string(master_boot_record.Mbr_partition[i].Part_size[:])
					s_size = strings.Trim(s_size, "\x00")
					i_size, err := strconv.Atoi(s_size)

					if err != nil {
						Mens_error(err)
					}

					espacio_usado += i_size
				}

				s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
				s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
				i_tamaño_disco, err2 := strconv.Atoi(s_tamaño_disco)

				if err2 != nil {
					Mens_error(err)
				}

				espacio_disponible := i_tamaño_disco - espacio_usado

				fmt.Println("Espacion disponible: ", espacio_disponible, " Bytes")
				fmt.Println("Espacio necesario: ", size_bytes, " Bytes")
				fmt.Println(num_particion)

				if espacio_disponible >= size_bytes {
					for i := 0; i < 4; i++ {
						s_part_typ := string(master_boot_record.Mbr_partition[i].Part_type[:])
						s_part_typ = strings.Trim(s_part_typ, "\x00")

						if s_part_typ == "e" {
							band_er = true
						}
					}

				}

				if band_er {
					fmt.Println("Error: No se puede crear la particion extendida porque ya existe una")
				} else {
					copy(master_boot_record.Mbr_partition[num_particion].Part_status[:], "0")
					copy(master_boot_record.Mbr_partition[num_particion].Part_type[:], "e")
					copy(master_boot_record.Mbr_partition[num_particion].Part_fit[:], []byte(fit))
					copy(master_boot_record.Mbr_partition[num_particion].Part_start[:], []byte(strconv.Itoa(espacio_disponible-size_bytes)))
					copy(master_boot_record.Mbr_partition[num_particion].Part_size[:], []byte(strconv.Itoa(size_bytes)))
					copy(master_boot_record.Mbr_partition[num_particion].Part_name[:], []byte(nombre))

					logis := estructuras.ListLogica{}
					logis.Part_status = "0"
					logis.Part_fit = "0"
					logis.Part_start = "-1"
					logis.Part_size = "0"
					logis.Part_next = "0"
					logis.Part_name = ""
					Listlogi = append(Listlogi, logis)
					fmt.Println(Listlogi)

					mbr_byte := Struct_a_bytes(master_boot_record)
					newpos, err := f.Seek(0, os.SEEK_SET)

					if err != nil {
						Mens_error(err)
					}

					_, err = f.WriteAt(mbr_byte, newpos)

					if err != nil {
						Mens_error(err)
					}

					fmt.Println("Particion extendida creada con exito")
				}

				if err != nil {
					Mens_error(err)
				}

			}
		}
		f.Close()
	}
}
func crear_particion_logixa(direccion string, nombre string, size int, fit string, unit string) {
	band_crea := false
	posicion := 0

	mbr_empty := estructuras.Mbr{}
	ebr_empty := estructuras.Ebr{}
	disco, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}

	mbr2 := Struct_a_bytes(mbr_empty)
	ebr2 := Struct_a_bytes(ebr_empty)
	sstruct := len(mbr2)
	estruct := len(ebr2)
	lectura := make([]byte, sstruct)
	lect := make([]byte, estruct)
	_, err = disco.ReadAt(lectura, 0)
	_, err = disco.ReadAt(lect, 0)
	if err != nil && err != io.EOF {
		Mens_error(err)
	}

	mbr := Bytes_a_struct_mbr(lectura)

	if err != nil {
		Mens_error(err)
	}
	for i := 0; i < 4; i++ {
		name := string(mbr.Mbr_partition[i].Part_type[:])
		name = strings.Trim(name, "\x00")
		if name == "e" {
			band_crea = true
			posicion = i
		}
	}

	if band_crea {
		aux_unit := ""
		aux_path := direccion
		size_bytes := 1024
		band_er := false
		mbr_empty := estructuras.Mbr{}
		var empty [100]byte

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
			band_particion := true
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

				if band_particion {
					espacio_usado := 0

					for i := 0; i < len(Listlogi); i++ {
						s_size := string(Listlogi[i].Part_size)
						s_size = strings.Trim(s_size, "\x00")
						i_size, err := strconv.Atoi(s_size)

						if err != nil {
							Mens_error(err)
						}

						espacio_usado += i_size

					}

					s_tamaño_disco := string(master_boot_record.Mbr_partition[posicion].Part_size[:])
					s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
					i_tamaño_disco, err2 := strconv.Atoi(s_tamaño_disco)

					if err2 != nil {
						Mens_error(err)
					}

					espacio_disponible := i_tamaño_disco - espacio_usado

					fmt.Println("Espacion disponible: ", espacio_disponible, " Bytes")
					fmt.Println("Espacio necesario: ", size_bytes, " Bytes")

					if espacio_disponible >= size_bytes {
						band_er = false

					}

					if band_er {
						fmt.Println("Error: No se puede crear la particion logica porque no hay espacio suficiente")
					} else {
						logi := estructuras.Ebr{}

						copy(logi.Part_status[:], "0")
						copy(logi.Part_fit[:], []byte(fit))
						copy(logi.Part_start[:], []byte(strconv.Itoa(espacio_disponible-size_bytes)))
						copy(logi.Part_size[:], []byte(strconv.Itoa(size_bytes)))
						copy(logi.Part_name[:], []byte(nombre))

						mbr_byte := Struct_a_bytes(logi)
						newpos, err := f.Seek(0, os.SEEK_SET)

						if err != nil {
							Mens_error(err)
						}

						_, err = f.WriteAt(mbr_byte, newpos)

						if err != nil {
							Mens_error(err)
						}

						fmt.Println("Particion logica creada con exito")
					}

					if err != nil {
						Mens_error(err)
					}

				}
			}
			f.Close()
		}

	} else {
		fmt.Println("Error: No existe una particion extendida por lo cual no se puede crear una logica.")
	}
	disco.Close()
}

func existe_particion(direccion string, nombre string) bool {
	mbr_empty := estructuras.Mbr{}
	disco, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	if err != nil {
		Mens_error(err)
	}

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
	for i := 0; i < 4; i++ {
		name := string(mbr.Mbr_partition[i].Part_name[:])
		name = strings.Trim(name, "\x00")
		if name == nombre {
			return true
		}
	}
	disco.Close()
	return false
}

func Bytes_a_struct_ebr(s []byte) estructuras.Ebr {
	p := estructuras.Ebr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		Mens_error(err)
	}

	return p
}
