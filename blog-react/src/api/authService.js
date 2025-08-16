// src/api/authService.js
import api from './config';

// 用户注册
export const register = (userData) => {
  return api.post('/auth/register', userData);
};

// 用户登录
export const login = (credentials) => {
  return api.post('/auth/login', credentials);
};

// 获取当前用户信息
export const getCurrentUser = () => {
  return api.get('/users/me');
};

// 更新用户信息
export const updateCurrentUser = (userData) => {
  return api.put('/users/me', userData);
};