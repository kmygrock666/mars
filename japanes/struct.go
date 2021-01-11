package main

import "time"

type gameTable struct {
	GameID    int64     `db:"game_id" json:",string"`
	State     int8      `db:"state"`
	UserID    int64     `db:"user_id"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt int64     `db:"created_at"`
}

type wordData struct {
	ID        int64     `db:"id" json:",string"`
	Foreign   string    `db:"foreign"`
	Native    string    `db:"native" json:",string"`
	Pinyin    string    `db:"pinyin"`
	ExtraInfo string    `db:"extra_info"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt int64     `db:"created_at"`
}

type orderTable struct {
	OrderID   int64     `db:"order_id" json:",string"`
	Result    int64     `db:"result"`
	GameID    int64     `db:"game_id" json:",string"`
	State     int8      `db:"state"`
	UpdatedAt time.Time `db:"updated_at"`
	CreatedAt int64     `db:"created_at"`
}

type checkOrder struct {
	GameID  int64 `json:",string"`
	Answer  string
	WordIdx int8 `json:",string"`
}
