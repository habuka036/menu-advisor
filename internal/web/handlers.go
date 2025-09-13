package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/habuka036/menu-advisor/internal/models"
	"github.com/habuka036/menu-advisor/internal/service"
)

// Handler manages HTTP requests
type Handler struct {
	menuService       *service.MenuAdvisorService
	documentProcessor *service.DocumentProcessor
	templates         *template.Template
}

// NewHandler creates a new HTTP handler
func NewHandler(menuService *service.MenuAdvisorService) *Handler {
	// Parse templates
	tmpl, err := template.ParseGlob("web/templates/*.html")
	if err != nil {
		// If templates don't exist yet, create a minimal one
		tmpl = template.New("base")
	}

	return &Handler{
		menuService:       menuService,
		documentProcessor: service.NewDocumentProcessor(menuService),
		templates:         tmpl,
	}
}

// HomeHandler serves the main page
func (h *Handler) HomeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get current school lunches
	lunches := h.menuService.GetAllSchoolLunches()
	
	// Create a simple HTML response
	htmlResponse := `
<!DOCTYPE html>
<html lang="ja">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>学校給食メニューアドバイザー</title>
    <style>
        body { font-family: Arial, sans-serif; margin: 40px; background-color: #f5f5f5; }
        .container { max-width: 800px; margin: 0 auto; background: white; padding: 30px; border-radius: 10px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); }
        h1 { color: #2c3e50; text-align: center; margin-bottom: 30px; }
        .section { margin: 20px 0; }
        .lunch-item { background: #ecf0f1; padding: 15px; margin: 10px 0; border-radius: 5px; }
        .date { font-weight: bold; color: #34495e; }
        .main-dish { color: #e74c3c; font-size: 1.1em; margin: 5px 0; }
        .side-dishes { color: #27ae60; }
        .soup { color: #3498db; }
        .form-section { background: #f8f9fa; padding: 20px; border-radius: 5px; margin: 20px 0; }
        input, select, button { padding: 10px; margin: 5px; border: 1px solid #ddd; border-radius: 3px; }
        button { background: #3498db; color: white; border: none; cursor: pointer; }
        button:hover { background: #2980b9; }
        .suggestion { background: #d5f4e6; padding: 15px; border-radius: 5px; margin: 10px 0; border-left: 4px solid #27ae60; }
    </style>
</head>
<body>
    <div class="container">
        <h1>🍱 学校給食メニューアドバイザー</h1>
        <p>学校給食メニューを参考に、自宅での食事メニューを提案します。</p>
        
        <div class="section">
            <h2>今週の学校給食メニュー</h2>
`

	for _, lunch := range lunches {
		htmlResponse += fmt.Sprintf(`
            <div class="lunch-item">
                <div class="date">📅 %s</div>
                <div class="main-dish">🍽️ メイン: %s</div>
                <div class="side-dishes">🥬 副菜: %s</div>
                <div class="soup">🍲 汁物: %s</div>
            </div>
        `, lunch.Date.Format("2006年01月02日"), lunch.MainDish, 
		   joinSlice(lunch.SideDishes), lunch.Soup)
	}

	htmlResponse += `
        </div>
        
        <div class="form-section">
            <h2>給食メニュー文書のアップロード</h2>
            <p>PDF、画像ファイル、JSONファイルから給食メニューを読み込むことができます。</p>
            <form id="uploadForm" enctype="multipart/form-data">
                <input type="file" id="document" name="document" accept=".pdf,.jpg,.jpeg,.png,.json" required>
                <button type="submit">文書をアップロード</button>
            </form>
            <div id="uploadResult"></div>
        </div>
        
        <div class="form-section">
            <h2>家庭メニュー提案</h2>
            <p>日付と食事タイプを選択して、おすすめメニューを取得してください。</p>
            <form id="menuForm">
                <input type="date" id="date" name="date" required>
                <select id="mealType" name="mealType" required>
                    <option value="">食事タイプを選択</option>
                    <option value="breakfast">朝食</option>
                    <option value="dinner">夕食</option>
                </select>
                <button type="submit">メニュー提案を取得</button>
            </form>
            <div id="result"></div>
        </div>
    </div>

    <script>
        // Handle document upload form
        document.getElementById('uploadForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const fileInput = document.getElementById('document');
            const file = fileInput.files[0];
            
            if (!file) {
                alert('ファイルを選択してください');
                return;
            }
            
            const formData = new FormData();
            formData.append('document', file);
            
            try {
                document.getElementById('uploadResult').innerHTML = '<p>アップロード中...</p>';
                
                const response = await fetch('/api/upload', {
                    method: 'POST',
                    body: formData
                });
                
                const data = await response.json();
                
                if (data.success) {
                    document.getElementById('uploadResult').innerHTML = ` + "`" + `
                        <div class="suggestion">
                            <h3>✅ アップロード成功</h3>
                            <p>${data.message}</p>
                            <p><small>文書ID: ${data.result.id}</small></p>
                            <p><small>処理状況: ${data.result.status}</small></p>
                        </div>
                    ` + "`" + `;
                    
                    // Reload the page to show new menu data
                    setTimeout(() => {
                        window.location.reload();
                    }, 2000);
                } else {
                    document.getElementById('uploadResult').innerHTML = ` + "`" + `
                        <div style="color: red; background: #ffebee; padding: 10px; border-radius: 5px;">
                            <h3>❌ アップロードエラー</h3>
                            <p>${data.error}</p>
                        </div>
                    ` + "`" + `;
                }
            } catch (error) {
                document.getElementById('uploadResult').innerHTML = ` + "`" + `
                    <div style="color: red; background: #ffebee; padding: 10px; border-radius: 5px;">
                        <h3>❌ エラー</h3>
                        <p>アップロードに失敗しました: ${error.message}</p>
                    </div>
                ` + "`" + `;
            }
        });

        // Handle menu suggestion form
        document.getElementById('menuForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            const date = document.getElementById('date').value;
            const mealType = document.getElementById('mealType').value;
            
            if (!date || !mealType) {
                alert('日付と食事タイプを選択してください');
                return;
            }
            
            try {
                const response = await fetch(` + "`" + `/api/suggest?date=${date}&meal_type=${mealType}` + "`" + `);
                const data = await response.json();
                
                if (response.ok) {
                    document.getElementById('result').innerHTML = ` + "`" + `
                        <div class="suggestion">
                            <h3>🌟 ${data.meal_type === 'breakfast' ? '朝食' : '夕食'}の提案</h3>
                            <p><strong>メイン:</strong> ${data.main_dish}</p>
                            <p><strong>副菜:</strong> ${data.side_dishes.join(', ')}</p>
                            ${data.soup ? ` + "`<p><strong>汁物:</strong> ${data.soup}</p>`" + ` : ''}
                            <p><strong>理由:</strong> ${data.reason}</p>
                            <p><small>参考給食: ${data.school_lunch_ref}</small></p>
                        </div>
                    ` + "`" + `;
                } else {
                    document.getElementById('result').innerHTML = ` + "`" + `<div style="color: red;">エラー: ${data.error}</div>` + "`" + `;
                }
            } catch (error) {
                document.getElementById('result').innerHTML = ` + "`" + `<div style="color: red;">エラーが発生しました: ${error.message}</div>` + "`" + `;
            }
        });
    </script>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(htmlResponse))
}

// SuggestHandler provides menu suggestions via API
func (h *Handler) SuggestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	dateStr := r.URL.Query().Get("date")
	mealType := r.URL.Query().Get("meal_type")

	if dateStr == "" || mealType == "" {
		http.Error(w, "Missing required parameters: date and meal_type", http.StatusBadRequest)
		return
	}

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	suggestion, err := h.menuService.GenerateHomeMenuSuggestion(date, mealType)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(suggestion)
}

// SchoolLunchHandler returns school lunch data
func (h *Handler) SchoolLunchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	lunches := h.menuService.GetAllSchoolLunches()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(lunches)
}

// Helper function to join string slices
func joinSlice(slice []string) string {
	if len(slice) == 0 {
		return "なし"
	}
	result := ""
	for i, item := range slice {
		if i > 0 {
			result += ", "
		}
		result += item
	}
	return result
}

// UploadHandler handles document uploads
func (h *Handler) UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(32 << 20) // 32MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Get uploaded file
	file, header, err := r.FormFile("document")
	if err != nil {
		http.Error(w, "Failed to get uploaded file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create processing request
	req := &models.DocumentProcessingRequest{
		File:   file,
		Header: header,
	}

	// Process document
	result, err := h.documentProcessor.ProcessDocument(req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"result":  result,
		})
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Document processed successfully",
		"result":  result,
	})
}