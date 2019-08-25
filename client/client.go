package main

import (
	"context"
	"fmt"
	"log"
	"time"

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

	t := time.Now()
	//fmt.Println(t.Format("20060102150405"))

	report1 := model.Report{
		Timestamp:     t.Format("2006-01-02 15:04:05.000-07:00"),
		Userid:        "ardiismail",
		Activities:    "Login",
		Applicationid: "userID-1234",
		Cif:           "123456",
		Accountno:     "000000",
		Ktpid:         "1723871273821",
		Status:        "Approved",
	}

	report := serviceReport()

	fmt.Println("===========> report test1")
	report.AddReport(context.Background(), &report1)
	fmt.Println("===========> report test2")

	// show all registered users
	// res1, err := report.List(context.Background(), &model.ApplicationUserId{UserId: "ardiismail"})
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	// res1String, _ := json.Marshal(res1.List)
	// log.Println(string(res1String))

}
