package evaluator

import (
	"../ast"
	"../object"
	"fmt"
	"math"
	"strings"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.BreakStatement:
		return &object.LoopControlStatement{Literal: "break"}
	case *ast.NextStatement:
		return &object.LoopControlStatement{Literal: "next"}
	case *ast.NumberLiteral:
		return &object.Number{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)
	case *ast.Null:
		return NULL
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		if node.Operator == "." {
			return evalObjectAccessExpression(left, node.Right)
		}

		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right, env)
	case *ast.DeclareExpression:
		right := Eval(node.Value, env)
		if isError(right) {
			return right
		}

		return evalDeclareExpression(node.Name, right, env)
	case *ast.AssignExpression:
		right := Eval(node.Value, env)
		if isError(right) {
			return right
		}

		return evalAssignExpression(node.Name, right, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Env: env, Body: body}
	case *ast.LambdaExpression:
		params := node.Parameters
		body := node.Body
		return &object.Lambda{Parameters: params, Env: env, Body: &body}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args, env)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.HashLiteral:
		return evalHashLiteral(node, env)
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index)
	case *ast.WhileExpression:
		return evalWhileExpression(node, env)
	case *ast.ForExpression:
		return evalForExpression(node, env)
	case *ast.ModelLiteral:
		return evalModelLiteral(node, env)
	default:
		return newError("evaluation for %T not yet implemented!", node)
	}

	return nil
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	if len(program.Statements) == 0 {
		return NULL
	}

	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	if _, ok := result.(*object.LoopControlStatement); ok {
		return NULL
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	if len(block.Statements) == 0 {
		return NULL
	}

	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ ||
				rt == object.ERROR_OBJ ||
				rt == object.LOOP_CONTROL_STATEMENT_OBJ {
				return result
			}
		}
	}

	return result
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "~":
		return evalBitwiseNotPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case FALSE, NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: -value}
}

func evalBitwiseNotPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.NUMBER_OBJ {
		return newError("unknown operator: ~%s", right.Type())
	}

	value := right.(*object.Number).Value
	return &object.Number{Value: float64(^uint64(value))}
}

func evalDeclareExpression(
	left ast.Expression,
	right object.Object,
	env *object.Environment,
) object.Object {
	switch left := left.(type) {
	case *ast.Identifier:
		return env.Declare(left.Value, right)
	case *ast.IndexExpression:
		return newError("cannot declare (:=) a hash field. try assigning (=)")
	default:
		return newError("cannot declare %v. expected an id or index expression",
			left.String())
	}
}

func evalAssignExpression(
	left ast.Expression,
	right object.Object,
	env *object.Environment,
) object.Object {
	switch left := left.(type) {
	case *ast.Identifier:
		return env.Assign(left.Value, right)
	case *ast.IndexExpression:
		return evalAssignIndexExpression(left, right, env)
	case *ast.InfixExpression:
		return evalAssignInfixExpression(left, right, env)
	default:
		return newError("not id")
	}
}

func evalAssignInfixExpression(
	left *ast.InfixExpression,
	right object.Object,
	env *object.Environment,
) object.Object {
	if left.Operator != "." {
		return newError("cannot assign any infix operator other than '.'")
	}

	obj := Eval(left.Left, env)
	if isError(obj) {
		return obj
	}

	fieldId, ok := left.Right.(*ast.Identifier)
	if !ok {
		return newError("expected an identifier")
	}

	switch obj := obj.(type) {
	case *object.Hash:
		obj.Set(fieldId.Value, right)
		return obj.Get(fieldId.Value)
	case *object.Model:
		fn := right
		if fn.Type() != "FUNCTION" {
			return newError("cannot assign a %v to a model field. expected a function",
				right.Type())
		}

		obj.Methods[fieldId] = fn
		return obj.Methods[fieldId]
	default:
		return newError("cannot assign fields of a %v. expected a hash or model",
			obj.Type())
	}
}

func evalAssignIndexExpression(
	left *ast.IndexExpression,
	right object.Object,
	env *object.Environment,
) object.Object {
	obj := Eval(left.Left, env)
	elem := Eval(left.Index, env)

	switch obj := obj.(type) {
	case *object.Array:
		index, ok := elem.(*object.Number)
		if !ok {
			return newError("expected a number for an array index, not %v",
				elem.Inspect())
		}

		if float64(int(index.Value)) != index.Value {
			return newError("expected an integral number for an array index. got a real")
		}

		return assignArrayIndex(obj, int(index.Value), right)
	case *object.Hash:
		key, ok := elem.(*object.String)
		if !ok {
			return newError("expected a string for a hash key, not %v",
				elem.Inspect())
		}

		return assignHashKey(obj, *key, right)
	default:
		return newError("cannot index %v\n", obj.Inspect())
	}

	return right
}

func assignArrayIndex(
	array *object.Array,
	index int,
	val object.Object,
) object.Object {
	for index >= len(array.Elements) {
		index %= len(array.Elements)
	}

	for index < 0 {
		index += len(array.Elements)
	}

	array.Elements[index] = val
	return array
}

func assignHashKey(
	hash *object.Hash,
	key object.String,
	val object.Object,
) object.Object {
	hash.Set(key.Value, val)
	return hash
}

func evalInfixExpression(operator string, left, right object.Object, env *object.Environment) object.Object {
	switch {
	case left.Type() == object.HASH_OBJ || right.Type() == object.HASH_OBJ:
		return evalHashInfixExpression(operator, left, right, env)
	case operator == "&&":
		return nativeBoolToBooleanObject(isTruthy(left) && isTruthy(right))
	case operator == "||":
		return nativeBoolToBooleanObject(isTruthy(left) || isTruthy(right))
	case operator == "==":
		return nativeBoolToBooleanObject(left.Equals(right))
	case operator == "!=":
		return nativeBoolToBooleanObject(!left.Equals(right))
	case operator == "in":
		return evalInOperator(operator, left, right)
	case left.Type() == object.NUMBER_OBJ && right.Type() == object.NUMBER_OBJ:
		return evalNumberInfixExpression(operator, left, right)
	case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s",
			left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalInOperator(operator string, left, right object.Object) object.Object {
	if right.Type() == object.ARRAY_OBJ {
		arr := right.(*object.Array)
		for _, obj := range arr.Elements {
			if left.Equals(obj) {
				return TRUE
			}
		}

		return FALSE
	} else if right.Type() == object.HASH_OBJ {
		if left.Type() != object.STRING_OBJ {
			return newError("expected a string for a hash key, got %v",
				left.Inspect())
		}

		key := left.(*object.String)

		hash := right.(*object.Hash)
		for k, _ := range hash.Pairs {
			if key.Value == k.Value {
				return TRUE
			}
		}

		return FALSE
	} else if right.Type() == object.STRING_OBJ {
		rightString := right.(*object.String).Value
		var s string

		switch left := left.(type) {
		case *object.Number:
			s = left.Inspect()
		case *object.String:
			s = left.Value
		default:
			return newError("expected a string or number to the left of 'in <string>'. got %v",
				left.Inspect())
		}

		return nativeBoolToBooleanObject(strings.Contains(rightString, s))
	} else if right.Type() == object.NUMBER_OBJ {
		if left.Type() != object.NUMBER_OBJ {
			return newError("expected a number to the left of 'in <number>'. got %v",
				left.Inspect())
		}

		leftVal := left.(*object.Number).Value
		rightVal := right.(*object.Number).Value

		return nativeBoolToBooleanObject(math.Mod(rightVal, leftVal) == 0)
	}

	return newError("expected a hash, array, string, or number to the right of 'in'. got %v",
		right.Inspect())
}

func evalNumberInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Number).Value
	rightVal := right.(*object.Number).Value

	switch operator {
	case "+":
		return &object.Number{Value: leftVal + rightVal}
	case "-":
		return &object.Number{Value: leftVal - rightVal}
	case "*":
		return &object.Number{Value: leftVal * rightVal}
	case "**":
		return &object.Number{Value: math.Pow(leftVal, rightVal)}
	case "/":
		return &object.Number{Value: leftVal / rightVal}
	case "%":
		return &object.Number{Value: math.Mod(leftVal, rightVal)}
	case "<":
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case ">":
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	case ">=":
		return nativeBoolToBooleanObject(leftVal >= rightVal)
	case "<=":
		return nativeBoolToBooleanObject(leftVal <= rightVal)
	case "<<":
		return &object.Number{Value: float64(int64(leftVal) << uint64(rightVal))}
	case ">>":
		return &object.Number{Value: float64(int64(leftVal) >> uint64(rightVal))}
	case "&":
		return &object.Number{Value: float64(int64(leftVal) & int64(rightVal))}
	case "|":
		return &object.Number{Value: float64(int64(leftVal) | int64(rightVal))}
	case "..":
		min := int64(leftVal)
		max := int64(rightVal)

		a := make([]object.Object, max-min+1)
		for i := range a {
			a[i] = &object.Number{Value: float64(min + int64(i))}
		}

		return &object.Array{Elements: a}
	case "..<":
		min := int64(leftVal)
		max := int64(rightVal)

		a := make([]object.Object, max-min)
		for i := range a {
			a[i] = &object.Number{Value: float64(min + int64(i))}
		}

		return &object.Array{Elements: a}
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalHashInfixExpression(
	operator string,
	left, right object.Object,
	env *object.Environment,
) object.Object {
	ops := map[string]string{
		"+":   "plus",
		"-":   "minus",
		"*":   "mul",
		"/":   "div",
		"**":  "exp",
		"%":   "mod",
		"<":   "lt",
		">":   "gt",
		"==":  "eq",
		"!=":  "n_eq",
		">=":  "gt_eq",
		"<=":  "lt_eq",
		"<<":  "bit_left",
		">>":  "bit_right",
		"..":  "range",
		"..<": "xrange",
		"&&":  "and",
		"||":  "or",
		"&":   "bit_and",
		"|":   "bit_or",
		"in":  "in",
	}

	f, ok := ops[operator]
	if !ok {
		return newError("operator %v cannot be overloaded", operator)
	}

	fnName := "_" + f

	var hash *object.Hash
	var operand object.Object

	if f == "in" {
		hash = right.(*object.Hash)
		operand = left
	} else {
		hash = left.(*object.Hash)
		operand = right
	}

	if o := hash.Get(fnName); o != NULL {
		if method, ok := o.(*object.MethodInstance); ok {
			result := applyFunctionWithThisValue(method, hash, []object.Object{operand}, env)
			return result
		} else {
			return newError("%v must be a method, not a property",
				fnName)
		}
	} else {
		return newError("operator %v not overloaded. to overload, use the special method %v",
			operator, fnName)
	}
}

func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case "!=":
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s",
			left.Type(), operator, right.Type())
	}
}

func evalObjectAccessExpression(left object.Object, right ast.Expression) object.Object {
	switch right := right.(type) {
	case *ast.Identifier:
		switch left := left.(type) {
		case *object.Hash:
			val := left.Get(right.Value)

			return val
		case *object.Model:
			id := right.Value

			if meth, ok := left.GetMethod(id); ok {
				return meth
			} else {
				return NULL
			}
		default:
			return newError("cannot access %v, expected a hash or a model", left.Inspect())
		}
	default:
		return newError("not an identifier")
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isError(condition) {
		return condition
	} else if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalWhileExpression(we *ast.WhileExpression, env *object.Environment) object.Object {
	result := &object.Array{Elements: []object.Object{}}

	for {
		condition := Eval(we.Condition, env)
		if isError(condition) {
			return condition
		}

		if isTruthy(condition) {
			res := Eval(we.Body, env)
			if isError(res) {
				return res
			}

			if ret, ok := res.(*object.ReturnValue); ok {
				return ret
			}

			if lcs, ok := res.(*object.LoopControlStatement); ok {
				if lcs.Literal == "break" {
					break
				} else {
					continue
				}
			}

			if res != nil {
				result.Elements = append(result.Elements, res)
			}
		} else {
			break
		}
	}

	return result
}

func evalForExpression(fe *ast.ForExpression, env *object.Environment) object.Object {
	set := Eval(fe.Set, env)

	varName, ok := fe.Var.(*ast.Identifier)
	if !ok {
		return newError("expected an identifier as the iterator name")
	}

	switch set := set.(type) {
	case *object.Array:
		return evalForExpressionOverArray(set, varName, fe.Body, env)
	case *object.Hash:
		return evalForExpressionOverHash(set, varName, fe.Body, env)
	case *object.String:
		return evalForExpressionOverString(set, varName, fe.Body, env)
	default:
		return newError("invalid set %v. expected an array or a hash", set.Inspect())
	}
}

func evalForExpressionOverArray(
	array *object.Array,
	varName *ast.Identifier,
	body *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	result := &object.Array{Elements: []object.Object{}}

	for elem, _ := range array.Elements {
		e := object.NewEnclosedEnvironment(env)
		e.Declare(varName.Value, &object.Number{Value: float64(elem)})

		res := Eval(body, e)
		if isError(res) {
			return res
		}

		if ret, ok := res.(*object.ReturnValue); ok {
			return ret
		}

		if lcs, ok := res.(*object.LoopControlStatement); ok {
			if lcs.Literal == "break" {
				break
			} else {
				continue
			}
		}

		if res != nil && res != NULL {
			result.Elements = append(result.Elements, res)
		}
	}

	return result
}

func evalForExpressionOverString(
	str *object.String,
	varName *ast.Identifier,
	body *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	result := &object.String{Value: ""}

	for elem := 0; elem < len(str.Value); elem++ {
		e := object.NewEnclosedEnvironment(env)
		e.Declare(varName.Value, &object.Number{Value: float64(elem)})

		res := Eval(body, e)
		if isError(res) {
			return res
		}

		if ret, ok := res.(*object.ReturnValue); ok {
			return ret
		}

		if lcs, ok := res.(*object.LoopControlStatement); ok {
			if lcs.Literal == "break" {
				break
			} else {
				continue
			}
		}

		if res != nil && res != NULL {
			result.Value += res.Inspect()
		}
	}

	return result
}

func evalForExpressionOverHash(
	hash *object.Hash,
	varName *ast.Identifier,
	body *ast.BlockStatement,
	env *object.Environment,
) object.Object {
	result := object.NewHash(object.OBJECT_MODEL)
	result.Pairs = make(map[object.String]object.Object)

	for key, _ := range hash.Pairs {
		e := object.NewEnclosedEnvironment(env)
		e.Declare(varName.Value, &key)

		res := Eval(body, e)
		if isError(res) {
			return res
		}

		if ret, ok := res.(*object.ReturnValue); ok {
			return ret
		}

		if lcs, ok := res.(*object.LoopControlStatement); ok {
			if lcs.Literal == "break" {
				break
			} else {
				continue
			}
		}

		if res != nil && res != NULL {
			result.Set(key.Value, res)
		}
	}

	return result
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ && index.Type() == object.STRING_OBJ:
		return evalHashIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ && index.Type() == object.NUMBER_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return newError("index operator not supported: %s[%s]",
			left.Type(), index.Type())
	}
}

func evalStringIndexExpression(str, index object.Object) object.Object {
	stringObject := str.(*object.String)
	idx := int64(index.(*object.Number).Value)
	max := int64(len(stringObject.Value) - 1)

	for idx > max {
		idx %= int64(len(stringObject.Value))
	}

	for idx < 0 {
		idx += int64(len(stringObject.Value))
	}

	return &object.String{Value: string(stringObject.Value[idx])}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	idx := int64(index.(*object.Number).Value)
	max := int64(len(arrayObject.Elements) - 1)

	for idx > max {
		idx %= int64(len(arrayObject.Elements))
	}

	for idx < 0 {
		idx += int64(len(arrayObject.Elements))
	}

	return arrayObject.Elements[idx]
}

func evalHashIndexExpression(hash, index object.Object) object.Object {
	hashObject := hash.(*object.Hash)
	key := index.(*object.String)

	return hashObject.Get(key.Value)
}

func evalHashLiteral(node *ast.HashLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.String]object.Object)

	for keyNode, valueNode := range node.Pairs {
		value := Eval(valueNode, env)
		if isError(value) {
			return value
		}

		switch keyNode := keyNode.(type) {
		case *ast.Identifier:
			str := object.String{Value: keyNode.Value}
			pairs[str] = value
		default:
			key := Eval(keyNode, env)
			if isError(key) {
				return key
			}

			str, ok := key.(*object.String)
			if !ok {
				return newError("unusable as hash key: %s. expected a string",
					key.Type())
			}

			pairs[*str] = value
		}
	}

	hash := object.NewHash(object.OBJECT_MODEL)
	hash.Pairs = pairs
	return hash
}

func evalModelLiteral(node *ast.ModelLiteral, env *object.Environment) object.Object {
	model := object.NewModel()
	model.Properties = node.Parameters

	if node.ParentName != nil {
		parent := Eval(*node.ParentName, env)
		if isError(parent) {
			return parent
		}
		model.Parent = parent.(*object.Model)

		model.ParentArgs = node.ParentArgs
	}

	return model
}

func applyFunction(
	fn object.Object,
	args []object.Object,
	env *object.Environment,
) object.Object {
	return applyFunctionWithThisValue(fn, nil, args, env)
}

func applyFunctionWithThisValue(
	fn object.Object,
	thisValue object.Object,
	args []object.Object,
	env *object.Environment,
) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		if len(fn.Parameters) != len(args) {
			return newError("invalid number of arguments. expected %v, got %v",
				len(fn.Parameters), len(args))
		}

		extendedEnv := extendFunctionEnv(fn, args)
		extendedEnv.Declare("this", thisValue)
		evaluated := Eval(fn.Body, extendedEnv)
		if isError(evaluated) {
			return evaluated
		}

		if evaluated == nil {
			return NULL
		}

		return unwrapReturnValue(evaluated)
	case *object.Lambda:
		if len(fn.Parameters) != len(args) {
			return newError("invalid number of arguments. expected %v, got %v",
				len(fn.Parameters), len(args))
		}

		extendedEnv := extendLambdaEnv(fn, args)
		extendedEnv.Declare("this", thisValue)
		evaluated := Eval(*fn.Body, extendedEnv)

		if evaluated == nil {
			return NULL
		}

		return unwrapReturnValue(evaluated)
	case *object.Model:
		m := fn

		hash := m.Instantiate(args).(*object.Hash)

		enclosedEnv := object.NewEnclosedEnvironment(env)
		for i, prop := range m.Properties {
			enclosedEnv.Declare(prop.Value, args[i])
		}

		if m.Parent != nil {
			for i, name := range m.Parent.Properties {
				val := Eval(m.ParentArgs[i], enclosedEnv)
				if isError(val) {
					return val
				}

				hash.Set(name.Value, val)
			}
		}

		if _new, ok := hash.Model.GetMethod("_new"); ok {
			_new.Hash = hash
			return applyFunction(_new, []object.Object{}, env)
		}
		return hash
	case *object.MethodInstance:
		res := applyFunctionWithThisValue(*fn.Function, fn.Hash, args, env)

		return res
	case *object.Builtin:
		return fn.Fn(thisValue, args...)
	default:
		return newError("cannot call a %s", fn.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Declare(param.Value, args[paramIdx])
	}

	return env
}

func extendLambdaEnv(fn *object.Lambda, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Declare(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return obj
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
