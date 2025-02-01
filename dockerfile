# Go rasmidan foydalanish
FROM golang:1.23-alpine

# Appni konteynerga nusxalash
WORKDIR /app
COPY . .

# Go ilovasini qurish
# RUN go mod tidy
RUN go build -o main .

# 8080 portini ochish
EXPOSE 8080
# Ilovani ishga tushurish
CMD ["./main"]
