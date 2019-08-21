package main

import (
	"context"
	"database/sql"
	"log"
	"net"

	"fmt"
	"ms-report/config"
	"ms-report/model"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"

	_ "github.com/go-sql-driver/mysql"
)

var localStorage *model.ReportList

func init() {
	localStorage = new(model.ReportList)
	localStorage.List = make([]*model.Report, 0)
}

func connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/poc")
	if err != nil {
		panic("dbpool init >> " + err.Error())
	}

	return db, nil
}

func insertReport(param *model.Report) {
	db, err := connect()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into tb_report values (?, ?, ?, ?, ?, ?, ?)")
	//stmt.Exec(timestamp, userid, activities, applicationid, cif, accountno, ktpid)
	stmt.Exec(param.Timestamp, param.Userid, param.Activities, param.Applicationid, param.Cif, param.Accountno, param.Ktpid)
	//stmt.Exec("1", "1", "1", "1", "1", "1", "1")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("insert success!")
}

type ReportsServer struct{}

func (ReportsServer) AddReport(ctx context.Context, param *model.Report) (*empty.Empty, error) {
	//localStorage.List = append(localStorage.List, param)

	log.Println("Adding new report", param.String())
	log.Println("Adding new report", param.Applicationid)

	insertReport(param)
	//insertReport(param.Timestamp, param.Userid, param.Activities, param.Applicationid, param.Cif, param.Accountno, param.Ktpid)

	return new(empty.Empty), nil
}

func (ReportsServer) List(ctx context.Context, void *model.UserId) (*model.ReportList, error) {
	return localStorage, nil
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
