package token_test

import (
	"time"

	"github.com/danyouknowme/awayfromus/pkg/token"
	"github.com/danyouknowme/awayfromus/pkg/util"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Payload", func() {
	lowercaseCharacters := "abcdefghijklmnopqrstuvwxyz"
	username := util.RandStringBytes(lowercaseCharacters, 7)
	duration := 24 * time.Hour
	payload, err := token.NewPayload(username, duration)

	It("Should not be nil and no error", func() {
		Expect(err).To(BeNil())
		Expect(payload).NotTo(BeNil())
		Expect(payload.ID).ShouldNot(BeZero())
		Expect(payload.Username).Should(Equal(username))
		Expect(payload.IssuedAt).ShouldNot(BeNil())
		Expect(payload.ExpiredAt).ShouldNot(BeNil())
	})

	It("Should not have error when verify payload", func() {
		err := payload.Valid()
		Expect(err).To(BeNil())
	})
})
