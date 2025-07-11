package parser

import (
	"compiler/internal/frontend/ast"
	"compiler/internal/frontend/lexer"
	"compiler/internal/report"
	"compiler/internal/source"
)

func parseAssignment(p *Parser, left ...ast.Expression) ast.Statement {
	assignees := ast.ExpressionList{}

	for _, expr := range left {
		assignees = append(assignees, expr)
	}

	expressions := ast.ExpressionList{}

	for p.peek().Kind == lexer.COMMA_TOKEN {
		p.advance()
		val := parseExpression(p)
		if val == nil {
			current := p.previous()
			p.ctx.Reports.Add(p.fullPath, source.NewLocation(&current.Start, &current.End), "Expected expression in assignment", report.PARSING_PHASE).SetLevel(report.SYNTAX_ERROR)
		}
		assignees = append(assignees, val)
	}

	p.consume(lexer.EQUALS_TOKEN, "Expected '=' in assignment")

	for {
		val := parseExpression(p)
		if val == nil {
			current := p.previous()
			p.ctx.Reports.Add(p.fullPath, source.NewLocation(&current.Start, &current.End), "Expected expression in assignment", report.PARSING_PHASE).SetLevel(report.SYNTAX_ERROR)
		}
		expressions = append(expressions, val)
		if p.peek().Kind == lexer.COMMA_TOKEN {
			p.advance()
		} else {
			break
		}
	}

	if len(assignees) < len(expressions) {
		current := p.previous()
		p.ctx.Reports.Add(p.fullPath, source.NewLocation(&current.Start, &current.End), "Mismatched number of variables and values", report.PARSING_PHASE).AddHint("Assignee count must be less than or equal to the number of expressions").SetLevel(report.SYNTAX_ERROR)
	}

	return &ast.AssignmentStmt{
		Left:     &assignees,
		Right:    &expressions,
		Location: *source.NewLocation(assignees[0].Loc().Start, expressions[len(expressions)-1].Loc().End),
	}
}
