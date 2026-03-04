import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  profile: {
    username: '',
    balance: 0,
    portfolioValue: 0,
  },
  wallet: [],
  isAuth: !!localStorage.getItem('token'),
  token: localStorage.getItem('token'),
  userId: localStorage.getItem('userId'),
};

export const userSlice = createSlice({
  name: 'user',
  initialState,
  reducers: {
    updateBalance: (state, action) => {
      state.profile.balance = action.payload;
    },
    login: (state, action) => {
      state.isAuth = true;
      state.token = action.payload.token;
      state.userId = action.payload.userId;
      state.profile.userId = action.payload.userId;
      localStorage.setItem('token', action.payload.token);
      localStorage.setItem('userId', action.payload.userId);
    },
    logout: (state) => {
      state.isAuth = false;
      state.token = null;
      state.userId = null;
      localStorage.removeItem('token');
      localStorage.removeItem('userId');
    },
  },
});

export const { updateBalance, login, logout } = userSlice.actions;
export default userSlice.reducer;
