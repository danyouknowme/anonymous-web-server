package util_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/danyouknowme/awayfromus/pkg/util"
)

var _ = Describe("Config", func() {
	util.LoadConfig()

	It("Should load config port from environment variables file", func() {
		Expect(util.AppConfig.Port).ShouldNot(BeNil())
		Expect(util.AppConfig.Port).Should(Equal(os.Getenv("PORT")))
	})

	It("Should load config mongo uri from environment variables file", func() {
		Expect(util.AppConfig.MongoUri).ShouldNot(BeNil())
		Expect(util.AppConfig.MongoUri).Should(Equal(os.Getenv("MONGO_URI")))
	})

	It("Should load config secret key from environment variables file", func() {
		Expect(util.AppConfig.SecretKey).ShouldNot(BeNil())
		Expect(util.AppConfig.SecretKey).Should(Equal(os.Getenv("SECRET_KEY")))
	})
})
