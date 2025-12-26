# Dmensions
Simple yet powerful playground for embeddings.
```
       /$$                                             /$$                              
      | $$                                            |__/                              
  /$$$$$$$ /$$$$$$/$$$$   /$$$$$$  /$$$$$$$   /$$$$$$$ /$$  /$$$$$$  /$$$$$$$   /$$$$$$$
 /$$__  $$| $$_  $$_  $$ /$$__  $$| $$__  $$ /$$_____/| $$ /$$__  $$| $$__  $$ /$$_____/
| $$  | $$| $$ \ $$ \ $$| $$$$$$$$| $$  \ $$|  $$$$$$ | $$| $$  \ $$| $$  \ $$|  $$$$$$ 
| $$  | $$| $$ | $$ | $$| $$_____/| $$  | $$ \____  $$| $$| $$  | $$| $$  | $$ \____  $$
|  $$$$$$$| $$ | $$ | $$|  $$$$$$$| $$  | $$ /$$$$$$$/| $$|  $$$$$$/| $$  | $$ /$$$$$$$/
 \_______/|__/ |__/ |__/ \_______/|__/  |__/|_______/ |__/ \______/ |__/  |__/|_______/ 
```



## Installation:
Program heavilly relies on [EmbeddingGemma](https://ollama.com/library/embeddinggemma) served by [Ollama API](https://github.com/ollama/ollama), please ensure they are up and running.

```bash
# Install Ollama from:
# https://ollama.com/download

# Get EmbeddingGemma
ollama pull embeddinggemma

# Downlaod the repository and enter it:
git clone https://github.com/umarbektokyo/dmensions
cd dmensions

# Run the app
go run ./cmd/dmensions/main.go
```

