package internal

import (
	"context"
	"time"
)

func GetContextWithTimeOut(timeOut time.Duration) (context.Context, context.CancelFunc){
	return context.WithTimeout(context.Background(), timeOut)
}