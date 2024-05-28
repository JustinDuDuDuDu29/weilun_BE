package controller

import (
	"database/sql"
	"fmt"
	"main/service"
	db "main/sql"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

type RevenueCtrl interface {
	RevenueDriver(c *gin.Context)
	RevenueExcel(c *gin.Context)
}

type RevenueCtrlImpl struct {
	svc *service.AppService
}

func (a *RevenueCtrlImpl) RevenueExcel(c *gin.Context) {

	year, err := strconv.Atoi(c.Query("year"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	month, err := strconv.Atoi(c.Query("month"))

	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	tDate := time.Now()

	var AppFrom sql.NullTime
	var AppEnd sql.NullTime
	fm, err := time.Parse(time.DateOnly, strings.Split(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	AppFrom.Scan(fm)

	me, err := time.Parse(time.DateOnly, strings.Split(time.Date(year, time.Month(month)+1, -1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	AppEnd.Scan(me)

	param := db.GetRevenueExcelParams{
		ApprovedDate:   AppFrom,
		ApprovedDate_2: AppEnd,
	}
	fmt.Println("AP1: ", AppFrom)
	fmt.Println("AP2: ", AppEnd)

	res, err := a.svc.RevenueServ.GetExcel(param)

	if err != nil {
		fmt.Println("err: ", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)

			c.AbortWithStatus(http.StatusBadRequest)
		}
	}()

	// for idx, row := range [][]interface{}{
	// 	{nil, "Apple", "Orange", "Pear"}, {"Small", 2, 3, 3},
	// 	{"Normal", 5, 2, 4}, {"Large", 6, 7, 8},
	// } {
	for idx, row := range [][]interface{}{
		{"日期", "車號", "駕駛", "發貨地", "中轉", "卸貨地", "趟次", "運費", "應收款", "甲方", "司機運費", "油資", "備註", "請款送單日"},
	} {
		cell, err := excelize.CoordinatesToCellName(1, idx+1)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetSheetRow("Sheet1", cell, &row)
	}
	for idx, row := range res {
		cell, err := excelize.CoordinatesToCellName(1, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		year := strconv.Itoa((row.Approveddate.Year()) - 1911)
		month := strconv.Itoa(int(row.Approveddate.Month()))
		day := strconv.Itoa(row.Approveddate.Day())

		f.SetCellValue("Sheet1", cell, (year + month + day))
		cell, err = excelize.CoordinatesToCellName(2, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Platenum.String)

		// cell, err = excelize.CoordinatesToCellName(3, idx+2)
		// if err != nil {
		// 	fmt.Println(err)
		// 	return
		// }
		// f.SetCellValue("Sheet1", cell, row.Name.String)

		cell, err = excelize.CoordinatesToCellName(3, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Username.String)

		cell, err = excelize.CoordinatesToCellName(4, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Fromloc)

		cell, err = excelize.CoordinatesToCellName(5, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Mid.String)

		cell, err = excelize.CoordinatesToCellName(6, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Toloc)

		cell, err = excelize.CoordinatesToCellName(7, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Count)

		cell, err = excelize.CoordinatesToCellName(8, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Price)

		cell, err = excelize.CoordinatesToCellName(9, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Totalprice)

		cell, err = excelize.CoordinatesToCellName(10, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Source)

		cell, err = excelize.CoordinatesToCellName(11, idx+2)
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetCellValue("Sheet1", cell, row.Togive)
	}
	targetPath := time.DateOnly + ".xlsx"
	if err := f.SaveAs("./excel/" + targetPath); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+targetPath)
	c.Header("Content-Type", "application/octet-stream")
	c.File("./excel/" + targetPath)

}
func (a *RevenueCtrlImpl) RevenueDriver(c *gin.Context) {

	uid := c.MustGet("UserID")
	role := c.MustGet("Role").(int16)
	bcmp := c.MustGet("belongCmp")

	qid := c.Query("id")
	qcmp := c.Query("cmp")

	if qid != "" {

		id, err := strconv.Atoi(qid)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		info, err := a.svc.UserServ.GetUserById(int64(id))
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if (role == 300 && id != uid) || (role == 200 && info.Belongcmp != bcmp) {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tDate := time.Now()
		fmt.Println("Today: ", tDate)
		y, m, _ := tDate.Date()
		today, err := time.Parse(time.DateOnly, strings.Split(tDate.String(), " ")[0])
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		fom, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var FinishedDate sql.NullTime
		FinishedDate.Scan(today)

		var FD2 sql.NullTime
		FD2.Scan(fom)

		param := db.GetDriverRevenueParams{
			Driverid:       int64(id),
			FinishedDate:   FD2,
			FinishedDate_2: FinishedDate,
		}
		res, err := a.svc.RevenueServ.GetRevenue(param)

		if err != nil {
			fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		lm, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m-1, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])

		FD2.Scan(lm)

		lme, err := time.Parse(time.DateOnly, strings.Split(lm.AddDate(0, 1, -1).String(), " ")[0])
		FinishedDate.Scan(lme)

		param = db.GetDriverRevenueParams{
			Driverid:       int64(id),
			FinishedDate:   FD2,
			FinishedDate_2: FinishedDate,
		}
		resp1, err := a.svc.RevenueServ.GetRevenue(param)

		if err != nil {
			fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		l2m, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m-2, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])

		FD2.Scan(l2m)

		l2me, err := time.Parse(time.DateOnly, strings.Split(l2m.AddDate(0, 1, -1).String(), " ")[0])
		FinishedDate.Scan(l2me)

		param = db.GetDriverRevenueParams{
			Driverid:       int64(id),
			FinishedDate:   FD2,
			FinishedDate_2: FinishedDate,
		}
		resp2, err := a.svc.RevenueServ.GetRevenue(param)

		if err != nil {
			fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resp2 = append(resp2, resp1...)
		resp2 = append(resp2, res...)

		c.JSON(http.StatusOK, resp2)
	}

	if qcmp != "" {

		id, err := strconv.Atoi(qcmp)
		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if role == 200 && bcmp != id || role >= 300 {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tDate := time.Now()
		y, m, _ := tDate.Date()
		today, err := time.Parse(time.DateOnly, strings.Split(tDate.String(), " ")[0])
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		fom, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
		if err != nil {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var FinishedDate sql.NullTime
		FinishedDate.Scan(today)

		var FD2 sql.NullTime
		FD2.Scan(fom)

		param := db.GetDriverRevenueByCmpParams{
			Belongcmp:      int64(id),
			FinishedDate:   FD2,
			FinishedDate_2: FinishedDate,
		}
		res, err := a.svc.RevenueServ.GetRevenueByCmp(param)

		if err != nil {
			fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		lm, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m-1, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])

		FD2.Scan(lm)

		lme, err := time.Parse(time.DateOnly, strings.Split(lm.AddDate(0, 1, -1).String(), " ")[0])
		FinishedDate.Scan(lme)

		param = db.GetDriverRevenueByCmpParams{

			Belongcmp:      int64(id),
			FinishedDate:   FD2,
			FinishedDate_2: FinishedDate,
		}
		resp1, err := a.svc.RevenueServ.GetRevenueByCmp(param)

		if err != nil {
			fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		l2m, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m-2, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])

		FD2.Scan(l2m)

		l2me, err := time.Parse(time.DateOnly, strings.Split(l2m.AddDate(0, 1, -1).String(), " ")[0])
		FinishedDate.Scan(l2me)

		param = db.GetDriverRevenueByCmpParams{

			Belongcmp:      int64(id),
			FinishedDate:   FD2,
			FinishedDate_2: FinishedDate,
		}
		resp2, err := a.svc.RevenueServ.GetRevenueByCmp(param)

		if err != nil {
			fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		resp2 = append(resp2, resp1...)
		resp2 = append(resp2, res...)

		c.JSON(http.StatusOK, resp2)
	}

	c.AbortWithStatus(http.StatusBadRequest)
	return
}

func RevenueCtrlInit(svc *service.AppService) *RevenueCtrlImpl {
	return &RevenueCtrlImpl{
		svc: svc,
	}
}
