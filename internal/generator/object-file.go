package generator

import (
	"fmt"
	"go/format"
	fp "path/filepath"
	"strings"
)

// createCppHeaderFile creates .h file and save it to obj.dirPath
func createCppHeaderFile(obj object) error {
	// Create initial result
	result := ""

	// Write include guard
	guardName := fmt.Sprintf("QAMEL_%s_H", strings.ToUpper(obj.name))
	result += fmt.Sprintf(""+
		"#pragma once\n\n"+
		"#ifndef %s\n"+
		"#define %s\n\n",
		guardName, guardName)

	// Write std library
	result += "" +
		"#include <stdint.h>\n" +
		"#include <stdbool.h>\n\n"

	// Write check for C++ and class declaration
	className := upperChar(obj.name, 0)
	result += fmt.Sprintf(""+
		"#ifdef __cplusplus\n\n"+
		"// Class\n"+
		"class %s;\n\n"+
		`extern "C" {`+"\n"+
		"#endif\n\n", className)

	// Write getter and setter for properties
	result += fmt.Sprintln("// Properties")
	for i, prop := range obj.properties {
		propName := upperChar(prop.name, 0)
		propType := mapGoType[prop.memberType].inC

		result += fmt.Sprintf(""+
			"%s %s_%s(void* ptr);\n"+
			"void %s_Set%s(void* ptr, %s %s);\n",
			propType, className, propName,
			className, propName, propType, prop.name)

		if i < len(obj.properties)-1 {
			result += "\n"
		}
	}

	// Write method for signals
	result += fmt.Sprintln()
	result += fmt.Sprintln("// Signals")
	for _, signal := range obj.signals {
		params := []string{"void* ptr"}
		for _, param := range signal.parameters {
			paramType := mapGoType[param.memberType].inC
			strParam := fmt.Sprintf("%s %s", paramType, param.name)
			params = append(params, strParam)
		}

		signalName := upperChar(signal.name, 0)
		result += fmt.Sprintf("void %s_%s(%s);\n",
			className, signalName, strings.Join(params, ", "))
	}

	// Write method for registering QML type
	result += fmt.Sprintf("\n"+
		"// Register\n"+
		"void %s_RegisterQML(char* uri, int versionMajor, int versionMinor, char* qmlName);\n\n",
		className)

	// Write #endifs
	result += fmt.Sprintln("" +
		"#ifdef __cplusplus\n" +
		"}\n" +
		"#endif\n\n" +
		"#endif")

	// Save result to file
	fileName := strings.ToLower(obj.name)
	fileName = fmt.Sprintf("qamel-%s.h", fileName)
	fileName = fp.Join(obj.dirPath, fileName)
	if err := saveToFile(fileName, result); err != nil {
		return fmt.Errorf("error creating header file for %s: %v", obj.name, err)
	}

	return nil
}

// createCppFile creates .cpp file and save it to obj.dirPath
func createCppFile(mocPath string, obj object) error {
	// Create initial result
	result := ""

	// Write #include list
	hFileName := strings.ToLower(obj.name)
	hFileName = fmt.Sprintf("qamel-%s.h", hFileName)
	result += fmt.Sprintf(""+
		"#include <QObject>\n"+
		"#include <QQuickItem>\n"+
		"#include <QString>\n"+
		"#include <QByteArray>\n"+
		"#include <QQmlEngine>\n"+
		"#include <QMetaObject>\n\n"+
		`#include "_cgo_export.h"`+"\n"+
		`#include "%s"`+"\n\n",
		hFileName)

	// Write class and property declaration
	className := upperChar(obj.name, 0)
	result += fmt.Sprintf(""+
		"class %s : public QQuickItem {\n"+
		"\tQ_OBJECT\n", className)

	for _, prop := range obj.properties {
		propType := mapGoType[prop.memberType].inCpp
		setterName := "set" + upperChar(prop.name, 0)
		notifierName := prop.name + "Changed"

		result += fmt.Sprintf("\tQ_PROPERTY(%s %s READ %s WRITE %s NOTIFY %s)\n",
			propType, prop.name, prop.name, setterName, notifierName)
	}

	// Write class's private member
	result += fmt.Sprintln("\nprivate:")
	for _, prop := range obj.properties {
		propType := mapGoType[prop.memberType].inCpp
		result += fmt.Sprintf("\t%s _%s;\n", propType, prop.name)
	}

	// Write class's public member
	result += fmt.Sprintln("\npublic:")

	// constructor
	result += fmt.Sprintf(""+
		"\t%s(QQuickItem* parent=Q_NULLPTR) : QQuickItem(parent) {\n"+
		"\t\tqamel%sConstructor(this);\n"+
		"\t}\n\n", className, className)

	// destroyer
	result += fmt.Sprintf(""+
		"\t~%s() {\n"+
		"\t\tqamelDestroy%s(this);\n"+
		"\t}\n\n", className, className)

	// getter and setter
	for i, prop := range obj.properties {
		propType := mapGoType[prop.memberType].inCpp
		setterName := "set" + upperChar(prop.name, 0)
		propNewName := "new" + upperChar(prop.name, 0)

		result += fmt.Sprintf(""+
			"\t%s %s() { return _%s; }\n"+
			"\tvoid %s(%s %s) { _%s = %s; }\n",
			propType, prop.name, prop.name,
			setterName, propType, propNewName, prop.name, propNewName)

		if i < len(obj.properties)-1 {
			result += "\n"
		}
	}

	// Write class's signals
	// properties signals
	result += fmt.Sprintln("signals:")
	for i, prop := range obj.properties {
		propType := mapGoType[prop.memberType].inCpp
		propNewName := "new" + upperChar(prop.name, 0)

		result += fmt.Sprintf("\tvoid %sChanged(%s %s);\n",
			prop.name, propType, propNewName)

		if i < len(obj.properties)-1 {
			result += "\n"
		}
	}

	// the real signals
	for i, signal := range obj.signals {
		var params []string
		for _, param := range signal.parameters {
			paramType := mapGoType[param.memberType].inCpp
			strParam := fmt.Sprintf("%s %s", paramType, param.name)
			params = append(params, strParam)
		}

		result += fmt.Sprintf("\tvoid %s(%s);\n",
			signal.name, strings.Join(params, ", "))

		if i < len(obj.signals)-1 {
			result += "\n"
		}
	}

	// Write class's slots
	result += fmt.Sprintln("\npublic slots:")
	for i, slot := range obj.slots {
		returnType := "void"
		if len(slot.returns) > 0 {
			returnType = mapGoType[slot.returns[0].memberType].inCpp
		}

		var params []string
		paramNames := []string{"this"}
		for _, param := range slot.parameters {
			paramType := mapGoType[param.memberType].inCpp
			paramName := param.name
			if param.memberType == "string" {
				paramName = fmt.Sprintf("%s.toLocal8Bit().data()", paramName)
			}

			strParam := fmt.Sprintf("%s %s", paramType, param.name)
			params = append(params, strParam)
			paramNames = append(paramNames, paramName)
		}

		result += fmt.Sprintf("\t%s %s(%s) {\n",
			returnType, slot.name, strings.Join(params, ", "))

		result += "\t\t"
		if returnType != "void" {
			result += "return "
		}

		result += fmt.Sprintf("qamel%s%s(%s);\n\t}\n",
			className, upperChar(slot.name, 0),
			strings.Join(paramNames, ", "))

		if i < len(obj.slots)-1 {
			result += "\n"
		}
	}

	// Finished writing definition of class
	result += "};\n"

	// Write public methods
	// for manipulating properties
	for i, prop := range obj.properties {
		propName := upperChar(prop.name, 0)
		propType := mapGoType[prop.memberType].inCpp
		propHeaderType := mapGoType[prop.memberType].inC

		// getter
		result += fmt.Sprintf(""+
			"%s %s_%s(void* ptr) {\n"+
			"\t%s *obj = static_cast<%s*>(ptr);\n"+
			"\treturn obj->%s()",
			propHeaderType, className, propName,
			className, className,
			prop.name)
		if prop.memberType == "string" {
			result += ".toLocal8Bit().data();\n"
		} else {
			result += ";\n"
		}
		result += "}\n\n"

		// setter
		result += fmt.Sprintf(""+
			"void %s_Set%s(void* ptr, %s %s) {\n"+
			"\t%s *obj = static_cast<%s*>(ptr);\n"+
			"\tobj->set%s(%s(%s));\n"+
			"}\n",
			className, propName, propHeaderType, prop.name,
			className, className,
			propName, propType, prop.name)

		if i < len(obj.properties)-1 {
			result += "\n"
		}
	}

	// for invoking signals
	result += "\n"
	for i, signal := range obj.signals {
		var invokerParams []string
		params := []string{"void* ptr"}
		for _, param := range signal.parameters {
			paramType := mapGoType[param.memberType].inCpp
			paramHeaderType := mapGoType[param.memberType].inC

			strParam := fmt.Sprintf("%s %s", paramHeaderType, param.name)
			strInvokerParam := fmt.Sprintf("%s(%s)", paramType, param.name)

			params = append(params, strParam)
			invokerParams = append(invokerParams, strInvokerParam)
		}

		signalName := upperChar(signal.name, 0)
		result += fmt.Sprintf(""+
			"void %s_%s(%s) {\n"+
			"\t%s *obj = static_cast<%s*>(ptr);\n"+
			"\tobj->%s(%s);\n"+
			"}\n", className, signalName, strings.Join(params, ", "),
			className, className,
			signal.name, strings.Join(invokerParams, ", "))

		if i < len(obj.signals)-1 {
			result += "\n"
		}
	}

	// for registering QML
	result += fmt.Sprintf("\n"+
		"void %s_RegisterQML(char* uri, int versionMajor, int versionMinor, char* qmlName) {\n"+
		"\tqmlRegisterType<%s>(uri, versionMajor, versionMinor, qmlName);\n"+
		"}\n", className, className)

	// Write #include moc file
	mocFileName := fmt.Sprintf("moc-%s", hFileName)
	result += fmt.Sprintf("\n"+`#include "%s"`+"\n", mocFileName)

	// Save result to file
	cppFileName := strings.ToLower(obj.name)
	cppFileName = fmt.Sprintf("qamel-%s.cpp", cppFileName)
	cppFileName = fp.Join(obj.dirPath, cppFileName)
	if err := saveToFile(cppFileName, result); err != nil {
		return fmt.Errorf("error creating C++ file for %s: %v", obj.name, err)
	}

	// Create moc file
	err := CreateMocFile(mocPath, cppFileName)
	if err != nil {
		return fmt.Errorf("error creating moc file for %s: %v", obj.name, err)
	}

	return nil
}

// createGoFile creates .go file and save it to obj.dirPath
func createGoFile(obj object) error {
	// Create file name
	baseName := strings.ToLower(obj.name)
	hFileName := fmt.Sprintf("qamel-%s.h", baseName)

	// Create initial result
	result := fmt.Sprintf("package %s\n", obj.packageName)
	result += fmt.Sprintln()

	// Write clause for importing C packages
	result += fmt.Sprintf(""+
		"// #include <stdlib.h>\n"+
		"// #include <stdint.h>\n"+
		"// #include <stdbool.h>\n"+
		"// #include <string.h>\n"+
		"// #include \"%s\"\n"+
		`import "C"`+"\n", hFileName)

	// Write clause for importing Go packages
	result += "" +
		"import (\n" +
		`"unsafe"` + "\n" +
		`"github.com/go-qamel/qamel"` + "\n" +
		")\n"

	// Write function for C constructor
	cClassName := upperChar(obj.name, 0)
	result += fmt.Sprintf(""+
		"//export qamel%sConstructor\n"+
		"func qamel%sConstructor(ptr unsafe.Pointer) {\n"+
		"obj := &%s{}\n"+
		"obj.Ptr = ptr\n"+
		"qamel.RegisterObject(ptr, obj)\n",
		cClassName, cClassName, obj.name)

	if len(obj.constructors) == 1 {
		result += fmt.Sprintf("obj.%s()\n", obj.constructors[0].name)
	}

	result += "}\n\n"

	// Write function for C destroyer
	result += fmt.Sprintf(""+
		"//export qamelDestroy%s\n"+
		"func qamelDestroy%s(ptr unsafe.Pointer) {\n"+
		"qamel.DeleteObject(ptr)\n"+
		"}\n\n", cClassName, cClassName)

	// Write function for C slots
	for _, slot := range obj.slots {
		returnType := ""
		cgoReturnType := ""
		if len(slot.returns) > 0 {
			returnType = slot.returns[0].memberType
			cgoReturnType = fmt.Sprintf("(result %s)", mapGoType[returnType].inCgo)
		}

		var castedNames []string
		var castedParams []string
		params := []string{"ptr unsafe.Pointer"}
		for _, param := range slot.parameters {
			cgoType := mapGoType[param.memberType].inCgo
			strParam := fmt.Sprintf("%s %s", param.name, cgoType)
			castedName := fmt.Sprintf("cgo%s", upperChar(param.name, 0))

			params = append(params, strParam)
			castedNames = append(castedNames, castedName)
			castedParams = append(castedParams, fmt.Sprintf("%s := %s",
				castedName, mapGoType[param.memberType].cgo2Go(param.name)))
		}

		slotName := upperChar(slot.name, 0)
		result += fmt.Sprintf(""+
			"//export qamel%s%s\n"+
			"func qamel%s%s(%s) %s {\n"+
			"obj := qamel.BorrowObject(ptr)\n"+
			"defer qamel.ReturnObject(ptr)\n"+
			"if obj == nil {\n"+
			"return\n"+
			"}\n\n"+
			"obj%s, ok := obj.(*%s)\n"+
			"if !ok {\n"+
			"return\n"+
			"}\n\n"+
			"%s\n",
			cClassName, slotName, cClassName, slotName,
			strings.Join(params, ", "), cgoReturnType,
			cClassName, obj.name,
			strings.Join(castedParams, "\n"))

		returnValue := fmt.Sprintf("obj%s.%s(%s)",
			cClassName, slot.name, strings.Join(castedNames, ", "))
		if returnType != "" {
			result += fmt.Sprintf("result = %s\n", mapGoType[returnType].go2C(returnValue))
		} else {
			result += returnValue + "\n"
		}
		result += "return\n}\n\n"
	}

	// Write struct member function
	// for manipulating properties
	result += "// getter and setter\n\n"
	for _, prop := range obj.properties {
		propName := upperChar(prop.name, 0)

		// getter
		result += fmt.Sprintf(""+
			"func (obj *%s) %s() (propValue %s) {\n"+
			"if obj.Ptr == nil || !qamel.ObjectExists(obj.Ptr) {\n"+
			"return\n"+
			"}\n\n"+
			"c%s := C.%s_%s(obj.Ptr)\n",
			obj.name, prop.name, prop.memberType,
			propName, cClassName, propName)

		result += fmt.Sprintf("propValue = %s\nreturn\n}\n\n",
			mapGoType[prop.memberType].cgo2Go("c"+propName))

		// setter
		result += fmt.Sprintf(""+
			"func (obj *%s) set%s(new%s %s) {\n"+
			"if obj.Ptr == nil || !qamel.ObjectExists(obj.Ptr) {\n"+
			"return\n"+
			"}\n\n"+
			"cNew%s := %s\n",
			obj.name, propName, propName, prop.memberType,
			propName, mapGoType[prop.memberType].go2C("new"+propName))

		if prop.memberType == "string" {
			result += fmt.Sprintf("defer C.free(unsafe.Pointer(cNew%s))\n", propName)
		}

		result += fmt.Sprintf("C.%s_Set%s(obj.Ptr, cNew%s)\n}\n\n",
			cClassName, propName, propName)
	}

	// for invoking signals
	result += "// signals invoker\n\n"
	for _, signal := range obj.signals {
		var params []string
		var castedParams []string
		castedNames := []string{"obj.Ptr"}
		for _, param := range signal.parameters {
			strParam := fmt.Sprintf("%s %s", param.name, param.memberType)
			castedName := fmt.Sprintf("c%s", upperChar(param.name, 0))

			params = append(params, strParam)
			castedNames = append(castedNames, castedName)
			castedParams = append(castedParams, fmt.Sprintf("%s := %s",
				castedName, mapGoType[param.memberType].go2C(param.name)))

			if param.memberType == "string" {
				castedParams = append(castedParams,
					fmt.Sprintf("defer C.free(unsafe.Pointer(%s))", castedName))
			}
		}

		result += fmt.Sprintf(""+
			"func (obj *%s) %s(%s) {\n"+
			"if obj.Ptr == nil || !qamel.ObjectExists(obj.Ptr) {\n"+
			"return\n"+
			"}\n\n"+
			"%s\n"+
			"C.%s_%s(%s)}\n\n",
			obj.name, signal.name, strings.Join(params, ", "),
			strings.Join(castedParams, "\n"),
			cClassName, upperChar(signal.name, 0),
			strings.Join(castedNames, ", "))
	}

	// Write function for registering QML object
	result += fmt.Sprintf(""+
		"// RegisterQml%s registers %s as QML object\n"+
		"func RegisterQml%s(uri string, versionMajor int, versionMinor int, qmlName string) {\n"+
		"cURI := C.CString(uri)\n"+
		"cQmlName := C.CString(qmlName)\n"+
		"cVersionMajor := C.int(int32(versionMajor))\n"+
		"cVersionMinor := C.int(int32(versionMinor))\n"+
		"defer func() {\n"+
		"C.free(unsafe.Pointer(cURI))\n"+
		"C.free(unsafe.Pointer(cQmlName))\n"+
		"}()\n\n"+
		"C.%s_RegisterQML(cURI, cVersionMajor, cVersionMinor, cQmlName)\n"+
		"}\n\n", cClassName, obj.name, cClassName, cClassName)

	// Format code
	fmtResult, err := format.Source([]byte(result))
	if err != nil {
		return fmt.Errorf("error formatting Go file for %s: %v", obj.name, err)
	}

	// Save result to file
	fileName := strings.ToLower(obj.name)
	fileName = fmt.Sprintf("qamel-%s.go", fileName)
	fileName = fp.Join(obj.dirPath, fileName)
	if err = saveToFile(fileName, string(fmtResult)); err != nil {
		return fmt.Errorf("error creating Go file for %s: %v", obj.name, err)
	}

	return nil
}
