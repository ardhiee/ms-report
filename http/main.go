package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	_ "github.com/go-sql-driver/mysql"
)

type activity struct {
	Timestamp  string
	UserID     string
	Activities string
}

// Query the database
func selectactivities(userID string) []activity {
	var result []activity

	// connect to db
	db, err := connectdb()
	if err != nil {
		fmt.Println("Can't connect DB")
		fmt.Println(err.Error())
	}
	defer db.Close()

	// query the activities
	var userid = userID
	rows, err := db.Query("select timestamp, userid, activities from tb_report where userid like ? order by timestamp desc", userid)
	if err != nil {
		fmt.Println("Can't query the table")
		fmt.Println(err.Error())
	}
	defer db.Close()

	for rows.Next() {
		var each = activity{}
		var err = rows.Scan(&each.Timestamp, &each.UserID, &each.Activities)

		if err != nil {
			fmt.Println(err.Error())
		}

		result = append(result, each)

	}

	return result
}

// JSON WEB API
func useractivities(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var id = r.FormValue("userid")
	if id == "" {
		id = "%"
	}
	var data = selectactivities(id)

	if r.Method == "GET" {
		var result, err = json.Marshal(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}
	http.Error(w, "", http.StatusBadRequest)

}

// JSON WEB API
func generatexls(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// GenerateExcel()
	var id = r.FormValue("userid")
	if id == "" {
		id = "%"
	}

	GenerateExcel(id)

	if r.Method == "GET" {
		var result, err = json.Marshal("success generate excel")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(result)
		return
	}
	http.Error(w, "", http.StatusBadRequest)

}

// HTTP Server
func main() {
	http.HandleFunc("/ms-report/useractivities", useractivities)
	http.HandleFunc("/ms-report/generatexls", generatexls)

	fmt.Println("Starting Web Server at http://localhost:8088/")
	http.ListenAndServe(":8088", nil)
}

// connect to database
func connectdb() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp("+os.Getenv("DBCONN")+")/poc")
	if err != nil {
		fmt.Println("Fail 1")
		panic("dbpool init >> " + err.Error())
	}

	return db, nil
}

// GenerateExcel for activities in form of xlsx
func GenerateExcel(UserID string) {
	// M type with interface, after that initialize data excelreportdata
	type M map[string]interface{}
	var excelreportdata = []M{}
	var id = UserID

	xlsx := excelize.NewFile()
	var activities = selectactivities(id)

	// convert the string Array from DB into Map type
	for _, activity := range activities {
		var temp = M{}
		temp["Timestamp"] = activity.Timestamp
		temp["UserID"] = activity.UserID
		temp["Activities"] = activity.Activities
		excelreportdata = append(excelreportdata, temp)
	}

	sheet1Name := "Sheet One"
	xlsx.SetSheetName(xlsx.GetSheetName(1), sheet1Name)

	xlsx.SetCellValue(sheet1Name, "A1", "Timestamp")
	xlsx.SetCellValue(sheet1Name, "B1", "UserID")
	xlsx.SetCellValue(sheet1Name, "C1", "Activities")

	err := xlsx.AutoFilter(sheet1Name, "A1", "C1", "")
	if err != nil {
		log.Fatal("ERROR", err.Error())
	}

	for i, each := range excelreportdata {
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("A%d", i+2), each["Timestamp"])
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("B%d", i+2), each["UserID"])
		xlsx.SetCellValue(sheet1Name, fmt.Sprintf("C%d", i+2), each["Activities"])
	}

	t := time.Now().Format("20060102-150405")

	err = xlsx.SaveAs("./Activities-Report-" + t + ".xlsx")
	if err != nil {
		fmt.Println(err)
	}

}
