package internal

import (
	"context"
	"time"
)

func GetTimeOutContext(duration ...time.Duration) (context.Context, context.CancelFunc){
	var m_duration time.Duration
	if len(duration)<=0{
		m_duration= time.Second*5
	}else{
		m_duration= duration[0]
	}

	return context.WithTimeout(context.Background(), m_duration)
}