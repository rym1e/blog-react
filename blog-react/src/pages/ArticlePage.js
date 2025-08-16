import React, { useState } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Typography, List, Input, Button, message, Space, Avatar } from 'antd';
import { UserOutlined, EyeOutlined, MessageOutlined } from '@ant-design/icons';
import ReactMarkdown from 'react-markdown';

const { Title, Paragraph } = Typography;
const { TextArea } = Input;

// 模拟文章数据
const mockArticle = {
  id: 1,
  title: '我的第一篇博客文章',
  content: `这是一篇示例文章的内容。在这里可以写任何你想写的内容。

## 标题二

这是一个段落，可以包含**粗体**和*斜体*文字。

1. 这是一个有序列表项
2. 这是另一个有序列表项

> 这是一个引用块

\`\`\`javascript
// 这是一段代码示例
function hello() {
  console.log('Hello, world!');
}
\`\`\``,
  author: '张三',
  createdAt: '2023-07-01 12:00:00',
  views: 128
};

// 模拟评论数据
const mockComments = [
  {
    id: 1,
    content: '这是一篇很棒的文章，学到了很多！',
    author: '李四',
    createdAt: '2023-07-01 15:30:00'
  },
  {
    id: 2,
    content: '感谢分享，期待更多好文章。',
    author: '王五',
    createdAt: '2023-07-02 09:15:00'
  }
];

const ArticlePage = () => {
  const { id } = useParams(); // 保留 useParams 以保持路由功能一致性
  const [comments, setComments] = useState(mockComments);
  const [newComment, setNewComment] = useState('');

  const handleCommentSubmit = () => {
    if (!newComment.trim()) {
      message.warning('请输入评论内容');
      return;
    }

    const comment = {
      id: comments.length + 1,
      content: newComment,
      author: '当前用户',
      createdAt: new Date().toLocaleString()
    };

    setComments([comment, ...comments]);
    setNewComment('');
    message.success('评论发表成功');
  };

  return (
    <div>
      <Card>
        <Title level={2}>{mockArticle.title}</Title>
        <Space style={{ marginBottom: '20px' }}>
          <span><UserOutlined /> {mockArticle.author}</span>
          <span>{mockArticle.createdAt}</span>
          <span><EyeOutlined /> {mockArticle.views}</span>
          <span><MessageOutlined /> {comments.length}</span>
        </Space>
        <Paragraph>
          <ReactMarkdown>{mockArticle.content}</ReactMarkdown>
        </Paragraph>
      </Card>

      <Card title="发表评论" style={{ marginTop: '20px' }}>
        <TextArea
          rows={4}
          value={newComment}
          onChange={(e) => setNewComment(e.target.value)}
          placeholder="请输入您的评论..."
        />
        <Button
          type="primary"
          onClick={handleCommentSubmit}
          style={{ marginTop: '10px' }}
        >
          发表评论
        </Button>
      </Card>

      <Card title="评论列表" style={{ marginTop: '20px' }}>
        <List
          dataSource={comments}
          renderItem={item => (
            <li key={item.id} style={{ marginBottom: '16px' }}>
              <div style={{ display: 'flex', marginBottom: '8px' }}>
                <Avatar icon={<UserOutlined />} style={{ marginRight: '12px' }} />
                <div>
                  <div>
                    <span style={{ fontWeight: 'bold', marginRight: '8px' }}>{item.author}</span>
                    <span style={{ fontSize: '12px', color: '#999' }}>{item.createdAt}</span>
                  </div>
                  <div style={{ marginTop: '4px' }}>{item.content}</div>
                </div>
              </div>
            </li>
          )}
        />
      </Card>
    </div>
  );
};

export default ArticlePage;