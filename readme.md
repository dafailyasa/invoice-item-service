## invoice-item-service

## 🏗️ How To Run

### Install Go (Windows, Mac, Linux)
```
https://go.dev/doc/install
```

### Copy Config
```bash
 cp ./config.example.yaml ./config.yaml
```

### Build Setup(MYSQL & Elasticsearch)
```bash
 make setup
```

### Run Migration into local env
```bash
make migrate-up    #  will run up migration scripts.
```

### Run the Application Locally
```bash
make run    # Start the application locally
```
