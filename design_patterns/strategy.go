package main

import "fmt"

type PasswordProtector struct {
	user          string
	passwordName  string
	hashAlgorithm HashAlgorithm
}

// HashAlgorithm es la clave del patron Strategy
type HashAlgorithm interface {
	Hash(p *PasswordProtector)
}

// Constructor de PasswordProtector
func NewPasswordProtector(user string, passwordName string, hash HashAlgorithm) *PasswordProtector {
	return &PasswordProtector{
		user:          user,
		passwordName:  passwordName,
		hashAlgorithm: hash,
	}
}

// Esta parte es clave tambien
// Permitimos cambiar el hash del objeto (Struct)
// Estoy vuelve muy reutilizable al codigo
func (p *PasswordProtector) SetHashAlgorithm(hash HashAlgorithm) {
	p.hashAlgorithm = hash
}

func (p *PasswordProtector) Hash() {
	p.hashAlgorithm.Hash(p)
}

type SHA struct{}

func (SHA) Hash(p *PasswordProtector) {
	fmt.Printf("Hashing using SHA for %s\n", p.passwordName)
}

type MD5 struct{}

func (MD5) Hash(p *PasswordProtector) {
	fmt.Printf("Hashing using MD5 for %s\n", p.passwordName)
}

func main() {
	sha := &SHA{}
	md5 := &MD5{}

	PasswordProtector := NewPasswordProtector("leandro", "gmail password", sha)
	PasswordProtector.Hash()
	PasswordProtector.SetHashAlgorithm(md5)
	PasswordProtector.Hash()
}
