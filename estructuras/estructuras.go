package estructuras

type Mbr = struct {
	Mbr_tamano         [100]byte
	Mbr_fecha_creacion [100]byte
	Mbr_dsk_signature  [100]byte
}

type Mkdisk struct {
	Path      string
	Fit, Unit byte
	Size      int
}
