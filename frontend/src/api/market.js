import axios from './axios';

export const getSymbols = async () => {
  const { data } = await axios.get('/market/symbols');
  return data;
};

export const getTicker = async (symbol) => {
  const { data } = await axios.get('/market/ticker', { params: { symbol } });
  return data;
};
