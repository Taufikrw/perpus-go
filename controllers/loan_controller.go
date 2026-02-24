package controllers

import (
	"belajar-go/config"
	"belajar-go/dto"
	"belajar-go/models"
	"belajar-go/resources"
	"belajar-go/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func IndexLoans(c *gin.Context) {
	var loans []models.Loan
	config.DB.Preload("Member.User.Role").Preload("BookItem.Book.Category").Find(&loans)
	utils.SendResponse(c, http.StatusOK, "Daftar peminjaman berhasil diambil!", resources.FormatLoans(loans))
}

func StoreLoan(c *gin.Context) {
	var input dto.CreateLoanDTO

	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Input tidak valid", errMsg)
		return
	}

	var member models.Member
	if err := config.DB.Where("id = ?", input.MemberID).Take(&member).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Member tidak ditemukan", nil)
		return
	}
	var bookItem models.BookItem
	if err := config.DB.Where("id = ?", input.BookItemID).Take(&bookItem).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Book item tidak ditemukan", nil)
		return
	}
	if bookItem.Status != "available" {
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Book item tidak tersedia untuk dipinjam", nil)
		return
	}

	tx := config.DB.Begin()
	bookItem.Status = "loaned"
	if err := tx.Save(&bookItem).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui status book item", nil)
		return
	}

	loanDate, _ := time.Parse("2006-01-02", input.LoanDate)
	dueDate, _ := time.Parse("2006-01-02", input.DueDate)

	if !dueDate.After(loanDate) {
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Validasi Gagal", []string{
			"Tanggal kembali (due date) harus lebih besar dari tanggal peminjaman (loan date)",
		})
		return
	}

	newLoan := models.Loan{
		MemberID:   member.ID,
		BookItemID: bookItem.ID,
		LoanDate:   loanDate,
		DueDate:    dueDate,
		Status:     "ongoing",
	}
	if err := config.DB.Create(&newLoan).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menyimpan data peminjaman", nil)
		return
	}
	tx.Commit()

	config.DB.Preload("Member.User.Role").Preload("BookItem.Book.Category").Take(&newLoan, newLoan.ID)
	utils.SendResponse(c, http.StatusCreated, "Peminjaman berhasil dibuat", resources.FormatLoan(newLoan))
}

func ShowLoan(c *gin.Context) {
	var loan models.Loan

	if err := config.DB.Where("id = ?", c.Param("id")).Preload("Member.User.Role").Preload("BookItem.Book.Category").Take(&loan).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Peminjaman tidak ditemukan!", nil)
		return
	}
	utils.SendResponse(c, http.StatusOK, "Peminjaman berhasil diambil!", resources.FormatLoan(loan))
}

func UpdateLoan(c *gin.Context) {
	var loan models.Loan
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&loan).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Peminjaman tidak ditemukan!", nil)
		return
	}

	var input dto.UpdateLoanDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		errMsg := utils.FormatError(err)
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Input tidak valid", errMsg)
		return
	}

	var member models.Member
	if err := config.DB.Where("id = ?", input.MemberID).Take(&member).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Member tidak ditemukan", nil)
		return
	}
	var bookItem models.BookItem
	if err := config.DB.Where("id = ?", input.BookItemID).Take(&bookItem).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Book item tidak ditemukan", nil)
		return
	}

	loanDate, _ := time.Parse("2006-01-02", input.LoanDate)
	dueDate, _ := time.Parse("2006-01-02", input.DueDate)
	var returnDate *time.Time
	if input.ReturnDate != "" {
		parsedDate, err := time.Parse("2006-01-02", input.ReturnDate)
		if err == nil {
			returnDate = &parsedDate
		}
	}

	if !dueDate.After(loanDate) {
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Validasi Gagal", []string{
			"Tanggal kembali (due date) harus lebih besar dari tanggal peminjaman (loan date)",
		})
		return
	}
	if input.ReturnDate != "" && !returnDate.After(loanDate) {
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Validasi Gagal", []string{
			"Tanggal pengembalian (return date) harus lebih besar dari tanggal peminjaman (loan date)",
		})
		return
	}

	newLoan := models.Loan{
		MemberID:   member.ID,
		BookItemID: bookItem.ID,
		LoanDate:   loanDate,
		DueDate:    dueDate,
		ReturnDate: returnDate,
		Status:     input.Status,
	}

	config.DB.Model(&loan).Updates(newLoan)
	config.DB.Preload("Member.User.Role").Preload("BookItem.Book.Category").Take(&newLoan, loan.ID)
	utils.SendResponse(c, http.StatusOK, "Peminjaman berhasil diupdate!", resources.FormatLoan(newLoan))
}

func ReturnLoan(c *gin.Context) {
	var loan models.Loan
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&loan).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Peminjaman tidak ditemukan!", nil)
		return
	}

	if loan.Status != "ongoing" {
		utils.SendErrorResponse(c, http.StatusUnprocessableEntity, "Peminjaman sudah dikembalikan atau melewati tanggal kembali", nil)
		return
	}

	tx := config.DB.Begin()
	now := time.Now()
	loan.ReturnDate = &now
	if now.After(loan.DueDate) {
		loan.Status = "overdue"
		newFine := models.Fine{
			LoanID: loan.ID,
			Amount: 5000 * float64(int(now.Sub(loan.DueDate).Hours()/24)),
		}
		tx.Create(&newFine)
	} else {
		loan.Status = "returned"
	}

	config.DB.Save(&loan)

	var bookItem models.BookItem
	config.DB.Where("id = ?", loan.BookItemID).Take(&bookItem)
	bookItem.Status = "available"
	config.DB.Save(&bookItem)

	tx.Commit()
	config.DB.Preload("Member.User.Role").Preload("BookItem.Book.Category").Take(&loan, loan.ID)
	utils.SendResponse(c, http.StatusOK, "Peminjaman berhasil dikembalikan!", resources.FormatLoan(loan))
}

func DeleteLoan(c *gin.Context) {
	var loan models.Loan
	if err := config.DB.Where("id = ?", c.Param("id")).Take(&loan).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Peminjaman tidak ditemukan!", nil)
		return
	}

	tx := config.DB.Begin()

	if err := tx.Where("id = ?", loan.ID).Delete(&loan).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal menghapus peminjaman", nil)
		return
	}

	var bookItem models.BookItem
	config.DB.Where("id = ?", loan.BookItemID).Take(&bookItem)
	bookItem.Status = "available"
	if err := tx.Save(&bookItem).Error; err != nil {
		tx.Rollback()
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui status book item", nil)
		return
	}

	tx.Commit()
	utils.SendResponse(c, http.StatusOK, "Peminjaman berhasil dihapus!", nil)
}

func PayFine(c *gin.Context) {
	var fine models.Fine
	if err := config.DB.Where("id = ? AND paid_at IS NULL", c.Param("id")).Take(&fine).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusNotFound, "Denda tidak ditemukan atau sudah dibayar!", nil)
		return
	}

	now := time.Now()
	fine.PaidAt = &now

	if err := config.DB.Save(&fine).Error; err != nil {
		utils.SendErrorResponse(c, http.StatusInternalServerError, "Gagal memperbarui status denda", nil)
		return
	}

	utils.SendResponse(c, http.StatusOK, "Denda berhasil dibayar!", nil)
}
