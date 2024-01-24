package code_generator

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
)

// CodeGenerator 提供了代码生成的功能
type CodeGenerator struct {
	template     *template.Template
	templateData *TemplateData
	templateFile []byte
	funcMap      template.FuncMap
}

// NewCodeGenerator 创建一个新的 CodeGenerator 实例
func NewCodeGenerator(templateData *TemplateData, templateFile []byte,
	funcMap template.FuncMap) (*CodeGenerator, error) {

	return &CodeGenerator{
		template:     template.New(templateData.TableName),
		templateData: templateData,
		templateFile: templateFile,
		funcMap:      funcMap,
	}, nil
}

// Generate 生成代码文件
func (g *CodeGenerator) Generate(destPath string) error {

	buffer := bytes.NewBuffer([]byte{})
	t, err := g.template.Funcs(g.funcMap).Parse(string(g.templateFile))
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	err = t.Execute(buffer, g.templateData)
	if err != nil {
		err = fmt.Errorf(fmt.Sprintf("Execute tp fail err:%v", err))
		return err
	}

	file, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
