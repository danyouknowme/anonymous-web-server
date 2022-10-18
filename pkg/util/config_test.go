package util_test

import (
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/danyouknowme/awayfromus/pkg/util"
)

var _ = Describe("Config", func() {
	util.LoadConfig()

	It("Should load config from environment variables file", func() {
		Expect(util.AppConfig.Port).Should(Equal(os.Getenv("PORT")))
		Expect(util.AppConfig.MongoUri).Should(Equal(os.Getenv("MONGO_URI")))
		Expect(util.AppConfig.SecretKey).Should(Equal(os.Getenv("SECRET_KEY")))
	})
})
