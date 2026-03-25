# 要件ドキュメント

## はじめに

本ドキュメントは、ブラウザで動作するHTML+JavaScript製のアラート風カウントダウンタイマーの要件を定義します。既存のGo版デスクトップアプリケーションをWeb版として再実装し、単一のHTMLファイルで完結する構成とします。プレゼンテーション用途を想定し、ゼロ到達後もマイナス方向へカウントダウンを継続する機能を備えます。

## 用語集

- **Timer_Application**: カウントダウンタイマーのWebアプリケーション（単一HTMLファイル）
- **Timer_Display**: カウントダウン時間を表示するUI要素
- **Timer_Controller**: タイマーの開始・停止・リセットを制御するJavaScriptモジュール
- **Countdown_Value**: 現在のカウントダウン値（正の値、ゼロ、または負の値）
- **Alert_Animation**: 警告感を演出するCSSアニメーション効果
- **Time_Input**: ユーザーがカウントダウン開始時間を入力するHTML入力要素

## 要件

### 要件 1: 起動時の時間入力画面

**ユーザーストーリー:** プレゼンターとして、ページ読み込み時にプレゼンテーションの時間を入力したい。これにより、すぐにカウントダウンの準備ができる。

#### 受け入れ基準

1. WHEN Timer_Application がブラウザで読み込まれた時、THE Timer_Application SHALL 最初に時間入力画面を表示する
2. THE Time_Input SHALL 分と秒を個別に入力できるフィールドを提供する
3. THE Time_Input SHALL 分フィールドと秒フィールドを明確に区別して表示する
4. THE Timer_Application SHALL 最大99分59秒までのカウントダウン時間をサポートする
5. WHEN ユーザーが Time_Input に値を入力した時、THE Timer_Application SHALL 入力値を検証し、有効な時間形式であることを確認する
6. IF 無効な時間形式が入力された場合、THEN THE Timer_Application SHALL エラーメッセージを表示し、再入力を促す

### 要件 2: スタートボタンによるカウントダウン開始

**ユーザーストーリー:** プレゼンターとして、スタートボタンを押してカウントダウンを開始したい。これにより、準備が整ったタイミングでプレゼンテーションを始められる。

#### 受け入れ基準

1. THE Timer_Application SHALL 時間入力画面にスタートボタンを表示する
2. WHILE Time_Input に有効な時間が入力されていない間、THE Timer_Application SHALL スタートボタンを無効化する
3. WHEN スタートボタンがクリックされた時、THE Timer_Application SHALL 時間入力画面からカウントダウン画面に遷移する
4. WHEN スタートボタンがクリックされた時、THE Timer_Controller SHALL 入力された時間でカウントダウンを開始する

### 要件 3: カウントダウン実行

**ユーザーストーリー:** プレゼンターとして、カウントダウンを停止・リセットしたい。これにより、プレゼンテーションの進行を柔軟に制御できる。

#### 受け入れ基準

1. WHILE カウントダウンが実行中の間、THE Timer_Display SHALL 1秒ごとに Countdown_Value を更新して表示する
2. WHEN 停止ボタンがクリックされた時、THE Timer_Controller SHALL カウントダウンを一時停止する
3. WHEN 再開ボタンがクリックされた時、THE Timer_Controller SHALL カウントダウンを再開する
4. WHEN リセットボタンがクリックされた時、THE Timer_Controller SHALL 時間入力画面に戻る
5. WHILE カウントダウンが一時停止中の間、THE Timer_Display SHALL 現在の Countdown_Value を維持して表示する

### 要件 4: マイナスカウントダウン継続

**ユーザーストーリー:** プレゼンターとして、タイマーがゼロになった後もマイナス方向にカウントダウンを継続させたい。これにより、予定時間を超過した時間を把握できる。

#### 受け入れ基準

1. WHEN Countdown_Value がゼロに到達した時、THE Timer_Controller SHALL カウントダウンを停止せずマイナス方向への計測を継続する
2. WHILE Countdown_Value が負の値の間、THE Timer_Display SHALL マイナス記号付きで経過時間を表示する（例：-00:15）
3. WHILE Countdown_Value が負の値の間、THE Timer_Display SHALL 超過状態であることを視覚的に強調表示する

### 要件 5: アラート風UIデザイン

**ユーザーストーリー:** プレゼンターとして、アラート風の警告感のあるデザインでタイマーを表示したい。これにより、視覚的にインパクトのあるプレゼンテーションができる。

#### 受け入れ基準

1. THE Timer_Application SHALL 赤と黒を基調としたカラースキームを使用する
2. THE Timer_Display SHALL セグメント表示風のフォントでカウントダウン値を表示する
3. WHILE カウントダウンが実行中の間、THE Timer_Display SHALL 警告感を演出する Alert_Animation を表示する
4. WHEN Countdown_Value が残り10秒以下になった時、THE Timer_Display SHALL 点滅アニメーションを開始する
5. WHILE Countdown_Value が負の値の間、THE Alert_Animation SHALL より強調された警告表示に切り替わる

### 要件 6: ブラウザWebアプリケーション

**ユーザーストーリー:** プレゼンターとして、ブラウザで動作するWebアプリケーションとしてタイマーを使用したい。これにより、OSを問わず利用でき、インストール不要で手軽に使える。

#### 受け入れ基準

1. THE Timer_Application SHALL 単一のHTMLファイルで完結する（CSS、JavaScriptを埋め込み）
2. THE Timer_Application SHALL モダンブラウザ（Chrome、Firefox、Safari、Edge）で動作する
3. THE Timer_Application SHALL ウィンドウのリサイズに対応し、Timer_Display を適切にスケーリングする
4. THE Timer_Application SHALL 外部依存なしで動作する（CDNやライブラリ不要）

### 要件 7: キーボードショートカット

**ユーザーストーリー:** プレゼンターとして、キーボードショートカットでタイマーを操作したい。これにより、プレゼンテーション中に素早く操作できる。

#### 受け入れ基準

1. WHEN 時間入力画面でEnterキーが押された時、THE Timer_Application SHALL スタートボタンと同じ動作を実行する
2. WHEN カウントダウン画面でスペースキーが押された時、THE Timer_Controller SHALL カウントダウンの一時停止または再開を切り替える
3. WHEN Escapeキーが押された時、THE Timer_Controller SHALL 時間入力画面に戻る
4. THE Timer_Application SHALL キーボードショートカットの一覧を表示するヘルプ機能を提供する
