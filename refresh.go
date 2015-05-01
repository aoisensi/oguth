package oguth

func DefaultRefreshTokenGenerator() (code string) {
	return SimpleRandomTokenGenerator(48)
}
