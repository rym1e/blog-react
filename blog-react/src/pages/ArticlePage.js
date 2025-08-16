// src/pages/ArticlePage.js
import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { Card, Typography, List, Input, Button, message, Space, Avatar, Spin, Comment } from 'antd';
import { UserOutlined, EyeOutlined, MessageOutlined } from '@ant-design/icons';
import ReactMarkdown from 'react-markdown';
import { getArticle, getComments, createComment } from '../api/articleService';

const { Title, Paragraph } = Typography;
const { TextArea } = Input;

const ArticlePage = () => {
  const { id } = useParams();
  const [article, setArticle] = useState(null);
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(true);
  const [commentLoading, setCommentLoading] = useState(false);
  const [newComment, setNewComment] = useState('');

  useEffect(() => {
    fetchArticle();
    fetchComments();
  }, [id]);

  const fetchArticle = async () => {
    try {
      setLoading(true);
      const response = await getArticle(id);
      if (response.success) {
        setArticle(response.data);
      } else {
        message.error(response.message || '获取文章失败');
      }
    } catch (error) {
      message.error('获取文章失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchComments = async () => {
    try {
      const response = await getComments(id);
      if (response.success) {
        setComments(response.data.comments);
      } else {
        message.error(response.message || '获取评论失败');
      }
    } catch (error) {
      message.error('获取评论失败');
    }
  };

  const handleCommentSubmit = async () => {
    if (!newComment.trim()) {
      message.warning('请输入评论内容');
      return;
    }

    try {
      setCommentLoading(true);
      const response = await createComment(id, { content: newComment });
      if (response.success) {
        message.success('评论发表成功');
        setNewComment('');
        fetchComments(); // 重新获取评论列表
      } else {
        message.error(response.message || '评论发表失败');
      }
    } catch (error) {
      message.error('评论发表失败');
    } finally {
      setCommentLoading(false);
    }
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  if (!article) {
    return <div>文章不存在</div>;
  }

  return (
    <div>
      <Card>
        <Title level={2}>{article.title}</Title>
        <Space style={{ marginBottom: '20px' }}>
          <span><UserOutlined /> {article.author?.username || article.author}</span>
          <span>{new Date(article.created_at).toLocaleString()}</span>
          <span><EyeOutlined /> {article.views}</span>
          <span><MessageOutlined /> {comments.length}</span>
        </Space>
        <Paragraph>
          <ReactMarkdown>{article.content}</ReactMarkdown>
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
          loading={commentLoading}
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
                    <span style={{ fontWeight: 'bold', marginRight: '8px' }}>{item.author?.username || item.author}</span>
                    <span style={{ fontSize: '12px', color: '#999' }}>{new Date(item.created_at).toLocaleString()}</span>
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