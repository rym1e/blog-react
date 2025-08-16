// src/pages/HomePage.js
import React, { useState, useEffect } from 'react';
import { List, Card, Button, Typography, Space, Spin, message } from 'antd';
import { EyeOutlined, MessageOutlined, UserOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { getArticles } from '../api/articleService';

const { Title, Paragraph } = Typography;

const HomePage = () => {
  const [articles, setArticles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    total_pages: 0
  });

  useEffect(() => {
    fetchArticles();
  }, []);

  const fetchArticles = async () => {
    try {
      setLoading(true);
      const response = await getArticles({
        page: pagination.page,
        limit: pagination.limit
      });
      
      if (response.success) {
        setArticles(response.data.articles);
        setPagination({
          ...pagination,
          ...response.data.pagination
        });
      } else {
        message.error(response.message || '获取文章列表失败');
      }
    } catch (error) {
      message.error('获取文章列表失败');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" />
      </div>
    );
  }

  return (
    <div>
      <Title level={2}>博客文章</Title>
      <List
        grid={{ gutter: 16, column: 1 }}
        dataSource={articles}
        renderItem={item => (
          <List.Item>
            <Card 
              title={<Title level={4}>{item.title}</Title>}
              extra={
                <Link to={`/article/${item.id}`}>
                  <Button type="primary">阅读更多</Button>
                </Link>
              }
            >
              <Paragraph ellipsis={{ rows: 2 }}>{item.content}</Paragraph>
              <Space>
                <span><UserOutlined /> {item.author?.username || item.author}</span>
                <span>{new Date(item.created_at).toLocaleDateString()}</span>
                <span><EyeOutlined /> {item.views}</span>
                <span><MessageOutlined /> {item.comments_count || 0}</span>
              </Space>
            </Card>
          </List.Item>
        )}
      />
    </div>
  );
};

export default HomePage;