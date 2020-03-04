package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/victorcmoura/go-redcoins-api/api/auth"
	"github.com/victorcmoura/go-redcoins-api/api/models"
	"github.com/victorcmoura/go-redcoins-api/api/responses"
	"github.com/victorcmoura/go-redcoins-api/api/utils/formaterror"
	"github.com/victorcmoura/go-redcoins-api/api/utils/timeutils"
)

func (server *Server) CreateTransaction(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	transaction := models.Transaction{}
	err = json.Unmarshal(body, &transaction)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = transaction.Prepare(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = transaction.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != transaction.OwnerID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	transactionCreated, err := transaction.SaveTransaction(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, transactionCreated.ID))
	responses.JSON(w, http.StatusCreated, transactionCreated)
}

func (server *Server) GetTransactions(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	transaction := models.Transaction{}

	transactions, err := transaction.FindAllTransactions(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, transactions)
}

func (server *Server) GetTransactionsByDay(w http.ResponseWriter, r *http.Request) {
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	var respMap map[string]string
	json.NewDecoder(r.Body).Decode(&respMap)

	layout := "2006-01-02"
	date, err := time.Parse(layout, respMap["date"])

	transactions := []models.Transaction{}
	err = server.DB.Where("created_at BETWEEN ? AND ?", timeutils.BeginningOfTheDay(date), timeutils.EndOfTheDay(date)).Find(&transactions).Error
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, transactions)
}
