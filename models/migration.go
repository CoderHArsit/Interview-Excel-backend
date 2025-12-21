package models


var modelsForMigration=[]interface{}{
	&User{},
	&Expert{},
	&AvailabilitySlot{},
	&Payment{},
	&Student{},
	&Session{},
}


func GetMigrationModel() []interface{} {
	return modelsForMigration
}
