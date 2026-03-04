import axios from './axios';

export const getOrders = async (userId) => {
  const response = await axios.get(`/trading/orders?user_id=${userId}`);
  return response.data;
};

export const createOrder = async (order) => {
  const response = await axios.post('/trading/orders', order);
  return response.data;
};

export const cancelOrder = async (orderId, userId) => {
  const response = await axios.delete(`/trading/orders/${orderId}?user_id=${userId}`);
  return response.data;
};

export const resetOrders = async (userId) => {
  const response = await axios.post(`/trading/reset?user_id=${userId}`);
  return response.data;
};
