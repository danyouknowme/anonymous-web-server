package util

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"golang.org/x/crypto/bcrypt"
)

var _ = Describe("Password", func() {
	password := "$Password1234"
	hashedPassword1, err := HashPassword(password)

	It("Should not be nil and no error", func() {
		Expect(err).To(BeNil())
		Expect(hashedPassword1).NotTo(BeNil())
	})

	It("Should not have an error when checking password", func() {
		err = CheckPassword(password, hashedPassword1)
		Expect(err).To(BeNil())
	})

	It("Should throw error mismatched when checking wrong password", func() {
		wrongPassword := "$Password12345"
		err = CheckPassword(wrongPassword, hashedPassword1)
		Expect(err).ShouldNot(BeNil())
		Expect(err).Should(Equal(bcrypt.ErrMismatchedHashAndPassword))
	})

	It("Should not equal when compare hash password with the same value", func() {
		hashedPassword2, err := HashPassword(password)
		Expect(err).To(BeNil())
		Expect(hashedPassword2).NotTo(BeNil())
		Expect(hashedPassword2).NotTo(Equal(hashedPassword1))
	})
})
