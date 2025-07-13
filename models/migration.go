package models


var modelsForMigration=[]interface{}{
	&Expert{},
	&User{},
	&AvailabilitySlot{},
}


func GetMigrationModel() []interface{} {
	return modelsForMigration
}
