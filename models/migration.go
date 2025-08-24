package models


var modelsForMigration=[]interface{}{
	&User{},
	&Expert{},
	&AvailabilitySlot{},
	&Payment{},
	&Student{},
}


func GetMigrationModel() []interface{} {
	return modelsForMigration
}
