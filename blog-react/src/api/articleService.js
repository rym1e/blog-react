// src/api/articleService.js
import api from './config';

// 获取文章列表
export const getArticles = (params = {}) => {
  return api.get('/articles', { params });
};

// 获取文章详情
export const getArticle = (id) => {
  return api.get(`/articles/${id}`);
};

// 创建文章
export const createArticle = (articleData) => {
  return api.post('/articles', articleData);
};

// 更新文章
export const updateArticle = (id, articleData) => {
  return api.put(`/articles/${id}`, articleData);
};

// 删除文章
export const deleteArticle = (id) => {
  return api.delete(`/articles/${id}`);
};

// 获取文章评论
export const getComments = (articleId, params = {}) => {
  return api.get(`/articles/${articleId}/comments`, { params });
};

// 发表评论
export const createComment = (articleId, commentData) => {
  return api.post(`/articles/${articleId}/comments`, commentData);
};

// 删除评论
export const deleteComment = (id) => {
  return api.delete(`/comments/${id}`);
};