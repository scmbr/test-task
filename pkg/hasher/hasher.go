package hasher

import "golang.org/x/crypto/bcrypt"

type Hasher struct {
	cost int
}

func NewHasher(cost int) *Hasher {
	return &Hasher{cost: cost}
}

func (h *Hasher) Hash(data []byte) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword(data, bcrypt.DefaultCost)
	return string(hashed), err
}
func (h *Hasher) Verify(token []byte, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), token)
	return err == nil
}
