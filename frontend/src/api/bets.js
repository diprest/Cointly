import axios from './axios';

export const getUserBets = async (userId) => {
  const response = await axios.get(`/bets?user_id=${userId}`);
  return response.data;
};

export const createBet = async (bet) => {
  const response = await axios.post('/bets', bet);
  return response.data;
};

export const resetBets = async (userId) => {
  const response = await axios.post(`/bets/reset?user_id=${userId}`);
  return response.data;
};
