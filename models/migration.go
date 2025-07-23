package models


var modelsForMigration=[]interface{}{
	&Expert{},
	&User{},
	&AvailabilitySlot{},
	&Payment{},
	&Student{},
}


func GetMigrationModel() []interface{} {
	return modelsForMigration
}
