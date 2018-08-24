package gas_test

import (
	"bytes"

	"github.com/GoASTScanner/gas"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Configuration", func() {
	var configuration gas.Config
	BeforeEach(func() {
		configuration = gas.NewConfig()
	})

	Context("when loading from disk", func() {

		It("should be possible to load configuration from a file", func() {
			json := `{"G101": {}}`
			buffer := bytes.NewBufferString(json)
			nread, err := configuration.ReadFrom(buffer)
			Expect(nread).Should(Equal(int64(len(json))))
			Expect(err).ShouldNot(HaveOccurred())
		})

		It("should return an error if configuration file is invalid", func() {
			var err error
			invalidBuffer := bytes.NewBuffer([]byte{0xc0, 0xff, 0xee})
			_, err = configuration.ReadFrom(invalidBuffer)
			Expect(err).Should(HaveOccurred())

			emptyBuffer := bytes.NewBuffer([]byte{})
			_, err = configuration.ReadFrom(emptyBuffer)
			Expect(err).Should(HaveOccurred())
		})

	})

	Context("when saving to disk", func() {
		It("should be possible to save an empty configuration to file", func() {
			expected := `{"global":{}}`
			buffer := bytes.NewBuffer([]byte{})
			nbytes, err := configuration.WriteTo(buffer)
			Expect(int(nbytes)).Should(Equal(len(expected)))
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).Should(Equal(expected))
		})

		It("should be possible to save configuration to file", func() {

			configuration.Set("G101", map[string]string{
				"mode": "strict",
			})

			buffer := bytes.NewBuffer([]byte{})
			nbytes, err := configuration.WriteTo(buffer)
			Expect(int(nbytes)).ShouldNot(BeZero())
			Expect(err).ShouldNot(HaveOccurred())
			Expect(buffer.String()).Should(Equal(`{"G101":{"mode":"strict"},"global":{}}`))

		})
	})

	Context("when configuring rules", func() {

		It("should be possible to get configuration for a rule", func() {
			settings := map[string]string{
				"ciphers": "AES256-GCM",
			}
			configuration.Set("G101", settings)

			retrieved, err := configuration.Get("G101")
			Expect(err).ShouldNot(HaveOccurred())
			Expect(retrieved).Should(HaveKeyWithValue("ciphers", "AES256-GCM"))
			Expect(retrieved).ShouldNot(HaveKey("foobar"))
		})
	})

	Context("when using global configuration options", func() {
		It("should have a default global section", func() {
			settings, err := configuration.Get("global")
			Expect(err).Should(BeNil())
			expectedType := make(map[string]string)
			Expect(settings).Should(BeAssignableToTypeOf(expectedType))
		})

		It("should save global settings to correct section", func() {
			configuration.SetGlobal("nosec", "enabled")
			settings, err := configuration.Get("global")
			Expect(err).Should(BeNil())
			if globals, ok := settings.(map[string]string); ok {
				Expect(globals["nosec"]).Should(MatchRegexp("enabled"))
			} else {
				Fail("globals are not defined as map")
			}

			setValue, err := configuration.GetGlobal("nosec")
			Expect(err).Should(BeNil())
			Expect(setValue).Should(MatchRegexp("enabled"))
		})
	})
})
