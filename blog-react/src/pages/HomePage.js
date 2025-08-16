import React from 'react';
import { List, Card, Button, Typography, Space } from 'antd';
import { EyeOutlined, MessageOutlined, UserOutlined } from '@ant-design/icons';

const { Title, Paragraph } = Typography;

const mockArticles = [
  {
    id: 1,
    title: '我的第一篇博客文章',
    content: '这是文章的摘要内容...',
    author: '张三',
    createdAt: '2023-07-01',
    views: 128,
    comments: 12
  },
  {
    id: 2,
    title: 'React学习心得',
    content: '这是文章的摘要内容...',
    author: '李四',
    createdAt: '2023-07-02',
    views: 96,
    comments: 8
  },
  {
    id: 3,
    title: 'Go语言并发编程',
    content: '这是文章的摘要内容...',
    author: '王五',
    createdAt: '2023-07-03',
    views: 256,
    comments: 15
  }
];

const HomePage = () => {
  return (
    <div>
      <Title level={2}>博客文章</Title>
      <List
        grid={{ gutter: 16, column: 1 }}
        dataSource={mockArticles}
        renderItem={item => (
          <List.Item>
            <Card 
              title={<Title level={4}>{item.title}</Title>}
              extra={<Button  type="primary">阅读更多</Button>}
            >
              <Paragraph>{item.content}</Paragraph>
              <Space>
                <span><UserOutlined /> {item.author}</span>
                <span>{item.createdAt}</span>
                <span><EyeOutlined /> {item.views}</span>
                <span><MessageOutlined /> {item.comments}</span>
              </Space>
            </Card>
          </List.Item>
        )}
      />
    </div>
  );
};

export default HomePage;