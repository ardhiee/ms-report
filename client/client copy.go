package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"ms-report/config"
	"ms-report/model"

	"google.golang.org/grpc"
)

func serviceReport() model.ReportsClient {
	port := config.SERVICE_REPORT_PORT
	conn, err := grpc.Dial(port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("could not connect to", port, err)
	}

	return model.NewReportsClient(conn)
}

func main() {

	report1 := model.Report{
		Timestamp:     "2019-08-20",
		Userid:        "user001",
		Activities:    "Login",
		Applicationid: "appID-1234",
		Cif:           "123456",
		Accountno:     "000000",
		Ktpid:         "1723871273821",
	}

	report2 := model.Report{
		Timestamp:     "2019-08-20",
		Userid:        "user002",
		Activities:    "Login",
		Applicationid: "appID-1234",
		Cif:           "123456",
		Accountno:     "000000",
		Ktpid:         "1723871273821",
	}

	report := serviceReport()

	fmt.Println("\n", "===========> report test")
	report.AddReport(context.Background(), &report1)
	report.AddReport(context.Background(), &report2)

	// show all garages of all user
	// fmt.Println("\n", "===========> report all")
	// res1, err := report.List(context.Background(), new(empty.Empty))
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// res1String, _ := json.Marshal(res1.List)
	// log.Println(string(res1String))

	// show specific report of user
	fmt.Println("\n", "===========> show only user 2")
	res2, err := report.List(context.Background(), &model.UserId{UserId: report2.Userid})
	if err != nil {
		log.Fatal(err.Error())
	}
	res2String, _ := json.Marshal(res2.List)
	log.Println(string(res2String))

}
