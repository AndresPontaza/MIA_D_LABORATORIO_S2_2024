package main

/*--------------------------Import--------------------------*/

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

/*--------------------------/Import--------------------------*/

/*--------------------------Structs--------------------------*/

// Master Boot Record
type mbr = struct {
	Mbr_tamano         [100]byte
	Mbr_fecha_creacion [100]byte
	Mbr_dsk_signature  [100]byte
	Dsk_fit            [100]byte
	Mbr_partition      [4]partition
}

// Partition
type partition = struct {
	Part_status [100]byte
	Part_type   [100]byte
	Part_fit    [100]byte
	Part_start  [100]byte
	Part_size   [100]byte
	Part_name   [100]byte
}

// Extended Boot Record
type ebr = struct {
	Part_status [100]byte
	Part_fit    [100]byte
	Part_start  [100]byte
	Part_size   [100]byte
	Part_next   [100]byte
	Part_name   [100]byte
}

// Super Bloque
type super_bloque = struct {
	S_filesystem_type   [100]byte
	S_inodes_count      [100]byte
	S_blocks_count      [100]byte
	S_free_blocks_count [100]byte
	S_free_inodes_count [100]byte
	S_mtime             [100]byte
	S_mnt_count         [100]byte
	S_magic             [100]byte
	S_inode_size        [100]byte
	S_block_size        [100]byte
	S_firts_ino         [100]byte
	S_first_blo         [100]byte
	S_bm_inode_start    [100]byte
	S_bm_block_start    [100]byte
	S_inode_start       [100]byte
	S_block_start       [100]byte
}

// Tablas de Inodos
type inodo = struct {
	I_uid   [100]byte
	I_gid   [100]byte
	I_size  [100]byte
	I_atime [100]byte
	I_ctime [100]byte
	I_mtime [100]byte
	I_block [100]byte
	I_type  [100]byte
	I_perm  [100]byte
}

// Bloques de Carpetas
type bloque_carpeta = struct {
	B_content [4]cotent
}

// Content
type cotent = struct {
	B_name  [100]byte
	B_inodo [100]byte
}

// Bloques de Archivos
type bloque_archivo = struct {
	B_content [100]byte
}

/*--------------------------/Structs--------------------------*/

/*--------------------------Metodos o Funciones--------------------------*/

// Metodo principal
func main() {
	analizar()
}

// Muestra el mensaje de error
func msg_error(err error) {
	fmt.Println("[ERROR] ", err)
}

// Obtiene y lee el comando
func analizar() {
	finalizar := false
	fmt.Println("Ejemplo 4: Fdisk 1.0")
	reader := bufio.NewReader(os.Stdin)

	//  Pide constantemente un comando
	for !finalizar {
		fmt.Print("Ingrese un comando: ")
		// Lee hasta que presione ENTER
		comando, _ := reader.ReadString('\n')

		if strings.Contains(comando, "exit") {
			/* SALIR */
			finalizar = true
		} else if strings.Contains(comando, "EXIT") {
			/* SALIR */
			finalizar = true
		} else {
			// Si no es vacio o el comando EXIT
			if comando != "" && comando != "exit\n" && comando != "EXIT\n" {
				// Obtener comando y parametros
				split_comando(comando)
			}
		}
	}
}

// Separa los diferentes comando con sus parametros si tienen
func split_comando(comando string) {
	var commandArray []string
	// Elimina los saltos de linea y retornos de carro
	comando = strings.Replace(comando, "\n", "", 1)
	comando = strings.Replace(comando, "\r", "", 1)

	// Banderas para verficar comentarios
	band_comentario := false

	if strings.Contains(comando, "#") {
		// Comentario
		band_comentario = true
		fmt.Println(comando)
	} else {
		// Comando con Parametros
		commandArray = strings.Split(comando, " -")
	}

	// Ejecuta el comando leido si no es un comentario
	if !band_comentario {
		ejecutar_comando(commandArray)
	}
}

// Identifica y ejecuta el comando encontrado
func ejecutar_comando(commandArray []string) {
	// Convierte el comando a minusculas
	data := strings.ToLower(commandArray[0])

	// Identifica el comando a ejecutar
	if data == "mkdisk" {
		/* MKDISK */
		mkdisk(commandArray)
	} else if data == "rmdisk" {
		/* RMDISK */
		rmdisk(commandArray)
	} else if data == "fdisk" { /* Ejemplo 4: Fdisk 1.0 */
		/* FDISK */
		fdisk(commandArray)
	} else if data == "rep" {
		/* REP */
		rep()
	} else {
		/* ERROR */
		fmt.Println("[ERROR] El comando no fue reconocido...")
	}
}

/*--------------------------Comandos--------------------------*/

/* MKDISK 2.0 */
func mkdisk(commandArray []string) {
	fmt.Println("[MENSAJE] El comando MKDISK aqui inicia")

	// Variables para los valores de los parametros
	val_size := 0
	val_fit := ""
	val_unit := ""
	val_path := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_size := false
	band_fit := false
	band_unit := false
	band_path := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> SIZE */
		case strings.Contains(data, "size="):
			// Valido si el parametro ya fue ingresado
			if band_size {
				fmt.Println("[ERROR] El parametro -size ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_size = true

			// Conversion a entero
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size

			// ERROR de conversion
			if err != nil {
				band_error = true
				msg_error(err)
			}

			// Valido que el tamaño sea positivo
			if val_size < 0 {
				band_error = true
				fmt.Println("[ERROR] El parametro -size es negativo...")
				break
			}
		/* PARAMETRO OPCIONAL -> FIT */
		case strings.Contains(data, "fit="):
			// Valido si el parametro ya fue ingresado
			if band_fit {
				fmt.Println("[ERROR] El parametro -fit ya fue ingresado...")
				band_error = true
				break
			}

			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "w"
			} else {
				fmt.Println("[ERROR] El Valor del parametro -fit no es valido...")
				band_error = true
				break
			}
		/* PARAMETRO OPCIONAL -> UNIT */
		case strings.Contains(data, "unit="):
			// Valido si el parametro ya fue ingresado
			if band_unit {
				fmt.Println("[ERROR] El parametro -unit ya fue ingresado...")
				band_error = true
				break
			}

			// Reemplaza comillas y lo paso a minusculas
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)

			if val_unit == "k" || val_unit == "m" {
				// Activo la bandera del parametro
				band_unit = true
			} else {
				// Parametro no valido
				fmt.Println("[ERROR] El Valor del parametro -unit no es valido...")
				band_error = true
				break
			}
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				fmt.Println("[ERROR] El parametro -path ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Println("[ERROR] Parametro no valido...")
		}
	}

	// Verifico si no hay errores
	if !band_error {
		// Verifico que el parametro "Path" (Obligatorio) este ingresado
		if band_path {
			// Verifico que el parametro "Size" (Obligatorio) este ingresado
			if band_size {
				total_size := 1024
				master_boot_record := mbr{}

				// Disco -> Archivo Binario
				crear_disco(val_path)

				// Fecha
				fecha := time.Now()
				str_fecha := fecha.Format("02/01/2006 15:04:05")

				// Copio valor al Struct
				copy(master_boot_record.Mbr_fecha_creacion[:], str_fecha)

				// Numero aleatorio
				rand.Seed(time.Now().UnixNano())
				min := 0
				max := 100
				num_random := rand.Intn(max-min+1) + min

				// Copio valor al Struct
				copy(master_boot_record.Mbr_dsk_signature[:], strconv.Itoa(int(num_random)))

				// Verifico si existe el parametro "Fit" (Opcional)
				if band_fit {
					// Copio valor al Struct
					copy(master_boot_record.Dsk_fit[:], val_fit)
				} else {
					// Si no especifica -> "Primer Ajuste"
					copy(master_boot_record.Dsk_fit[:], "f")
				}

				// Verifico si existe el parametro "Unit" (Opcional)
				if band_unit {
					// Megabytes
					if val_unit == "m" {
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
						total_size = val_size * 1024
					} else {
						// Kilobytes
						copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024)))
						total_size = val_size
					}
				} else {
					// Si no especifica -> Megabytes
					copy(master_boot_record.Mbr_tamano[:], strconv.Itoa(int(val_size*1024*1024)))
					total_size = val_size * 1024
				}

				// Inicializar Parcticiones
				for i := 0; i < 4; i++ {
					copy(master_boot_record.Mbr_partition[i].Part_status[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_type[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_fit[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_start[:], "-1")
					copy(master_boot_record.Mbr_partition[i].Part_size[:], "0")
					copy(master_boot_record.Mbr_partition[i].Part_name[:], "")
				}

				// Convierto de entero a string
				str_total_size := strconv.Itoa(total_size)

				// Comando para definir el tamaño (Kilobytes) y llenarlo de ceros
				cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero of=\""+val_path+"\" bs=1024 count="+str_total_size)
				cmd.Dir = "/"
				_, err := cmd.Output()

				// ERROR
				if err != nil {
					msg_error(err)
				}

				// Se escriben los datos en disco

				// Apertura del archivo
				disco, err := os.OpenFile(val_path, os.O_RDWR, 0660)

				// ERROR
				if err != nil {
					msg_error(err)
				}

				// Conversion de struct a bytes
				mbr_byte := struct_a_bytes(master_boot_record)

				// Se posiciona al inicio del archivo para guardar la informacion del disco
				newpos, err := disco.Seek(0, os.SEEK_SET)

				// ERROR
				if err != nil {
					msg_error(err)
				}

				// Escritura de struct en archivo binario
				_, err = disco.WriteAt(mbr_byte, newpos)

				// ERROR
				if err != nil {
					msg_error(err)
				}

				disco.Close()
			}
		}
	}

	fmt.Println("[MENSAJE] El comando MKDISK aqui finaliza")
}

/* RMDISK 1.0 */
func rmdisk(commandArray []string) {
	fmt.Println("[MENSAJE] El comando RMDISK aqui inicia")

	// Variables para los valores de los parametros
	val_path := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_path := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				fmt.Println("[ERROR] El parametro -path ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Println("[ERROR] Parametro no valido...")
		}
	}

	// Verifico si no hay errores
	if !band_error {
		// Verifico que el parametro "Path" (Obligatorio) este ingresado
		for band_path {

			// Si existe el archivo binario
			_, e := os.Stat(val_path)

			if e != nil {
				// Si no existe
				if os.IsNotExist(e) {
					fmt.Println("[ERROR] No existe el disco que desea eliminar...")
					band_path = false
				}
			} else {
				// Si existe
				fmt.Print("[MENSAJE] ¿Desea eliminar el disco [S/N]?: ")

				// Obtengo la opcion ingresada por el usuario
				var opcion string
				fmt.Scanln(&opcion)

				// Verificando entrada
				if opcion == "s" || opcion == "S" {

					// Elimina el archivo
					cmd := exec.Command("/bin/sh", "-c", "rm \""+val_path+"\"")
					cmd.Dir = "/"
					_, err := cmd.Output()

					// ERROR
					if err != nil {
						msg_error(err)
					} else {
						fmt.Println("[SUCCES] El Disco fue eliminado!")
					}

					band_path = false
				} else if opcion == "n" || opcion == "N" {
					fmt.Println("[MENSAJE] EL disco no fue eliminado")
					band_path = false
				} else {
					fmt.Println("[ERROR] Opcion no valida intentalo de nuevo...")
				}
			}
		}
	}

	fmt.Println("[MENSAJE] El comando RMDISK aqui finaliza")
}

/* Ejemplo 4: Fdisk 1.0 */
/* FDISK 1.0 */
func fdisk(commandArray []string) {
	fmt.Println("[MENSAJE] El comando FDISK aqui inicia")

	// Variables para los valores de los parametros
	val_size := 0
	val_unit := ""
	val_path := ""
	val_type := ""
	val_fit := ""
	val_name := ""

	// Banderas para verificar los parametros y ver si se repiten
	band_size := false
	band_unit := false
	band_path := false
	band_type := false
	band_fit := false
	band_name := false
	band_error := false

	// Obtengo solo los parametros validos
	for i := 1; i < len(commandArray); i++ {
		aux_data := strings.SplitAfter(commandArray[i], "=")
		data := strings.ToLower(aux_data[0])
		val_data := aux_data[1]

		// Identifica los parametos
		switch {
		/* PARAMETRO OBLIGATORIO -> SIZE */
		case strings.Contains(data, "size="):
			// Valido si el parametro ya fue ingresado
			if band_size {
				fmt.Println("[ERROR] El parametro -size ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_size = true

			// Conversion a entero
			aux_size, err := strconv.Atoi(val_data)
			val_size = aux_size
			fmt.Println("Size: ", val_size)
			// ERROR de conversion
			if err != nil {
				msg_error(err)
				band_error = true
			}

			// Valido que el tamaño sea positivo
			if val_size < 0 {
				band_error = true
				fmt.Println("[ERROR] El parametro -size es negativo...")
				break
			}
		/* PARAMETRO OPCIONAL -> UNIT */
		case strings.Contains(data, "unit="):
			// Valido si el parametro ya fue ingresado
			if band_unit {
				fmt.Println("[ERROR] El parametro -unit ya fue ingresado...")
				band_error = true
				break
			}

			// Reemplaza comillas y lo paso a minusculas
			val_unit = strings.Replace(val_data, "\"", "", 2)
			val_unit = strings.ToLower(val_unit)
			fmt.Println("Unit: ", val_unit)
			if val_unit == "b" || val_unit == "k" || val_unit == "m" {
				// Activo la bandera del parametro
				band_unit = true
			} else {
				// Parametro no valido
				fmt.Println("[ERROR] El Valor del parametro -unit no es valido...")
				band_error = true
				break
			}
		/* PARAMETRO OBLIGATORIO -> PATH */
		case strings.Contains(data, "path="):
			if band_path {
				fmt.Println("[ERROR] El parametro -path ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_path = true

			// Reemplaza comillas
			val_path = strings.Replace(val_data, "\"", "", 2)
			fmt.Println("Path: ", val_path)
		/* PARAMETRO OPCIONAL -> TYPE */
		case strings.Contains(data, "type="):
			if band_type {
				fmt.Println("[ERROR] El parametro -type ya fue ingresado...")
				band_error = true
				break
			}

			// Reemplaza comillas y lo paso a minusculas
			val_type = strings.Replace(val_data, "\"", "", 2)
			val_type = strings.ToLower(val_type)
			fmt.Println("Type: ", val_type)
			if val_type == "p" || val_type == "e" || val_type == "l" {
				// Activo la bandera del parametro
				band_type = true
			} else {
				// Parametro no valido
				fmt.Println("[ERROR] El Valor del parametro -type no es valido...")
				band_error = true
				break
			}
		/* PARAMETRO OPCIONAL -> FIT */
		case strings.Contains(data, "fit="):
			// Valido si el parametro ya fue ingresado
			if band_fit {
				fmt.Println("[ERROR] El parametro -fit ya fue ingresado...")
				band_error = true
				break
			}

			// Le quito las comillas y lo paso a minusculas
			val_fit = strings.Replace(val_data, "\"", "", 2)
			val_fit = strings.ToLower(val_fit)

			if val_fit == "bf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "b"
			} else if val_fit == "ff" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "f"
			} else if val_fit == "wf" {
				// Activo la bandera del parametro y obtengo el caracter que me interesa
				band_fit = true
				val_fit = "w"
			} else {
				fmt.Println("[ERROR] El Valor del parametro -fit no es valido...")
				band_error = true
				break
			}
			fmt.Println("fit: ", val_fit)
		/* PARAMETRO OBLIGATORIO -> NAME */
		case strings.Contains(data, "name="):
			// Valido si el parametro ya fue ingresado
			if band_name {
				fmt.Println("[ERROR] El parametro -name ya fue ingresado...")
				band_error = true
				break
			}

			// Activo la bandera del parametro
			band_name = true

			// Reemplaza comillas
			val_name = strings.Replace(val_data, "\"", "", 2)
			fmt.Println("Name: ", val_name)
		/* PARAMETRO NO VALIDO */
		default:
			fmt.Println("[ERROR] Parametro no valido...")
		}
	}

	// Verifico si no hay errores
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
						} else {
							// Logica
						}
					} else {
						// Si no lo indica se tomara como Primaria
						crear_particion_primaria(val_path, val_name, val_size, val_fit, val_unit)
					}
				} else {
					fmt.Println("[ERROR] El parametro -name no fue ingresado")
				}
			} else {
				fmt.Println("[ERROR] El parametro -path no fue ingresado")
			}
		} else {
			fmt.Println("[ERROR] El parametro -size no fue ingresado")
		}
	}

	fmt.Println("[MENSAJE] El comando FDISK aqui finaliza")
}

/* REP 1.0 */
func rep() {
	fin_archivo := false
	var empty [100]byte
	mbr_empty := mbr{}
	cont := 0

	fmt.Println("* Reporte de MKDISK: *")

	// Apertura de archivo
	disco, err := os.OpenFile("/home/andres-pontaza/Documentos/Laboratorio MIA/Discos/A.dsk", os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		msg_error(err)
	}

	// Calculo del tamano de struct en bytes
	mbr2 := struct_a_bytes(mbr_empty)
	sstruct := len(mbr2)

	// RECORRE CADA STRUCT DEL ARCHIVO
	for !fin_archivo {
		// Lectrura de conjunto de bytes en archivo binario
		lectura := make([]byte, sstruct)
		_, err = disco.ReadAt(lectura, int64(sstruct*cont))

		// ERROR
		if err != nil && err != io.EOF {
			msg_error(err)
		}

		// Conversion de bytes a struct
		mbr := bytes_a_struct_mbr(lectura)
		sstruct = len(lectura)

		// ERROR
		if err != nil {
			msg_error(err)
		}

		if mbr.Mbr_tamano == empty {
			fin_archivo = true
		} else {
			fmt.Print("Tamaño: ")
			fmt.Print(string(mbr.Mbr_tamano[:]))
			fmt.Println(" bytes ")
			fmt.Print("Fecha: ")
			fmt.Println(string(mbr.Mbr_fecha_creacion[:]))
			fmt.Print("Signature: ")
			fmt.Println(string(mbr.Mbr_dsk_signature[:]))
		}

		cont++
	}
	disco.Close()
}

/*--------------------------/Comandos--------------------------*/

// Crea el archivo que simula ser un disco duro
func crear_disco(ruta string) {
	aux, err := filepath.Abs(ruta)

	// ERROR
	if err != nil {
		msg_error(err)
	}

	// Crea el directiorio de forma recursiva
	cmd1 := exec.Command("/bin/sh", "-c", "sudo mkdir -p '"+filepath.Dir(aux)+"'")
	cmd1.Dir = "/"
	_, err1 := cmd1.Output()

	// ERROR
	if err1 != nil {
		msg_error(err)
	}

	// Da los permisos al directorio
	cmd2 := exec.Command("/bin/sh", "-c", "sudo chmod -R 777 '"+filepath.Dir(aux)+"'")
	cmd2.Dir = "/"
	_, err2 := cmd2.Output()

	// ERROR
	if err2 != nil {
		msg_error(err)
	}

	// Verifica si existe la ruta para el archivo
	if _, err := os.Stat(filepath.Dir(aux)); errors.Is(err, os.ErrNotExist) {
		if err != nil {
			fmt.Println("[FAILURE] No se pudo crear el disco...")
		}
	}
}

/* Ejemplo 4: Fdisk 1.0 */
// Crea la Particion Primaria
func crear_particion_primaria(direccion string, nombre string, size int, fit string, unit string) {
	//aux_fit := ""
	aux_unit := ""
	aux_path := direccion
	size_bytes := 1024
	//buffer := "1"

	mbr_empty := mbr{}
	var empty [100]byte

	/* Pendiente */
	// Verifico si tiene Ajuste
	if fit != "" {
		//aux_fit = fit
	} else {
		// Por default es Peor ajuste
		//aux_fit = "w"
	}

	// Verifico si tiene Unidad
	if unit != "" {
		aux_unit = unit

		// *Bytes
		if aux_unit == "b" {
			size_bytes = size
		} else if aux_unit == "k" {
			// *Kilobytes
			size_bytes = size * 1024
		} else {
			// *Megabytes
			size_bytes = size * 1024 * 1024
		}
	} else {
		// Por default Kilobytes
		size_bytes = size * 1024
	}

	// Abro el archivo para lectura con opcion a modificar
	// * OpenFile(name string, flag int, perm FileMode)
	// * O_RDWR -> Lectura y Escritura
	// * 0660 -> Permisos de lectura y escritura
	f, err := os.OpenFile(aux_path, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		msg_error(err)
	} else {
		// Procede a leer el archivo
		band_particion := false
		num_particion := 0

		// Calculo del tamano de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		// make -> Crea un slice de bytes con el tamaño indicado (sstruct)
		// ReadAt -> Lee el archivo binario desde la posicion indicada (0) y lo guarda en el slice de bytes
		// Slice de byte es un arreglo de bytes que se puede modificar y con ReadAt se llena con los bytes del archivo
		lectura := make([]byte, sstruct)
		_, err = f.ReadAt(lectura, 0)

		// ERROR
		if err != nil && err != io.EOF {
			msg_error(err)
		}

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)

		// ERROR
		if err != nil {
			msg_error(err)
		}

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_start := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				s_part_start = string(master_boot_record.Mbr_partition[i].Part_start[:])
				// Le quito los caracteres null
				s_part_start = strings.Trim(s_part_start, "\x00")

				// Verifico si en las particiones hay espacio
				if s_part_start == "-1" && band_particion == false {
					band_particion = true
					num_particion = i
				}
			}

			// Verifico si hay espacio
			if band_particion {
				espacio_usado := 0

				// Recorro las 4 particiones
				for i := 0; i < 4; i++ {
					// Obtengo el espacio utilizado
					s_size := string(master_boot_record.Mbr_partition[i].Part_size[:])
					// Le quito los caracteres null
					s_size = strings.Trim(s_size, "\x00")
					i_size, err := strconv.Atoi(s_size)

					// ERROR
					if err != nil {
						msg_error(err)
					}

					// Le sumo el valor al espacio
					espacio_usado += i_size
				}

				/* Tamaño del disco */

				// Obtengo el tamaño del disco
				s_tamaño_disco := string(master_boot_record.Mbr_tamano[:])
				// Le quito los caracteres null
				s_tamaño_disco = strings.Trim(s_tamaño_disco, "\x00")
				i_tamaño_disco, err2 := strconv.Atoi(s_tamaño_disco)

				// ERROR
				if err2 != nil {
					msg_error(err)
				}

				espacio_disponible := i_tamaño_disco - espacio_usado

				fmt.Println("[ESPACIO DISPONIBLE] ", espacio_disponible, " Bytes")
				fmt.Println("[ESPACIO NECESARIO] ", size_bytes, " Bytes")
				fmt.Println(num_particion)

				// Verifico que haya espacio suficiente
				if espacio_disponible >= size_bytes {
					fmt.Println("Si cumple " + nombre + " !")
				}
			}
		}
		f.Close()
	}
}

/* Ejemplo 4: Fdisk 1.0 */
// Verifica si el nombre de la particion esta disponible
func existe_particion(direccion string, nombre string) bool {
	extendida := -1
	mbr_empty := mbr{}
	ebr_empty := ebr{}
	var empty [100]byte
	cont := 0
	fin_archivo := false

	// Abro el archivo para lectura con opcion a modificar
	f, err := os.OpenFile(direccion, os.O_RDWR, 0660)

	// ERROR
	if err != nil {
		msg_error(err)
	} else {
		// Procedo a leer el archivo

		// Calculo del tamano de struct en bytes
		mbr2 := struct_a_bytes(mbr_empty)
		sstruct := len(mbr2)

		// Lectrura del archivo binario desde el inicio
		// make -> Crea un slice de bytes con el tamaño indicado (sstruct)
		// ReadAt -> Lee el archivo binario desde la posicion indicada (0) y lo guarda en el slice de bytes
		// Slice de byte es un arreglo de bytes que se puede modificar y con ReadAt se llena con los bytes del archivo
		lectura := make([]byte, sstruct)
		_, err = f.ReadAt(lectura, 0)

		// ERROR
		if err != nil && err != io.EOF {
			msg_error(err)
		}

		// Conversion de bytes a struct
		master_boot_record := bytes_a_struct_mbr(lectura)
		sstruct = len(lectura)

		// ERROR
		if err != nil {
			msg_error(err)
		}

		// Si el disco esta creado
		if master_boot_record.Mbr_tamano != empty {
			s_part_name := ""
			s_part_type := ""

			// Recorro las 4 particiones
			for i := 0; i < 4; i++ {
				// Antes de comparar limpio la cadena
				// Obtengo el nombre de la particion
				// [:] -> Convierte el arreglo de bytes a cadena
				s_part_name = string(master_boot_record.Mbr_partition[i].Part_name[:])
				s_part_name = strings.Trim(s_part_name, "\x00")

				/* Pendiente */
				// Verifico si ya existe una particion con ese nombre
				if s_part_name == nombre {

				}

				// Antes de comparar limpio la cadena
				s_part_type = string(master_boot_record.Mbr_partition[i].Part_type[:])
				s_part_type = strings.Trim(s_part_type, "\x00")

				// Verifico si de tipo extendida
				if s_part_type == "E" {
					extendida = i
				}
			}

			// Lo busco en las extendidas
			if extendida != -1 {
				// Obtengo el inicio de la particion
				s_part_start := string(master_boot_record.Mbr_partition[extendida].Part_start[:])
				// Le quito los caracteres null
				s_part_start = strings.Trim(s_part_start, "\x00")
				i_part_start, err := strconv.Atoi(s_part_start)

				// ERROR
				if err != nil {
					msg_error(err)
					fin_archivo = true
				}

				// Obtengo el espacio de la partcion
				s_part_size := string(master_boot_record.Mbr_partition[extendida].Part_size[:])
				// Le quito los caracteres null
				s_part_size = strings.Trim(s_part_size, "\x00")
				i_part_size, err := strconv.Atoi(s_part_size)

				// ERROR
				if err != nil {
					msg_error(err)
					fin_archivo = true
				}

				// Calculo del tamano de struct en bytes
				ebr2 := struct_a_bytes(ebr_empty)
				sstruct := len(ebr2)

				// Lectrura de conjunto de bytes desde el inicio de la particion
				for !fin_archivo {
					// Lectrura de conjunto de bytes en archivo binario
					lectura := make([]byte, sstruct)
					n_leidos, err := f.ReadAt(lectura, int64(sstruct*cont+i_part_start))

					// ERROR
					if err != nil {
						msg_error(err)
						fin_archivo = true
					}

					// Posicion actual en el archivo
					// Seek -> Cambia la posicion del puntero de lectura/escritura
					// Seek(offset int64, whence int) (int64, error)
					// whence -> 0: desde el inicio, 1: desde la posicion actual, 2: desde el final
					// os.SEEK_CUR -> Desde la posicion actual
					pos_actual, err := f.Seek(0, os.SEEK_CUR)

					// ERROR
					if err != nil {
						msg_error(err)
						fin_archivo = true
					}

					// Si no lee nada y ya se paso del tamaño de la particion
					if n_leidos == 0 && pos_actual < int64(i_part_start+i_part_size) {
						fin_archivo = true
						break
					}

					// Conversion de bytes a struct
					extended_boot_record := bytes_a_struct_ebr(lectura)
					sstruct = len(lectura)

					if err != nil {
						msg_error(err)
					}

					if extended_boot_record.Part_size == empty {
						fin_archivo = true
					} else {
						fmt.Print(" Nombre: ")
						fmt.Print(string(extended_boot_record.Part_name[:]))

						// Antes de comparar limpio la cadena
						s_part_name = string(extended_boot_record.Part_name[:])
						s_part_name = strings.Trim(s_part_name, "\x00")

						// Verifico si ya existe una particion con ese nombre
						if s_part_name == nombre {
							f.Close()
							return true
						}

						// Obtengo el espacio utilizado
						s_part_next := string(extended_boot_record.Part_next[:])
						// Le quito los caracteres null
						s_part_next = strings.Trim(s_part_next, "\x00")
						i_part_next, err := strconv.Atoi(s_part_next)

						// ERROR
						if err != nil {
							msg_error(err)
						}

						// Si ya termino
						if i_part_next != -1 {
							f.Close()
							return false
						}
					}
					cont++
				}
			}
		}
	}
	f.Close()
	return false
}

// Codifica de Struct a []Bytes
func struct_a_bytes(p interface{}) []byte {
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(p)

	// ERROR
	if err != nil && err != io.EOF {
		msg_error(err)
	}

	return buf.Bytes()
}

/* Ejemplo 4: Fdisk 1.0 */
// Decodifica de [] Bytes a Struct
func bytes_a_struct_mbr(s []byte) mbr {
	// Descodificacion
	p := mbr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		msg_error(err)
	}

	return p
}

/* Ejemplo 4: Fdisk 1.0 */
// Decodifica de [] Bytes a Struct
func bytes_a_struct_ebr(s []byte) ebr {
	// Descodificacion
	p := ebr{}
	dec := gob.NewDecoder(bytes.NewReader(s))
	err := dec.Decode(&p)

	// ERROR
	if err != nil && err != io.EOF {
		msg_error(err)
	}

	return p
}

/*--------------------------/Metodos o Funciones--------------------------*/
