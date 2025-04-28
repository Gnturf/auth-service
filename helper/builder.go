package helper

import "fmt"

func PostgresDSNBuilder(username, password, ip, port, dbName string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, ip, port, dbName)
}