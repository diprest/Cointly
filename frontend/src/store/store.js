import { configureStore } from '@reduxjs/toolkit';
import userReducer from './userSlice';
import marketReducer from './marketSlice';

export const store = configureStore({
  reducer: {
    user: userReducer,
    market: marketReducer,
  },
});
