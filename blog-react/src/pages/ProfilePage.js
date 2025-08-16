import React from 'react';
import { Card, Typography, Button, List, Avatar } from 'antd';
import { EditOutlined, DeleteOutlined } from '@ant-design/icons';

const { Title, Text } = Typography;

// 模拟用户数据
const user = {
  username: '张三',
  email: 'zhangsan@example.com',
  joinDate: '2023-01-01'
};

// 模拟用户文章数据
const userArticles = [
  {
    id: 1,
    title: '我的第一篇博客文章',
    createdAt: '2023-07-01'
  },
  {
    id: 2,
    title: 'React学习心得',
    createdAt: '2023-07-02'
  }
];

const ProfilePage = () => {
  return (
    <div>
      <Card>
        <Title level={2}>个人资料</Title>
        <Text strong>用户名: </Text> {user.username}<br/>
        <Text strong>邮箱: </Text> {user.email}<br/>
        <Text strong>注册时间: </Text> {user.joinDate}<br/>
        <Button type="primary" style={{ marginTop: '20px' }}>编辑资料</Button>
      </Card>

      <Card title="我的文章" style={{ marginTop: '20px' }}>
        <List
          itemLayout="horizontal"
          dataSource={userArticles}
          renderItem={item => (
            <List.Item
              actions={[
                <Button type="link" icon={<EditOutlined />}>编辑</Button>,
                <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
              ]}
            >
              <List.Item.Meta
                title={<a href={`/article/${item.id}`}>{item.title}</a>}
                description={item.createdAt}
              />
            </List.Item>
          )}
        />
        {userArticles.length === 0 && <Text>暂无文章</Text>}
      </Card>
    </div>
  );
};

export default ProfilePage;