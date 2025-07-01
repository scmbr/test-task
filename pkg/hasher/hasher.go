package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	cost int
}

func NewHasher(cost int) *Hasher {
	return &Hasher{cost: cost}
}

func (h *Hasher) Hash(token string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(token), h.cost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}
func (h *Hasher) Verify(token, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}
