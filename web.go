package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// クエリパラメータを解析する
	query := r.URL.Query()
	name := query.Get("name") // "名前"から"name"へ修正

	// レスポンス用のマップを作成
	response := map[string]string{
		"message": "Hello " + name, // “message”： "Hello " + name、から修正
	}

	// Content-Typeヘッダーをapplication/jsonに設定
	w.Header().Set("Content-Type", "application/json")

	// マップをJSONにエンコードしてレスポンスとして送信
	json.NewEncoder(w).Encode(response)
}

// api/categories
func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories := []string{"Sprint", "Mile", "Medium", "Long"} // 競馬の距離区分(ウマ娘より)

	response := map[string][]string{
		"categories": categories,
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(response)
}

func calculatorHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	operator := query.Get("o")
	x, err1 := strconv.ParseFloat(query.Get("x"), 64)
	y, err2 := strconv.ParseFloat(query.Get("y"), 64)

	if err1 != nil || err2 != nil {
		// 400 Bad Request を返す
		http.Error(w, "Invalid value for x or y. Must be a number.", http.StatusBadRequest)
		// return で処理を終了
		return
	}

	// ゼロ除算チェック
	if operator == "/" && y == 0 {
		http.Error(w, "Cannot divide by zero.", http.StatusBadRequest)
		return
	}

	response := map[string]float64{
		"Answer": calculator(operator, x, y),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func calculator(operator string, x, y float64) float64 {
	var result float64
	switch operator {
	case "+":
		result = x + y
	case "-":
		result = x - y
	case "*":
		result = x * y
	case "/":
		result = x / y
	default:
		result = 0.0
	}

	return result
}

func main() {
	fmt.Println("Starting the server!")

	// ルートとハンドラ関数を定義
	http.HandleFunc("/api/hello", helloHandler)

	// http://localhost:8000/api/categories -> {"categories":["Sprint","Mile","Medium","Long"]}
	http.HandleFunc("/api/categories", categoriesHandler)

	// e.g., http://localhost:8000/api/calculator?o=%2B&x=5&y=2 (operatorに+を与えると空白として処理されるため、%2Bにした) -> {"Answer":7}
	// 引き算と掛け算の場合はURLエンコードは不要だけど、割り算はパス階層の区切りとして認識されるから、%2Fを使う
	http.HandleFunc("/api/calculator", calculatorHandler)

	// 8000番ポートでサーバを開始
	http.ListenAndServe(":8000", nil)
}
