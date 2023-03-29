package main

// Writes the header for the source file
func (g *Generator) WriteHeader() {
	g.buf.WriteString("// Code generated by parsley for scanning JSON strings. DO NOT EDIT.\n")
	g.buf.WriteString("package " + g.pkgName + "\n\n")
}

// Writes imports for the source file
func (g *Generator) WriteImports(pkgs map[string]string) {
	g.buf.WriteString("import (\n")
	for name, path := range pkgs {
		g.buf.WriteString(name + path + "\n")
	}
	g.buf.WriteString(")\n\n")
	g.buf.WriteString("var _ *reader.Reader\n")
	g.buf.WriteString("var _ *writer.Writer\n\n")
}

// Implements functions for a define and writes it to the buffer
func (g *Generator) WriteDefine(df Define) {
	name := df.name
	di := NewDefineInfo(df)

	code := "func (o *" + name + ")UnmarshalParsleyJSON(r *reader.Reader) (err error) {\n" +
		createUnmarshalDefineBody(di) +
		"    return\n" +
		"}\n\n" +
		"func (o *" + name + ")sequenceParsleyJSON(r *reader.Reader, idx int) (res []" + name + ", err error) {\n" +
		"    var e " + name + "\n" +
		"    if err = e.UnmarshalParsleyJSON(r); err == nil {\n" +
		"        if !r.Next() {\n" +
		"            res = make([]" + name + ", idx+1)\n" +
		"            res[idx] = e\n" +
		"            return\n" +
		"        } else if res, err = o.sequenceParsleyJSON(r, idx + 1); err == nil {\n" +
		"            res[idx] = e\n" +
		"        }\n" +
		"    }\n" +
		"    return\n" +
		"}\n\n" +
		"func (o *" + name + ") UnmarshalParsleyJSONSlice(r *reader.Reader) (res []" + name + ", err error) {\n" +
		"    if err = r.OpenArray(); err == nil {\n" +
		"        if res, err = o.sequenceParsleyJSON(r, 0); err == nil {\n" +
		"            err = r.CloseArray()\n" +
		"        }\n" +
		"    }\n" +
		"    return\n" +
		"}\n\n" +
		"func (o *" + name + ") MarshalParsleyJSON(dst []byte) (ln int) {\n" +
		"    if o == nil {\n" +
		"        return writer.WriteNull(dst)\n" +
		"    }\n" +
		"    return " + createMarshalDefineBody(di) + "\n" +
		"}\n\n" +
		"func (o *" + name + ") MarshalParsleyJSONSlice(dst []byte, slc []" + name + ") (ln int) {\n" +
		"    if slc == nil {\n" +
		"        return writer.WriteNull(dst)\n" +
		"    }\n" +
		"    dst[0] = '['\n" +
		"    ln++\n" +
		"    if len(slc) > 0 {\n" +
		"        ln += slc[0].MarshalParsleyJSON(dst[1:])\n" +
		"        for _, o := range slc[1:] {\n" +
		"            dst[ln] = ','\n" +
		"            ln++\n" +
		"            ln += o.MarshalParsleyJSON(dst[ln:])\n" +
		"        }\n" +
		"    }\n" +
		"    dst[ln] = ']'\n" +
		"    return ln + 1\n" +
		"}\n\n"

	g.buf.WriteString(code)
}

func (g *Generator) WriteStruct(st Struct) {
	name := st.name
	fis := []FieldInfo{}
	for _, f := range st.fields {
		fis = append(fis, NewFieldInfo(f, g.defaultCase))
	}

	code := "func (o *" + name + ")UnmarshalParsleyJSON(r *reader.Reader) (err error) {\n" +
		"    var key []byte\n" +
		"    err = r.OpenObject()\n" +
		"    if r.GetType() != reader.TerminatorToken {\n" +
		"        for err == nil {\n" +
		"            if key, err = r.GetKey(); err == nil {\n" +
		"                if r.IsNull() {\n" +
		"                    r.SkipNull()\n" +
		"                } else {\n" +
		"                    switch string(key) {\n" +
		createUnmarshalStructBody(fis) +
		"                    default:\n" +
		"                        err = r.Skip()\n" +
		"                    }\n" +
		"                }\n" +
		"                if err == nil && !r.Next() {\n" +
		"                    break\n" +
		"                }\n" +
		"            }\n" +
		"        }\n" +
		"    }\n" +
		"    if err == nil {\n" +
		"        err = r.CloseObject()\n" +
		"    }\n" +
		"    return\n" +
		"}\n\n" +
		"func (o *" + name + ")sequenceParsleyJSON(r *reader.Reader, idx int) (res []" + name + ", err error) {\n" +
		"    var e " + name + "\n" +
		"    if err = e.UnmarshalParsleyJSON(r); err == nil {\n" +
		"        if !r.Next() {\n" +
		"            res = make([]" + name + ", idx+1)\n" +
		"            res[idx] = e\n" +
		"            return\n" +
		"        } else if res, err = o.sequenceParsleyJSON(r, idx + 1); err == nil {\n" +
		"            res[idx] = e\n" +
		"        }\n" +
		"    }\n" +
		"    return\n" +
		"}\n\n" +
		"func (o *" + name + ") UnmarshalParsleyJSONSlice(r *reader.Reader) (res []" + name + ", err error) {\n" +
		"    if err = r.OpenArray(); err == nil {\n" +
		"        if res, err = o.sequenceParsleyJSON(r, 0); err == nil {\n" +
		"            err = r.CloseArray()\n" +
		"        }\n" +
		"    }\n" +
		"    return\n" +
		"}\n\n" +
		"func (o *" + name + ") MarshalParsleyJSON(dst []byte) (ln int) {\n" +
		"    if o == nil {\n" +
		"        return writer.WriteNull(dst)\n" +
		"    }\n" +
		"    off := 1\n" +
		"    dst[0] = '{'\n" +
		"    ln++\n" +
		createMarshalStructBody(fis) +
		"    dst[ln] = '}'\n" +
		"    ln++\n" +
		"    return ln\n" +
		"}\n\n" +
		"func (o *" + name + ") MarshalParsleyJSONSlice(dst []byte, slc []" + name + ") (ln int) {\n" +
		"    if slc == nil {\n" +
		"        return writer.WriteNull(dst)\n" +
		"    }\n" +
		"    dst[0] = '['\n" +
		"    ln++\n" +
		"    if len(slc) > 0 {\n" +
		"        ln += slc[0].MarshalParsleyJSON(dst[1:])\n" +
		"        for _, o := range slc[1:] {\n" +
		"            dst[ln] = ','\n" +
		"            ln++\n" +
		"            ln += o.MarshalParsleyJSON(dst[ln:])\n" +
		"        }\n" +
		"    }\n" +
		"    dst[ln] = ']'\n" +
		"    return ln + 1\n" +
		"}\n\n"

	g.buf.WriteString(code)
}
