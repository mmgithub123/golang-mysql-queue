package main

import (
	"fmt"
	"strconv"
	"strings"

	//"mysqltaskdirdemo/makeData"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type VsTask struct {
	ID       uint
	Sha256   string
	FileType uint8
}

type VsDir struct {
	ID       uint
	VsTaskId string
}

func CreateVsDir(db *gorm.DB) error {
	// Note the use of tx as the database handle once you are within a transaction
	//get min id from vs_tasks
	//get five records from vs_tasks
	//insert into this five records to vs_dirs
	//delete this five records from vs_tasks
	//commit
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	var ID string
	if err := tx.Raw("SELECT MIN(ID) FROM vs_tasks").Scan(&ID).Error; err != nil {
		tx.Rollback()
		return err
	}

	var IDs []string
	IDnum, err := strconv.Atoi(ID)
	if err != nil {
		panic(err)
	}
	upID := IDnum + 4
	if err := tx.Raw("select ID from vs_tasks WHERE ID >= ? and ID <= ? for update", ID, upID).Scan(&IDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	VsTaskIdStr := strings.Join(IDs, ", ")
	vsdir := VsDir{VsTaskId: VsTaskIdStr}
	if err := tx.Create(&vsdir).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Exec("delete from vs_tasks WHERE ID IN ?", IDs).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func main() {
	fmt.Println("start")
	dsn := ":@tcp(192:3306)/?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	//makeData.InsertTableRand(db)
	for {
		CreateVsDir(db)
	}
}
