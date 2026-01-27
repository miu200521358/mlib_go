//go:build windows
// +build windows

// 指示: miu200521358
package controller

import "github.com/miu200521358/walk/pkg/walk"

// ProgressBar は進捗バーを表す。
type ProgressBar struct {
	*walk.ProgressBar
}

// NewProgressBar はProgressBarを生成する。
func NewProgressBar(parent walk.Container) (*ProgressBar, error) {
	pb := new(ProgressBar)

	var err error
	pb.ProgressBar, err = walk.NewProgressBar(parent)
	if err != nil {
		return nil, NewProgressBarInitFailed("ProgressBarの初期化に失敗しました", err)
	}

	return pb, nil
}

// SetMax は最大値を設定する。
func (pb *ProgressBar) SetMax(max int) {
	pb.SetRange(0, max)
}

// SetValue は現在値を設定する。
func (pb *ProgressBar) SetValue(v int) {
	pb.ProgressBar.SetValue(v)
}

// Increment は現在値を1増やす。
func (pb *ProgressBar) Increment() {
	pb.SetValue(pb.Value() + 1)
}
