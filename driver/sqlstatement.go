package driver

const (
	SQL_DROP_TABLE                 = "DROP TABLE IF EXISTS VIRTUAL_RESOURCE;"
	SQL_CREATE_TABLE               = "CREATE TABLE VIRTUAL_RESOURCE (DEVICE_NAME string, COMMAND_NAME String, DEVICE_RESOURCE_NAME string, ENABLE_RANDOMIZATION bool, DATA_TYPE String, VALUE string);"
	SQL_SELECT                     = "SELECT DEVICE_NAME, COMMAND_NAME, DEVICE_RESOURCE_NAME, ENABLE_RANDOMIZATION, DATA_TYPE, VALUE FROM VIRTUAL_RESOURCE where DEVICE_NAME==$1 and DEVICE_RESOURCE_NAME==$2;"
	SQL_INSERT                     = "INSERT INTO VIRTUAL_RESOURCE VALUES ($1, $2, $3, $4, $5, $6);"
	SQL_UPDATE_ENABLERANDOMIZATION = "UPDATE VIRTUAL_RESOURCE SET ENABLE_RANDOMIZATION=$1 WHERE DEVICE_NAME==$2 AND DEVICE_RESOURCE_NAME==$3;"
	SQL_UPDATE_VALUE               = "UPDATE VIRTUAL_RESOURCE SET VALUE=$1 WHERE DEVICE_NAME==$2 AND DEVICE_RESOURCE_NAME==$3;"
)
