import axios from './axios';

export const getBalance = async (userId, asset = 'USDT') => {
  const response = await axios.get('/portfolio/balance', {
    params: { user_id: userId, asset },
  });
  return response.data;
};

export const getPortfolio = async (userId) => {
  const response = await axios.get(`/portfolio/list?user_id=${userId}`);
  return response.data;
};

export const resetPortfolio = async (userId) => {
  const response = await axios.post(`/portfolio/reset?user_id=${userId}`);
  return response.data;
};
