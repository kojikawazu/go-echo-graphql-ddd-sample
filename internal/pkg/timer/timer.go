package pkg_timer

import (
	"time"
)

// タイマーの構造体
type TimerPkg struct {
	StartTime time.Time
	EndTime   time.Time
}

// タイマーのインスタンス化
func NewTimerPkg() *TimerPkg {
	return &TimerPkg{}
}

// タイマーの開始
func (t *TimerPkg) Start() {
	t.StartTime = time.Now()
}

// タイマーの終了
func (t *TimerPkg) End() {
	t.EndTime = time.Now()
}

// タイム時間を取得
func (t *TimerPkg) GetDuration() time.Duration {
	t.End()
	return t.EndTime.Sub(t.StartTime)
}
