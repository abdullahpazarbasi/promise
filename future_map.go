package promise

type FutureMap[T any] interface {
	Commit() ProgressMap[T]
	Await() *map[interface{}]Output[T]
	Race() (key interface{}, pay T, err error)
	commit() ProgressMap[T]
}

type futureMap[T any] map[interface{}]Future[T]

func (fm *futureMap[T]) Commit() ProgressMap[T] {
	return fm.commit()
}

func (fm *futureMap[T]) Await() *map[interface{}]Output[T] {
	return fm.commit().await()
}

func (fm *futureMap[T]) Race() (key interface{}, pay T, err error) {
	//TODO implement me
	panic("implement me")
}

func (fm *futureMap[T]) commit() ProgressMap[T] {
	pm := make(progressMap[T])
	for k, p := range *fm {
		pm[k] = p.commit()
	}

	return &pm
}
