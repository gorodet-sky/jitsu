// @Libs
import React, {memo, useCallback, useMemo, useState} from 'react';
import { Form, FormInstance, Input } from 'antd';
import { set } from 'lodash';
import MonacoEditor from 'react-monaco-editor';
// @Types
import { SourceFormConfigFieldProps as Props } from './SourceForm.types';
// @Components
import { LabelWithTooltip } from "../../../../../../lib/components/components";
import * as monacoEditor from "monaco-editor";

const SourceFormConfigFieldComponent = ({ displayName, initialValue, required, id, type, documentation }: Props) => {
  const fieldName = useMemo(() => `config.${id}`, [id]);

  const [jsonValue, setJsonValue] = useState<string>();

  const handleChange = useCallback(
    (getFieldsValue, setFieldsValue) => (e: React.ChangeEvent<HTMLInputElement>) => {
      const formValues = getFieldsValue();
      const value = e.target.value.replace(/\D/g, '');

      set(formValues, fieldName, value);

      setFieldsValue(formValues);
    },
    [fieldName]
  );

  const handleMonacoChange = useCallback((getFieldsValue, setFieldsValue) => (value: string, e: monacoEditor.editor.IModelContentChangedEvent) => {
      setJsonValue(JSON.stringify(value));

      const formValues = getFieldsValue();

      set(formValues, fieldName, value);

      setFieldsValue(formValues);
  }, [fieldName]);

  const formItemChild = useCallback(
    (getFieldsValue, setFieldsValue) => {
      switch (type) {
      case 'string':
      default:
        return <Input />;

      case 'json':
        return <MonacoEditor
            height="300"
            language="json"
            theme="vs-dark"
            options={{
                selectOnLineNumbers: true,
                lineNumbers: 'off'
            }}
            onChange={handleMonacoChange(getFieldsValue, setFieldsValue)}
        />;

      case 'int':
        return <Input onChange={handleChange(getFieldsValue, setFieldsValue)} />;
      }
    },
    [type, handleChange]
  );

  return (
    <div className="test-test">
      <Form.Item noStyle shouldUpdate={(prevValues, currentValues) => prevValues[fieldName] !== currentValues[fieldName]}>
        {({ getFieldsValue, setFieldsValue }: FormInstance) => (
          <Form.Item
            initialValue={initialValue}
            className="form-field_fixed-label"
            label={documentation ? <LabelWithTooltip documentation={documentation}>{displayName}</LabelWithTooltip> : <span className="field-label">{displayName}</span>}
            name={fieldName}
            rules={required ? [{ required, message: `${displayName} is required` }] : undefined}
          >
            {formItemChild(getFieldsValue, setFieldsValue)}
          </Form.Item>
        )}
      </Form.Item>
    </div>
  );
};

SourceFormConfigFieldComponent.displayName = 'SourceFormConfigField';

export const SourceFormConfigField = memo(SourceFormConfigFieldComponent);
