package database

import (
	"database/sql" //gives access to db for fns like db connections and running queries
	"log"          //used to log msgs or errors

	_ "modernc.org/sqlite" //this is sqlite driver it allows GO to talk w SQLite db
)

var DB *sql.DB //'DB' is a globbal var. that'll store db connection

func InitDB() *sql.DB { //initializes the SQLite db connection
	db, err := sql.Open("sqlite", "./todo.db") //todo.db file will store in the same folder as ur GO project
	if err != nil {

		log.Fatal("Cannot open database:", err)
	} //if smtg goes wrong while opening the db, it'll show and error and stop the prog.
	DB = db   //if everything's alright save this open connection in 'DB' global var
	return DB //return DB conn. so it can be used by other parts of  project
}
