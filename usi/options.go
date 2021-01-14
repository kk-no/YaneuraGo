package usi

type Option func(u *usi)

func SetEnginePath(path string) Option {
	return func(u *usi) {
		u.path = path
	}
}

func WithDebug(isDebug bool) Option {
	return func(u *usi) {
		u.isDebug = isDebug
	}
}
