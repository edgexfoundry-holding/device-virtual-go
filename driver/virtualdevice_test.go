package driver

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

const (
	DeviceName                                        = "Random-Value-Generator01"
	DeviceCommandName_Bool, DeviceResource_Bool       = "RandomValue_Bool", "RandomValue_Bool"
	DeviceCommandName_Int8, DeviceResource_Int8       = "RandomValue_Int8", "RandomValue_Int8"
	DeviceCommandName_Int16, DeviceResource_Int16     = "RandomValue_Int16", "RandomValue_Int16"
	DeviceCommandName_Int32, DeviceResource_Int32     = "RandomValue_Int32", "RandomValue_Int32"
	DeviceCommandName_Int64, DeviceResource_Int64     = "RandomValue_Int64", "RandomValue_Int64"
	DeviceCommandName_Uint8, DeviceResource_Uint8     = "RandomValue_Uint8", "RandomValue_Uint8"
	DeviceCommandName_Uint16, DeviceResource_Uint16   = "RandomValue_Uint16", "RandomValue_Uint16"
	DeviceCommandName_Uint32, DeviceResource_Uint32   = "RandomValue_Uint32", "RandomValue_Uint32"
	DeviceCommandName_Uint64, DeviceResource_Uint64   = "RandomValue_Uint64", "RandomValue_Uint64"
	DeviceCommandName_Float32, DeviceResource_Float32 = "RandomValue_Float32", "RandomValue_Float32"
	DeviceCommandName_Float64, DeviceResource_Float64 = "RandomValue_Float64", "RandomValue_Float64"
	EnableRandomization_True                          = "true"
)

func init() {
	if _, err := os.Stat(qlDatabaseDriverName); os.IsNotExist(err) {
		os.Mkdir(qlDatabaseDir, os.ModeDir)
	}

	db := getDb()
	if err := db.openDb(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer func() {
		if err := db.closeDb(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	err := db.startTransaction()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := db.exec(SQL_DROP_TABLE); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := db.exec(SQL_CREATE_TABLE); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	ds := [][]string{
		{DeviceName, DeviceCommandName_Bool, DeviceResource_Bool, EnableRandomization_True, typeBool, "true"},
		{DeviceName, DeviceCommandName_Int8, DeviceResource_Int8, EnableRandomization_True, typeInt8, "0"},
		{DeviceName, DeviceCommandName_Int16, DeviceResource_Int16, EnableRandomization_True, typeInt16, "0"},
		{DeviceName, DeviceCommandName_Int32, DeviceResource_Int32, EnableRandomization_True, typeInt32, "0"},
		{DeviceName, DeviceCommandName_Int64, DeviceResource_Int64, EnableRandomization_True, typeInt64, "0"},
		{DeviceName, DeviceCommandName_Uint8, DeviceResource_Uint8, EnableRandomization_True, typeUint8, "0"},
		{DeviceName, DeviceCommandName_Uint16, DeviceResource_Uint16, EnableRandomization_True, typeUint16, "0"},
		{DeviceName, DeviceCommandName_Uint32, DeviceResource_Uint32, EnableRandomization_True, typeUint32, "0"},
		{DeviceName, DeviceCommandName_Uint64, DeviceResource_Uint64, EnableRandomization_True, typeUint64, "0"},
		{DeviceName, DeviceCommandName_Float32, DeviceResource_Float32, EnableRandomization_True, typeFloat32, "0"},
		{DeviceName, DeviceCommandName_Float64, DeviceResource_Float64, EnableRandomization_True, typeFloat64, "0"},
	}
	for _, d := range ds {
		b, _ := strconv.ParseBool(d[3])
		if err := db.exec(SQL_INSERT, d[0], d[1], d[2], b, d[4], d[5]); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	if err = db.commit(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func TestValue_Bool(t *testing.T) {
	db := getDb()
	if err := db.openDb(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	rd := newVirtualDevice()
	v1, err := rd.value(DeviceName, DeviceResource_Bool, "", "", db)
	if err != nil {
		t.Fatal(err)
	}

	//the return string must be convertible to boolean
	b1, err := strconv.ParseBool(v1)
	if err != nil {
		t.Fatal(err)
	}

	rounds := 20
	//EnableRandomization = true
	for x := 1; x <= rounds; x++ {
		v2, _ := rd.value(DeviceName, DeviceResource_Bool, "", "", db)
		b2, _ := strconv.ParseBool(v2)
		if b1 != b2 {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got same value in %d rounds", rounds)
		}
	}

	//EnableRandomization = false
	if err := db.startTransaction(); err != nil {
		t.Fatalf("Start a transaction failed: %v", err)
	}
	if err := db.exec(SQL_UPDATE_ENABLERANDOMIZATION, false, DeviceName, DeviceResource_Bool); err != nil {
		t.Fatal(err)
	}
	if err := db.commit(); err != nil {
		t.Fatalf("Commit the transaction failed: %v", err)
	}

	v1, _ = rd.value(DeviceName, DeviceResource_Bool, "", "", db)
	b1, _ = strconv.ParseBool(v1)
	for x := 0; x <= rounds; x++ {
		v2, _ := rd.value(DeviceName, DeviceResource_Bool, "", "", db)
		b2, _ := strconv.ParseBool(v2)
		if b1 != b2 {
			t.Fatalf("EnableRandomization is false, but got different value")
		}
	}
}

func TestValueIntx(t *testing.T) {
	SubTestValueIntx(t, DeviceResource_Int8)
	SubTestValueIntx(t, DeviceResource_Int16)
	SubTestValueIntx(t, DeviceResource_Int32)
	SubTestValueIntx(t, DeviceResource_Int64)
}

func TestValueUintx(t *testing.T) {
	SubTestValueUintx(t, DeviceResource_Uint8)
	SubTestValueUintx(t, DeviceResource_Uint16)
	SubTestValueUintx(t, DeviceResource_Uint32)
	SubTestValueUintx(t, DeviceResource_Uint64)
}

func TestValueFloatx(t *testing.T) {
	SubTestValueFloatx(t, DeviceResource_Float32)
	SubTestValueFloatx(t, DeviceResource_Float64)
}

func SubTestValueIntx(t *testing.T, dr string) {
	db := getDb()
	if err := db.openDb(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	rd := newVirtualDevice()

	var minInt, maxInt int64
	var bitSize int
	switch dr {
	case DeviceResource_Int8:
		minInt, maxInt = rd.minInt8, rd.maxInt8
		bitSize = 8
	case DeviceResource_Int16:
		minInt, maxInt = rd.minInt16, rd.maxInt16
		bitSize = 16
	case DeviceResource_Int32:
		minInt, maxInt = rd.minInt32, rd.maxInt32
		bitSize = 32
	case DeviceResource_Int64:
		minInt, maxInt = rd.minInt64, rd.maxInt64
		bitSize = 64
	default:
		t.Fatal("unknown device resource:", dr)
	}

	if _, err := rd.value(DeviceName, dr, strconv.FormatInt(minInt-1, 10), strconv.FormatInt(maxInt, 10), db); err == nil {
		t.Fatalf("when using the minimum value of %s to minus 1, expected to get an error", dr)
	}
	if _, err := rd.value(DeviceName, dr, strconv.FormatInt(minInt, 10), strconv.FormatInt(maxInt+1, 10), db); err == nil {
		t.Fatalf("when using the maximum value of %s to plus 1, expected to get an error", dr)
	}

	rounds := 100

	//EnableRandomization = true
	for x := 1; x <= rounds; x++ {
		vn, err := rd.value(DeviceName, dr, strconv.FormatInt(minInt, 10), strconv.FormatInt(maxInt, 10), db)
		if err != nil {
			t.Fatal(err)
		}
		in, err := strconv.ParseInt(vn, 10, bitSize)
		if err != nil {
			t.Fatal(err)
		}
		var i1 int64
		if x == 1 {
			i1 = in
		}
		if i1 != in {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got same value in %d rounds", rounds)
		}
	}

	//generate value 100 times
	for x := 1; x <= rounds; x++ {
		v, err := rd.value(DeviceName, dr, strconv.FormatInt(minInt, 10), strconv.FormatInt(maxInt, 10), db)
		if err != nil {
			t.Fatal(err)
		}
		i, err := strconv.ParseInt(v, 10, bitSize)
		if err != nil {
			t.Fatal(err)
		}
		if i < minInt || i > maxInt {
			t.Fatalf("random value: %d,  out of range: %d ~ %d", i, minInt, maxInt)
		}
	}

	//EnableRandomization = false
	if err := db.startTransaction(); err != nil {
		t.Fatalf("Start a transaction failed: %v", err)
	}
	if err := db.exec(SQL_UPDATE_ENABLERANDOMIZATION, false, DeviceName, dr); err != nil {
		t.Fatal(err)
	}
	if err := db.commit(); err != nil {
		t.Fatal(err)
	}

	v1, _ := rd.value(DeviceName, dr, strconv.FormatInt(minInt, 10), strconv.FormatInt(maxInt, 10), db)
	i1, _ := strconv.ParseInt(v1, 10, bitSize)
	for x := 1; x <= rounds; x++ {
		v2, _ := rd.value(DeviceName, dr, strconv.FormatInt(minInt, 10), strconv.FormatInt(maxInt, 10), db)
		i2, _ := strconv.ParseInt(v2, 10, bitSize)
		if i1 != i2 {
			t.Fatalf("EnableRandomization is false, but got different value")
		}
	}
}

func SubTestValueUintx(t *testing.T, dr string) {
	db := getDb()
	if err := db.openDb(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	rd := newVirtualDevice()

	var minUint, maxUint uint64
	var bitSize int
	switch dr {
	case DeviceResource_Uint8:
		minUint, maxUint = rd.minUint8, rd.maxUint8
		bitSize = 8
	case DeviceResource_Uint16:
		minUint, maxUint = rd.minUint16, rd.maxUint16
		bitSize = 16
	case DeviceResource_Uint32:
		minUint, maxUint = rd.minUint32, rd.maxUint32
		bitSize = 32
	case DeviceResource_Uint64:
		minUint, maxUint = rd.minUint64, rd.maxUint64
		bitSize = 64
	default:
		t.Fatal("unknown device resource:", dr)
	}

	if _, err := rd.value(DeviceName, dr, strconv.FormatUint(minUint-1, 10), strconv.FormatUint(maxUint, 10), db); err == nil {
		t.Fatalf("when using the minimum value of %s to minus 1, expected to get an error", dr)
	}
	if _, err := rd.value(DeviceName, dr, strconv.FormatUint(minUint, 10), strconv.FormatUint(maxUint+1, 10), db); err == nil {
		t.Fatalf("when using the maximum value of %s to plus 1, expected to get an error", dr)
	}

	rounds := 100

	//EnableRandomization = true
	for x := 1; x <= rounds; x++ {
		vn, _ := rd.value(DeviceName, dr, strconv.FormatUint(minUint, 10), strconv.FormatUint(maxUint, 10), db)
		in, _ := strconv.ParseUint(vn, 10, bitSize)
		var i1 uint64
		if x == 1 {
			i1 = in
		}
		if i1 != in {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got same value in %d rounds", rounds)
		}
	}

	//generate value 100 times
	for x := 1; x <= rounds; x++ {
		v, err := rd.value(DeviceName, dr, strconv.FormatUint(minUint, 10), strconv.FormatUint(maxUint, 10), db)
		if err != nil {
			t.Fatal(err)
		}
		i, err := strconv.ParseUint(v, 10, bitSize)
		if err != nil {
			t.Fatal(err)
		}
		if i < minUint || i > maxUint {
			t.Fatalf("random value: %d,  out of range: %d ~ %d", i, minUint, maxUint)
		}
	}

	//EnableRandomization = false
	if err := db.startTransaction(); err != nil {
		t.Fatalf("Start a transaction failed: %v", err)
	}
	if err := db.exec(SQL_UPDATE_ENABLERANDOMIZATION, false, DeviceName, dr); err != nil {
		t.Fatal(err)
	}
	if err := db.commit(); err != nil {
		t.Fatal(err)
	}

	v1, _ := rd.value(DeviceName, dr, strconv.FormatUint(minUint, 10), strconv.FormatUint(maxUint, 10), db)
	i1, _ := strconv.ParseUint(v1, 10, bitSize)
	for x := 1; x <= rounds; x++ {
		v2, _ := rd.value(DeviceName, dr, strconv.FormatUint(minUint, 10), strconv.FormatUint(maxUint, 10), db)
		i2, _ := strconv.ParseUint(v2, 10, bitSize)
		if i1 != i2 {
			t.Fatalf("EnableRandomization is false, but got different value")
		}
	}
}

func SubTestValueFloatx(t *testing.T, dr string) {
	db := getDb()
	if err := db.openDb(); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := db.closeDb(); err != nil {
			t.Fatal(err)
		}
	}()

	rd := newVirtualDevice()

	var minFloat, maxFloat float64
	var bitSize, prec int
	switch dr {
	case DeviceResource_Float32:
		minFloat, maxFloat = rd.minFloat32, rd.maxFloat32
		bitSize, prec = 32, 6
	case DeviceResource_Float64:
		minFloat, maxFloat = rd.minFloat64, rd.maxFloat64
		bitSize, prec = 64, 15
	default:
		t.Fatal("unknown device resource:", dr)
	}

	if _, err := rd.value(DeviceName, dr, strconv.FormatFloat(minFloat-1, 'f', prec, bitSize), strconv.FormatFloat(maxFloat, 'f', prec, bitSize), db); err == nil {
		t.Fatalf("when using the minimum value of %s to minus 1, expected to get an error", dr)
	}
	if _, err := rd.value(DeviceName, dr, strconv.FormatFloat(minFloat, 'f', prec, bitSize), strconv.FormatFloat(maxFloat+1, 'f', prec, bitSize), db); err == nil {
		t.Fatalf("when using the maximum value of %s to plus 1, expected to get an error", dr)
	}

	rounds := 100

	//EnableRandomization = true
	for x := 1; x <= rounds; x++ {
		vn, _ := rd.value(DeviceName, dr, strconv.FormatFloat(minFloat, 'f', prec, bitSize), strconv.FormatFloat(maxFloat, 'f', prec, bitSize), db)
		fn, _ := strconv.ParseFloat(vn, bitSize)
		var f1 float64
		if x == 1 {
			f1 = fn
		}
		if f1 != fn {
			break
		}
		if x == rounds {
			t.Fatalf("EnableRandomization is true, but got same value in %d rounds", rounds)
		}
	}

	//generate value 100 times
	for x := 1; x <= rounds; x++ {
		v, err := rd.value(DeviceName, dr, strconv.FormatFloat(minFloat, 'f', prec, bitSize), strconv.FormatFloat(maxFloat, 'f', prec, bitSize), db)
		if err != nil {
			t.Fatal(err)
		}
		f, err := strconv.ParseFloat(v, bitSize)
		if err != nil {
			t.Fatal(err)
		}
		if f < minFloat || f > maxFloat {
			t.Fatalf("random value: %f,  out of range: %f ~ %f", f, minFloat, maxFloat)
		}
	}

	//EnableRandomization = false
	if err := db.startTransaction(); err != nil {
		t.Fatalf("Start a transaction failed: %v", err)
	}
	if err := db.exec(SQL_UPDATE_ENABLERANDOMIZATION, false, DeviceName, dr); err != nil {
		t.Fatal(err)
	}
	if err := db.commit(); err != nil {
		t.Fatal(err)
	}

	v1, _ := rd.value(DeviceName, dr, strconv.FormatFloat(minFloat, 'f', prec, bitSize), strconv.FormatFloat(maxFloat, 'f', prec, bitSize), db)
	f1, _ := strconv.ParseFloat(v1, bitSize)
	for x := 1; x <= rounds; x++ {
		v2, _ := rd.value(DeviceName, dr, strconv.FormatFloat(minFloat, 'f', prec, bitSize), strconv.FormatFloat(maxFloat, 'f', prec, bitSize), db)
		f2, _ := strconv.ParseFloat(v2, bitSize)
		if f1 != f2 {
			t.Fatalf("EnableRandomization is false, but got different value")
		}
	}
}
