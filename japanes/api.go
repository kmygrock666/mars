package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//APIv1 APIv1
type APIv1 struct{}

var japan JapanAccent

// Ping Ping
func (api *APIv1) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("pong"))
}

//LearnAccent LearnAccent
func (api *APIv1) LearnAccent(w http.ResponseWriter, r *http.Request) {
	resp := NewResponse(w, r)
	defer resp.HTTPResponse()

	vars := mux.Vars(r)
	userID, err := strconv.ParseInt(vars["user_id"], 10, 64)
	if err != nil {
		resp.SetError(RequsetParamsError("user_id", vars["user_id"]))
		return
	}

	//TODO 先查詢是否曾經已開新局
	game, err := mysql.getGamesByUserId(userID)
	if err != nil {
		resp.SetError(SystemError(err))
		return
	}

	if game == nil {
		var newGame *gameTable
		newGame.UserID = userID
		_, err := mysql.InsertGames(newGame)
		if err != nil {
			resp.SetError(SystemError(err))
			return
		}
		//TODO 新增拼音
		// if rowsAffected >= 1) {
		// 	resp.AddData("isSuccess", true)
		// }
	}

	newRoundshuffle()
	res := getAccentWord()

	resp.addData("games", res)

}

// CheckAnswer CheckAnswer
func (api *APIv1) CheckAnswer(w http.ResponseWriter, r *http.Request) {
	resp := NewResponse(w, r)
	defer resp.HTTPResponse()

	var pushAns *checkOrder
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pushAns)
	if err != nil {
		resp.SetError(SystemError(err))
		return
	}

	orders, err := mysql.getOrders(pushAns.GameID)
	if err != nil {
		resp.SetError(SystemError(err))
		return
	}

	if len(orders) <= 0 {
		resp.SetError(NoOrdersError(pushAns.GameID, pushAns.WordIdx))
		return
	}
	//TODO 檢查回答是否正確
	var aim *JapanAccent
	var chkres bool
	next := true
	updateOrders := make([]*orderTable, 0)
	for _, oldOrder := range orders {
		aim.getCurrentAccent(oldOrder.Result)
		chkres = aim.CheckPinyin(pushAns.Answer)
		order := oldOrder
		if chkres {
			order.State = 2
		} else {
			order.State = 1
		}

		updateOrders = append(updateOrders, order)
	}

	_, err = mysql.InsertOrders(updateOrders)
	if err != nil {
		resp.SetError(SystemError(err))
		return
	}

	if len(tmpIndex) == 0 {
		var game *gameTable
		game.State = 2
		//TODO 取得用戶id
		mysql.updateStatus(game)
		next = false
	}

	resp.addData("result", chkres)
	resp.addData("next", next)
}
