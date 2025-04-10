package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"golang.org/x/sync/errgroup"
)	

// run関数で期待通りにHTTPサーバーが起動するか、
// テストコードから意図通りに終了するかを検証する
func TestRun(t *testing.T) {
	// キャンセル可能なcontext.Contextのオブジェクトを作成する
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)

	// 別ゴルーチンでテスト対象のrun関数を実行してHTTPサーバーを起動する
	eg.Go(func() error {
		return run(ctx)
	})

	in := "message"
	// エンドポイントに対してGETリクエストを送信する
	rsp, err := http.Get("http://localhost:18080/" + in)
	if err != nil {
		t.Fatalf("failed to get: %+v", err)
	}
	defer rsp.Body.Close()

	got, err := io.ReadAll(rsp.Body)
	if err != nil {
		t.Fatalf("failed to read body: %v", err)
	}

	// HTTP serverの戻り値を検証する
	want := fmt.Sprintf("Hello, %s!", in)
	if string(got) != want {
		t.Errorf("want %q, but got %q", want, got)
	}

	// run関数に終了通知を送信する
	cancel()
	
	// *errgroup.Group.Waitメソッド経由でrun関数の戻り値を検証する
	if err := eg.Wait(); err != nil {
		t.Fatal(err)
	}
}