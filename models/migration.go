package models


var modelsForMigration=[]interface{}{
	&Expert{},
	&User{},

}


func GetMigrationModel() []interface{} {
	return modelsForMigration
}
