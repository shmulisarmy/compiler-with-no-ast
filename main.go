package main

import (
	"fmt"
	"no-ast/tokenizer"
	"no-ast/utils/assert"
	"no-ast/utils/displayStruct"
	"os"
	"strconv"
)

type Opcode int

const (
	OPCODE_ADD Opcode = iota
	OPCODE_SUB
	OPCODE_MUL
	OPCODE_DIV
	OPCODE_EQ
	OPCODE_GT
	OPCODE_LT
	Pop
	Push
	LoadVar
	Assign
	StackTopType
	Blank
	Invoke_function_on_stack_top
	Return
	LoadLocal
	JumpIfZero
	SetLocal
)

type VarInfo struct {
	Name       string
	Type       string
	mem_offset int
}

type Instruction struct {
	Opcode   Opcode
	Operands []any
}

func (this Instruction) String() string {
	res := ""
	switch this.Opcode {
	case Pop:
		res += "POP"
	case Push:
		res += "PUSH"
	case LoadVar:
		res += "LOADVAR"
	case Assign:
		res += "ASSIGN"
	case StackTopType:
		res += "STACKTOPTYPE"
	case Blank:
		res += "BLANK"
	case OPCODE_ADD:
		res += "ADD"
	case OPCODE_SUB:
		res += "SUB"
	case OPCODE_MUL:
		res += "MUL"
	case OPCODE_DIV:
		res += "DIV"
	case OPCODE_EQ:
		res += "EQ"
	case Invoke_function_on_stack_top:
		res += "INVOKE_FUNCTION_ON_STACK_TOP"
	case Return:
		res += "RETURN"
	case JumpIfZero:
		res += "JUMP_IF_ZERO"
	case LoadLocal:
		res += "LOAD_LOCAL"
	case OPCODE_GT:
		res += "LOAD_LOCAL"
	case OPCODE_LT:
		res += "LOAD_LOCAL"
	default:
		panic(fmt.Sprintf("Unknown opcode: %d", this.Opcode))
	}
	for _, operand := range this.Operands {
		res += " " + fmt.Sprint(operand)
	}
	return res

}

type Parser struct {
	tokens []tokenizer.Token
	index  int
}

func (p *Parser) in_range() bool {
	return p.index < len(p.tokens)
}

func (p *Parser) cur_token() tokenizer.Token {
	return p.tokens[p.index]
}

func (p *Parser) parse_wrapped_term() []Instruction {
	exp_byte_code := p.parse_term()
	was_ident := exp_byte_code[len(exp_byte_code)-1].Opcode != StackTopType
	if !was_ident {
		return exp_byte_code
	}
	fmt.Print(p.cur_token().Value, "p.tokens[p.index].Value")
	for p.in_range() {
		state_changed := false
		if p.tokens[p.index].Value == "(" {
			p.index++
			arg_count := 0
			for p.in_range() && p.cur_token().Value != ")" {
				exp_byte_code = append(exp_byte_code, p.parse_expression()...)
				arg_count++
				if p.cur_token().Value == "," {
					p.index++
				} else {
					break
				}
			}
			p.expect_token(")")
			exp_byte_code = append(exp_byte_code, Instruction{Opcode: Invoke_function_on_stack_top, Operands: []any{arg_count}})
			state_changed = true
		}
		if !state_changed {
			break
		}
	}
	return exp_byte_code
}

func (p *Parser) expect_token(token_value string) {
	if p.NextToken().Value != token_value {
		panic(fmt.Sprintf("Expected token %s, got %s", token_value, p.NextToken().Value))
	}
}
func (p *Parser) parse_term() []Instruction {
	t := p.NextToken()
	switch t.Type {
	case tokenizer.TOKEN_NUMBER:
		n, _ := strconv.Atoi(t.Value)
		return []Instruction{Instruction{Opcode: Push, Operands: []any{n}}} // Instruction{Opcode: StackTopType, Operands: []any{"int"}}

	case tokenizer.TOKEN_STRING:
		return []Instruction{Instruction{Opcode: Push, Operands: []any{t.Value}}} //  Instruction{Opcode: StackTopType, Operands: []any{"string"}}

	case tokenizer.TOKEN_IDENTIFIER:
		if in_function {
			if v, ok := current_parsing_function.local_vars[t.Value]; ok {
				return []Instruction{Instruction{Opcode: LoadLocal, Operands: []any{
					v.mem_offset, v.Type,
				}}}
			}
		}
		return []Instruction{Instruction{Opcode: LoadVar, Operands: []any{t.Value}}}
	default:
		panic("Unexpected token: " + t.String())
	}

}
func (p *Parser) parse_expression() []Instruction {
	instructions := p.parse_wrapped_term()
	for p.in_range() && p.tokens[p.index].Type == tokenizer.TOKEN_OPERATOR {
		t := p.NextToken()
		operation_byte_code := Instruction{Opcode: Blank}
		switch t.Value {
		case "+":
			operation_byte_code = Instruction{Opcode: OPCODE_ADD}
		case "-":
			operation_byte_code = Instruction{Opcode: OPCODE_SUB}
		case "*":
			operation_byte_code = Instruction{Opcode: OPCODE_MUL}
		case "==":
			operation_byte_code = Instruction{Opcode: OPCODE_EQ}
		case ">":
			operation_byte_code = Instruction{Opcode: OPCODE_GT}
		case "<":
			operation_byte_code = Instruction{Opcode: OPCODE_LT}
		default:
			panic("Unexpected operator: " + t.String())
		}
		right := p.parse_expression()
		instructions = append(instructions, right...)
		instructions = append(instructions, operation_byte_code)
	}
	return instructions
}

func (p *Parser) parse_statement(previous_instruction_amount int) []Instruction {
	instructions := []Instruction{}
	t := p.NextToken()
	if t.Value == "return" {
		p.index++
		return append(instructions, Instruction{Opcode: Return, Operands: []any{}})
	}
	if t.Value == "if" {
		instructions = append(instructions, p.parse_expression()...)
		instructions = append(instructions, Instruction{Opcode: JumpIfZero, Operands: []any{}})
		conditional_jump_instruction_index := len(instructions) - 1
		p.expect_token("{")
		for p.in_range() && p.cur_token().Value != "}" {
			instructions = append(instructions, p.parse_statement(previous_instruction_amount+len(instructions))...)
		}
		p.expect_token("}")
		instructions[conditional_jump_instruction_index].Operands = append(instructions[conditional_jump_instruction_index].Operands, previous_instruction_amount+len(instructions))
		return instructions
	}
	if t.Value == "while" {
		start_index := len(instructions)
		instructions = append(instructions, p.parse_expression()...)
		instructions = append(instructions, Instruction{Opcode: JumpIfZero, Operands: []any{}})
		conditional_jump_instruction_index := len(instructions) - 1
		p.expect_token("{")
		for p.in_range() && p.cur_token().Value != "}" {
			instructions = append(instructions, p.parse_statement(previous_instruction_amount+len(instructions))...)
		}
		p.expect_token("}")
		instructions = append(instructions, Instruction{Opcode: Push, Operands: []any{0}})
		instructions = append(instructions, Instruction{Opcode: JumpIfZero, Operands: []any{previous_instruction_amount + start_index}})
		instructions[conditional_jump_instruction_index].Operands = append(instructions[conditional_jump_instruction_index].Operands, previous_instruction_amount+len(instructions))
		return instructions
	}
	if p.in_range() && p.cur_token().Value == "=" {
		p.index++
		if in_function {
			if v, ok := current_parsing_function.local_vars[t.Value]; ok {
				instructions = append(instructions, p.parse_expression()...)
				return append(instructions, Instruction{Opcode: SetLocal, Operands: []any{
					v.mem_offset, v.Type,
				}})
			}
		}
		instructions = append(instructions, p.parse_expression()...)
		return append(instructions, Instruction{Opcode: Assign, Operands: []any{t.Value}})
	}
	if p.in_range() && p.cur_token().Value == "(" {
		p.index--
		return append(instructions, p.parse_wrapped_term()...)
	}
	panic("Unexpected statement: " + t.String())

}
func (p *Parser) NextToken() tokenizer.Token {
	if !p.in_range() {
		return tokenizer.Token{Type: tokenizer.TOKEN_EOF, Value: ""}
	}
	token := p.tokens[p.index]
	p.index++
	return token
}

func print_sum(nums []any) {
	sum := 0
	for _, num := range nums {
		sum += num.(int)
	}
	fmt.Println(sum)
}

func print_all(args []any) {
	for _, arg := range args {
		fmt.Println(arg)
	}
}

func print_one(args []any) {
	fmt.Println(args[0].(int))
}

var vars = map[string]VarInfo{
	"x":           VarInfo{Name: "x", Type: "int", mem_offset: 0},
	"y":           VarInfo{Name: "y", Type: "int", mem_offset: 1},
	"print_all":   VarInfo{Name: "print_all", Type: "builtin-function", mem_offset: 2},
	"print_one":   VarInfo{Name: "print_one", Type: "builtin-function", mem_offset: 3},
	"print_added": VarInfo{Name: "print_added", Type: "function", mem_offset: 4},
	"done":        VarInfo{Name: "done", Type: "builtin-function", mem_offset: 5},
}

type StackFrame struct {
	return_address              int
	function_locals_start_index int
}
type Function struct {
	Name                    string
	param_types             []string
	return_type             string
	instruction_start_index int
	local_vars              map[string]VarInfo
}

var memory = make([]any, 100)
var bytecode []Instruction

var current_parsing_function = Function{}
var in_function = false

func init() {

	source := `y = 1+90-7 
				x=y+2 
				print_one(y)
				print_one(x)
				print_one(y) 
				print_added(x) 
				done()
	`
	bytecode = block_instructions(source, 0)
	memory[vars["x"].mem_offset] = 0
	memory[vars["y"].mem_offset] = 0
	memory[vars["print_all"].mem_offset] = print_all
	memory[vars["print_one"].mem_offset] = print_one
	memory[vars["done"].mem_offset] = func([]any) {
		fmt.Println("done the program")
		os.Exit(0)
	}
	function_header := Function{Name: "print_added", param_types: []string{"int"}, return_type: "void", instruction_start_index: len(bytecode), local_vars: map[string]VarInfo{
		"x": VarInfo{Name: "x", Type: "int", mem_offset: 0},
	}}
	makeFunction(function_header, "print_added", `print_one(x-1) 
												if x>20{ print_added(x-1)}
												x=0
												while x<10{x=x+1 
													print_one(x) 
													if x == 5{x=x+2}
												} 
												return`)

	function_header = Function{Name: "defers", param_types: []string{"int"}, return_type: "void", instruction_start_index: len(bytecode), local_vars: map[string]VarInfo{
		"x": VarInfo{Name: "x", Type: "int", mem_offset: 0},
	}}
	makeFunction(function_header, "defers", `
					print_one(x)
				`)

}

func makeFunction(function_header Function, function_name string, block_code string) {
	current_parsing_function = function_header
	in_function = true
	memory[vars[function_name].mem_offset] = current_parsing_function
	print_added_instructions := block_instructions(block_code, len(bytecode))
	for _, instruction := range print_added_instructions {
		bytecode = append(bytecode, instruction)
	}
	current_parsing_function = Function{}
	in_function = false
}

type DataAndHeader struct {
	Type string
	Data any
}

var stack = make([]DataAndHeader, 0)
var frames = make([]StackFrame, 0)

func block_instructions(source string, previous_instruction_amount int) []Instruction {
	tokens := tokenizer.Tokenize(source)
	for _, token := range tokens {
		fmt.Println(token.String())
	}
	p := Parser{tokens: tokens, index: 0}
	return p.parse_block(previous_instruction_amount)
}

func (p *Parser) parse_block(previous_instruction_amount int) []Instruction {
	bytecode := []Instruction{}
	for p.in_range() && p.cur_token().Type != tokenizer.TOKEN_EOF {
		for _, instruction := range p.parse_statement(len(bytecode) + previous_instruction_amount) {
			fmt.Println(instruction)
			bytecode = append(bytecode, instruction)
		}
	}
	return bytecode
}

func main() {
	// source := "x = 42 + 3.14 * (y - 1);"
	instruction_ptr := 0
	for instruction_ptr < len(bytecode) {
		instruction := bytecode[instruction_ptr]

		displayStruct.Print(stack)
		displayStruct.Print(instruction)
		switch instruction.Opcode {
		case LoadVar:
			name := instruction.Operands[0].(string)
			// fmt.Println(name, "is name")
			data := memory[vars[name].mem_offset]
			stack = append(stack, DataAndHeader{Type: vars[name].Type, Data: data})
		case Assign:
			name := instruction.Operands[0].(string)
			data := stack_pop()
			memory[vars[name].mem_offset] = data.Data
		case OPCODE_ADD:
			right := stack_pop()
			left := stack_pop()
			stack = append(stack, DataAndHeader{Type: "int", Data: left.Data.(int) + right.Data.(int)})
		case OPCODE_SUB:
			right := stack_pop()
			left := stack_pop()
			if left.Type != "int" || right.Type != "int" {
				panic("sub on non-int")
			}
			stack = append(stack, DataAndHeader{Type: "int", Data: left.Data.(int) - right.Data.(int)})
		case OPCODE_MUL:
			right := stack_pop()
			left := stack_pop()
			stack = append(stack, DataAndHeader{Type: "int", Data: left.Data.(int) * right.Data.(int)})
		case OPCODE_DIV:
			right := stack_pop()
			left := stack_pop()
			stack = append(stack, DataAndHeader{Type: "int", Data: left.Data.(int) / right.Data.(int)})
		case LoadLocal:
			offset := instruction.Operands[0].(int)
			type_ := instruction.Operands[1].(string)
			var_stack_index := frames[len(frames)-1].function_locals_start_index + offset
			stack = append(stack, DataAndHeader{Type: type_, Data: stack[var_stack_index].Data})
		case SetLocal:
			offset := instruction.Operands[0].(int)
			type_ := instruction.Operands[1].(string)
			var_stack_index := frames[len(frames)-1].function_locals_start_index + offset
			if stack[var_stack_index].Type != type_ {
				assert.Assert(stack[len(stack)-1].Type != type_, "something has gone wrong within the compiler")
				panic(fmt.Sprintf("expected %s, got %s", type_, stack[var_stack_index].Type))
			}
			stack[var_stack_index].Data = stack_pop().Data
		case Push:
			d := instruction.Operands[0]
			type_ := ""
			switch d.(type) {
			case int:
				type_ = "int"
			case float64:
				type_ = "int"
			case string:
				type_ = "string"
			default:
				panic(fmt.Sprintf("unhandled type %T", d))
			}
			stack = append(stack, DataAndHeader{Type: type_, Data: d})
		case Invoke_function_on_stack_top:
			arg_count := instruction.Operands[0].(int)
			// println(arg_count, "arg_count")
			function := stack[len(stack)-1-arg_count]
			if function.Type == "builtin-function" {
				args := make([]any, arg_count)
				for i := 0; i < arg_count; i++ {
					args[i] = stack[len(stack)-arg_count+i].Data
				}
				function.Data.(func([]any))(args)
				stack = stack[:len(stack)-1-arg_count]

			} else if function.Type == "function" {
				if len(function.Data.(Function).param_types) != arg_count {
					panic(fmt.Sprintf("expected %d args, got %d", len(function.Data.(Function).param_types), arg_count))
				}
				for i := 0; i < arg_count; i++ {
					if function.Data.(Function).param_types[i] != stack[len(stack)-arg_count+i].Type {
						panic(fmt.Sprintf("expected %s, got %s for arg %d", function.Data.(Function).param_types[i], stack[len(stack)-arg_count+i].Type, i))
					}
				}
				frames = append(frames, StackFrame{return_address: len(bytecode), function_locals_start_index: len(stack) - arg_count})
				instruction_ptr = function.Data.(Function).instruction_start_index
				continue
			} else {
				panic("unhandled Invoke_function_on_stack_top")
			}

		case JumpIfZero:
			if stack[len(stack)-1].Type != "int" {
				panic("jump_if_zero on non-int")
			}
			if stack_pop().Data.(int) == 0 {
				instruction_ptr = instruction.Operands[0].(int)
				continue
			}
		case Return:
			if len(frames) == 0 {
				panic("return from main")
			}
			frame := frames[len(frames)-1]
			frames = frames[:len(frames)-1]
			stack = stack[:frame.function_locals_start_index]
			instruction_ptr = frame.return_address
			continue
		case OPCODE_GT:
			right := stack_pop()
			left := stack_pop()
			if left.Type != "int" || right.Type != "int" {
				panic("sub on non-int")
			}
			if left.Data.(int) > right.Data.(int) {
				stack = append(stack, DataAndHeader{Type: "int", Data: 1})
			} else {
				stack = append(stack, DataAndHeader{Type: "int", Data: 0})
			}
		case OPCODE_LT:
			right := stack_pop()
			left := stack_pop()
			if left.Type != "int" || right.Type != "int" {
				panic("sub on non-int")
			}
			if left.Data.(int) < right.Data.(int) {
				stack = append(stack, DataAndHeader{Type: "int", Data: 1})
			} else {
				stack = append(stack, DataAndHeader{Type: "int", Data: 0})
			}
		case OPCODE_EQ:
			right := stack_pop()
			left := stack_pop()
			if left.Type != "int" || right.Type != "int" {
				panic("sub on non-int")
			}
			if left.Data.(int) == right.Data.(int) {
				stack = append(stack, DataAndHeader{Type: "int", Data: 1})
			} else {
				stack = append(stack, DataAndHeader{Type: "int", Data: 0})
			}
		default:
			panic(fmt.Sprintf("unhandled " + instruction.String()))
		}
		instruction_ptr++
	}
}

func stack_pop() DataAndHeader {
	v := stack[len(stack)-1]
	stack = stack[:len(stack)-1]
	return v
}
