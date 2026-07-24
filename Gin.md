*gin.Context (c)
│
├── Request (*http.Request)
│   │
│   ├── Method
│   ├── URL
│   ├── Proto
│   ├── Host
│   ├── RemoteAddr
│   ├── Header
│   ├── Body
│   ├── ContentLength
│   ├── TransferEncoding
│   ├── Form
│   ├── PostForm
│   ├── MultipartForm
│   ├── Trailer
│   ├── TLS
│   └── Context()
│
├── Writer (gin.ResponseWriter)
│   │
│   ├── Header()
│   ├── Write()
│   ├── WriteHeader()
│   ├── Status()
│   ├── Size()
│   └── Written()
│
├── Params
│
├── Keys
│
├── Errors
│
├── Accepted
│
├── FullPath()
│
├── ClientIP()
│
├── GetHeader()
│
├── Query()
│
├── Param()
│
├── PostForm()
│
├── BindJSON()
│
├── ShouldBindJSON()
│
├── JSON()
│
├── String()
│
├── XML()
│
├── File()
│
├── Redirect()
│
├── Abort()
│
├── Next()
│
└── Set()/Get()