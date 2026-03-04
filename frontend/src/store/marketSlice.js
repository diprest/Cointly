import { createSlice } from '@reduxjs/toolkit';

const initialState = {
  coins: [],
  loading: false,
};

export const marketSlice = createSlice({
  name: 'market',
  initialState,
  reducers: {
    setCoins: (state, action) => {
      state.coins = action.payload;
    },
  },
});

export const { setCoins } = marketSlice.actions;
export default marketSlice.reducer;
