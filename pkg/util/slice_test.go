package util_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/danyouknowme/awayfromus/pkg/util"
)

var _ = Describe("Slice", func() {
	pets := []string{"dog", "cat", "rabbit"}

	It("Should return true when slice contains specific string", func() {
		Expect(util.Contains(pets, "dog")).Should(BeTrue())
	})

	It("Should return false when slice not contains specific string", func() {
		Expect(util.Contains(pets, "tiger")).Should(BeFalse())
	})

	It("Should return the correct new slice when remove element from slice", func() {
		removedDogFromSlice := []string{"cat", "rabbit"}
		Expect(util.Remove(pets, "dog")).Should(Equal(removedDogFromSlice))
		remainingSlice := []string{"rabbit"}
		Expect(util.Remove(removedDogFromSlice, "cat")).Should(Equal(remainingSlice))
	})

	It("Should return the same slice when element cannot remove from slice", func() {
		Expect(util.Remove(pets, "shrimp")).Should(Equal(pets))
	})
})
