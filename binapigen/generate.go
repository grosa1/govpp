//  Copyright (c) 2020 Cisco and/or its affiliates.
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at:
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package binapigen

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"git.fd.io/govpp.git/version"
)

// generatedCodeVersion indicates a version of the generated code.
// It is incremented whenever an incompatibility between the generated code and
// GoVPP api package is introduced; the generated code references
// a constant, api.GoVppAPIPackageIsVersionN (where N is generatedCodeVersion).
const generatedCodeVersion = 2

// message field names
const (
	msgIdField       = "_vl_msg_id"
	clientIndexField = "client_index"
	contextField     = "context"
	retvalField      = "retval"
)

const (
	outputFileExt = ".ba.go" // file extension of the Go generated files
	rpcFileSuffix = "_rpc"   // file name suffix for the RPC services

	constModuleName = "ModuleName" // module name constant
	constAPIVersion = "APIVersion" // API version constant
	constVersionCrc = "VersionCrc" // version CRC constant

	unionDataField = "XXX_UnionData" // name for the union data field

	serviceApiName    = "RPCService"    // name for the RPC service interface
	serviceImplName   = "serviceClient" // name for the RPC service implementation
	serviceClientName = "ServiceClient" // name for the RPC service client

	// TODO: register service descriptor
	//serviceDescType = "ServiceDesc"             // name for service descriptor type
	//serviceDescName = "_ServiceRPC_serviceDesc" // name for service descriptor var
)

// MessageType represents the type of a VPP message
type MessageType int

const (
	requestMessage MessageType = iota // VPP request message
	replyMessage                      // VPP reply message
	eventMessage                      // VPP event message
	otherMessage                      // other VPP message
)

type GenFile struct {
	*Generator
	filename   string
	file       *File
	packageDir string
	buf        bytes.Buffer
}

func generatePackage(ctx *GenFile, w io.Writer) {
	logf("----------------------------")
	logf("generating binapi package: %q", ctx.file.PackageName)
	logf("----------------------------")

	generateHeader(ctx, w)
	generateImports(ctx, w)

	// generate module desc
	fmt.Fprintln(w, "const (")
	fmt.Fprintf(w, "\t// %s is the name of this module.\n", constModuleName)
	fmt.Fprintf(w, "\t%s = \"%s\"\n", constModuleName, ctx.file.Name)

	if ctx.IncludeAPIVersion {
		fmt.Fprintf(w, "\t// %s is the API version of this module.\n", constAPIVersion)
		fmt.Fprintf(w, "\t%s = \"%s\"\n", constAPIVersion, ctx.file.Version())
		fmt.Fprintf(w, "\t// %s is the CRC of this module.\n", constVersionCrc)
		fmt.Fprintf(w, "\t%s = %v\n", constVersionCrc, ctx.file.CRC)
	}
	fmt.Fprintln(w, ")")
	fmt.Fprintln(w)

	// generate enums
	if len(ctx.file.Enums) > 0 {
		for _, enum := range ctx.file.Enums {
			if imp, ok := ctx.file.imports[enum.Name]; ok {
				generateImportedAlias(ctx, w, enum.GoName, imp)
				continue
			}
			generateEnum(ctx, w, enum)
		}
	}

	// generate aliases
	if len(ctx.file.Aliases) > 0 {
		for _, alias := range ctx.file.Aliases {
			if imp, ok := ctx.file.imports[alias.Name]; ok {
				generateImportedAlias(ctx, w, alias.GoName, imp)
				continue
			}
			generateAlias(ctx, w, alias)
		}
	}

	// generate types
	if len(ctx.file.Structs) > 0 {
		for _, typ := range ctx.file.Structs {
			if imp, ok := ctx.file.imports[typ.Name]; ok {
				generateImportedAlias(ctx, w, typ.GoName, imp)
				continue
			}
			generateStruct(ctx, w, typ)
		}
	}

	// generate unions
	if len(ctx.file.Unions) > 0 {
		for _, union := range ctx.file.Unions {
			if imp, ok := ctx.file.imports[union.Name]; ok {
				generateImportedAlias(ctx, w, union.GoName, imp)
				continue
			}
			generateUnion(ctx, w, union)
		}
	}

	// generate messages
	if len(ctx.file.Messages) > 0 {
		for _, msg := range ctx.file.Messages {
			generateMessage(ctx, w, msg)
		}

		initFnName := fmt.Sprintf("file_%s_binapi_init", ctx.file.PackageName)

		// generate message registrations
		fmt.Fprintf(w, "func init() { %s() }\n", initFnName)
		fmt.Fprintf(w, "func %s() {\n", initFnName)
		for _, msg := range ctx.file.Messages {
			fmt.Fprintf(w, "\tapi.RegisterMessage((*%s)(nil), \"%s\")\n",
				msg.GoName, ctx.file.Name+"."+msg.GoName)
		}
		fmt.Fprintln(w, "}")
		fmt.Fprintln(w)

		// generate list of messages
		fmt.Fprintf(w, "// Messages returns list of all messages in this module.\n")
		fmt.Fprintln(w, "func AllMessages() []api.Message {")
		fmt.Fprintln(w, "\treturn []api.Message{")
		for _, msg := range ctx.file.Messages {
			fmt.Fprintf(w, "\t(*%s)(nil),\n", msg.GoName)
		}
		fmt.Fprintln(w, "}")
		fmt.Fprintln(w, "}")
	}

	generateFooter(ctx, w)

}

func generateHeader(ctx *GenFile, w io.Writer) {
	fmt.Fprintln(w, "// Code generated by GoVPP's binapi-generator. DO NOT EDIT.")
	fmt.Fprintln(w, "// versions:")
	fmt.Fprintf(w, "//  binapi-generator: %s\n", version.Version())
	if ctx.IncludeVppVersion {
		fmt.Fprintf(w, "//  VPP:              %s\n", ctx.VPPVersion)
	}
	fmt.Fprintf(w, "// source: %s\n", ctx.file.Path)
	fmt.Fprintln(w)

	fmt.Fprintln(w, "/*")
	fmt.Fprintf(w, "Package %s contains generated code for VPP binary API defined by %s.api (version %s).\n",
		ctx.file.PackageName, ctx.file.Name, ctx.file.Version())
	fmt.Fprintln(w)
	fmt.Fprintln(w, "It consists of:")
	printObjNum := func(obj string, num int) {
		if num > 0 {
			if num > 1 {
				if strings.HasSuffix(obj, "s") {

					obj += "es"
				} else {
					obj += "s"
				}
			}
			fmt.Fprintf(w, "\t%3d %s\n", num, obj)
		}
	}
	//printObjNum("RPC", len(ctx.file.Service.RPCs))
	printObjNum("alias", len(ctx.file.Aliases))
	printObjNum("enum", len(ctx.file.Enums))
	printObjNum("message", len(ctx.file.Messages))
	printObjNum("type", len(ctx.file.Structs))
	printObjNum("union", len(ctx.file.Unions))
	fmt.Fprintln(w, "*/")
	fmt.Fprintf(w, "package %s\n", ctx.file.PackageName)
	fmt.Fprintln(w)
}

func generateImports(ctx *GenFile, w io.Writer) {
	fmt.Fprintln(w, "import (")
	fmt.Fprintln(w, `	"bytes"`)
	fmt.Fprintln(w, `	"context"`)
	fmt.Fprintln(w, `	"encoding/binary"`)
	fmt.Fprintln(w, `	"io"`)
	fmt.Fprintln(w, `	"math"`)
	fmt.Fprintln(w, `	"strconv"`)
	fmt.Fprintln(w)
	fmt.Fprintf(w, "\tapi \"%s\"\n", "git.fd.io/govpp.git/api")
	fmt.Fprintf(w, "\tcodec \"%s\"\n", "git.fd.io/govpp.git/codec")
	fmt.Fprintf(w, "\tstruc \"%s\"\n", "github.com/lunixbochs/struc")
	if len(ctx.file.Imports) > 0 {
		fmt.Fprintln(w)
		for _, imp := range ctx.file.Imports {
			importPath := path.Join(ctx.ImportPrefix, imp)
			if ctx.ImportPrefix == "" {
				importPath = getImportPath(ctx.packageDir, imp)
			}
			fmt.Fprintf(w, "\t%s \"%s\"\n", imp, strings.TrimSpace(importPath))
		}
	}
	fmt.Fprintln(w, ")")
	fmt.Fprintln(w)

	fmt.Fprintln(w, "// This is a compile-time assertion to ensure that this generated file")
	fmt.Fprintln(w, "// is compatible with the GoVPP api package it is being compiled against.")
	fmt.Fprintln(w, "// A compilation error at this line likely means your copy of the")
	fmt.Fprintln(w, "// GoVPP api package needs to be updated.")
	fmt.Fprintf(w, "const _ = api.GoVppAPIPackageIsVersion%d // please upgrade the GoVPP api package\n", generatedCodeVersion)
	fmt.Fprintln(w)
}

func getImportPath(outputDir string, pkg string) string {
	absPath, err := filepath.Abs(filepath.Join(outputDir, "..", pkg))
	if err != nil {
		panic(err)
	}
	cmd := exec.Command("go", "list", absPath)
	var errbuf, outbuf bytes.Buffer
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf
	if err := cmd.Run(); err != nil {
		fmt.Printf("ERR: %v\n", errbuf.String())
		panic(err)
	}
	return outbuf.String()
}

func generateFooter(ctx *GenFile, w io.Writer) {
	fmt.Fprintf(w, "// Reference imports to suppress errors if they are not otherwise used.\n")
	fmt.Fprintf(w, "var _ = api.RegisterMessage\n")
	fmt.Fprintf(w, "var _ = codec.DecodeString\n")
	fmt.Fprintf(w, "var _ = bytes.NewBuffer\n")
	fmt.Fprintf(w, "var _ = context.Background\n")
	fmt.Fprintf(w, "var _ = io.Copy\n")
	fmt.Fprintf(w, "var _ = strconv.Itoa\n")
	fmt.Fprintf(w, "var _ = struc.Pack\n")
	fmt.Fprintf(w, "var _ = binary.BigEndian\n")
	fmt.Fprintf(w, "var _ = math.Float32bits\n")
}

func generateComment(ctx *GenFile, w io.Writer, goName string, vppName string, objKind string) {
	if objKind == "service" {
		fmt.Fprintf(w, "// %s represents RPC service API for %s module.\n", goName, ctx.file.Name)
	} else {
		fmt.Fprintf(w, "// %s represents VPP binary API %s '%s'.\n", goName, objKind, vppName)
	}
}

func generateEnum(ctx *GenFile, w io.Writer, enum *Enum) {
	name := enum.GoName
	typ := binapiTypes[enum.Type]

	logf(" writing ENUM %q (%s) with %d entries", enum.Name, name, len(enum.Entries))

	// generate enum comment
	generateComment(ctx, w, name, enum.Name, "enum")

	// generate enum definition
	fmt.Fprintf(w, "type %s %s\n", name, typ)
	fmt.Fprintln(w)

	// generate enum entries
	fmt.Fprintln(w, "const (")
	for _, entry := range enum.Entries {
		fmt.Fprintf(w, "\t%s %s = %v\n", entry.Name, name, entry.Value)
	}
	fmt.Fprintln(w, ")")
	fmt.Fprintln(w)

	// generate enum conversion maps
	fmt.Fprintln(w, "var (")
	fmt.Fprintf(w, "%s_name = map[%s]string{\n", name, typ)
	for _, entry := range enum.Entries {
		fmt.Fprintf(w, "\t%v: \"%s\",\n", entry.Value, entry.Name)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintf(w, "%s_value = map[string]%s{\n", name, typ)
	for _, entry := range enum.Entries {
		fmt.Fprintf(w, "\t\"%s\": %v,\n", entry.Name, entry.Value)
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w, ")")
	fmt.Fprintln(w)

	fmt.Fprintf(w, "func (x %s) String() string {\n", name)
	fmt.Fprintf(w, "\ts, ok := %s_name[%s(x)]\n", name, typ)
	fmt.Fprintf(w, "\tif ok { return s }\n")
	fmt.Fprintf(w, "\treturn \"%s(\" + strconv.Itoa(int(x)) + \")\"\n", name)
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}

func generateImportedAlias(ctx *GenFile, w io.Writer, name string, imp string) {
	fmt.Fprintf(w, "type %s = %s.%s\n", name, imp, name)
	fmt.Fprintln(w)
}

func generateAlias(ctx *GenFile, w io.Writer, alias *Alias) {
	name := alias.GoName

	logf(" writing ALIAS %q (%s), length: %d", alias.Name, name, alias.Length)

	// generate struct comment
	generateComment(ctx, w, name, alias.Name, "alias")

	// generate struct definition
	fmt.Fprintf(w, "type %s ", name)

	if alias.Length > 0 {
		fmt.Fprintf(w, "[%d]", alias.Length)
	}

	dataType := convertToGoType(ctx.file, alias.Type)
	fmt.Fprintf(w, "%s\n", dataType)

	fmt.Fprintln(w)
}

func generateStruct(ctx *GenFile, w io.Writer, typ *Struct) {
	name := typ.GoName

	logf(" writing STRUCT %q (%s) with %d fields", typ.Name, name, len(typ.Fields))

	// generate struct comment
	generateComment(ctx, w, name, typ.Name, "type")

	// generate struct definition
	fmt.Fprintf(w, "type %s struct {\n", name)

	// generate struct fields
	for i := range typ.Fields {
		// skip internal fields
		switch strings.ToLower(typ.Name) {
		case msgIdField:
			continue
		}

		generateField(ctx, w, typ.Fields, i)
	}

	// generate end of the struct
	fmt.Fprintln(w, "}")

	// generate name getter
	generateTypeNameGetter(w, name, typ.Name)

	fmt.Fprintln(w)
}

// generateUnionMethods generates methods that implement struc.Custom
// interface to allow having XXX_uniondata field unexported
// TODO: do more testing when unions are actually used in some messages
/*func generateUnionMethods(w io.Writer, structName string) {
	// generate struc.Custom implementation for union
	fmt.Fprintf(w, `
func (u *%[1]s) Pack(p []byte, opt *struc.Options) (int, error) {
	var b = new(bytes.Buffer)
	if err := struc.PackWithOptions(b, u.union_data, opt); err != nil {
		return 0, err
	}
	copy(p, b.Bytes())
	return b.Len(), nil
}
func (u *%[1]s) Unpack(r io.Reader, length int, opt *struc.Options) error {
	return struc.UnpackWithOptions(r, u.union_data[:], opt)
}
func (u *%[1]s) Size(opt *struc.Options) int {
	return len(u.union_data)
}
func (u *%[1]s) String() string {
	return string(u.union_data[:])
}
`, structName)
}*/

/*func generateUnionGetterSetterNew(w io.Writer, structName string, getterField, getterStruct string) {
	fmt.Fprintf(w, `
func %[1]s%[2]s(a %[3]s) (u %[1]s) {
	u.Set%[2]s(a)
	return
}
func (u *%[1]s) Set%[2]s(a %[3]s) {
	copy(u.%[4]s[:], a[:])
}
func (u *%[1]s) Get%[2]s() (a %[3]s) {
	copy(a[:], u.%[4]s[:])
	return
}
`, structName, getterField, getterStruct, unionDataField)
}*/

func generateUnion(ctx *GenFile, w io.Writer, union *Union) {
	name := union.GoName

	logf(" writing UNION %q (%s) with %d fields", union.Name, name, len(union.Fields))

	// generate struct comment
	generateComment(ctx, w, name, union.Name, "union")

	// generate struct definition
	fmt.Fprintln(w, "type", name, "struct {")

	// maximum size for union
	maxSize := getUnionSize(ctx.file, union)

	// generate data field
	fmt.Fprintf(w, "\t%s [%d]byte\n", unionDataField, maxSize)

	// generate end of the struct
	fmt.Fprintln(w, "}")

	// generate name getter
	generateTypeNameGetter(w, name, union.Name)

	// generate getters for fields
	for _, field := range union.Fields {
		fieldType := convertToGoType(ctx.file, field.Type)
		generateUnionGetterSetter(w, name, field.GoName, fieldType)
	}

	// generate union methods
	//generateUnionMethods(w, name)

	fmt.Fprintln(w)
}

func generateUnionGetterSetter(w io.Writer, structName string, getterField, getterStruct string) {
	fmt.Fprintf(w, `
func %[1]s%[2]s(a %[3]s) (u %[1]s) {
	u.Set%[2]s(a)
	return
}
func (u *%[1]s) Set%[2]s(a %[3]s) {
	var b = new(bytes.Buffer)
	if err := struc.Pack(b, &a); err != nil {
		return
	}
	copy(u.%[4]s[:], b.Bytes())
}
func (u *%[1]s) Get%[2]s() (a %[3]s) {
	var b = bytes.NewReader(u.%[4]s[:])
	struc.Unpack(b, &a)
	return
}
`, structName, getterField, getterStruct, unionDataField)
}

func generateMessage(ctx *GenFile, w io.Writer, msg *Message) {
	name := msg.GoName

	logf(" writing MESSAGE %q (%s) with %d fields", msg.Name, name, len(msg.Fields))

	// generate struct comment
	generateComment(ctx, w, name, msg.Name, "message")

	// generate struct definition
	fmt.Fprintf(w, "type %s struct {", name)

	msgType := otherMessage
	wasClientIndex := false

	// generate struct fields
	n := 0
	for i, field := range msg.Fields {
		if i == 1 {
			if field.Name == clientIndexField {
				// "client_index" as the second member,
				// this might be an event message or a request
				msgType = eventMessage
				wasClientIndex = true
			} else if field.Name == contextField {
				// reply needs "context" as the second member
				msgType = replyMessage
			}
		} else if i == 2 {
			if wasClientIndex && field.Name == contextField {
				// request needs "client_index" as the second member
				// and "context" as the third member
				msgType = requestMessage
			}
		}

		// skip internal fields
		switch strings.ToLower(field.Name) {
		case /*crcField,*/ msgIdField:
			continue
		case clientIndexField, contextField:
			if n == 0 {
				continue
			}
		}
		n++
		if n == 1 {
			fmt.Fprintln(w)
		}

		generateField(ctx, w, msg.Fields, i)
	}

	// generate end of the struct
	fmt.Fprintln(w, "}")

	// generate message methods
	generateMessageResetMethod(w, name)
	generateMessageNameGetter(w, name, msg.Name)
	generateCrcGetter(w, name, msg.CRC)
	generateMessageTypeGetter(w, name, msgType)
	generateMessageSize(ctx, w, name, msg.Fields)
	generateMessageMarshal(ctx, w, name, msg.Fields)
	generateMessageUnmarshal(ctx, w, name, msg.Fields)

	fmt.Fprintln(w)
}

func generateMessageSize(ctx *GenFile, w io.Writer, name string, fields []*Field) {
	fmt.Fprintf(w, "func (m *%[1]s) Size() int {\n", name)

	fmt.Fprintf(w, "\tif m == nil { return 0 }\n")
	fmt.Fprintf(w, "\tvar size int\n")

	encodeBaseType := func(typ, name string, length int, sizefrom string) bool {
		t, ok := BaseTypeNames[typ]
		if !ok {
			return false
		}

		var s = BaseTypeSizes[t]
		switch t {
		case STRING:
			if length > 0 {
				s = length
				fmt.Fprintf(w, "\tsize += %d\n", s)
			} else {
				s = 4
				fmt.Fprintf(w, "\tsize += %d + len(%s)\n", s, name)
			}
		default:
			if sizefrom != "" {
				//fmt.Fprintf(w, "\tsize += %d * int(%s)\n", s, sizefrom)
				fmt.Fprintf(w, "\tsize += %d * len(%s)\n", s, name)
			} else {
				if length > 0 {
					s = BaseTypeSizes[t] * length
				}
				fmt.Fprintf(w, "\tsize += %d\n", s)
			}
		}

		return true
	}

	lvl := 0
	var encodeFields func(fields []*Field, parentName string)
	encodeFields = func(fields []*Field, parentName string) {
		lvl++
		defer func() { lvl-- }()

		n := 0
		for _, field := range fields {
			// skip internal fields
			switch strings.ToLower(field.Name) {
			case /*crcField,*/ msgIdField:
				continue
			case clientIndexField, contextField:
				if n == 0 {
					continue
				}
			}
			n++

			fieldName := field.GoName //camelCaseName(strings.TrimPrefix(field.Name, "_"))
			name := fmt.Sprintf("%s.%s", parentName, fieldName)
			sizeFrom := camelCaseName(strings.TrimPrefix(field.SizeFrom, "_"))
			var sizeFromName string
			if sizeFrom != "" {
				sizeFromName = fmt.Sprintf("%s.%s", parentName, sizeFrom)
			}

			fmt.Fprintf(w, "\t// field[%d] %s\n", lvl, name)

			if encodeBaseType(field.Type, name, field.Length, sizeFromName) {
				continue
			}

			char := fmt.Sprintf("s%d", lvl)
			index := fmt.Sprintf("j%d", lvl)

			if field.Array {
				if field.Length > 0 {
					fmt.Fprintf(w, "\tfor %[2]s := 0; %[2]s < %[1]d; %[2]s ++ {\n", field.Length, index)
				} else if field.SizeFrom != "" {
					//fmt.Fprintf(w, "\tfor %[1]s := 0; %[1]s < int(%[2]s.%[3]s); %[1]s++ {\n", index, parentName, sizeFrom)
					fmt.Fprintf(w, "\tfor %[1]s := 0; %[1]s < len(%[2]s); %[1]s++ {\n", index, name)
				}

				fmt.Fprintf(w, "\tvar %[1]s %[2]s\n_ = %[1]s\n", char, convertToGoType(ctx.file, field.Type))
				fmt.Fprintf(w, "\tif %[1]s < len(%[2]s) { %[3]s = %[2]s[%[1]s] }\n", index, name, char)
				name = char
			}

			if enum := getEnumByRef(ctx.file, field.Type); enum != nil {
				if encodeBaseType(enum.Type, name, 0, "") {
				} else {
					fmt.Fprintf(w, "\t// ??? ENUM %s %s\n", name, enum.Type)
				}
			} else if alias := getAliasByRef(ctx.file, field.Type); alias != nil {
				if encodeBaseType(alias.Type, name, alias.Length, "") {
				} else if typ := getTypeByRef(ctx.file, alias.Type); typ != nil {
					encodeFields(typ.Fields, name)
				} else {
					fmt.Fprintf(w, "\t// ??? ALIAS %s %s\n", name, alias.Type)
				}
			} else if typ := getTypeByRef(ctx.file, field.Type); typ != nil {
				encodeFields(typ.Fields, name)
			} else if union := getUnionByRef(ctx.file, field.Type); union != nil {
				maxSize := getUnionSize(ctx.file, union)
				fmt.Fprintf(w, "\tsize += %d\n", maxSize)
			} else {
				fmt.Fprintf(w, "\t// ??? buf[pos] = (%s)\n", name)
			}

			if field.Array {
				fmt.Fprintf(w, "\t}\n")
			}
		}
	}

	encodeFields(fields, "m")

	fmt.Fprintf(w, "return size\n")

	fmt.Fprintf(w, "}\n")
}

func generateMessageMarshal(ctx *GenFile, w io.Writer, name string, fields []*Field) {
	fmt.Fprintf(w, "func (m *%[1]s) Marshal(b []byte) ([]byte, error) {\n", name)

	fmt.Fprintf(w, "\to := binary.BigEndian\n")
	fmt.Fprintf(w, "\t_ = o\n")
	fmt.Fprintf(w, "\tpos := 0\n")
	fmt.Fprintf(w, "\t_ = pos\n")

	var buf = new(strings.Builder)

	encodeBaseType := func(typ, name string, length int, sizefrom string) bool {
		t, ok := BaseTypeNames[typ]
		if !ok {
			return false
		}

		isArray := length > 0 || sizefrom != ""

		switch t {
		case I8, U8, I16, U16, I32, U32, I64, U64, F64:
			if isArray {
				if length != 0 {
					fmt.Fprintf(buf, "\tfor i := 0; i < %d; i++ {\n", length)
				} else if sizefrom != "" {
					//fmt.Fprintf(buf, "\tfor i := 0; i < int(%s); i++ {\n", sizefrom)
					fmt.Fprintf(buf, "\tfor i := 0; i < len(%s); i++ {\n", name)
				}
			}
		}

		switch t {
		case I8, U8:
			if isArray {
				fmt.Fprintf(buf, "\tvar x uint8\n")
				fmt.Fprintf(buf, "\tif i < len(%s) { x = uint8(%s[i]) }\n", name, name)
				name = "x"
			}
			fmt.Fprintf(buf, "\tbuf[pos] = uint8(%s)\n", name)
			fmt.Fprintf(buf, "\tpos += 1\n")
			if isArray {
				fmt.Fprintf(buf, "\t}\n")
			}
		case I16, U16:
			if isArray {
				fmt.Fprintf(buf, "\tvar x uint16\n")
				fmt.Fprintf(buf, "\tif i < len(%s) { x = uint16(%s[i]) }\n", name, name)
				name = "x"
			}
			fmt.Fprintf(buf, "\to.PutUint16(buf[pos:pos+2], uint16(%s))\n", name)
			fmt.Fprintf(buf, "\tpos += 2\n")
			if isArray {
				fmt.Fprintf(buf, "\t}\n")
			}
		case I32, U32:
			if isArray {
				fmt.Fprintf(buf, "\tvar x uint32\n")
				fmt.Fprintf(buf, "\tif i < len(%s) { x = uint32(%s[i]) }\n", name, name)
				name = "x"
			}
			fmt.Fprintf(buf, "\to.PutUint32(buf[pos:pos+4], uint32(%s))\n", name)
			fmt.Fprintf(buf, "\tpos += 4\n")
			if isArray {
				fmt.Fprintf(buf, "\t}\n")
			}
		case I64, U64:
			if isArray {
				fmt.Fprintf(buf, "\tvar x uint64\n")
				fmt.Fprintf(buf, "\tif i < len(%s) { x = uint64(%s[i]) }\n", name, name)
				name = "x"
			}
			fmt.Fprintf(buf, "\to.PutUint64(buf[pos:pos+8], uint64(%s))\n", name)
			fmt.Fprintf(buf, "\tpos += 8\n")
			if isArray {
				fmt.Fprintf(buf, "\t}\n")
			}
		case F64:
			if isArray {
				fmt.Fprintf(buf, "\tvar x float64\n")
				fmt.Fprintf(buf, "\tif i < len(%s) { x = float64(%s[i]) }\n", name, name)
				name = "x"
			}
			fmt.Fprintf(buf, "\to.PutUint64(buf[pos:pos+8], math.Float64bits(float64(%s)))\n", name)
			fmt.Fprintf(buf, "\tpos += 8\n")
			if isArray {
				fmt.Fprintf(buf, "\t}\n")
			}
		case BOOL:
			fmt.Fprintf(buf, "\tif %s { buf[pos] = 1 }\n", name)
			fmt.Fprintf(buf, "\tpos += 1\n")
		case STRING:
			if length != 0 {
				fmt.Fprintf(buf, "\tcopy(buf[pos:pos+%d], %s)\n", length, name)
				fmt.Fprintf(buf, "\tpos += %d\n", length)
			} else {
				fmt.Fprintf(buf, "\to.PutUint32(buf[pos:pos+4], uint32(len(%s)))\n", name)
				fmt.Fprintf(buf, "\tpos += 4\n")
				fmt.Fprintf(buf, "\tcopy(buf[pos:pos+len(%s)], %s[:])\n", name, name)
				fmt.Fprintf(buf, "\tpos += len(%s)\n", name)
			}
		default:
			fmt.Fprintf(buf, "\t// ??? %s %s\n", name, typ)
			return false
		}
		return true
	}

	lvl := 0
	var encodeFields func(fields []*Field, parentName string)
	encodeFields = func(fields []*Field, parentName string) {
		lvl++
		defer func() { lvl-- }()

		n := 0
		for _, field := range fields {
			// skip internal fields
			switch strings.ToLower(field.Name) {
			case /*crcField,*/ msgIdField:
				continue
			case clientIndexField, contextField:
				if n == 0 {
					continue
				}
			}
			n++

			getFieldName := func(name string) string {
				fieldName := camelCaseName(strings.TrimPrefix(name, "_"))
				return fmt.Sprintf("%s.%s", parentName, fieldName)
			}

			fieldName := camelCaseName(strings.TrimPrefix(field.Name, "_"))
			name := fmt.Sprintf("%s.%s", parentName, fieldName)
			sizeFrom := camelCaseName(strings.TrimPrefix(field.SizeFrom, "_"))
			var sizeFromName string
			if sizeFrom != "" {
				sizeFromName = fmt.Sprintf("%s.%s", parentName, sizeFrom)
			}

			fmt.Fprintf(buf, "\t// field[%d] %s\n", lvl, name)

			getSizeOfField := func() *Field {
				for _, f := range fields {
					if f.SizeFrom == field.Name {
						return f
					}
				}
				return nil
			}
			if f := getSizeOfField(); f != nil {
				if encodeBaseType(field.Type, fmt.Sprintf("len(%s)", getFieldName(f.Name)), field.Length, "") {
					continue
				}
				panic(fmt.Sprintf("failed to encode base type of sizefrom field: %s", field.Name))
			}

			if encodeBaseType(field.Type, name, field.Length, sizeFromName) {
				continue
			}

			char := fmt.Sprintf("v%d", lvl)
			index := fmt.Sprintf("j%d", lvl)

			if field.Array {
				if field.Length > 0 {
					fmt.Fprintf(buf, "\tfor %[2]s := 0; %[2]s < %[1]d; %[2]s ++ {\n", field.Length, index)
				} else if field.SizeFrom != "" {
					//fmt.Fprintf(buf, "\tfor %[1]s := 0; %[1]s < int(%[2]s.%[3]s); %[1]s++ {\n", index, parentName, sizeFrom)
					fmt.Fprintf(buf, "\tfor %[1]s := 0; %[1]s < len(%[2]s); %[1]s++ {\n", index, name)
				}

				fmt.Fprintf(buf, "\tvar %s %s\n", char, convertToGoType(ctx.file, field.Type))
				fmt.Fprintf(buf, "\tif %[1]s < len(%[2]s) { %[3]s = %[2]s[%[1]s] }\n", index, name, char)
				name = char
			}

			if enum := getEnumByRef(ctx.file, field.Type); enum != nil {
				if encodeBaseType(enum.Type, name, 0, "") {
				} else {
					fmt.Fprintf(buf, "\t// ??? ENUM %s %s\n", name, enum.Type)
				}
			} else if alias := getAliasByRef(ctx.file, field.Type); alias != nil {
				if encodeBaseType(alias.Type, name, alias.Length, "") {
				} else if typ := getTypeByRef(ctx.file, alias.Type); typ != nil {
					encodeFields(typ.Fields, name)
				} else {
					fmt.Fprintf(buf, "\t// ??? ALIAS %s %s\n", name, alias.Type)
				}
			} else if typ := getTypeByRef(ctx.file, field.Type); typ != nil {
				encodeFields(typ.Fields, name)
			} else if union := getUnionByRef(ctx.file, field.Type); union != nil {
				maxSize := getUnionSize(ctx.file, union)
				fmt.Fprintf(buf, "\tcopy(buf[pos:pos+%d], %s.%s[:])\n", maxSize, name, unionDataField)
				fmt.Fprintf(buf, "\tpos += %d\n", maxSize)
			} else {
				fmt.Fprintf(buf, "\t// ??? buf[pos] = (%s)\n", name)
			}

			if field.Array {
				fmt.Fprintf(buf, "\t}\n")
			}
		}
	}

	encodeFields(fields, "m")

	fmt.Fprintf(w, "\tvar buf []byte\n")
	fmt.Fprintf(w, "\tif b == nil {\n")
	fmt.Fprintf(w, "\tbuf = make([]byte, m.Size())\n")
	fmt.Fprintf(w, "\t} else {\n")
	fmt.Fprintf(w, "\tbuf = b\n")
	fmt.Fprintf(w, "\t}\n")
	fmt.Fprint(w, buf.String())

	fmt.Fprintf(w, "return buf, nil\n")

	fmt.Fprintf(w, "}\n")
}

func generateMessageUnmarshal(ctx *GenFile, w io.Writer, name string, fields []*Field) {
	fmt.Fprintf(w, "func (m *%[1]s) Unmarshal(tmp []byte) error {\n", name)

	fmt.Fprintf(w, "\to := binary.BigEndian\n")
	fmt.Fprintf(w, "\t_ = o\n")
	fmt.Fprintf(w, "\tpos := 0\n")
	fmt.Fprintf(w, "\t_ = pos\n")

	decodeBaseType := func(typ, orig, name string, length int, sizefrom string, alloc bool) bool {
		t, ok := BaseTypeNames[typ]
		if !ok {
			return false
		}

		isArray := length > 0 || sizefrom != ""

		switch t {
		case I8, U8, I16, U16, I32, U32, I64, U64, F64:
			if isArray {
				if alloc {
					if length != 0 {
						fmt.Fprintf(w, "\t%s = make([]%s, %d)\n", name, orig, length)
					} else if sizefrom != "" {
						fmt.Fprintf(w, "\t%s = make([]%s, %s)\n", name, orig, sizefrom)
					}
				}
				fmt.Fprintf(w, "\tfor i := 0; i < len(%s); i++ {\n", name)
			}
		}

		switch t {
		case I8, U8:
			if isArray {
				fmt.Fprintf(w, "\t%s[i] = %s(tmp[pos])\n", name, convertToGoType(ctx.file, typ))
			} else {
				fmt.Fprintf(w, "\t%s = %s(tmp[pos])\n", name, orig)
			}
			fmt.Fprintf(w, "\tpos += 1\n")
			if isArray {
				fmt.Fprintf(w, "\t}\n")
			}
		case I16, U16:
			if isArray {
				fmt.Fprintf(w, "\t%s[i] = %s(o.Uint16(tmp[pos:pos+2]))\n", name, orig)
			} else {
				fmt.Fprintf(w, "\t%s = %s(o.Uint16(tmp[pos:pos+2]))\n", name, orig)
			}
			fmt.Fprintf(w, "\tpos += 2\n")
			if isArray {
				fmt.Fprintf(w, "\t}\n")
			}
		case I32, U32:
			if isArray {
				fmt.Fprintf(w, "\t%s[i] = %s(o.Uint32(tmp[pos:pos+4]))\n", name, orig)
			} else {
				fmt.Fprintf(w, "\t%s = %s(o.Uint32(tmp[pos:pos+4]))\n", name, orig)
			}
			fmt.Fprintf(w, "\tpos += 4\n")
			if isArray {
				fmt.Fprintf(w, "\t}\n")
			}
		case I64, U64:
			if isArray {
				fmt.Fprintf(w, "\t%s[i] = %s(o.Uint64(tmp[pos:pos+8]))\n", name, orig)
			} else {
				fmt.Fprintf(w, "\t%s = %s(o.Uint64(tmp[pos:pos+8]))\n", name, orig)
			}
			fmt.Fprintf(w, "\tpos += 8\n")
			if isArray {
				fmt.Fprintf(w, "\t}\n")
			}
		case F64:
			if isArray {
				fmt.Fprintf(w, "\t%s[i] = %s(math.Float64frombits(o.Uint64(tmp[pos:pos+8])))\n", name, orig)
			} else {
				fmt.Fprintf(w, "\t%s = %s(math.Float64frombits(o.Uint64(tmp[pos:pos+8])))\n", name, orig)
			}
			fmt.Fprintf(w, "\tpos += 8\n")
			if isArray {
				fmt.Fprintf(w, "\t}\n")
			}
		case BOOL:
			fmt.Fprintf(w, "\t%s = tmp[pos] != 0\n", name)
			fmt.Fprintf(w, "\tpos += 1\n")
		case STRING:
			if length != 0 {
				fmt.Fprintf(w, "\t{\n")
				fmt.Fprintf(w, "\tnul := bytes.Index(tmp[pos:pos+%d], []byte{0x00})\n", length)
				fmt.Fprintf(w, "\t%[1]s = codec.DecodeString(tmp[pos:pos+nul])\n", name)
				fmt.Fprintf(w, "\tpos += %d\n", length)
				fmt.Fprintf(w, "\t}\n")
			} else {
				fmt.Fprintf(w, "\t{\n")
				fmt.Fprintf(w, "\tsiz := o.Uint32(tmp[pos:pos+4])\n")
				fmt.Fprintf(w, "\tpos += 4\n")
				fmt.Fprintf(w, "\t%[1]s = codec.DecodeString(tmp[pos:pos+int(siz)])\n", name)
				fmt.Fprintf(w, "\tpos += len(%s)\n", name)
				fmt.Fprintf(w, "\t}\n")
			}
		default:
			fmt.Fprintf(w, "\t// ??? %s %s\n", name, typ)
			return false
		}
		return true
	}

	lvl := 0
	var decodeFields func(fields []*Field, parentName string)
	decodeFields = func(fields []*Field, parentName string) {
		lvl++
		defer func() { lvl-- }()

		n := 0
		for _, field := range fields {
			// skip internal fields
			switch strings.ToLower(field.Name) {
			case /*crcField,*/ msgIdField:
				continue
			case clientIndexField, contextField:
				if n == 0 {
					continue
				}
			}
			n++

			fieldName := camelCaseName(strings.TrimPrefix(field.Name, "_"))
			name := fmt.Sprintf("%s.%s", parentName, fieldName)
			sizeFrom := camelCaseName(strings.TrimPrefix(field.SizeFrom, "_"))
			var sizeFromName string
			if sizeFrom != "" {
				sizeFromName = fmt.Sprintf("%s.%s", parentName, sizeFrom)
			}

			fmt.Fprintf(w, "\t// field[%d] %s\n", lvl, name)

			if decodeBaseType(field.Type, convertToGoType(ctx.file, field.Type), name, field.Length, sizeFromName, true) {
				continue
			}

			//char := fmt.Sprintf("v%d", lvl)
			index := fmt.Sprintf("j%d", lvl)

			if field.Array {
				if field.Length > 0 {
					fmt.Fprintf(w, "\tfor %[2]s := 0; %[2]s < %[1]d; %[2]s ++ {\n", field.Length, index)
				} else if field.SizeFrom != "" {
					fieldType := getFieldType(ctx, field)
					if strings.HasPrefix(fieldType, "[]") {
						fmt.Fprintf(w, "\t%s = make(%s, int(%s.%s))\n", name, fieldType, parentName, sizeFrom)
					}
					fmt.Fprintf(w, "\tfor %[1]s := 0; %[1]s < int(%[2]s.%[3]s); %[1]s++ {\n", index, parentName, sizeFrom)
				}

				/*fmt.Fprintf(w, "\tvar %s %s\n", char, convertToGoType(ctx, field.Type))
				fmt.Fprintf(w, "\tif %[1]s < len(%[2]s) { %[3]s = %[2]s[%[1]s] }\n", index, name, char)
				name = char*/
				name = fmt.Sprintf("%s[%s]", name, index)
			}

			if enum := getEnumByRef(ctx.file, field.Type); enum != nil {
				if decodeBaseType(enum.Type, convertToGoType(ctx.file, field.Type), name, 0, "", false) {
				} else {
					fmt.Fprintf(w, "\t// ??? ENUM %s %s\n", name, enum.Type)
				}
			} else if alias := getAliasByRef(ctx.file, field.Type); alias != nil {
				if decodeBaseType(alias.Type, convertToGoType(ctx.file, field.Type), name, alias.Length, "", false) {
				} else if typ := getTypeByRef(ctx.file, alias.Type); typ != nil {
					decodeFields(typ.Fields, name)
				} else {
					fmt.Fprintf(w, "\t// ??? ALIAS %s %s\n", name, alias.Type)
				}
			} else if typ := getTypeByRef(ctx.file, field.Type); typ != nil {
				decodeFields(typ.Fields, name)
			} else if union := getUnionByRef(ctx.file, field.Type); union != nil {
				maxSize := getUnionSize(ctx.file, union)
				fmt.Fprintf(w, "\tcopy(%s.%s[:], tmp[pos:pos+%d])\n", name, unionDataField, maxSize)
				fmt.Fprintf(w, "\tpos += %d\n", maxSize)
			} else {
				fmt.Fprintf(w, "\t// ??? buf[pos] = (%s)\n", name)
			}

			if field.Array {
				fmt.Fprintf(w, "\t}\n")
			}
		}
	}

	decodeFields(fields, "m")

	fmt.Fprintf(w, "return nil\n")

	fmt.Fprintf(w, "}\n")
}

func getFieldType(ctx *GenFile, field *Field) string {
	//fieldName := strings.TrimPrefix(field.Name, "_")
	//fieldName = camelCaseName(fieldName)
	//fieldName := field.GoName

	dataType := convertToGoType(ctx.file, field.Type)
	fieldType := dataType

	// check if it is array
	if field.Length > 0 || field.SizeFrom != "" {
		if dataType == "uint8" {
			dataType = "byte"
		}
		if dataType == "string" && field.Array {
			fieldType = "string"
			dataType = "byte"
		} else if _, ok := BaseTypeNames[field.Type]; !ok && field.SizeFrom == "" {
			fieldType = fmt.Sprintf("[%d]%s", field.Length, dataType)
		} else {
			fieldType = "[]" + dataType
		}
	}

	return fieldType
}

func generateField(ctx *GenFile, w io.Writer, fields []*Field, i int) {
	field := fields[i]

	//fieldName := strings.TrimPrefix(field.Name, "_")
	//fieldName = camelCaseName(fieldName)
	fieldName := field.GoName

	dataType := convertToGoType(ctx.file, field.Type)
	fieldType := dataType

	// generate length field for strings
	if field.Type == "string" && field.Length == 0 {
		fmt.Fprintf(w, "\tXXX_%sLen uint32 `struc:\"sizeof=%s\"`\n", fieldName, fieldName)
	}

	// check if it is array
	if field.Length > 0 || field.SizeFrom != "" {
		if dataType == "uint8" {
			dataType = "byte"
		}
		if dataType == "string" && field.Array {
			fieldType = "string"
			dataType = "byte"
		} else if _, ok := BaseTypeNames[field.Type]; !ok && field.SizeFrom == "" {
			fieldType = fmt.Sprintf("[%d]%s", field.Length, dataType)
		} else {
			fieldType = "[]" + dataType
		}
	}
	fmt.Fprintf(w, "\t%s %s", fieldName, fieldType)

	fieldTags := map[string]string{}

	if field.Length > 0 && field.Array {
		// fixed size array
		fieldTags["struc"] = fmt.Sprintf("[%d]%s", field.Length, dataType)
	} else {
		for _, f := range fields {
			if f.SizeFrom == field.Name {
				// variable sized array
				//sizeOfName := camelCaseName(f.Name)
				fieldTags["struc"] = fmt.Sprintf("sizeof=%s", f.GoName)
			}
		}
	}

	if ctx.IncludeBinapiNames {
		typ := fromApiType(field.Type)
		if field.Array {
			if field.Length > 0 {
				fieldTags["binapi"] = fmt.Sprintf("%s[%d],name=%s", typ, field.Length, field.Name)
			} else if field.SizeFrom != "" {
				fieldTags["binapi"] = fmt.Sprintf("%s[%s],name=%s", typ, field.SizeFrom, field.Name)
			}
		} else {
			fieldTags["binapi"] = fmt.Sprintf("%s,name=%s", typ, field.Name)
		}
	}
	if limit, ok := field.Meta["limit"]; ok && limit.(int) > 0 {
		fieldTags["binapi"] = fmt.Sprintf("%s,limit=%d", fieldTags["binapi"], limit)
	}
	if def, ok := field.Meta["default"]; ok && def != nil {
		actual := getActualType(ctx.file, fieldType)
		if t, ok := binapiTypes[actual]; ok && t != "float64" {
			defnum := int(def.(float64))
			fieldTags["binapi"] = fmt.Sprintf("%s,default=%d", fieldTags["binapi"], defnum)
		} else {
			fieldTags["binapi"] = fmt.Sprintf("%s,default=%v", fieldTags["binapi"], def)
		}
	}

	fieldTags["json"] = fmt.Sprintf("%s,omitempty", field.Name)

	if len(fieldTags) > 0 {
		fmt.Fprintf(w, "\t`")
		var keys []string
		for k := range fieldTags {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		var n int
		for _, tt := range keys {
			t, ok := fieldTags[tt]
			if !ok {
				continue
			}
			if n > 0 {
				fmt.Fprintf(w, " ")
			}
			n++
			fmt.Fprintf(w, `%s:"%s"`, tt, t)
		}
		fmt.Fprintf(w, "`")
	}

	fmt.Fprintln(w)
}

func generateMessageResetMethod(w io.Writer, structName string) {
	fmt.Fprintf(w, "func (m *%[1]s) Reset() { *m = %[1]s{} }\n", structName)
}

func generateMessageNameGetter(w io.Writer, structName, msgName string) {
	fmt.Fprintf(w, "func (*%s) GetMessageName() string {	return %q }\n", structName, msgName)
}

func generateTypeNameGetter(w io.Writer, structName, msgName string) {
	fmt.Fprintf(w, "func (*%s) GetTypeName() string { return %q }\n", structName, msgName)
}

func generateCrcGetter(w io.Writer, structName, crc string) {
	crc = strings.TrimPrefix(crc, "0x")
	fmt.Fprintf(w, "func (*%s) GetCrcString() string { return %q }\n", structName, crc)
}

func generateMessageTypeGetter(w io.Writer, structName string, msgType MessageType) {
	fmt.Fprintf(w, "func (*"+structName+") GetMessageType() api.MessageType {")
	if msgType == requestMessage {
		fmt.Fprintf(w, "\treturn api.RequestMessage")
	} else if msgType == replyMessage {
		fmt.Fprintf(w, "\treturn api.ReplyMessage")
	} else if msgType == eventMessage {
		fmt.Fprintf(w, "\treturn api.EventMessage")
	} else {
		fmt.Fprintf(w, "\treturn api.OtherMessage")
	}
	fmt.Fprintln(w, "}")
	fmt.Fprintln(w)
}
