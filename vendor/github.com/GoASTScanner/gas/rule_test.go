package gas_test

import (
	"fmt"
	"go/ast"

	"github.com/GoASTScanner/gas"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type mockrule struct {
	issue    *gas.Issue
	err      error
	callback func(n ast.Node, ctx *gas.Context) bool
}

func (m *mockrule) ID() string {
	return "MOCK"
}

func (m *mockrule) Match(n ast.Node, ctx *gas.Context) (*gas.Issue, error) {
	if m.callback(n, ctx) {
		return m.issue, nil
	}
	return nil, m.err
}

var _ = Describe("Rule", func() {

	Context("when using a ruleset", func() {

		var (
			ruleset        gas.RuleSet
			dummyErrorRule gas.Rule
			dummyIssueRule gas.Rule
		)

		JustBeforeEach(func() {
			ruleset = gas.NewRuleSet()
			dummyErrorRule = &mockrule{
				issue:    nil,
				err:      fmt.Errorf("An unexpected error occurred"),
				callback: func(n ast.Node, ctx *gas.Context) bool { return false },
			}
			dummyIssueRule = &mockrule{
				issue: &gas.Issue{
					Severity:   gas.High,
					Confidence: gas.High,
					What:       `Some explanation of the thing`,
					File:       "main.go",
					Code:       `#include <stdio.h> int main(){ puts("hello world"); }`,
					Line:       "42",
				},
				err:      nil,
				callback: func(n ast.Node, ctx *gas.Context) bool { return true },
			}
		})
		It("should be possible to register a rule for multiple ast.Node", func() {
			registeredNodeA := (*ast.CallExpr)(nil)
			registeredNodeB := (*ast.AssignStmt)(nil)
			unregisteredNode := (*ast.BinaryExpr)(nil)

			ruleset.Register(dummyIssueRule, registeredNodeA, registeredNodeB)
			Expect(ruleset.RegisteredFor(unregisteredNode)).Should(BeEmpty())
			Expect(ruleset.RegisteredFor(registeredNodeA)).Should(ContainElement(dummyIssueRule))
			Expect(ruleset.RegisteredFor(registeredNodeB)).Should(ContainElement(dummyIssueRule))

		})

		It("should not register a rule when no ast.Nodes are specified", func() {
			ruleset.Register(dummyErrorRule)
			Expect(ruleset).Should(BeEmpty())
		})

		It("should be possible to retrieve a list of rules for a given node type", func() {
			registeredNode := (*ast.CallExpr)(nil)
			unregisteredNode := (*ast.AssignStmt)(nil)
			ruleset.Register(dummyErrorRule, registeredNode)
			ruleset.Register(dummyIssueRule, registeredNode)
			Expect(ruleset.RegisteredFor(unregisteredNode)).Should(BeEmpty())
			Expect(ruleset.RegisteredFor(registeredNode)).Should(HaveLen(2))
			Expect(ruleset.RegisteredFor(registeredNode)).Should(ContainElement(dummyErrorRule))
			Expect(ruleset.RegisteredFor(registeredNode)).Should(ContainElement(dummyIssueRule))
		})

	})

})
