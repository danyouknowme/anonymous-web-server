package token_test

import (
	"time"

	"github.com/danyouknowme/awayfromus/pkg/token"
	"github.com/danyouknowme/awayfromus/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("JwtMaker", func() {
	lowercaseCharacters := "abcdefghijklmnopqrstuvwxyz"
	username := util.RandStringBytes(lowercaseCharacters, 7)
	duration := time.Minute

	It("Should return the new token and not have error", func() {
		recievedToken, err := token.CreateToken(username, duration)
		Expect(err).To(BeNil())
		Expect(recievedToken).ToNot(BeEmpty())

		payload, err := token.VerifyToken(recievedToken)
		Expect(err).To(BeNil())
		Expect(payload).ToNot(BeNil())
		Expect(payload.ID).ShouldNot(BeZero())
		Expect(payload.Username).Should(Equal(username))
		Expect(payload.IssuedAt).ShouldNot(BeNil())
		Expect(payload.ExpiredAt).ShouldNot(BeNil())
	})

	It("Should throw error when token is expired", func() {
		recievedToken, err := token.CreateToken(username, -time.Minute)
		Expect(err).To(BeNil())
		Expect(recievedToken).ToNot(BeEmpty())

		payload, err := token.VerifyToken(recievedToken)
		Expect(err).Should(Equal(token.ErrExpiredToken))
		Expect(payload).Should(BeNil())
	})
})
