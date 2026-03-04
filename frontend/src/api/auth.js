import axios from './axios';

export const authApi = {
  login: async (login, password) => {
    const response = await axios.post('/auth/login', { login, password });
    return response.data;
  },
  register: async (login, password) => {
    const response = await axios.post('/auth/register', { login, password });
    return response.data;
  },
  changeLogin: async (newLogin) => {
    const response = await axios.post('/auth/user/change-login', { new_login: newLogin });
    return response.data;
  },
  changePassword: async (oldPassword, newPassword) => {
    const response = await axios.post('/auth/user/change-password', { old_password: oldPassword, new_password: newPassword });
    return response.data;
  },
};
