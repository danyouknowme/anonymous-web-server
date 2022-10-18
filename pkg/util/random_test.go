package util_test

import (
	"regexp"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/danyouknowme/awayfromus/pkg/util"
)

var _ = Describe("Random", func() {
	n := 5
	lowercaseCharacters := "abcdefghijklmnopqrstuvwxyz"
	randomString := util.RandStringBytes(lowercaseCharacters, n)

	It("Should return random string with length n", func() {
		Expect(randomString).Should(HaveLen(n))
	})

	username := "adminuser1"
	remainingLicenseLength := 7
	license := util.GenerateLicense(username)

	It("Should return string with length of username length plus remaining license length", func() {
		Expect(license).Should(HaveLen(len(username) + remainingLicenseLength))
	})

	It("Should contain username in generated license", func() {
		match1, err := regexp.MatchString(username, license)
		Expect(err).To(BeNil())
		Expect(match1).Should(BeTrue())
	})

	secretCodeSliceLength := 9
	secretCodeLength := 5
	secretCode := util.GenerateSecretCode()

	It("Should return slice of string with length of secret code length", func() {
		Expect(secretCode).Should(HaveLen(secretCodeSliceLength))
	})

	It("Should ", func() {
		Expect(secretCode[0]).Should(HaveLen(secretCodeLength))
	})
})
