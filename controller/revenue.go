package controller

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"main/apptypes"
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
	// SimpleExcel(c *gin.Context)
}

type RevenueCtrlImpl struct {
	svc *service.AppService
}

// func (a *RevenueCtrlImpl) SimpleExcel(c *gin.Context) {
// 	bcmp := c.MustGet("belongCmp")
// 	year, _ := strconv.Atoi(c.Query("year"))
// 	month, _ := strconv.Atoi(c.Query("month"))

// 	// Parsing the date range
// 	loc := time.Now().Location()
// 	fromDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)
// 	toDate := fromDate.AddDate(0, 1, 0)

// 	param := db.GetRevenueExcelParams{
// 		ApprovedDate:   sql.NullTime{Time: fromDate, Valid: true},
// 		ApprovedDate_2: sql.NullTime{Time: toDate, Valid: true},
// 		Belongcmp:      bcmp.(int64),
// 	}

// 	res, err := a.svc.RevenueServ.GetExcel(param)
// 	if err != nil {
// 		fmt.Println(err)
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		return
// 	}

// 	f := excelize.NewFile()
// 	for _, userRecord := range res {
// 		var record []db.GetRevenueExcelRow
// 		json.Unmarshal(userRecord, &record)
// 		sheetname := record.Username + strconv.Itoa(month) + "月報表"
// 		f.NewSheet(sheetname)

// 		header := []interface{}{"日期", "車號", "駕駛", "發貨地", "中轉", "卸貨地", "趟次", "運費", "應收款", "甲方", "司機運費", "油資", "維修", "備註", "業主"}
// 		f.SetSheetRow(sheetname, "A1", &header)

// 		rowIdx := 2
// 		for _, ls := range record.List {
// 			date, _ := time.Parse("2006-01-02", ls.Date)
// 			for _, job := range ls.Data {
// 				values := []interface{}{
// 					date.Format("2006/01/02"),
// 					job.Platenum,
// 					record.Username,
// 					job.FromLoc,
// 					job.Mid,
// 					job.Toloc,
// 					job.Count,
// 					job.Jp,
// 					job.Total,
// 					job.CmpName,
// 					"", // Driver freight
// 					ls.Gas,
// 					ls.Repair,
// 					"", // Notes
// 					job.Ss,
// 				}
// 				cell, _ := excelize.CoordinatesToCellName(1, rowIdx)
// 				f.SetSheetRow(sheetname, cell, &values)
// 				rowIdx++
// 			}
// 		}
// 	}

// 	filePath := fmt.Sprintf("./excel/RevenueReport_%d_%d.xlsx", year, month)
// 	if err := f.SaveAs(filePath); err != nil {
// 		c.AbortWithStatus(http.StatusInternalServerError)
// 		return
// 	}

//		c.File(filePath)
//	}
func (a *RevenueCtrlImpl) RevenueExcel(c *gin.Context) {
	// Get the company ID (belongCmp) from the context
	bcmp := c.MustGet("belongCmp")

	// Parse year and month from the query parameters
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

	// Get the current date
	tDate := time.Now()

	// Initialize variables for date ranges
	var AppFrom sql.NullTime
	var AppEnd sql.NullTime

	// Parse the start and end dates
	fm, err := time.Parse(time.DateOnly, strings.Split(time.Date(year, time.Month(month), 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	AppFrom.Scan(fm)

	me, err := time.Parse(time.DateOnly, strings.Split(time.Date(year, time.Month(month)+1, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	AppEnd.Scan(me)

	// Parse the company ID from the query, or use the one from the context
	var qBcmp int64
	if c.Query("cmp") != "" {
		tmp, err := strconv.Atoi(c.Query("cmp"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		qBcmp = int64(tmp)
	} else {
		qBcmp = bcmp.(int64)
	}

	// Prepare the parameters for the database query
	param := db.GetRevenueExcelParams{
		ApprovedDate:   AppFrom,
		ApprovedDate_2: AppEnd,
		Belongcmp:      qBcmp,
	}

	// Call the service to get the revenue data
	res, err := a.svc.RevenueServ.GetExcel(param)
	if err != nil {
		fmt.Println("Error fetching revenue data:", err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Create a new Excel file using the excelize library
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
		}
	}()

	// Loop over the data returned from the SQL query
	for _, userRecord := range res {
		fmt.Println(userRecord)
		// Marshal the user record to JSON bytes first
		jsonData, err := json.Marshal(userRecord)
		if err != nil {
			fmt.Println("Error marshaling userRecord to JSON:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Now unmarshal the JSON data into the desired struct
		var record apptypes.Excel
		err = json.Unmarshal(jsonData, &record)
		if err != nil {
			fmt.Println("Error unmarshaling JSON data:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Extract the date and create a sheet name
		date, err := time.Parse(time.DateOnly, record.List[0].Date)
		if err != nil {
			fmt.Println("Error parsing date:", err)
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		sheetname := record.Username + strconv.Itoa(int(date.Month())) + "月報表"

		// Create a new sheet in the Excel file
		_, err = f.NewSheet(sheetname)
		if err != nil {
			fmt.Println("Error creating new sheet:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Set the header row for the Excel sheet
		for _, row := range [][]interface{}{
			{"日期", "車號", "駕駛", "發貨地", "中轉", "卸貨地", "趟次", "運費", "應收款", "甲方", "司機運費", "油資", "維修", "備註", "業主"},
		} {
			cell, err := excelize.CoordinatesToCellName(1, 1)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			f.SetSheetRow(sheetname, cell, &row)
		}

		// Row counter for the Excel sheet
		rr := 0
		// Loop through the data for each day (record.List)
		for _, ls := range record.List {
			date, err := time.Parse(time.DateOnly, ls.Date)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}

			year := strconv.Itoa((date.Year()) - 1911)
			month := strconv.Itoa(int(date.Month()))
			day := strconv.Itoa(date.Day())

			// Loop through each job data entry (ls.Data)
			for ids, row := range ls.Data {
				// Set the date, platenum, and other fields
				cell, err := excelize.CoordinatesToCellName(1, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, year+month+day)

				// Platenum
				cell, err = excelize.CoordinatesToCellName(2, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, row.Platenum)

				// Driver Name (Username)
				cell, err = excelize.CoordinatesToCellName(3, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, record.Username)

				// From Location
				cell, err = excelize.CoordinatesToCellName(4, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, row.FromLoc)

				// Mid Location
				cell, err = excelize.CoordinatesToCellName(5, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, row.Mid)

				// To Location
				cell, err = excelize.CoordinatesToCellName(6, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, row.ToLoc)

				// Count (Trip count)
				cell, err = excelize.CoordinatesToCellName(7, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, row.Count)

				// Freight (運費)
				cell, err = excelize.CoordinatesToCellName(8, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, row.Jp)

				// Total amount (應收款)
				cell, err = excelize.CoordinatesToCellName(9, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, row.Total)

				// CMPT (甲方)
				cell, err = excelize.CoordinatesToCellName(10, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, row.CmpName)

				// Driver Freight (司機運費) - Empty for now
				cell, err = excelize.CoordinatesToCellName(11, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, "")

				// Gas (油資) - Only fill in the last row for each date
				cell, err = excelize.CoordinatesToCellName(12, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				if ids == len(ls.Data)-1 {
					f.SetCellValue(sheetname, cell, ls.Gas.Gas)
				} else {
					f.SetCellValue(sheetname, cell, "")
				}

				// Repair (維修) - Only fill in the last row for each date
				cell, err = excelize.CoordinatesToCellName(13, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				if ids == len(ls.Data)-1 {
					f.SetCellValue(sheetname, cell, ls.Repair.Repair)
				} else {
					f.SetCellValue(sheetname, cell, "")
				}

				// Remarks (備註)
				cell, err = excelize.CoordinatesToCellName(14, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, "")

				// Owner (業主)
				cell, err = excelize.CoordinatesToCellName(15, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, row.Ss)

				rr++
			}
		}
	}

	// Save the Excel file
	targetPath := time.DateOnly + ".xlsx"
	if err := f.SaveAs("./excel/" + targetPath); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Set the headers for downloading the Excel file
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
		y, m, _ := tDate.Date()
		today, err := time.Parse(time.DateOnly, strings.Split(tDate.String(), " ")[0])
		if err != nil {
			// fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		fom, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
		if err != nil {
			// fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var FinishedDate sql.NullTime
		FinishedDate.Scan(today)

		var FD2 sql.NullTime
		FD2.Scan(fom)

		param := db.GetDriverRevenueParams{
			Driverid: int64(id),
			Date:     FD2,
			Date_2:   FinishedDate,
		}
		res, err := a.svc.RevenueServ.GetRevenue(param)

		if err != nil {
			// fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		lm, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m-1, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])

		FD2.Scan(lm)

		lme, err := time.Parse(time.DateOnly, strings.Split(lm.AddDate(0, 1, -1).String(), " ")[0])
		FinishedDate.Scan(lme)

		param = db.GetDriverRevenueParams{
			Driverid: int64(id),
			Date:     FD2,
			Date_2:   FinishedDate,
		}
		resp1, err := a.svc.RevenueServ.GetRevenue(param)

		if err != nil {
			// fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		l2m, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m-2, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])

		FD2.Scan(l2m)

		l2me, err := time.Parse(time.DateOnly, strings.Split(l2m.AddDate(0, 1, -1).String(), " ")[0])
		FinishedDate.Scan(l2me)

		param = db.GetDriverRevenueParams{
			Driverid: int64(id),
			Date:     FD2,
			Date_2:   FinishedDate,
		}
		resp2, err := a.svc.RevenueServ.GetRevenue(param)

		if err != nil {
			// fmt.Print(err)
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
			// fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		fom, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])
		if err != nil {
			// fmt.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var FinishedDate sql.NullTime
		FinishedDate.Scan(today)

		var FD2 sql.NullTime
		FD2.Scan(fom)

		param := db.GetDriverRevenueByCmpParams{
			Belongcmp: int64(id),
			Date:      FD2,
			Date_2:    FinishedDate,
		}
		res, err := a.svc.RevenueServ.GetRevenueByCmp(param)

		if err != nil {
			// fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		lm, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m-1, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])

		FD2.Scan(lm)

		lme, err := time.Parse(time.DateOnly, strings.Split(lm.AddDate(0, 1, -1).String(), " ")[0])
		FinishedDate.Scan(lme)

		param = db.GetDriverRevenueByCmpParams{

			Belongcmp: int64(id),
			Date:      FD2,
			Date_2:    FinishedDate,
		}
		resp1, err := a.svc.RevenueServ.GetRevenueByCmp(param)

		if err != nil {
			// fmt.Print(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		l2m, err := time.Parse(time.DateOnly, strings.Split(time.Date(y, m-2, 1, 0, 0, 0, 0, tDate.Location()).String(), " ")[0])

		FD2.Scan(l2m)

		l2me, err := time.Parse(time.DateOnly, strings.Split(l2m.AddDate(0, 1, -1).String(), " ")[0])
		FinishedDate.Scan(l2me)

		param = db.GetDriverRevenueByCmpParams{

			Belongcmp: int64(id),
			Date:      FD2,
			Date_2:    FinishedDate,
		}
		resp2, err := a.svc.RevenueServ.GetRevenueByCmp(param)

		if err != nil {
			// fmt.Print(err)
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
