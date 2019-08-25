package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"os"

	"fmt"
	"ms-report/config"
	"ms-report/model"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql"
)

// connect to database
func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp("+os.Getenv("DBCONN")+")/poc")
	if err != nil {
		fmt.Println("Fail 1")
		panic("dbpool init >> " + err.Error())
	}

	return db, nil
}

// insert the report function
func insertReport(param *model.Report) {
	db, err := connect()
	if err != nil {
		fmt.Println("Cant connect to DB.")
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into tb_report values (?, ?, ?, ?, ?, ?, ?, ?)")
	stmt.Exec(param.Timestamp, param.Userid, param.Activities, param.Applicationid, param.Cif, param.Accountno, param.Ktpid, param.Status)

	if err != nil {
		fmt.Println("Cant insert to DB.")
		fmt.Println(err.Error())
		return
	}
	fmt.Println("insert success!")
}

// query the data
func selectReport(appUserID string) *model.ReportList {

	var result = &model.ReportList{}

	db, err := connect()
	if err != nil {
		fmt.Println("Cant connect to DB.")
		fmt.Println(err.Error())
	}
	defer db.Close()

	var userid = appUserID
	rows, err := db.Query("select timestamp, userid, activities, applicationid, cif, accountno, ktpid, status from tb_report where userid = ?", userid)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var each = &model.Report{}
		var err = rows.Scan(&each.Timestamp, &each.Userid, &each.Activities, &each.Applicationid, &each.Cif, &each.Accountno, &each.Ktpid, &each.Status)

		if err != nil {
			fmt.Println(err.Error())
		}
		result.List = append(result.List, each)
	}
	return result
}

// ReportServer data type
type ReportsServer struct{}

// AddReport grpc function server to insert into database
func (ReportsServer) AddReport(ctx context.Context, param *model.Report) (*empty.Empty, error) {
	log.Println("Adding new report", param.String())
	go insertReport(param)
	return new(empty.Empty), nil
}

// List grpc function server to query from database
func (ReportsServer) List(ctx context.Context, param *model.ApplicationUserId) (*model.ReportList, error) {
	appUserID := param.UserId

	fmt.Println(appUserID)

	return selectReport(appUserID), nil
}

func main() {

	srv := grpc.NewServer()
	var userSrv ReportsServer
	model.RegisterReportsServer(srv, userSrv)

	log.Println("Starting RPC server at", config.SERVICE_REPORT_PORT)

	l, err := net.Listen("tcp", config.SERVICE_REPORT_PORT)
	if err != nil {
		log.Fatalf("could not listen to %s: %v", config.SERVICE_REPORT_PORT, err)
	}

	log.Fatal(srv.Serve(l))

}
