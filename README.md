# LaTeX 目錄中文字數統計工具

這個工具可以遍歷指定目錄下的所有 LaTeX 文件，計算每個文件的中文字數，並以表格形式顯示結果。

## 安裝

確保你已經安裝了 Go（版本 1.22 或更高）。

### 方法 1：直接從 GitHub 安裝

```bash
go install github.com/leon123858/latex-directory-counter/cmd/latex-directory-counter@latest
```

### 方法 2：克隆倉庫後安裝

```bash
git clone https://github.com/your_github_username/latex-directory-counter.git
cd latex-directory-counter
go install ./cmd/latex-directory-counter
```

安裝後，確保 `$GOPATH/bin` 在你的 PATH 中。你可以通過將以下行添加到你的 `.bashrc` 或 `.zshrc` 文件來實現：

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

## 使用方法

```bash
latex-directory-counter <LaTeX目錄路徑>
```

例如：

```bash
latex-directory-counter /path/to/your/latex/project
```

這將顯示一個包含每個 .tex 文件中文字數的表格，以及總字數。

## 功能

- 遍歷指定目錄下的所有 .tex 文件
- 計算每個文件中的中文字數
- 忽略 LaTeX 註釋、命令和環境
- 按字數降序排列文件
- 以表格形式顯示結果，包括總字數

## 貢獻

歡迎提交 issues 和 pull requests 來改進這個工具。

## 許可證

本項目採用 MIT 許可證。