package generator

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	fp "path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/RadhiFadlillah/qamel/qamel/config"
)

var (
	qamelImportName = "qamel"
	qamelImportPath = "github.com/RadhiFadlillah/qamel"
	qamelObjectName = "QmlObject"

	rxNumber = regexp.MustCompile(`\d`)
	rxSymbol = regexp.MustCompile(`[^A-Za-z0-9]`)
)

type object struct {
	name         string
	dirPath      string
	fileName     string
	packageName  string
	structNode   *ast.StructType
	constructors []objectMethod
	properties   []objectMember
	signals      []objectMethod
	slots        []objectMethod
}

type objectMember struct {
	name       string
	memberType string
}

type objectMethod struct {
	name       string
	parameters []objectMember
	returns    []objectMember
}

// CreateQmlObjectCode generates Go code and C++ code for all QmlObject
// in specified directory
func CreateQmlObjectCode(profile config.Profile, projectDir string, buildTags ...string) []error {
	// Make sure project directory is exists
	if !dirExists(projectDir) {
		err := fmt.Errorf("directory %s doesn't exist", projectDir)
		return []error{err}
	}

	// Find all sub directories, including the root dir
	subDirs, err := getSubDirs(projectDir)
	if err != nil {
		return []error{err}
	}

	// For each sub directories, find Go files with matching build tags
	goFiles := []string{}
	buildCtx := build.Default
	buildCtx.BuildTags = append([]string{}, buildTags...)
	for _, dir := range subDirs {
		// If a directory doesn't have any Go files, this method will throw error.
		// That's why in this part, when error happened just continue.
		dirPkg, err := buildCtx.ImportDir(dir, 0)
		if err != nil {
			continue
		}

		for _, goFile := range dirPkg.GoFiles {
			goFiles = append(goFiles, fp.Join(dir, goFile))
		}
	}

	// From each Go files, find struct with qamel.QmlObject embedded to it
	qmlObjects := []object{}
	for _, goFile := range goFiles {
		objects, err := getQmlObjectStructs(goFile)
		if err != nil {
			return []error{err}
		}
		qmlObjects = append(qmlObjects, objects...)
	}

	// Parse each qml objects
	errors := []error{}
	for i, obj := range qmlObjects {
		tmpObj, tmpErrors := parseQmlObject(obj)
		errors = append(errors, tmpErrors...)
		qmlObjects[i] = tmpObj
	}

	if len(errors) > 0 {
		return errors
	}

	// Create code for each QML objects
	mapDirPackage := map[string]string{}
	for _, obj := range qmlObjects {
		mapDirPackage[obj.dirPath] = obj.packageName

		err = createCppHeaderFile(obj)
		if err != nil {
			return []error{err}
		}

		err = createCppFile(profile.Moc, obj)
		if err != nil {
			return []error{err}
		}

		err = createGoFile(obj)
		if err != nil {
			return []error{err}
		}
	}

	// Create cgo file for each package
	for dir, packageName := range mapDirPackage {
		err = CreateCgoFile(profile, dir, packageName)
		if err != nil {
			return []error{err}
		}
	}

	return nil
}

// getQmlObjectStruct fetch structs that embedded with qamel.QmlObject inside specified Go file
func getQmlObjectStructs(goFile string) ([]object, error) {
	// Parse file
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, goFile, nil, 0)
	if err != nil {
		return nil, err
	}

	// Check if this file importing qamel package
	usingQamel := false
	importName := qamelImportName
	for _, importSpec := range f.Imports {
		importPath := strings.Trim(importSpec.Path.Value, `"`)
		if importPath != qamelImportPath {
			continue
		}

		if importSpec.Name != nil {
			importName = importSpec.Name.Name
		}

		usingQamel = true
	}

	if !usingQamel {
		return []object{}, nil
	}

	// Traverse file tree to find struct that embedding QmlObject
	result := []object{}
	embeddedName := fmt.Sprintf("%s.%s", importName, qamelObjectName)
	ast.Inspect(f, func(node ast.Node) bool {
		// Make sure this node is a declaration of `type`
		typeNode, isTypeSpec := node.(*ast.TypeSpec)
		if !isTypeSpec {
			return true
		}

		// Make sure this type is declaring a struct
		structNode, isStruct := typeNode.Type.(*ast.StructType)
		if !isStruct {
			return false
		}

		// Make sure this struct has fields
		if structNode.Fields == nil || len(structNode.Fields.List) == 0 {
			return false
		}

		// Check this struck embedding QmlObject
		embeddingQamel := false
		for _, field := range structNode.Fields.List {
			selector, isSelector := field.Type.(*ast.SelectorExpr)
			if !isSelector {
				continue
			}

			selectorValue := fmt.Sprintf("%s.%s", selector.X, selector.Sel.Name)
			if selectorValue == embeddedName {
				embeddingQamel = true
				break
			}
		}

		if !embeddingQamel {
			return false
		}

		// Save this struct to list of object
		result = append(result, object{
			name:        typeNode.Name.Name,
			dirPath:     fp.Dir(goFile),
			fileName:    goFile,
			packageName: f.Name.Name,
			structNode:  structNode,
		})

		return false
	})

	return result, nil
}

// parseNode parse struct nodes inside object and find the property, signal and slots
func parseQmlObject(obj object) (object, []error) {
	errors := []error{}
	slots := []objectMethod{}
	signals := []objectMethod{}
	properties := []objectMember{}
	constructors := []objectMethod{}

	nPropName := map[string]int{}
	nSignalName := map[string]int{}
	nSlotName := map[string]int{}

	for _, structField := range obj.structNode.Fields.List {
		// Make sure this field either identity or function
		identField, isIdent := structField.Type.(*ast.Ident)
		funcField, isFunc := structField.Type.(*ast.FuncType)
		if !isIdent && !isFunc {
			continue
		}

		// Get and check field tag
		if structField.Tag == nil {
			continue
		}

		fieldTag := strings.Trim(structField.Tag.Value, "`")
		structTag := reflect.StructTag(fieldTag)

		propName := strings.TrimSpace(structTag.Get("property"))
		signalName := strings.TrimSpace(structTag.Get("signal"))
		slotName := strings.TrimSpace(structTag.Get("slot"))
		constructorName := strings.TrimSpace(structTag.Get("constructor"))
		mergedName := propName + signalName + slotName + constructorName

		if mergedName == "" {
			continue
		}

		if mergedName != propName && mergedName != signalName && mergedName != slotName && mergedName != constructorName {
			err := fmt.Errorf("object %s: a field must be only used for one purpose", obj.name)
			errors = append(errors, err)
			continue
		}

		// Check whether this field is blank or not
		isBlankField := len(structField.Names) == 1 && structField.Names[0].String() == "_"

		// Check if it's property
		if isIdent && propName != "" {
			if !isBlankField {
				err := fmt.Errorf("object %s, property %s: must be a single blank field", obj.name, propName)
				errors = append(errors, err)
				continue
			}

			if err := validateTagName(propName); err != nil {
				err = fmt.Errorf("object %s, property %s: %v", obj.name, propName, err)
				errors = append(errors, err)
				continue
			}

			if nPropName[propName] > 0 {
				err := fmt.Errorf("object %s, property %s: property has been declared before", obj.name, propName)
				errors = append(errors, err)
				continue
			}

			properties = append(properties, objectMember{
				name:       propName,
				memberType: identField.String(),
			})

			nPropName[propName]++
			continue
		}

		// Check if it's constructor
		if isFunc && constructorName != "" {
			if !isBlankField {
				err := fmt.Errorf("object %s, constructor %s: must be a single blank field", obj.name, constructorName)
				errors = append(errors, err)
				continue
			}

			if err := validateTagName(constructorName); err != nil {
				err = fmt.Errorf("object %s, constructor %s: %v", obj.name, constructorName, err)
				errors = append(errors, err)
				continue
			}

			if len(constructors) > 0 {
				err := fmt.Errorf("object %s, constructor %s: other constructor has been declared before", obj.name, constructorName)
				errors = append(errors, err)
				continue
			}

			if funcField.Results != nil {
				err := fmt.Errorf("object %s, constructor %s: must not have return value", obj.name, constructorName)
				errors = append(errors, err)
				continue
			}

			parameters := parseAstFuncParams(funcField.Params)
			if len(parameters) > 0 {
				err := fmt.Errorf("object %s, constructor %s: must not have any parameter", obj.name, constructorName)
				errors = append(errors, err)
				continue
			}

			constructors = append(constructors, objectMethod{
				name: constructorName,
			})

			continue
		}

		// Check if it's signal
		if isFunc && signalName != "" {
			if !isBlankField {
				err := fmt.Errorf("object %s, signal %s: must be a single blank field", obj.name, signalName)
				errors = append(errors, err)
				continue
			}

			if err := validateTagName(signalName); err != nil {
				err = fmt.Errorf("object %s, signal %s: %v", obj.name, signalName, err)
				errors = append(errors, err)
				continue
			}

			if nSignalName[signalName] > 0 {
				err := fmt.Errorf("object %s, signal %s: signal has been declared before", obj.name, signalName)
				errors = append(errors, err)
				continue
			}

			if funcField.Results != nil {
				err := fmt.Errorf("object %s, signal %s: must not have return value", obj.name, signalName)
				errors = append(errors, err)
				continue
			}

			signalParameters := parseAstFuncParams(funcField.Params)
			err := validateMethodType(signalParameters)
			if err != nil {
				err1 := fmt.Errorf("object %s, signal %s: %v", obj.name, signalName, err)
				errors = append(errors, err1)
				continue
			}

			signals = append(signals, objectMethod{
				name:       signalName,
				parameters: signalParameters,
			})

			nSignalName[signalName]++
			continue
		}

		// Check if it's slot
		if isFunc && slotName != "" {
			if !isBlankField {
				err := fmt.Errorf("object %s, slot %s: must be a single blank field", obj.name, slotName)
				errors = append(errors, err)
				continue
			}

			if err := validateTagName(slotName); err != nil {
				err = fmt.Errorf("object %s, slot %s: %v", obj.name, slotName, err)
				errors = append(errors, err)
				continue
			}

			if nSlotName[slotName] > 0 {
				err := fmt.Errorf("object %s, slot %s: slot has been declared before", obj.name, slotName)
				errors = append(errors, err)
				continue
			}

			slotReturns := parseAstFuncParams(funcField.Results)
			if len(slotReturns) > 1 {
				err := fmt.Errorf("object %s, slot %s: only allowed max one return value", obj.name, slotName)
				errors = append(errors, err)
				continue
			}

			err := validateMethodType(slotReturns)
			if err != nil {
				err1 := fmt.Errorf("object %s, slot %s: %v", obj.name, slotName, err)
				errors = append(errors, err1)
				continue
			}

			slotParameters := parseAstFuncParams(funcField.Params)
			err = validateMethodType(slotParameters)
			if err != nil {
				err1 := fmt.Errorf("object %s, slot %s: %v", obj.name, slotName, err)
				errors = append(errors, err1)
				continue
			}

			slots = append(slots, objectMethod{
				name:       slotName,
				parameters: slotParameters,
				returns:    slotReturns,
			})

			nSlotName[slotName]++
		}
	}

	if len(errors) == 0 {
		obj.slots = slots
		obj.signals = signals
		obj.properties = properties
		obj.constructors = constructors
	}

	return obj, errors
}

// parseAstFuncParams converts field list to object member
func parseAstFuncParams(fieldList *ast.FieldList) []objectMember {
	// Make sure field list exists
	if fieldList == nil {
		return []objectMember{}
	}

	// Get name and type of each parameter
	paramIdx := 0
	result := []objectMember{}
	for _, param := range fieldList.List {
		if len(param.Names) == 0 {
			result = append(result, objectMember{
				name:       fmt.Sprintf("p%d", paramIdx),
				memberType: fmt.Sprint(param.Type),
			})
			paramIdx++
			continue
		}

		for _, nameIdent := range param.Names {
			result = append(result, objectMember{
				name:       nameIdent.String(),
				memberType: fmt.Sprint(param.Type),
			})
			paramIdx++
		}
	}

	return result
}

// validateMethodType check if member has unknown type
func validateMethodType(members []objectMember) error {
	for _, member := range members {
		if _, known := mapGoType[member.memberType]; !known {
			return fmt.Errorf("unknown type %s", member.memberType)
		}
	}

	return nil
}

// validateTagName check if member has a valid tag name
func validateTagName(tagName string) error {
	if tagName == "" {
		return nil
	}

	firstChar := tagName[0:1]
	if rxNumber.MatchString(firstChar) {
		return fmt.Errorf("name must not started with number")
	}

	if firstChar == strings.ToUpper(firstChar) {
		return fmt.Errorf("name must be unexported")
	}

	if rxSymbol.MatchString(tagName) {
		return fmt.Errorf("name must be only consisted of letters and numbers")
	}

	return nil
}
