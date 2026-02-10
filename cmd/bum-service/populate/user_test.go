package populate_test

import (
	"golang.org/x/crypto/bcrypt"
	"testing"
)

var password = []byte("myStrongPassword123!")

func benchmarkBcrypt(cost int, b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := bcrypt.GenerateFromPassword(password, cost)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkBcryptCost2(b *testing.B)  { benchmarkBcrypt(2, b) }
func BenchmarkBcryptCost4(b *testing.B)  { benchmarkBcrypt(4, b) }
func BenchmarkBcryptCost8(b *testing.B)  { benchmarkBcrypt(8, b) }
func BenchmarkBcryptCost10(b *testing.B) { benchmarkBcrypt(10, b) }
func BenchmarkBcryptCost12(b *testing.B) { benchmarkBcrypt(12, b) }
func BenchmarkBcryptCost14(b *testing.B) { benchmarkBcrypt(14, b) }
