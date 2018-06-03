package rx

func NewObservable(subscribe func(this Observable, subscriber Subscriber))
