# 実装計画: アラート風カウントダウンタイマー

## 概要

Go言語とGioフレームワークを使用し、cgo不要でmacOSとWindowsの両方で動作するアラート風カウントダウンタイマーを実装します。MVCパターンに基づき、モデル層、コントローラー層、ビュー層を順次実装していきます。

## タスク

- [ ] 1. プロジェクト構造とコアインターフェースのセットアップ
  - [x] 1.1 Go modulesの初期化とGio依存関係の追加
    - `go mod init evangelion-timer` を実行
    - `gioui.org` パッケージを追加
    - `github.com/leanovate/gopter` をテスト用に追加
    - _Requirements: 6.1_

  - [x] 1.2 パッケージ構造の作成
    - `internal/model/`, `internal/controller/`, `internal/view/`, `internal/animation/` ディレクトリを作成
    - 各パッケージの基本ファイルを作成
    - _Requirements: 6.1_

- [ ] 2. データモデルの実装
  - [x] 2.1 タイマー状態モデルの実装 (`internal/model/timer.go`)
    - `TimerStatus` 列挙型（Idle, Running, Paused）を定義
    - `TimerModel` 構造体を実装（InitialSeconds, RemainingSeconds, Status）
    - `GetRemainingSeconds()`, `IsRunning()`, `IsPaused()`, `IsNegative()` メソッドを実装
    - _Requirements: 3.1, 4.1_

  - [x] 2.2 タイマーモデルのプロパティテスト
    - **Property 4: タイマー状態に基づく更新動作**
    - **Validates: Requirements 3.1, 3.5**

  - [x] 2.3 アプリ設定モデルの実装 (`internal/model/config.go`)
    - `AppConfig` 構造体を実装（AlwaysOnTop, WindowWidth, WindowHeight）
    - _Requirements: 6.3_

- [x] 3. チェックポイント - モデル層の検証
  - すべてのテストが通ることを確認し、質問があればユーザーに確認する

- [ ] 4. コントローラー層の実装
  - [x] 4.1 入力バリデーターの実装 (`internal/controller/validator.go`)
    - `InputValidator` 構造体を実装
    - `ValidateMinutes()`, `ValidateSeconds()`, `ValidateTotal()` メソッドを実装
    - エラーメッセージ定数を定義（ErrEmptyInput, ErrInvalidMinutes, ErrInvalidSeconds, ErrZeroTime, ErrInvalidFormat）
    - _Requirements: 1.4, 1.5, 1.6_

  - [x] 4.2 入力バリデーターのプロパティテスト
    - **Property 1: 入力バリデーションの正確性**
    - **Validates: Requirements 1.4, 1.5**

  - [x] 4.3 入力バリデーターのユニットテスト
    - 空入力、非数値入力、範囲外入力のエッジケースをテスト
    - _Requirements: 1.5, 1.6_

  - [x] 4.4 タイマーコントローラーの実装 (`internal/controller/timer.go`)
    - `TimerController` 構造体を実装
    - `Start()`, `Stop()`, `Resume()`, `Reset()`, `Toggle()` メソッドを実装
    - `time.Ticker` を使用した1秒ごとの更新処理を実装
    - ゼロ到達後のマイナス継続ロジックを実装
    - _Requirements: 2.4, 3.1, 3.2, 3.3, 3.4, 4.1_

  - [x] 4.5 タイマーコントローラーのプロパティテスト
    - **Property 3: カウントダウン開始時の初期値設定**
    - **Property 5: 停止・再開のラウンドトリップ**
    - **Property 6: ゼロ到達後のマイナス継続**
    - **Validates: Requirements 2.4, 3.2, 3.3, 4.1**

- [x] 5. チェックポイント - コントローラー層の検証
  - すべてのテストが通ることを確認し、質問があればユーザーに確認する

- [ ] 6. 時間フォーマットとアニメーション
  - [x] 6.1 時間フォーマット関数の実装 (`internal/view/format.go`)
    - `FormatTime()` 関数を実装（MM:SS:mm形式、負の値対応）
    - ミリ秒表示用のフォーマット関数を実装
    - _Requirements: 4.2_

  - [x] 6.2 時間フォーマットのプロパティテスト
    - **Property 7: 負の値の表示フォーマット**
    - **Validates: Requirements 4.2**

  - [x] 6.3 アニメーションコントローラーの実装 (`internal/animation/alert.go`)
    - `AnimationController` 構造体を実装
    - `AnimationState` 列挙型（Normal, Blink, Critical）を定義
    - `GetAnimationState()` 関数を実装（残り秒数に基づく状態判定）
    - 点滅アニメーション用のタイミング計算を実装
    - _Requirements: 5.3, 5.4, 5.5_

  - [x] 6.4 アニメーション状態のプロパティテスト
    - **Property 9: アニメーション状態の閾値遷移**
    - **Validates: Requirements 5.4, 5.5**

- [x] 7. チェックポイント - フォーマットとアニメーションの検証
  - すべてのテストが通ることを確認し、質問があればユーザーに確認する

- [ ] 8. アラートテーマの実装
  - [x] 8.1 テーマ定義の実装 (`internal/view/theme.go`)
    - 赤と黒を基調としたカラーパレットを定義
    - セグメント表示風フォント設定を定義
    - 通常状態、点滅状態、超過状態の色定義
    - _Requirements: 5.1, 5.2_

- [ ] 9. ビュー層の実装
  - [x] 9.1 メインアプリケーションウィンドウの実装 (`internal/view/app.go`)
    - Gioの `app.Window` を使用したウィンドウ作成
    - 画面状態管理（Input, Countdown）を実装
    - ウィンドウリサイズ対応を実装
    - 常に最前面表示オプションを実装
    - _Requirements: 1.1, 6.2, 6.3_

  - [x] 9.2 時間入力画面の実装 (`internal/view/input.go`)
    - 分と秒の入力フィールドを実装
    - スタートボタンを実装
    - 入力バリデーションとエラーメッセージ表示を実装
    - 有効な入力時のみスタートボタンを有効化
    - _Requirements: 1.2, 1.3, 1.5, 1.6, 2.1, 2.2_

  - [x] 9.3 スタートボタン無効化のプロパティテスト
    - **Property 2: 無効入力時のスタートボタン無効化**
    - **Validates: Requirements 2.2**

  - [x] 9.4 カウントダウン画面の実装 (`internal/view/countdown.go`)
    - セグメント表示風のタイマー表示を実装
    - ミリ秒アニメーション表示を実装（約60fps）
    - 停止/再開ボタン、リセットボタンを実装
    - マイナス表示と超過状態の視覚的強調を実装
    - 警告アニメーション（点滅、強調）を実装
    - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5, 4.2, 4.3, 5.2, 5.3, 5.4, 5.5_

- [x] 10. チェックポイント - ビュー層の検証
  - すべてのテストが通ることを確認し、質問があればユーザーに確認する

- [ ] 11. キーボードショートカットの実装
  - [x] 11.1 キーボードハンドラーの実装
    - 時間入力画面でのEnterキー処理（スタート）を実装
    - カウントダウン画面でのスペースキー処理（一時停止/再開）を実装
    - Escapeキー処理（リセット）を実装
    - _Requirements: 7.1, 7.2, 7.3_

  - [x] 11.2 スペースキートグルのプロパティテスト
    - **Property 11: スペースキーによる状態トグル**
    - **Validates: Requirements 7.2**

  - [x] 11.3 ヘルプ機能の実装
    - キーボードショートカット一覧を表示するヘルプダイアログを実装
    - _Requirements: 7.4_

- [ ] 12. エントリーポイントと統合
  - [x] 12.1 main.goの実装
    - アプリケーション初期化処理を実装
    - Gioイベントループの開始を実装
    - _Requirements: 6.1_

  - [x] 12.2 コンポーネントの統合
    - モデル、コントローラー、ビューの接続を実装
    - 画面遷移ロジックの統合を実装
    - _Requirements: 2.3, 2.4, 3.4_

- [x] 13. 最終チェックポイント - 全体の検証
  - すべてのテストが通ることを確認し、質問があればユーザーに確認する
  - `CGO_ENABLED=0 go build` でビルドが成功することを確認

## 備考

- `*` マークのタスクはオプションであり、MVPを早く完成させるためにスキップ可能
- 各タスクは特定の要件を参照しており、トレーサビリティを確保
- チェックポイントで段階的な検証を実施
- プロパティテストは普遍的な正確性プロパティを検証
- ユニットテストは特定の例とエッジケースを検証
