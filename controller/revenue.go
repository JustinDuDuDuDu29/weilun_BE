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
	SimpleExcel(c *gin.Context)
}

type RevenueCtrlImpl struct {
	svc *service.AppService
}

func (a *RevenueCtrlImpl) SimpleExcel(c *gin.Context) {
	// Get the company ID (belongCmp) from the context
	bcmp := c.MustGet("belongCmp")
	role := c.MustGet("Role")

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

	// Initialize variables for date ranges
	var AppFrom, AppEnd sql.NullTime

	// Parse the start and end dates for the given month
	fm, err := time.Parse(time.DateOnly, fmt.Sprintf("%04d-%02d-01", year, month))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	AppFrom.Scan(fm)

	// Calculate the end date for the given month
	nextMonth := month + 1
	if nextMonth > 12 {
		year++
		nextMonth = 1
	}
	me, err := time.Parse(time.DateOnly, fmt.Sprintf("%04d-%02d-01", year, nextMonth))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	AppEnd.Scan(me)

	// Parse the company ID from the query or use the context value
	var qBcmp sql.NullInt64
	if cmp := c.Query("cmpid"); cmp != "" {
		tmp, err := strconv.Atoi(c.Query("cmpid"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		qBcmp.Scan(tmp)
		if role == 200 {
			qBcmp.Scan(bcmp)
		}
	} else {
		c.AbortWithStatus(http.StatusBadRequest)
	}

	var userid sql.NullInt64

	userid.Valid = false

	// Prepare parameters for the database query
	param := db.GetJobCmpParams{
		ApprovedDate:   AppFrom,
		ApprovedDate_2: AppEnd,
		CmpId:          qBcmp,
		UserId:         userid,
	}

	// Fetch the data from the service
	res, err := a.svc.RevenueServ.GetSimpleExcel(param)
	if err != nil {
		fmt.Println("Error fetching revenue data:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Create a new Excel file using the excelize library
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}
	}()

	// Loop over the data returned from the SQL query
	for _, userRecord := range res {
		// Marshal the user record to JSON bytes
		jsonData, err := json.Marshal(userRecord)
		if err != nil {
			fmt.Println("Error marshaling userRecord to JSON:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Now unmarshal the JSON data into the desired struct
		var record apptypes.AutoGenerated
		err = json.Unmarshal(jsonData, &record)
		if err != nil {
			fmt.Println("Error unmarshaling JSON data:", err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Loop through each company in the record
		sheetname := record.Name + strconv.Itoa(int(year)) + "年" + strconv.Itoa(int(month)) + "月報表"

		// Create a new sheet for each company
		f.NewSheet(sheetname)

		// Set the header row for the Excel sheet
		headers := []interface{}{"駕駛", "趟次", "應收款", "油資", "維修", "應收款 - 油資 + 維修"}
		if err := f.SetSheetRow(sheetname, "A1", &headers); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		// Row counter for the Excel sheet
		rowIndex := 2

		// Loop through the data for each user in the company
		for _, user := range record.Users {
			// Calculate the new volume: GasTotal + RepairTotal
			newVolume := user.JobTotal - user.GasTotal + user.RepairTotal

			// Prepare the row data, including the new volume
			row := []interface{}{
				user.UserName,    // 駕駛 (Username)
				user.JobCount,    // 趟次 (Job count)
				user.JobTotal,    // 應收款 (Job total)
				user.GasTotal,    // 油資 (Gas total)
				user.RepairTotal, // 維修 (Repair total)
				newVolume,        // 新增總額: 油資 + 維修
			}

			// Set each row of data in the sheet
			cell, err := excelize.CoordinatesToCellName(1, rowIndex)
			if err != nil {
				c.AbortWithStatus(http.StatusBadRequest)
				return
			}
			f.SetSheetRow(sheetname, cell, &row)
			rowIndex++
		}
	}

	// Save the Excel file to the server
	targetPath := fmt.Sprintf("%04d-%02d-Report.xlsx", year, month)
	if err := f.SaveAs("./excel/" + targetPath); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// Set the headers for downloading the Excel file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+targetPath)
	c.Header("Content-Type", "application/octet-stream")
	c.File("./excel/" + targetPath)
}

func (a *RevenueCtrlImpl) RevenueExcel(c *gin.Context) {
	// Get the company ID (belongCmp) from the context
	bcmp := c.MustGet("belongCmp")

	// Parse year and month from the query parameters
	year, err := strconv.Atoi(c.Query("year"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	total := 0

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
	fmt.Println(AppFrom)

	fmt.Println(AppEnd)
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

			year := date.Year() - 1911 // 民國年
			month := int(date.Month())
			day := date.Day()

			yearStr := strconv.Itoa(year)
			monthStr := fmt.Sprintf("%02d", month) // 補前導0
			dayStr := fmt.Sprintf("%02d", day)     // 補前導0

			formattedDate := yearStr + monthStr + dayStr
			// fmt.Println(formattedDate)

			// Loop through each job data entry (ls.Data)
			for _, row := range ls.Data {
				// fmt.Println(row)
				// Set the date, platenum, and other fields
				cell, err := excelize.CoordinatesToCellName(1, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				f.SetCellValue(sheetname, cell, formattedDate)

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
				total += int(row.Total)
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
				// fmt.Println(ls.Gas.Gas)
				// if ids == len(ls.Data)-1 {
				f.SetCellValue(sheetname, cell, ls.Gas.Gas)
				// } else {
				// f.SetCellValue(sheetname, cell, "")
				// }

				// Repair (維修) - Only fill in the last row for each date
				cell, err = excelize.CoordinatesToCellName(13, rr+2)
				if err != nil {
					c.AbortWithStatus(http.StatusBadRequest)
					return
				}
				// if ids == len(ls.Data)-1 {
				f.SetCellValue(sheetname, cell, ls.Repair.Repair)
				// } else {
				// f.SetCellValue(sheetname, cell, "")
				// }

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

	sheetname := "總和"
	_, err = f.NewSheet(sheetname)
	if err != nil {
		fmt.Println("Error creating new sheet:", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, row := range [][]interface{}{
		{"應收款總和"},
	} {
		cell, err := excelize.CoordinatesToCellName(1, 1)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		f.SetSheetRow(sheetname, cell, &row)
	}
	cell, err := excelize.CoordinatesToCellName(1, 2)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	f.SetCellValue(sheetname, cell, total)

	now := time.Now()
	targetPath := fmt.Sprintf("%d%02d%02d.xlsx", now.Year()-1911, now.Month(), now.Day())

	if err := f.SaveAs("./excel/" + targetPath); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	// Set the headers for downloading the Excel file
	c.Header("Content-Disposition", "attachment; filename=\""+targetPath+"\"")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Expires", "0")
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
