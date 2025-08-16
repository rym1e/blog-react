import React, { useState } from 'react';
import { Form, Input, Button, Card, message, Typography } from 'antd';
import ReactMarkdown from 'react-markdown';

const { Title } = Typography;

const CreateArticlePage = () => {
  const [loading, setLoading] = useState(false);
  const [contentPreview, setContentPreview] = useState('');

  const onFinish = (values) => {
    setLoading(true);
    // 模拟发布文章请求
    setTimeout(() => {
      setLoading(false);
      message.success('文章发布成功');
      console.log('Received values of form: ', values);
    }, 1000);
  };

  const handleContentChange = (e) => {
    setContentPreview(e.target.value);
  };

  return (
    <div>
      <Title level={2}>发布文章</Title>
      <Card>
        <Form
          name="create-article"
          onFinish={onFinish}
          layout="vertical"
        >
          <Form.Item
            label="文章标题"
            name="title"
            rules={[{ required: true, message: '请输入文章标题!' }]}
          >
            <Input placeholder="请输入文章标题" />
          </Form.Item>
          <Form.Item
            label="文章内容"
            name="content"
            rules={[{ required: true, message: '请输入文章内容!' }]}
          >
            <Input.TextArea
              rows={10}
              placeholder="请输入文章内容，支持Markdown语法"
              onChange={handleContentChange}
            />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={loading}>
              发布文章
            </Button>
          </Form.Item>
        </Form>
      </Card>
      
      {contentPreview && (
        <Card title="内容预览" style={{ marginTop: '20px' }}>
          <ReactMarkdown>{contentPreview}</ReactMarkdown>
        </Card>
      )}
    </div>
  );
};

export default CreateArticlePage;