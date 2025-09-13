# 🍱 学校給食メニューアドバイザー (School Lunch Menu Advisor)

学校給食メニューを参考にした自宅での食事メニューの提案サービス

## 概要

このサービスは学校で提供される給食メニューを分析し、栄養バランスを考慮して家庭での朝食・夕食メニューを提案します。子どもが学校で摂取する栄養を補完し、1日全体での栄養バランスを最適化することを目的としています。

## 機能

- 📊 学校給食メニューの表示
- 🍳 給食内容に基づく朝食・夕食メニューの提案
- 🥗 栄養バランスを考慮した補完的なメニュー推奨
- 🌐 ウェブインターフェースでの簡単操作
- 📱 レスポンシブデザイン対応

## 技術スタック

- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, JavaScript
- **Data**: JSON形式での給食データ管理

## 使用方法

### 1. サービスの起動

```bash
# リポジトリをクローン
git clone https://github.com/habuka036/menu-advisor.git
cd menu-advisor

# 依存関係の取得
go mod tidy

# サービスの起動
go run cmd/main.go
```

### 2. ウェブインターフェースへのアクセス

ブラウザで `http://localhost:8080` にアクセス

### 3. API使用例

```bash
# 全ての給食メニューを取得
curl http://localhost:8080/api/school-lunches

# 特定日の朝食メニュー提案を取得
curl "http://localhost:8080/api/suggest?date=2025-01-13&meal_type=breakfast"

# 特定日の夕食メニュー提案を取得
curl "http://localhost:8080/api/suggest?date=2025-01-13&meal_type=dinner"
```

## APIエンドポイント

- `GET /` - メインのウェブインターフェース
- `GET /api/school-lunches` - 学校給食データの取得
- `GET /api/suggest?date=YYYY-MM-DD&meal_type=breakfast|dinner` - メニュー提案

## プロジェクト構造

```
menu-advisor/
├── cmd/
│   └── main.go              # メインアプリケーション
├── internal/
│   ├── models/
│   │   └── menu.go          # データモデル
│   ├── service/
│   │   ├── menu_advisor.go  # メニュー提案ロジック
│   │   └── menu_advisor_test.go # テスト
│   └── web/
│       └── handlers.go      # HTTPハンドラー
├── data/
│   └── school_lunch_sample.json # サンプル給食データ
├── go.mod
└── README.md
```

## データ形式

### 学校給食メニュー

```json
{
  "date": "2025-01-13T00:00:00Z",
  "main_dish": "鶏肉の照り焼き",
  "side_dishes": ["野菜炒め", "白米"],
  "soup": "味噌汁（わかめ）",
  "nutrition": {
    "calories": 650,
    "protein_g": 28.5,
    "carbs_g": 85.2,
    "fat_g": 18.3,
    "fiber_g": 4.2,
    "sodium_mg": 850,
    "vegetables_servings": 2
  }
}
```

## 開発

### テストの実行

```bash
go test ./...
```

### 依存関係の更新

```bash
go mod tidy
```

## ライセンス

MIT License - 詳細は [LICENSE](LICENSE) ファイルをご確認ください。

## 貢献

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/AmazingFeature`)
3. 変更をコミット (`git commit -m 'Add some AmazingFeature'`)
4. ブランチにプッシュ (`git push origin feature/AmazingFeature`)
5. プルリクエストを作成

## 作者

Osamu Habuka (@habuka036)